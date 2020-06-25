package interfaces

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sslab-archive/key_custody_provider/application"
	"github.com/sslab-archive/key_custody_provider/domain/entity"
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	"github.com/sslab-archive/key_custody_provider/domain/service"
	"net/http"
	"strconv"
	"time"
)

type PhoneAuthentication struct {
	repository repository.AuthenticationRepository
	authApp    application.AuthenticationApp
	userApp    application.UserApp
	keyService service.KeyManagementService
}

func NewPhoneAuthentication(
	repository repository.AuthenticationRepository, keyService service.KeyManagementService,
	authApp application.AuthenticationApp,userApp application.UserApp) *PhoneAuthentication {
	return &PhoneAuthentication{
		repository: repository,
		authApp:    authApp,
		userApp:    userApp,
		keyService: keyService,
	}
}

func (au *PhoneAuthentication) StartAuthenticationPage(c *gin.Context) {
	// check required params
	//queryParams := c.Request.URL.Query()
	//requiredParams := []string{"partial_key", "user_public_key", "redirect_url"}
	//for _, requiredParam := range requiredParams {
	//	if _, found := queryParams[requiredParam]; !found {
	//		c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
	//		return
	//	}
	//}
	//
	//_, err := hex.DecodeString(queryParams.Get("partial_key"))
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, "err : "+err.Error())
	//	return
	//}

	// process
	c.HTML(http.StatusOK, "authentication_phone.tmpl", nil)
}

func (au *PhoneAuthentication) SendVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	requiredParams := []string{"phone",}
	for _, requiredParam := range requiredParams {
		if _, found := queryParams[requiredParam]; !found {
			c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
			return
		}
	}
	phone := queryParams.Get("phone")
	code, err := au.authApp.SendVerificationCode(phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err : "+err.Error())
		return
	}

	existedAuthentication, err := au.repository.GetAuthenticationByPayload(phone)
	if err != nil {
		au.repository.DeleteAuthentication(existedAuthentication.ID)
	}

	authentication := entity.Authentication{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Payload:    phone,
		AuthCode:   code,
		IsVerified: false,
	}
	_, _ = au.repository.SaveAuthentication(&authentication)
	c.JSON(http.StatusOK, nil)
	return
}

func (au *PhoneAuthentication) CheckVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	fmt.Println(queryParams)
	// check purpose params
	purpose := queryParams.Get("purpose")
	phone := queryParams.Get("phone")
	code := queryParams.Get("code")
	if purpose == "" || phone == "" || code == "" {
		c.JSON(http.StatusBadRequest, "query param required : "+"purpose/phone/code")
		return
	}

	if purpose == "encrypt" {
		requiredParams := []string{"phone", "code", "user_public_key","partial_key", "partial_key_index", "purpose"}
		for _, requiredParam := range requiredParams {
			if _, found := queryParams[requiredParam]; !found {
				c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
				return
			}
		}

		phone, code := queryParams.Get("phone"), queryParams.Get("code")
		partialKey, partialKeyIndex := queryParams.Get("partial_key"), queryParams.Get("partial_key_index")
		userPublicKey := queryParams.Get("user_public_key")

		err := au.authApp.CheckVerificationCode(phone, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		pubKey, privKey := au.keyService.GetServerPublicKey(), au.keyService.GetServerPrivateKey()

		encryptedPayload, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pubKey, []byte(phone), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		bArr, err := hex.DecodeString(partialKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		encryptedPartialKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pubKey, bArr, nil)
		rawData, _ := json.Marshal(gin.H{
			"encrypted_payload":     hex.EncodeToString(encryptedPayload),
			"encrypted_partial_key": hex.EncodeToString(encryptedPartialKey),
			"credential_type":       "phone",
		})
		h := sha256.New()
		h.Write(rawData)

		signData, err := rsa.SignPKCS1v15(rand.Reader, &privKey, crypto.SHA256, h.Sum(nil))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		intPartialKeyIndex,_ := strconv.ParseUint(partialKeyIndex, 10, 64)
		err = au.userApp.CreateUser(userPublicKey,partialKey,phone,intPartialKeyIndex)
		if err != nil{
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"encrypted_payload":     hex.EncodeToString(encryptedPayload),
			"payload":               phone,
			"encrypted_partial_key": hex.EncodeToString(encryptedPartialKey),
			"partial_key":           partialKey,
			"partial_key_index":     partialKeyIndex,
			"credential_type":       "phone",
			"provider_id":           1,
			"public_key":            pubKey.N.String(),
			"signed_by_private_key": signData,
			"purpose":               purpose,
		})
		return
	} else if purpose == "decrypt" {
		phone, code := queryParams.Get("phone"), queryParams.Get("code")
		err := au.authApp.CheckVerificationCode(phone, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		partialKey,partialKeyIdx,err := au.userApp.GetPartialKeyByPayload(phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payload":           phone,
			"credential_type":   "phone",
			"provider_id":       1,
			"purpose":           purpose,
			"partial_key":       partialKey,
			"partial_key_index": partialKeyIdx,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, "err : invalid purpose")
		return
	}
}

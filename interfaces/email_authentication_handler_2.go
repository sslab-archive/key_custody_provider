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

type SecondEmailAuthentication struct {
	repository repository.AuthenticationRepository
	authApp    application.AuthenticationApp
	userApp    application.UserApp
	keyService service.KeyManagementService
}

func NewSecondEmailAuthentication(
	repository repository.AuthenticationRepository, keyService service.KeyManagementService,
	authApp application.AuthenticationApp,userApp application.UserApp) *SecondEmailAuthentication {
	return &SecondEmailAuthentication{
		repository: repository,
		authApp:    authApp,
		userApp:    userApp,
		keyService: keyService,
	}
}

func (au *SecondEmailAuthentication) StartAuthenticationPage(c *gin.Context) {
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
	c.HTML(http.StatusOK, "mail_authentication_2.tmpl", nil)
}

func (au *SecondEmailAuthentication) SendVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	requiredParams := []string{"email",}
	for _, requiredParam := range requiredParams {
		if _, found := queryParams[requiredParam]; !found {
			c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
			return
		}
	}
	email := queryParams.Get("email")
	code, err := au.authApp.SendVerificationCode(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err : "+err.Error())
		return
	}

	existedAuthentication, err := au.repository.GetAuthenticationByPayload(email)
	if err != nil {
		au.repository.DeleteAuthentication(existedAuthentication.ID)
	}

	authentication := entity.Authentication{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Payload:    email,
		AuthCode:   code,
		IsVerified: false,
	}
	_, _ = au.repository.SaveAuthentication(&authentication)
	c.JSON(http.StatusOK, nil)
	return
}

func (au *SecondEmailAuthentication) CheckVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	fmt.Println(queryParams)
	// check purpose params
	purpose := queryParams.Get("purpose")
	email := queryParams.Get("email")
	code := queryParams.Get("code")
	if purpose == "" || email == "" || code == "" {
		c.JSON(http.StatusBadRequest, "query param required : "+"purpose/email/code")
		return
	}

	if purpose == "encrypt" {
		requiredParams := []string{"email", "code", "user_public_key","partial_key", "partial_key_index", "purpose"}
		for _, requiredParam := range requiredParams {
			if _, found := queryParams[requiredParam]; !found {
				c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
				return
			}
		}

		email, code := queryParams.Get("email"), queryParams.Get("code")
		partialKey, partialKeyIndex := queryParams.Get("partial_key"), queryParams.Get("partial_key_index")
		userPublicKey := queryParams.Get("user_public_key")

		err := au.authApp.CheckVerificationCode(email, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		pubKey, privKey := au.keyService.GetServerPublicKey(), au.keyService.GetServerPrivateKey()

		encryptedPayload, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pubKey, []byte(email), nil)
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
			"credential_type":       "email",
		})
		h := sha256.New()
		h.Write(rawData)

		signData, err := rsa.SignPKCS1v15(rand.Reader, &privKey, crypto.SHA256, h.Sum(nil))
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		intPartialKeyIndex,_ := strconv.ParseUint(partialKeyIndex, 10, 64)
		err = au.userApp.CreateUser(userPublicKey,partialKey,email,intPartialKeyIndex)
		if err != nil{
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"encrypted_payload":     hex.EncodeToString(encryptedPayload),
			"payload":               email,
			"encrypted_partial_key": hex.EncodeToString(encryptedPartialKey),
			"partial_key":           partialKey,
			"partial_key_index":     partialKeyIndex,
			"credential_type":       "email",
			"provider_id":           1,
			"public_key":            pubKey.N.String(),
			"signed_by_private_key": signData,
			"purpose":               purpose,
		})
		return
	} else if purpose == "decrypt" {
		email, code := queryParams.Get("email"), queryParams.Get("code")
		err := au.authApp.CheckVerificationCode(email, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}
		partialKey,partialKeyIdx,err := au.userApp.GetPartialKeyByPayload(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "err : "+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payload":           email,
			"credential_type":   "email",
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

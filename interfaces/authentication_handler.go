package interfaces

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sslab-archive/key_custody_provider/application"
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	"github.com/sslab-archive/key_custody_provider/domain/service"
	"net/http"
)

type Authentication struct {
	repository repository.AuthenticationRepository
	app        application.AuthenticationApp
	keyService service.KeyManagementService
}

func NewAuthentication(repository repository.AuthenticationRepository, autApp application.AuthenticationApp, keyService service.KeyManagementService) *Authentication {
	return &Authentication{
		repository: repository,
		app:        autApp,
		keyService: keyService,
	}
}

func (au *Authentication) StartAuthenticationPage(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	requiredParams := []string{"partial_key", "user_public_key", "redirect_url"}
	for _, requiredParam := range requiredParams {
		if _, found := queryParams[requiredParam]; !found {
			c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
			return
		}
	}

	_, err := hex.DecodeString(queryParams.Get("partial_key"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "err : "+err.Error())
		return
	}

	// process
	c.HTML(http.StatusOK, "authentication.tmpl", nil)
}

func (au *Authentication) SendVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	requiredParams := []string{"email",}
	for _, requiredParam := range requiredParams {
		if _, found := queryParams[requiredParam]; !found {
			c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
			return
		}
	}
	_, err := au.app.SendVerificationCode(queryParams.Get("email"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err : "+err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
	return
}

func (au *Authentication) CheckVerificationCodeAPI(c *gin.Context) {
	// check required params
	queryParams := c.Request.URL.Query()
	requiredParams := []string{"email", "code", "partial_key", "partial_key_index", "purpose"}
	for _, requiredParam := range requiredParams {
		if _, found := queryParams[requiredParam]; !found {
			c.JSON(http.StatusBadRequest, "query param required : "+requiredParam)
			return
		}
	}
	email, code := queryParams.Get("email"), queryParams.Get("code")
	partialKey, partialKeyIndex := queryParams.Get("partial_key"), queryParams.Get("partial_key_index")
	purpose := queryParams.Get("purpose")
	err := au.app.CheckVerificationCode(email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "err : "+err.Error())
		return
	}

	if purpose == "encrypt"{

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
	}else if purpose == "decrypt"{

	}else{
		c.JSON(http.StatusInternalServerError, "err : invalid purpose")
		return
	}
}

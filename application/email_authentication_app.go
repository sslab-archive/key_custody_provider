package application

import (
	"github.com/sslab-archive/key_custody_provider/domain/entity"
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	"github.com/sslab-archive/key_custody_provider/util"
	"log"
	"net/smtp"
	"time"
)

type EmailAuthenticationApp struct {
	authenticationRepository repository.AuthenticationRepository
}

func NewEmailAuthenticationApp(ar repository.AuthenticationRepository) *EmailAuthenticationApp {
	return &EmailAuthenticationApp{authenticationRepository: ar}
}

func (aa *EmailAuthenticationApp) SendVerificationCode(email string) (code string, err error) {
	d, err := aa.authenticationRepository.GetAuthenticationByPayload(email)
	if err != nil {
		aa.authenticationRepository.DeleteAuthentication(d.ID)
	}
	authCode := util.GetRandomString(20)
	authInfo := entity.Authentication{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Payload:    email,
		AuthCode:   authCode,
		IsVerified: false,
	}

	err = sendEmail(email, authCode)
	if err != nil {
		return "", err
	}

	_, err = aa.authenticationRepository.SaveAuthentication(&authInfo)
	if err != nil {
		return "", err
	}

	return authCode, nil
}

func (aa *EmailAuthenticationApp) CheckVerificationCode(email string, code string) error {
	return nil
}

func sendEmail(to string, authCode string) error {
	from := util.GetConfigInstance().GmailID
	pass := util.GetConfigInstance().GmailPW

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Authentication Code\n\n" +
		"Authentication code : " + authCode

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	return nil
}

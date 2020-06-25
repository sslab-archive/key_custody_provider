package application

import (
	"errors"
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
	if err == nil {
		aa.authenticationRepository.DeleteAuthentication(d.ID)
	}
	authCode := util.GetRandomString(10)
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
	// check authentication repository
	existedAuthentication , err := aa.authenticationRepository.GetAuthenticationByPayload(email)
	if err != nil{
		return errors.New("인증 내역이 존재하지 않습니다.")
	}

	// 이미 사용된것
	if existedAuthentication.IsVerified == true{
		return errors.New("이미 사용된 인증입니다.")
	}

	// 코드 확인
	if existedAuthentication.AuthCode != code{
		return errors.New("인증 코드가 잘못됬습니다.")
	}

	existedAuthentication.IsVerified=true

	_, err = aa.authenticationRepository.SaveAuthentication(&existedAuthentication)
	return err
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

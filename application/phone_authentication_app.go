package application

import (
	"errors"
	"fmt"
	"github.com/sslab-archive/key_custody_provider/domain/entity"
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	"github.com/sslab-archive/key_custody_provider/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PhoneAuthenticationApp struct {
	authenticationRepository repository.AuthenticationRepository
}

func NewPhoneAuthenticationApp(ar repository.AuthenticationRepository) *PhoneAuthenticationApp {
	return &PhoneAuthenticationApp{authenticationRepository: ar}
}

func (aa *PhoneAuthenticationApp) SendVerificationCode(phoneNumber string) (code string, err error) {
	d, err := aa.authenticationRepository.GetAuthenticationByPayload(phoneNumber)
	if err == nil {
		aa.authenticationRepository.DeleteAuthentication(d.ID)
	}
	authCode := util.GetRandomString(10)
	authInfo := entity.Authentication{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Payload:    phoneNumber,
		AuthCode:   authCode,
		IsVerified: false,
	}

	err = sendSms(phoneNumber, authCode)
	if err != nil {
		return "", err
	}

	_, err = aa.authenticationRepository.SaveAuthentication(&authInfo)
	if err != nil {
		return "", err
	}

	return authCode, nil
}

func (aa *PhoneAuthenticationApp) CheckVerificationCode(email string, code string) error {
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
	aa.authenticationRepository.DeleteAuthentication(existedAuthentication.ID)
	return nil
}

func sendSms(to string, authCode string) error {
	form := url.Values{}
	form.Add("send_phone","01049315539")
	form.Add("dest_phone", to)
	form.Add("msg_body",  "code : " + authCode)


	req, err := http.NewRequest("POST", "http://api.apistore.co.kr/ppurio/1/message/sms/thefamilylab", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("x-waple-authorization", util.GetConfigInstance().PhoneAPI)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err)
		return err
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)
	return nil
}

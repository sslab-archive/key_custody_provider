package application

type AuthenticationApp interface {
	SendVerificationCode(payload string) (code string, err error)
	CheckVerificationCode(payload string, code string) error
}

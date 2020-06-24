package repository

import "github.com/sslab-archive/key_custody_provider/domain/entity"

type AuthenticationRepository interface {
	SaveAuthentication(authentication *entity.Authentication) (*entity.Authentication, error)
	GetAuthenticationById(uint64) (entity.Authentication, error)
	GetAuthenticationByPayload(string) (entity.Authentication, error)
	GetAuthenticationByAuthCode(string) (entity.Authentication, error)
	DeleteAuthentication(uint64)
}

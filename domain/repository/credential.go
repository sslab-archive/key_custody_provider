package repository

import "github.com/sslab-archive/key_custody_provider/domain/entity"

type CredentialRepository interface {
	SaveCredential(credential *entity.Credential) (*entity.Credential, error)
	GetAuthenticationById(uint64) (entity.Authentication, error)
	GetAuthenticationByPayload(string) (entity.Authentication, error)
	GetAuthenticationByAuthCode(string) (entity.Authentication, error)
	DeleteAuthentication(uint64)
}
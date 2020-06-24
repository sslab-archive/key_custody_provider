package persistence

import (
	"github.com/sslab-archive/key_custody_provider/domain/entity"
)

type AuthenticationMemoryRepository struct {
	authentication  map[uint64]entity.Authentication // id-user object map
	autoIncreasedId uint64
}

func (ar *AuthenticationMemoryRepository) DeleteAuthentication(uint64) {
	panic("implement me")
}

func (ar *AuthenticationMemoryRepository) SaveAuthentication(au *entity.Authentication) (*entity.Authentication, error) {
	// 아이디가 없는경우 새로 생성
	if au.ID == 0 {
		ar.autoIncreasedId++
		au.ID = ar.autoIncreasedId
		// 해당 아이디로 없는경우 역시 새로 생성, 아이디 부여
	} else if _, err := ar.GetAuthenticationById(au.ID); err == nil {
		ar.autoIncreasedId++
		au.ID = ar.autoIncreasedId
	}
	ar.authentication[au.ID] = *au
	return au, nil
}

func (ar *AuthenticationMemoryRepository) GetAuthenticationById(uint64) (entity.Authentication, error) {
	panic("implement me")
}

func (ar *AuthenticationMemoryRepository) GetAuthenticationByPayload(string) (entity.Authentication, error) {
	panic("implement me")
}

func (ar *AuthenticationMemoryRepository) GetAuthenticationByAuthCode(string) (entity.Authentication, error) {
	panic("implement me")
}

func NewAuthenticationMemoryRepository() *AuthenticationMemoryRepository {
	return &AuthenticationMemoryRepository{
		authentication:  make(map[uint64]entity.Authentication),
		autoIncreasedId: 0,
	}
}

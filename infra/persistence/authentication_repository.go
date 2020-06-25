package persistence

import (
	"errors"
	"github.com/sslab-archive/key_custody_provider/domain/entity"
)

type AuthenticationMemoryRepository struct {
	authentication  map[uint64]entity.Authentication // id-user object map
	autoIncreasedId uint64
}

func (ar *AuthenticationMemoryRepository) DeleteAuthentication(id uint64) {
	delete(ar.authentication,id)
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

func (ar *AuthenticationMemoryRepository) GetAuthenticationById(id uint64) (entity.Authentication, error) {
	if data, found := ar.authentication[id]; found{
		return data, nil
	}
	return entity.Authentication{},errors.New("해당하는 authentication 을 발견하지 못했습니다.")
}

func (ar *AuthenticationMemoryRepository) GetAuthenticationByPayload(payload string) (entity.Authentication, error) {
	for _,auth := range ar.authentication{
		if auth.Payload == payload{
			return auth,nil
		}
	}
	return entity.Authentication{}, errors.New("해당하는 authentication 을 찾지 못했습니다")
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

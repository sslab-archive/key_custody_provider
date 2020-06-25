package persistence

import (
	"errors"
	"fmt"
	"github.com/sslab-archive/key_custody_provider/domain/entity"
)

type UserMemoryRepository struct {
	users map[uint64]entity.User // id-user object map
	autoIncreasedId uint64
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{users: make(map[uint64]entity.User)}
}

func (u *UserMemoryRepository) SaveUser(inputUser *entity.User) (*entity.User, error) {
	if inputUser.ID == 0 {
		u.autoIncreasedId++
		inputUser.ID = u.autoIncreasedId
		// 해당 아이디로 없는경우 역시 새로 생성, 아이디 부여
	} else if _, err := u.GetUserById(inputUser.ID); err == nil {
		u.autoIncreasedId++
		inputUser.ID = u.autoIncreasedId
	}
	u.users[inputUser.ID] = *inputUser
	return inputUser, nil
}

func (u *UserMemoryRepository) GetUserById(id uint64) (entity.User, error) {
	if user, found := u.users[id]; found {
		return user, nil
	}
	return entity.User{}, errors.New(fmt.Sprintf("user id %d not found", id))
}

func (u *UserMemoryRepository) GetUserByPubKey(pbKey string) (entity.User, error) {
	for _, v := range u.users {
		if v.PublicKey == pbKey {
			return v, nil
		}
	}
	return entity.User{}, errors.New(fmt.Sprintf("user pubkey %s not found", pbKey))
}

func (u *UserMemoryRepository) GetUsers() ([]entity.User, error) {
	var values []entity.User
	for _, value := range u.users {
		values = append(values, value)
	}
	return values, nil
}

func (u *UserMemoryRepository) DeleteUser(id uint64) {
	delete(u.users,id)
}
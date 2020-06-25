package repository

import "github.com/sslab-archive/key_custody_provider/domain/entity"

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUserById(uint64) (entity.User, error)
	GetUserByPubKey(string) (entity.User, error)
	GetUsers() ([]entity.User, error)
	DeleteUser(id uint64)
}

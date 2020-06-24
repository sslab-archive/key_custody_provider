package persistence

import (
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	_ "gorm.io/driver/postgres"
)

type Repositories struct {
	User           repository.UserRepository
	Authentication repository.AuthenticationRepository
}

func NewRepositories() (*Repositories, error) {
	return &Repositories{
		User:           NewUserMemoryRepository(),
		Authentication: NewAuthenticationMemoryRepository(),
	}, nil
}

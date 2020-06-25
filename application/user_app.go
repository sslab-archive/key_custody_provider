/* 
 * Copyright 2019 hea9549
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package application

import (
	"errors"
	"github.com/sslab-archive/key_custody_provider/domain/entity"
	"github.com/sslab-archive/key_custody_provider/domain/repository"
	"time"
)

type UserApp interface {
	CreateUser(pubKey, partialKey, payload string,partialKeyIndex uint64) error
	GetPartialKeyByPayload(payload string) (string, uint64, error)
}

type DefaultUserApp struct {
	userRepository repository.UserRepository
}

func NewDefaultUserApp(userRepository repository.UserRepository) *DefaultUserApp {
	return &DefaultUserApp{
		userRepository: userRepository,
	}
}

func (dua *DefaultUserApp) CreateUser(pubKey, partialKey, payload string, partialKeyIndex uint64) error {
	u, err := dua.userRepository.GetUserByPubKey(pubKey)
	if err == nil {
		dua.userRepository.DeleteUser(u.ID)
	}

	newUser := entity.User{
		ID:              0,
		PublicKey:       pubKey,
		PartialKey:      partialKey,
		PartialKeyIndex: partialKeyIndex,
		Payload:         payload,
		RegisteredTx:    "",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_, err = dua.userRepository.SaveUser(&newUser)
	return err
}

func (dua *DefaultUserApp) GetPartialKeyByPayload(payload string) (string, uint64, error) {
	users, _ := dua.userRepository.GetUsers()

	for _, u := range users {
		if u.Payload == payload {
			return u.PartialKey, u.PartialKeyIndex, nil
		}
	}
	return "", 0, errors.New("해당하는 payload 에 user 가 없습니다.")

}

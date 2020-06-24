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

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sslab-archive/key_custody_provider/application"
	"github.com/sslab-archive/key_custody_provider/domain/service"
	"github.com/sslab-archive/key_custody_provider/infra/persistence"
	"github.com/sslab-archive/key_custody_provider/interfaces"
	"github.com/sslab-archive/key_custody_provider/util"
	"log"
)

func main() {
	util.InitConfig("C:\\Users\\user\\Desktop\\sslab-archive\\key_custody_provider\\config\\server.json")
	repositories, err := persistence.NewRepositories()
	if err != nil {
		panic(err)
	}

	authApp := application.NewEmailAuthenticationApp(repositories.Authentication)
	keyService,err := service.NewStoredKeyManagementService()
	if err != nil {
		panic(err)
	}

	authentications := interfaces.NewAuthentication(repositories.Authentication, authApp, keyService)

	r := gin.Default()

	r.GET("/authentication", authentications.StartAuthenticationPage)
	r.POST("/api/authentication/send_code", authentications.SendVerificationCodeAPI)
	r.POST("/api/authentication/check", authentications.CheckVerificationCodeAPI)
	r.Static("static", "template/static")
	r.LoadHTMLGlob("template/html/*")

	log.Fatal(r.Run(":8888"))
}

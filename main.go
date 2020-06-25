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


	//p1
	p1Repositories, err := persistence.NewRepositories()
	if err != nil {
		panic(err)
	}

	p1AuthApp := application.NewEmailAuthenticationApp(p1Repositories.Authentication)
	p1UserApp := application.NewDefaultUserApp(p1Repositories.User)
	keyService, err := service.NewStoredKeyManagementService()
	if err != nil {
		panic(err)
	}
	p1Authentications := interfaces.NewAuthentication(p1Repositories.Authentication, keyService, p1AuthApp, p1UserApp)

	//p2
	p2Repositories, err := persistence.NewRepositories()
	if err != nil {
		panic(err)
	}

	p2AuthApp := application.NewPhoneAuthenticationApp(p2Repositories.Authentication)
	p2UserApp := application.NewDefaultUserApp(p2Repositories.User)
	if err != nil {
		panic(err)
	}
	p2Authentications := interfaces.NewPhoneAuthentication(p2Repositories.Authentication, keyService, p2AuthApp, p2UserApp)

	// p3
	p3Repositories, err := persistence.NewRepositories()
	if err != nil {
		panic(err)
	}

	p3AuthApp := application.NewEmailAuthenticationApp(p3Repositories.Authentication)
	p3UserApp := application.NewDefaultUserApp(p3Repositories.User)
	if err != nil {
		panic(err)
	}
	p3Authentications := interfaces.NewSecondEmailAuthentication(p3Repositories.Authentication, keyService, p3AuthApp, p3UserApp)

	// router setting
	r := gin.Default()

	// p1
	r.GET("/p1/authentication", p1Authentications.StartAuthenticationPage)
	r.POST("/p1/api/authentication/send_code", p1Authentications.SendVerificationCodeAPI)
	r.POST("/p1/api/authentication/check", p1Authentications.CheckVerificationCodeAPI)

	// p2
	r.GET("/p2/authentication", p2Authentications.StartAuthenticationPage)
	r.POST("/p2/api/authentication/send_code", p2Authentications.SendVerificationCodeAPI)
	r.POST("/p2/api/authentication/check", p2Authentications.CheckVerificationCodeAPI)

	// p3
	r.GET("/p3/authentication", p3Authentications.StartAuthenticationPage)
	r.POST("/p3/api/authentication/send_code", p3Authentications.SendVerificationCodeAPI)
	r.POST("/p3/api/authentication/check", p3Authentications.CheckVerificationCodeAPI)

	r.Static("static", "template/static")
	r.LoadHTMLGlob("template/html/*")

	log.Fatal(r.Run(":8888"))
}

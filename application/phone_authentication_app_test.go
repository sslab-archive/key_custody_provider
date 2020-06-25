package application

import (
	"github.com/sslab-archive/key_custody_provider/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendSMS(t *testing.T) {
	util.InitConfig("C:\\Users\\user\\Desktop\\sslab-archive\\key_custody_provider\\config\\server.json")
	err := sendSms("01049315539","1234")
	assert.NoError(t,err,"err!")
}

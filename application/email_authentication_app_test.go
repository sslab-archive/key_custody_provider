package application

import (
	"github.com/sslab-archive/key_custody_provider/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendEmail(t *testing.T) {
	util.InitConfig("/Users/hea9549/Desktop/sslab-archive/key_custody_provider/config/server.json")
	err := sendEmail("hea9549@gmail.com","1234")
	assert.NoError(t,err,"err!")
}

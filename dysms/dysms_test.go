package dysms

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"fmt"
)

const accessKeyId = "accessKeyId"
const accessKeySecret = "accessKeySecret"

func TestNew(t *testing.T) {
	dysms := New(accessKeyId, accessKeySecret)
	assert.Equal(t, accessKeyId, dysms.accessKeyId, "accessKeyId 保存错误")
	assert.Equal(t, accessKeySecret, dysms.accessKeySecret, "accessKeySecret 保存错误")
}

func TestDysms_SendSms(t *testing.T) {
	dysms := New(os.Getenv("accessKeyId"), os.Getenv("accessKeySecret"))

	requestId, message, bizId, err := dysms.SendSms([]string{"18686861234", "18686895678"},
		"xxx", "SMS_11221122",
		TemplateParam{"name": "param1", "time": "param2"}, "")
	fmt.Println(requestId, message, bizId, err)
	assert.NoError(t, err, "sendsms error")
}

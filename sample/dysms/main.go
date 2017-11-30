package main

import (
	"github.com/zcx2001/aliyun-go/dysms"
	"fmt"
)

func main() {
	sms := dysms.New("accessKeyId", "accessKeySecret")
	requestId, message, bizId, err := sms.SendSms([]string{"18686861234", "18686895678"},
		"xxx", "SMS_11221122",
		dysms.TemplateParam{"name": "param1", "time": "param2"}, "")
	fmt.Println(requestId, message, bizId, err)
}

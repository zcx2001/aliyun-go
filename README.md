# aliyun-go
通过go语言实现阿里云服务的实现

## 短信服务（dysms）
```
sms := dysms.New("accessKeyId", "accessKeySecret")
requestId, message, bizId, err := sms.SendSms([]string{"18686861234", "18686895678"},
    "xxx", "SMS_11221122",
    dysms.TemplateParam{"name": "param1", "time": "param2"}, "")
```
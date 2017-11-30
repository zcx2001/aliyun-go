package dysms

import (
	"strings"
	"errors"
	"time"
	"github.com/satori/go.uuid"
	"encoding/json"
	"net/url"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

const (
	dysms_url              = "http://dysmsapi.aliyuncs.com/"
	dysms_format           = "JSON"
	dysms_signatureMethod  = "HMAC-SHA1"
	dysms_signatureVersion = "1.0"

	dysms_action   = "SendSms"
	dysms_version  = "2017-05-25"
	dysms_regionId = "cn-hangzhou"
)

type TemplateParam map[string]string

type Dysms struct {
	accessKeyId     string
	accessKeySecret string
}

type result struct {
	Code      string
	Message   string
	RequestId string
	BizId     string
}

func New(accessKeyId, accessKeySecret string) *Dysms {
	return &Dysms{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}

func (sms *Dysms) SendSms(phoneNumbers []string, signName, templateCode string,
	templateParam TemplateParam, outId string) (Message, RequestId, BizId string, err error) {
	if len(phoneNumbers) > 1000 {
		return "", "", "", errors.New("phoneNumbers size > 1000")
	}

	params := map[string]string{
		"AccessKeyId":      sms.accessKeyId,
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           dysms_format,
		"SignatureMethod":  dysms_signatureMethod,
		"SignatureVersion": dysms_signatureVersion,
		"SignatureNonce":   uuid.NewV1().String(),

		"Action":       dysms_action,
		"Version":      dysms_version,
		"RegionId":     dysms_regionId,
		"PhoneNumbers": strings.Join(phoneNumbers, ","),
		"SignName":     signName,
		"TemplateCode": templateCode,
	}

	if templateParam != nil {
		j, _ := json.Marshal(templateParam)
		params["TemplateParam"] = string(j)
	}

	if len(outId) > 0 {
		params["OutId"] = outId
	}

	signature := sign(sms.accessKeySecret, "GET", "/", params)
	urlParams := url.Values{}
	urlParams.Set("Signature", signature)

	for k, v := range params {
		urlParams.Set(k, v)
	}

	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/?%s", dysms_url, urlParams.Encode()), nil)
	if err != nil {
		return "", "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", "", "", err
	}

	var r result
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", "", "", err
	}

	if !strings.EqualFold("OK", r.Code) {
		return r.RequestId, r.Message, r.BizId, errors.New(r.Code)
	}

	return r.Message, r.RequestId, r.BizId, nil
}

func sign(accessKeySecret, httpMethod, httpUrl string, params map[string]string) (signature string) {
	originStr1 := httpMethod + "&" + url.QueryEscape(httpUrl) + "&"
	originStr2 := ""
	params_k := make([]string, 0)
	for k := range params {
		params_k = append(params_k, k)
	}
	sort.Strings(params_k)
	for _, k := range params_k {
		originStr2 += "&" + k + "=" + url.QueryEscape(params[k])
	}
	originStr3 := originStr1 + url.QueryEscape(originStr2[1:])

	mac := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	mac.Write([]byte(originStr3))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

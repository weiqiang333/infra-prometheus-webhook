// 发送云片短信通知
// 通知模板：
// 	申请签名: 【systemName】
// 	申请模板: 故障级别 #receiver#,当前状态 #status#,故障名称 #alertname#,故障数量#alertsum#个
// 报警 SMS 信息案例: 【systemName】故障级别 p0,当前状态 firing,故障名称 Kube Pod Container Restart,故障数量3个
// 语音 Voice 信息案例: 你的验证码是 001101. 验证码位置解释: 1-2位的00代表故障等级p0=00,p1=11,p2=22,p3=33; 3-4位的11代表当前状态firing=11,resolved=00; 5-6位的01代表故障数量最多可承载99个故障

package yunpian_clnt

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/spf13/viper"
	model2 "github.com/weiqiang333/infra-prometheus-webhook/internal/model"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/utils/http_client"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/utils/notification_process"
)

type YunPianModule struct {
	ApiKey    string `json:"apikey"`
	Mobile    string `json:"mobile"` // 手机号逗号分割
	Text      string `json:"text"`   // sms 发送内容: 【systemName】故障级别 p0,当前状态 firing,故障名称 Kube Pod Container Restart,故障数量3个
	TplId     int64  `json:"tpl_id"` // 模板编号
	VoiceCode string `json:"code"`   // 语音验证码
	Url       string `json:"url"`    // API url 地址
}

func NewYunPianModule() *YunPianModule {
	apiKey := viper.GetString("yunpian.apikey")
	mobile := viper.GetString("yunpian.sms.mobile.p0")
	return &YunPianModule{
		ApiKey:    apiKey,
		Mobile:    mobile,
		Text:      "【systemName】故障级别 p0,当前状态 firing,故障名称 Kube Pod Container Restart,故障数量3个",
		TplId:     0,
		VoiceCode: "1234",
		Url:       "https://sms.yunpian.com/v2/sms/tpl_single_send.json",
	}
}

func (y *YunPianModule) PostSms(notification model2.Notification, priority string) error {
	tplId := viper.GetInt64("yunpian.sms.tpl_id")
	mobile := viper.GetString(fmt.Sprintf("yunpian.sms.mobile.%s", priority))
	y.TplId = tplId
	y.Mobile = mobile
	y.Url = "https://sms.yunpian.com/v2/sms/tpl_batch_send.json"

	var receiver = "yunpianSms"
	var status string
	switch notification.Status {
	case "firing":
		status = "PROBLEM"
	case "resolved":
		status = "OK"
	}
	alertname := notification.CommonLabels["alertname"]
	description := notification_process.GetDescriptionList(notification)
	descriptionSum := len(notification.Alerts)
	summary := notification.CommonAnnotations["summary"]

	tplValue := url.Values{
		"#receiver#":  {priority},
		"#status#":    {status},
		"#alertname#": {alertname},
		"#alertsum#":  {strconv.Itoa(descriptionSum)},
	}.Encode()
	dataTplSms := url.Values{
		"apikey":    {y.ApiKey},
		"mobile":    {mobile},
		"tpl_id":    {fmt.Sprintf("%d", tplId)},
		"tpl_value": {tplValue}}
	if err := http_client.HttpsPostForm(y.Url, dataTplSms); err != nil {
		msg := fmt.Sprintln("Failed http.Post", receiver, status, priority, alertname, summary, descriptionSum, description, err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	log.Println("INFO http.Post", receiver, status, priority, alertname, summary, descriptionSum, description)
	return nil
}

func (y *YunPianModule) PostVoice(notification model2.Notification, priority string) error {
	mobile := viper.GetString(fmt.Sprintf("yunpian.voice.mobile.%s", priority))
	y.Mobile = mobile
	y.Url = "https://voice.yunpian.com/v2/voice/send.json"

	var receiver = "yunpianVoice"
	var status string
	switch notification.Status {
	case "firing":
		status = "11"
	case "resolved":
		status = "00"
	}
	alertname := notification.CommonLabels["alertname"]
	description := notification_process.GetDescriptionList(notification)
	descriptionSum := len(notification.Alerts)
	if descriptionSum > 99 {
		descriptionSum = 99
	}
	var priorityCode string
	switch priority {
	case "p0":
		priorityCode = "00"
	case "p1":
		priorityCode = "11"
	case "p2":
		priorityCode = "22"
	case "p3":
		priorityCode = "33"
	}
	y.VoiceCode = fmt.Sprintf("%s%s%02d", priorityCode, status, descriptionSum)
	dataSendVoice := url.Values{
		"apikey": {y.ApiKey},
		"mobile": {mobile},
		"code":   {y.VoiceCode}}
	if err := http_client.HttpsPostForm(y.Url, dataSendVoice); err != nil {
		msg := fmt.Sprintln("Failed http.Post", receiver, status, priority, alertname, descriptionSum, description, "VoiceCode: ", y.VoiceCode, err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	log.Println("INFO http.Post", receiver, status, priority, alertname, descriptionSum, description, "VoiceCode: ", y.VoiceCode)
	return nil
}

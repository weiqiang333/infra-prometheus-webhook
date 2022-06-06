package phonecall

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
	"github.com/weiqiang333/infra-prometheus-webhook/model"
)

// Phonecall 发送电话报警给值班人员
func Phonecall(notification model.Notification, role string) (int ,error) {
	var receiver = "phonecall"
	if ! selectTime() {
		fmt.Println("非电话报警时间段")
		// 响应 code 200,避免 Alertmanager 误以为发送失败，频发发送
		return http.StatusOK, fmt.Errorf("非电话报警时间段")
	}

	var status string
	switch notification.Status {
	case "firing":
		status = "PROBLEM"
	case "resolved":
		status = "OK"
	}
	grade := notification.CommonLabels["priority"]
	alertname := notification.CommonLabels["alertname"]

	// 获取值班用户、电话
	user = GetOncallUser(role)

	// 发送云片电话
	clnt := ypclnt.New(model.Config.PhoneCall["apikey"])
	voice := clnt.Voice()
	param := ypclnt.NewParam(2)
	param[ypclnt.MOBILE] = user.Mobile
	param[ypclnt.CODE] = "33333"
	r := voice.Send(param)
	// Result 内容: Code Msg Detail Data
	// Fail：2 请求参数格式错误 参数 code 格式不正确，code:验证码必须是数字 <nil>
	// Success: 0   map[count:1 fee:0.05 sid:86039ca95f5e47b68ec6db8034263600]
	if r.Code == 0 {
		log.Println(http.StatusOK, receiver, status, grade, alertname, role, user.UserName)
		return http.StatusOK, nil
	} else {
		log.Println(http.StatusInternalServerError, receiver, status, grade, alertname, role, user.UserName, r.Detail)
		return http.StatusInternalServerError, fmt.Errorf(r.Detail)
	}
}

// selectTime 凌晨至早8点时间确认
func selectTime() bool {
	now := time.Now().UTC()
	date, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	startTime := date.Add(time.Hour * 16)
	stopTime := date.Add(time.Hour * 24)
	if now.After(startTime) && now.Before(stopTime) {
		return true
	}
	return false
}

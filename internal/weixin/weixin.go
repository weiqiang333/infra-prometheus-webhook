package weixin

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	model2 "github.com/weiqiang333/infra-prometheus-webhook/internal/model"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/utils/notification_process"
)

// Weixin 发送企业微信消息程序
func Weixin(notification model2.Notification, priority string) error {
	var receiver = "weixin"
	var status string
	switch notification.Status {
	case "firing":
		status = "PROBLEM"
	case "resolved":
		status = "OK"
	}
	grade := notification.CommonLabels["priority"]
	alertname := notification.CommonLabels["alertname"]
	description := notification_process.GetDescriptionList(notification)
	summary := notification.CommonAnnotations["summary"]
	weixinKey := viper.GetString(fmt.Sprintf("Weixin.%s", priority))

	content := fmt.Sprintf(`状态: %s

等级: %s

告警: %s

Item values: 

%s

故障: %s`,
		status, grade, alertname, description, summary)

	data := fmt.Sprintf(`{
        "msgtype": "text",
		"text": {
			"content": "%s",
		},
    }`, content)
	bodys := strings.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", weixinKey), "application/json", bodys)
	if err != nil {
		log.Println("Failed http.Post", http.StatusInternalServerError, receiver, status, grade, alertname, summary, description)
		return err
	}
	log.Println("INFO http.Post", resp.StatusCode, receiver, status, grade, alertname, summary, description)
	return nil
}

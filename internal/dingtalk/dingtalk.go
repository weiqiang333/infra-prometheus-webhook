package dingtalk

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/model"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/utils/notification_process"
)

// Dingtalk 发送钉钉消息程序
func Dingtalk(notification model.Notification, priority string) error {
	var receiver = "dingtalk"
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
	accessToken := viper.GetString(fmt.Sprintf("Dingtalk.%s", priority))

	content := fmt.Sprintf(`状态: %s

等级: %s

告警: %s

Item values: 

%s


故障修复: %s`,
		status, grade, alertname, description, summary)

	data := fmt.Sprintf(`{
        "msgtype": "text",
            "text": {
            "content": "%s",
        },
        "at": {
            "isAtAll": "true",
        },
    }`, content)
	bodys := strings.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken), "application/json", bodys)
	if err != nil {
		log.Println("Failed http.Post", http.StatusInternalServerError, receiver, status, grade, alertname, summary, description)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		msg := fmt.Sprintf("Failed http.Post StatusCode is %v, body %s", resp.StatusCode, string(body))
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	log.Println("INFO http.Post", resp.StatusCode, receiver, status, grade, alertname, summary, description)
	return nil
}

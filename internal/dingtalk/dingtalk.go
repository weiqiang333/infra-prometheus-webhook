package dingtalk

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"infra-prometheus-webhook/model"
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
	description := notification.CommonAnnotations["description"]
	alertSummary := summary(notification)

	content := fmt.Sprintf(`状态: %s

等级: %s

告警: %s



Item values: 

%s



故障修复: %s`,
		status, grade, alertname, alertSummary, description)

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
	resp, err := http.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s",
		model.Config.Dingtalk[priority]), "application/json", bodys)
	if err != nil {
		log.Println(http.StatusInternalServerError, receiver, status, grade, alertname, alertSummary)
		return err
	}
	log.Println(resp.StatusCode, receiver, status, grade, alertname, alertSummary)
	return nil
}

// summary 将多条报警内容总结为一条
func summary(notification model.Notification) string {
	var annotations bytes.Buffer
	for i, alert := range notification.Alerts {
		annotations.WriteString(strconv.Itoa(i+1) + ". " + alert.Annotations["summary"])
		if i + 1 != len(notification.Alerts) {
			annotations.WriteString("\n")
		}
	}
	fmt.Print(annotations.String())
	return annotations.String()
}
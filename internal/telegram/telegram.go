package telegram

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/weiqiang333/infra-prometheus-webhook/model"
)

// Telegram 发送消息
func Telegram(notification model.Notification, priority string) error {
	var receiver = "telegram"
	var status string
	switch notification.Status {
	case "firing":
		status = "PROBLEM"
	case "resolved":
		status = "OK"
	}
	grade := notification.CommonLabels["priority"]
	alertname := notification.CommonLabels["alertname"]
	description := getDescriptionList(notification)
	summary := notification.CommonAnnotations["summary"]

	content := fmt.Sprintf(`状态: %s

等级: %s

告警: %s

Item values: 

%s

故障: %s`,
		status, grade, alertname, description, summary)

	data := fmt.Sprintf(`{
        "parse_mode": "Markdown",
		"chat_id": "%s",
		"text": "%s",
    }`, model.Config.Telegram[priority], content)
	bodys := strings.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage",
		model.Config.Telegram["bot_token"]), "application/json", bodys)
	if err != nil {
		log.Println(http.StatusInternalServerError, receiver, status, grade, alertname, summary, description)
		return err
	}
	log.Println(resp.StatusCode, receiver, status, grade, alertname, summary, description)
	return nil
}

// getDescriptionList 将多条报警内容总结为一条信息清单
func getDescriptionList(notification model.Notification) string {
	var annotations bytes.Buffer
	for i, alert := range notification.Alerts {
		annotations.WriteString(strconv.Itoa(i+1) + ". " + alert.Annotations["description"])
		if i+1 != len(notification.Alerts) {
			annotations.WriteString("\n")
		}
	}
	fmt.Print(annotations.String())
	return annotations.String()
}

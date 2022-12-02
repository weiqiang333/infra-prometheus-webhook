package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	model2 "github.com/weiqiang333/infra-prometheus-webhook/internal/model"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/utils/notification_process"
)

// Telegram 发送消息
func Telegram(notification model2.Notification, priority string) error {
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
	description := notification_process.GetDescriptionList(notification)
	summary := notification.CommonAnnotations["summary"]
	chatId := viper.GetString(fmt.Sprintf("Telegram.%s", priority))
	botToken := viper.GetString("Telegram.bot_token")

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
    }`, chatId, content)
	bodys := strings.NewReader(data)
	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken), "application/json", bodys)
	if err != nil {
		log.Println("Failed http.Post", http.StatusInternalServerError, receiver, status, grade, alertname, summary, description, err.Error())
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

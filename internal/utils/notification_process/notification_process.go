package notification_process

import (
	"bytes"
	"strconv"

	"github.com/weiqiang333/infra-prometheus-webhook/internal/model"
)

// GetDescriptionList 将多条报警内容总结为一条信息清单
func GetDescriptionList(notification model.Notification) string {
	var annotations bytes.Buffer
	for i, alert := range notification.Alerts {
		annotations.WriteString(strconv.Itoa(i+1) + ". " + alert.Annotations["description"])
		if i+1 != len(notification.Alerts) {
			annotations.WriteString("\n")
		}
	}
	return annotations.String()
}

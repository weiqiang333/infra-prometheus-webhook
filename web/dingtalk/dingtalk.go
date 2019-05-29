package dingtalk

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"infra-prometheus-webhook/model"
	"infra-prometheus-webhook/internal/dingtalk"
)

var notification  = model.Notification{}

// Dingtalk 路由入口、响应
func Dingtalk(c *gin.Context) {
	err := c.BindJSON(&notification)
	priority, _ := c.Params.Get("priority")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 发送钉钉消息
	err = dingtalk.Dingtalk(notification, priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
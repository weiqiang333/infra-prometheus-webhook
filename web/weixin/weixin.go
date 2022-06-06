package weixin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/weixin"
	"github.com/weiqiang333/infra-prometheus-webhook/model"
)

// Weixin 路由入口、响应
func Weixin(c *gin.Context) {
	notification := model.Notification{}
	err := c.BindJSON(&notification)
	priority, _ := c.Params.Get("priority")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 发送钉钉消息
	err = weixin.Weixin(notification, priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
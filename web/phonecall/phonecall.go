package phonecall

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"infra-prometheus-webhook/internal/phonecall"
	"infra-prometheus-webhook/model"
)

var notification  = model.Notification{}

// Phonecall 路由入口、响应
func Phonecall(c *gin.Context)  {
	err := c.BindJSON(&notification)
	role, _ := c.Params.Get("role")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 发送电话
	statusCode, err := phonecall.Phonecall(notification, role)
	if err != nil {
		// StatusNotImplemented 501 非报警时间段, 或发送电话报警失败
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(statusCode, gin.H{"message": "Message sent successfully"})
}
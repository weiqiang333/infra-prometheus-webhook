package web

import (
	"github.com/gin-gonic/gin"
	"infra-prometheus-webhook/model"

	"infra-prometheus-webhook/web/dingtalk"
	"infra-prometheus-webhook/web/phonecall"
)

// Webhook 路由入口文件
func Webhook() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	alerts := router.Group("alerts")
	{
		alerts.POST("/dingtalk/:priority", dingtalk.Dingtalk)
		alerts.POST("/phonecall/:role", phonecall.Phonecall)
	}
	router.Run(model.Config.ListenPort)
}

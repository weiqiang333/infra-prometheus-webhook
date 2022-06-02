package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/infra-prometheus-webhook/model"
	"github.com/weiqiang333/infra-prometheus-webhook/web/dingtalk"
	"github.com/weiqiang333/infra-prometheus-webhook/web/phonecall"
	"github.com/weiqiang333/infra-prometheus-webhook/web/weixin"
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
		alerts.POST("/weixin/:priority", weixin.Weixin)
	}
	err := router.Run(model.Config.ListenPort)
	if err != nil {
		fmt.Println("service run failed: %s", err.Error())
	}
}

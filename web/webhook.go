package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/weiqiang333/infra-prometheus-webhook/model"
	"github.com/weiqiang333/infra-prometheus-webhook/web/dingtalk"
	"github.com/weiqiang333/infra-prometheus-webhook/web/phonecall"
	"github.com/weiqiang333/infra-prometheus-webhook/web/telegram"
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
	router.POST("/-/reload", reloadConfig)

	alerts := router.Group("alerts")
	{
		alerts.POST("/dingtalk/:priority", dingtalk.Dingtalk)
		alerts.POST("/phonecall/:role", phonecall.Phonecall)
		alerts.POST("/weixin/:priority", weixin.Weixin)
		alerts.POST("/telegram/:priority", telegram.Telegram)
	}
	err := router.Run(model.Config.ListenPort)
	if err != nil {
		fmt.Println("service run failed: ", err.Error())
	}
}

// reloadConfig 127.0.0.1:8080/-/reload
func reloadConfig(c *gin.Context) {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error config file: %w \n", err))
		c.String(http.StatusOK, fmt.Sprintf("Failed reload config file: %s, err: %s", viper.ConfigFileUsed(), err.Error()))
		return
	}
	fmt.Println("reload config file: ", viper.ConfigFileUsed())
	c.String(http.StatusOK, fmt.Sprintf("reload config file: %s", viper.ConfigFileUsed()))
}

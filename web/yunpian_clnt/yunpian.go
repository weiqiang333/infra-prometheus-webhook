package yunpian_clnt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/model"
	"github.com/weiqiang333/infra-prometheus-webhook/internal/yunpian_clnt"
)

// YunPianClnt 云片路由入口、响应
func YunPianClnt(c *gin.Context) {
	notification := model.Notification{}
	err := c.BindJSON(&notification)
	sendtype, _ := c.Params.Get("sendtype")
	priority, _ := c.Params.Get("priority")
	if err != nil {
		log.Println("Failed YunPianClnt BindJSON err: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 发送 yunpian 消息
	ypClnt := yunpian_clnt.NewYunPianModule()
	switch sendtype {
	case "sms":
		err = ypClnt.PostSms(notification, priority)
	case "voice":
		err = ypClnt.PostVoice(notification, priority)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Parameter exception, %s", sendtype)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

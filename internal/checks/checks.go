package checks

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var (
	ready []string
	notready []string
	client = http.Client{
		Timeout: 2 * time.Second,
	}
)

// Cron asd
func Cron()  {
	c := cron.New()
	c.AddFunc("00 00 10 * * *", checks())
	c.Start()
}

func checks() func() {
	checksUrl := viper.GetStringMapString("checks")
	return func() {
		fmt.Println("The scheduled task that will run the check")
		ready = []string{}
		notready = []string{}
		for service, url := range checksUrl {
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println(http.StatusInternalServerError, service, err)
			}
			if resp != nil && resp.StatusCode == 200 {
				ready = append(ready, service +  ":\tis ready")
			} else {
				notready = append(notready, service + ":\tis not ready")
			}

		}
		content := content()
		dingtalk(content)
	}
}

func content() string {
	var content bytes.Buffer
	content.WriteString("亲，我是 Prometheus's Webhook.\n")
	content.WriteString("\t我检查到 " + strconv.Itoa(len(ready)) + " 个服务正常，" + strconv.Itoa(len(notready)) + "个服务异常.\n")
	if len(notready) > 0 {
		content.WriteString("\t亲，你可能要加班了!!!\n\n")
	} else {
		content.WriteString("\t亲，祝你加班愉快!!!\n\n")
	}
	if len(ready) != 0 {
		content.WriteString("\tService ready:")
		for _, v := range ready {
			content.WriteString("\n\t\t" + v)
		}
	}
	if len(notready) != 0 {
		content.WriteString("\n\tService not ready:")
		for _, v := range notready {
			content.WriteString("\n\t\t" + v)
		}
	}
	return content.String()
}

func dingtalk(content string)  {
	data := fmt.Sprintf(`{
        "msgtype": "text",
            "text": {
            "content": "%s",
        },
        "at": {
            "isAtAll": "",
        },
    }`, content)

	bodys := strings.NewReader(data)
	resp, err := client.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s",
		viper.GetString("dingtalk.p0")), "application/json", bodys)
	if err != nil {
		fmt.Println(resp.StatusCode, content)
	}
}
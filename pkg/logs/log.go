package logs

import (
	"fmt"
	"log"
	"os"
)

// access log: $data $statusCode $receiver $status $grade $alertname $[alertSummary|msg]
func init() {
	file, err := os.OpenFile("log/infra-prometheus-webhook.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)
	log.SetOutput(file)
}

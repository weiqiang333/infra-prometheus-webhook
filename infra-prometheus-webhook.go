// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/weiqiang333/infra-prometheus-webhook/internal/checks"
	"github.com/weiqiang333/infra-prometheus-webhook/web"
)

func init() {
	pflag.String("config", "configs/production.yaml", "config file")
	pflag.String("listen_address", "0.0.0.0:8099", "server listen address.")
	pflag.Bool("check", false, "check's cron: Used to check the infrastructure.(It is a temporary solution and not recommended. - deprecating -)")

	file, err := os.OpenFile("log/infra-prometheus-webhook.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Failed open log file: ", err.Error())
		os.Exit(1)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)
	log.SetOutput(file)
}

func main() {
	loadConfig()

	if viper.GetBool("check") {
		fmt.Println("Run check Cron")
		checks.Cron()
	}

	web.Webhook()
}

// load config and flag config
func loadConfig() {
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println(err.Error())
		panic(fmt.Errorf("Fatal error BindPFlags: %w \n", err))
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("config"))
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	log.Println("load config file:", viper.GetString("config"))
}

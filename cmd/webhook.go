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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"infra-prometheus-webhook/internal/checks"
	"infra-prometheus-webhook/web"
)

// webhookCmd represents the webhook command
var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "prometheus alert receiver webhook",
	Long: `prometheus alert receiver webhook:

		Here, the DingTalk group robot receiver and the yunpian voice receiver are implemented..`,
	Run: func(cmd *cobra.Command, args []string) {
		check, _ := cmd.Flags().GetBool("check")
		if check == true {
			fmt.Println("Run check Cron")
			checks.Cron()
		}
		web.Webhook()
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webhookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webhookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	webhookCmd.Flags().BoolP("check", "c", false, "check's cron: Used to check the infrastructure")
}

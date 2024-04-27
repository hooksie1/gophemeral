/*
Copyright © 2024 John Hooks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hooksie1/gophemeral/service"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "Get a secret",
	RunE:         get,
	SilenceUsage: true,
}

func init() {
	clientCmd.AddCommand(getCmd)
	getCmd.Flags().String("id", "", "The ID of the secret")
	viper.BindPFlag("id", getCmd.Flags().Lookup("id"))
	getCmd.Flags().String("password", "", "The password for the secret")
	viper.BindPFlag("password", getCmd.Flags().Lookup("password"))
	getCmd.Flags().String("get-subject", "gophemeral.secrets.get", "The subject to get a secret")
	viper.BindPFlag("get_subject", getCmd.Flags().Lookup("get-subject"))
}

func get(cmd *cobra.Command, args []string) error {
	var data []byte
	var tv service.TextViews
	nc, err := newNatsConnection("gophemeral-client")
	if err != nil {
		return err
	}

	if len(args) != 0 {
		data = []byte(args[0])
	} else {
		req := service.IDPassword{
			ID:       viper.GetString("id"),
			Password: viper.GetString("password"),
		}

		data, err = json.Marshal(req)
		if err != nil {
			return err
		}

	}

	msg := nats.NewMsg(viper.GetString("get_subject"))
	msg.Data = data
	resp, err := nc.RequestMsg(msg, 1*time.Second)
	if err != nil {
		return err
	}

	if resp.Header.Get("Nats-Service-Error-Code") != "" {
		return fmt.Errorf(string(resp.Data))
	}

	if err := json.Unmarshal(resp.Data, &tv); err != nil {
		return err
	}

	if viper.GetBool("json") {
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("Text: %s\nViews: %d\n", tv.Text, tv.Views)
	if tv.Views == 0 && !viper.GetBool("json") {
		fmt.Println("This is the last time you can view this message")
	}

	return nil
}

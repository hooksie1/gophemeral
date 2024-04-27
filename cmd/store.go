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

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:          "store",
	Short:        "Store a secret",
	RunE:         store,
	SilenceUsage: true,
}

func init() {
	clientCmd.AddCommand(storeCmd)
	storeCmd.Flags().String("text", "", "The text to store")
	viper.BindPFlag("text", storeCmd.Flags().Lookup("text"))
	storeCmd.Flags().Int("views", 1, "The number of views for this secret")
	viper.BindPFlag("views", storeCmd.Flags().Lookup("views"))
	storeCmd.Flags().String("store-subject", "gophemeral.secrets.store", "The subject to store a secret")
	viper.BindPFlag("store_subject", storeCmd.Flags().Lookup("store-subject"))

}

func store(cmd *cobra.Command, args []string) error {
	var idp service.IDPassword
	var data []byte
	nc, err := newNatsConnection("gophemeral-client")
	if err != nil {
		return err
	}

	if len(args) != 0 {
		data = []byte(args[0])
	} else {
		req := service.TextViews{
			Text:  viper.GetString("text"),
			Views: viper.GetInt("views"),
		}

		data, err = json.Marshal(req)
		if err != nil {
			return err
		}

	}

	msg := nats.NewMsg(viper.GetString("store_subject"))
	msg.Data = data
	resp, err := nc.RequestMsg(msg, 1*time.Second)
	if err != nil {
		return err
	}

	if resp.Header.Get("Nats-Service-Error-Code") != "" {
		return fmt.Errorf(string(resp.Data))
	}

	if err := json.Unmarshal(resp.Data, &idp); err != nil {
		return err
	}

	if viper.GetBool("json") {
		fmt.Println(string(resp.Data))
		return nil
	}

	fmt.Printf("ID: %s\nPassword: %s\n", idp.ID, idp.Password)

	return nil

}

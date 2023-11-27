/*
Copyright © 2021 John Hooks john@hooks.technology

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a secret",
	Run:   create,
}

type Record struct {
	Text     string `json:"text,omitempty"`
	Views    int    `json:"views,omitempty"`
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringP("secret", "s", "", "the secret to store")
	createCmd.MarkPersistentFlagRequired("secret")
	createCmd.PersistentFlags().IntP("views", "v", 1, "the number of times the secret can be viewed")
	viper.BindPFlag("secString", createCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("numViews", createCmd.PersistentFlags().Lookup("views"))
}

func create(cmd *cobra.Command, args []string) {
	secret := viper.GetString("secString")
	views := viper.GetInt("numViews")

	rec := Record{
		Text:  secret,
		Views: views,
	}

	data, err := json.Marshal(rec)
	if err != nil {
		log.Printf("Error marshaling data: %s", err)
		os.Exit(1)
	}

	server := viper.GetString("server-address")

	resp, err := http.Post(server+"/api/message", "Content-Type: application/json", bytes.NewReader(data))

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)
	body := buf.Bytes()

	if err := json.Unmarshal(body, &rec); err != nil {
		log.Printf("Error unmarshaling data: %s", err)
		os.Exit(1)
	}

	fmt.Printf("ID: %s\nPassword: %s\n", rec.ID, rec.Password)

}

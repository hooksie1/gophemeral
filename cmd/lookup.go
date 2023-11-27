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

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup a stored secret",
	Run:   lookup,
}

func init() {
	rootCmd.AddCommand(lookupCmd)
	lookupCmd.PersistentFlags().StringP("id", "i", "", "the ID of the secret")
	lookupCmd.PersistentFlags().StringP("pass", "p", "", "the password for the secret")
	viper.BindPFlag("secID", lookupCmd.PersistentFlags().Lookup("id"))
	viper.BindPFlag("secPass", lookupCmd.PersistentFlags().Lookup("pass"))

}

func lookup(cmd *cobra.Command, args []string) {
	id := viper.GetString("secID")
	pass := viper.GetString("secPass")

	client := http.Client{}

	query := map[string]string{
		"id": id,
	}

	server := viper.GetString("server-address")

	req, err := NewRequest(
		SetURL(server+"/api/message"),
		SetMethod("GET"),
		SetHeader("X-Password", pass),
		SetQuery(query),
	)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to server: %s", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var rec Record

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)
	body := buf.Bytes()

	if err := json.Unmarshal(body, &rec); err != nil {
		log.Printf("Error unmarshaling data: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Secret: %s\nViews: %d\n", rec.Text, rec.Views)

}

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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophemeral",
	Short: "gophemeral holds secrets for a limited time",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	_, ok := os.LookupEnv("GOPHEMERAL_SERVER")
	if !ok {
		viper.Set("server-address", "https://api.gophemeral.com")
	} else {
		viper.Set("server-address", os.Getenv("GOPHEMERAL_SERVER"))
	}

	if !strings.HasPrefix(viper.GetString("server-address"), "http") {
		log.Println("server must start with http:// or https://")
		os.Exit(1)
	}

}

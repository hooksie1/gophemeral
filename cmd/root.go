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
	"os"
	"strings"

	"github.com/CoverWhale/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg Config

var rootCmd = &cobra.Command{
	Use:   "gophemeral",
	Short: "The app description",
}
var replacer = strings.NewReplacer("-", "_")

type Config struct {
	Port int `mapstructure:"port"`
}

func Execute() {
	viper.SetDefault("service-name", "gophemeral-local")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gophemeral.json)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".gophemeral")
	}

	viper.SetEnvPrefix("gophemeral")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	logger := logr.NewLogger()
	if err := viper.ReadInConfig(); err == nil {
		logger.Debugf("using config %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		cobra.CheckErr(err)
	}
}

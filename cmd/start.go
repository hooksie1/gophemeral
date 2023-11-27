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
	bc "github.com/hooksie1/bclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/hooksie1/gophemeral/app"
	"gitlab.com/hooksie1/gophemeral/rest"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	RunE:  startServer,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().String("path", "gophemeral.db", "path to local db")
	viper.BindPFlag("path", startCmd.PersistentFlags().Lookup("path"))
	startCmd.PersistentFlags().String("token", "", "token for backend")
	viper.BindPFlag("token", startCmd.PersistentFlags().Lookup("token"))
	startCmd.PersistentFlags().String("collection", "", "collection for backend")
	viper.BindPFlag("collection", startCmd.PersistentFlags().Lookup("collection"))
	startCmd.PersistentFlags().IntP("port", "p", 8080, "port for server")
	viper.BindPFlag("port", startCmd.PersistentFlags().Lookup("port"))

}

func setBackend() (app.Backend, error) {
	//if viper.GetString("collection") != "" {
	//	return app.NewFaunaConnection(viper.GetString("token"), viper.GetString("collection")), nil
	//}

	client := bc.NewClient()
	if err := client.NewDB(viper.GetString("path")); err != nil {
		return nil, err
	}

	bc := app.BoltDB{
		Client: client,
	}

	if err := bc.Init(); err != nil {
		return nil, err
	}

	return bc, nil
}

func startServer(cmd *cobra.Command, args []string) error {
	backend, err := setBackend()
	if err != nil {
		return err
	}

	s := rest.NewServer(backend)

	s.Serve(viper.GetInt("port"))

	return nil
}

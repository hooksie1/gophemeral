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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Flags are defined here. Because of the way Viper binds values, if the same flag name is called
// with viper.BindPFlag multiple times during init() the value will be overwritten. For example if
// two subcommands each have a flag called name but they each have their own default values,
// viper can overwrite any value passed in for one subcommand with the default value of the other subcommand.
// The answer here is to not use init() and instead use something like PersistentPreRun to bind the
// viper values. Using init for the cobra flags is ok, they are only in here to limit duplication of names.

// bindNatsFlags binds nats flag values to viper
func bindNatsFlags(cmd *cobra.Command) {
	viper.BindPFlag("nats_urls", cmd.Flags().Lookup("nats-urls"))
	viper.BindPFlag("nats_seed", cmd.Flags().Lookup("nats-seed"))
	viper.BindPFlag("nats_jwt", cmd.Flags().Lookup("nats-jwt"))
	viper.BindPFlag("nats_secret", cmd.Flags().Lookup("nats-secret"))
	viper.BindPFlag("credentials_file", cmd.Flags().Lookup("credentials-file"))
}

// natsFlags adds the nats flags to the passed in cobra command
func natsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("nats-jwt", "", "NATS JWT as a string")
	cmd.PersistentFlags().String("nats-seed", "", "NATS seed as a string")
	cmd.PersistentFlags().String("credentials-file", "", "Path to NATS user credentials file")
	cmd.PersistentFlags().String("nats-urls", "nats://localhost:4222", "NATS URLs")
}

// bindServiceFlags binds the secret flag values to viper
func bindServiceFlags(cmd *cobra.Command) {
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("max_characters", cmd.Flags().Lookup("max-characters"))
}

// sererFlags adds the service flags to the passed in command
func serviceFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP("port", "p", 8080, "Server port")
	cmd.PersistentFlags().IntP("max-characters", "m", 200, "Maximum characters for a secret")
}

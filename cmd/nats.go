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

	"github.com/CoverWhale/logr"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func newNatsConnection(name string) (*nats.Conn, error) {
	opts := []nats.Option{nats.Name(name)}

	_, ok := os.LookupEnv("USER")

	if viper.GetString("credentials_file") == "" && viper.GetString("nats_jwt") == "" && ok {
		logr.Debug("using NATS context")
		return natscontext.Connect("", opts...)
	}

	if viper.GetString("nats_jwt") != "" {
		opts = append(opts, nats.UserJWTAndSeed(viper.GetString("nats_jwt"), viper.GetString("nats_seed")))
	}
	if viper.GetString("credentials_file") != "" {
		opts = append(opts, nats.UserCredentials(viper.GetString("credentials_file")))
	}

	return nats.Connect(viper.GetString("nats_urls"), opts...)
}

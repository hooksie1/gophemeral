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
	"context"

	cwnats "github.com/CoverWhale/coverwhale-go/transports/nats"
	"github.com/CoverWhale/logr"
	"github.com/hooksie1/gophemeral/rest"
	"github.com/hooksie1/gophemeral/secrets"
	"github.com/hooksie1/gophemeral/service"
	"github.com/invopop/jsonschema"
	"github.com/nats-io/nats.go/micro"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "starts the service",
	RunE:         start,
	SilenceUsage: true,
}

func init() {
	// attach start subcommand to service subcommand
	serviceCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	logger := logr.NewLogger()

	config := micro.Config{
		Name:        "gophemeral",
		Version:     "0.0.1",
		Description: "Secrets sharing for everyone",
	}

	nc, err := newNatsConnection("gophemeral-server")
	if err != nil {
		return err
	}
	defer nc.Close()

	backend, err := secrets.NewNatsBackend(nc)
	if err != nil {
		return err
	}

	// uncomment for config watching
	//js, err := nc.JetStream()
	//if err != nil {
	//    return err
	//}

	// uncomment to enable logging over NATS
	//logger.SetOutput(cwnats.NewNatsLogger("prime.logs.gophemeral", nc))

	svc, err := micro.AddService(nc, config)
	if err != nil {
		logr.Fatal(err)
	}

	// add a handler group
	grp := svc.AddGroup("gophemeral.secrets")
	grp.AddEndpoint("store",
		service.SecretHandler(backend, logger, service.StoreSecret),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "stores a secret",
			"format":          "application/json",
			"request_schema":  schemaString(&secrets.Secret{}),
			"response_schema": schemaString(&secrets.Secret{}),
		}),
		micro.WithEndpointSubject("store"),
	)
	grp.AddEndpoint("get",
		service.SecretHandler(backend, logger, service.GetSecret),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "gets a secret",
			"format":          "application/json",
			"request_schema":  schemaString(&secrets.Secret{}),
			"response_schema": schemaString(&secrets.Secret{}),
		}),
		micro.WithEndpointSubject("get"),
	)

	// uncomment to enable config watching
	//go service.WatchForConfig(logger, js)
	logger.Infof("service %s %s started", svc.Info().Name, svc.Info().ID)
	go cwnats.HandleNotify(svc)

	errChan := make(chan error)

	s := rest.NewServer(backend, logger, viper.GetInt("port"))

	logger.Infof("starting HTTP server on port %d", viper.GetInt("port"))
	go s.Serve(errChan)
	s.AutoHandleErrors(ctx, errChan)
	return nil
}

func schemaString(s any) string {
	schema := jsonschema.Reflect(s)
	data, err := schema.MarshalJSON()
	if err != nil {
		logr.Fatal(err)
	}

	return string(data)
}

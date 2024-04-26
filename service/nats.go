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

package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	cwnats "github.com/CoverWhale/coverwhale-go/transports/nats"
	"github.com/CoverWhale/logr"
	"github.com/hooksie1/gophemeral/secrets"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Handler func(secrets.Backend, *logr.Logger, micro.Request) error

type TextViews struct {
	Text  string `json:"text"`
	Views int    `json:"views"`
}

type IDPassword struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func StoreSecret(b secrets.Backend, logger *logr.Logger, r micro.Request) error {
	var tv TextViews
	if err := json.Unmarshal(r.Data(), &tv); err != nil {
		return cwnats.NewClientError(err, 400)
	}

	s := secrets.Secret{
		Text:  tv.Text,
		Views: tv.Views,
	}

	secret, err := secrets.AddSecret(b, s)
	if err != nil {
		return err
	}

	r.RespondJSON(IDPassword{ID: secret.ID, Password: secret.Password})

	return nil
}

func GetSecret(b secrets.Backend, logger *logr.Logger, r micro.Request) error {
	var idp IDPassword
	if err := json.Unmarshal(r.Data(), &idp); err != nil {
		return cwnats.NewClientError(err, 400)
	}

	s := secrets.Secret{
		ID:       idp.ID,
		Password: idp.Password,
	}

	secret, err := secrets.GetSecret(s, b)
	if err != nil {
		return err
	}

	r.RespondJSON(TextViews{Text: secret.Text, Views: secret.Views})

	return nil
}

func WatchForConfig(logger *logr.Logger, js nats.JetStreamContext) {
	kv, err := js.KeyValue("configs")
	if err != nil {
		logr.Fatal(err)
	}

	w, err := kv.Watch("gophemeral.log_level")
	if err != nil {
		logr.Fatal(err)
	}

	for val := range w.Updates() {
		if val == nil {
			continue
		}

		level := string(val.Value())
		if level == "info" {
			logger.Level = logr.InfoLevel
		}

		if level == "error" {
			logger.Level = logr.ErrorLevel
		}

		if level == "debug" {
			logger.Level = logr.DebugLevel
		}

		logger.Infof("set log level to %s", level)
	}

	time.Sleep(5 * time.Second)
}

func SecretHandler(b secrets.Backend, logger *logr.Logger, h Handler) micro.HandlerFunc {
	return func(r micro.Request) {
		start := time.Now()
		reqLogger := logger.WithContext(map[string]string{"subject": r.Subject()})
		defer func() {
			reqLogger.Infof("duration %dms", time.Since(start).Milliseconds())
		}()

		err := h(b, reqLogger, r)
		if err == nil {
			return
		}

		handleRequestError(reqLogger, err, r)
	}
}

func handleRequestError(logger *logr.Logger, err error, r micro.Request) {
	var ce cwnats.ClientError
	var re secrets.RecordError
	if errors.As(err, &re) {
		r.Error(strconv.Itoa(re.Code()), http.StatusText(re.Code()), []byte(re.Body()))
	}
	if errors.As(err, &ce) {
		r.Error(ce.CodeString(), http.StatusText(ce.Code), ce.Body())
		return
	}

	logger.Error(err)

	r.Error("500", "internal server error", []byte(`{"error": "internal server error"}`))
}

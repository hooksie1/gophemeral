/*
Copyright © 2023 John Hooks

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

package secrets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nats-io/nats.go"
)

type ValidateFunc func(s Secret) error

type NATS struct {
	urls      string
	bucket    string
	conn      *nats.Conn
	opts      []nats.Option
	js        nats.JetStreamContext
	kv        nats.KeyValue
	validator ValidateFunc
}

func NewNatsBackend(nc *nats.Conn, v ValidateFunc) (*NATS, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	kv, err := js.KeyValue("secrets")
	if err != nil {
		return nil, err
	}

	return &NATS{
		conn:      nc,
		bucket:    "secrets",
		js:        js,
		kv:        kv,
		validator: v,
	}, nil
}

func DefaultValidator(length int) ValidateFunc {
	return func(s Secret) error {
		if len(s.Text) > length {
			return NewSecretError(http.StatusBadRequest, fmt.Sprintf("secret length cannot be greater than %d", length))
		}

		return nil
	}
}

func (n *NATS) Validate(s Secret) error {
	return n.validator(s)
}

func (n *NATS) Write(s Secret) error {
	data, err := json.Marshal(s)
	if err != nil {
		return NewSecretError(400, err.Error())
	}
	_, err = n.kv.Put(s.ID, data)
	if err != nil {
		return err
	}

	return nil
}

func (n *NATS) Read(id string) (Secret, error) {
	var secret Secret
	v, err := n.kv.Get(id)
	if err != nil && errors.Is(err, nats.ErrKeyNotFound) {
		return secret, NewSecretError(404, errSecretNotFound.Error())
	}

	if err != nil {
		return secret, err
	}

	if err := json.Unmarshal(v.Value(), &secret); err != nil {
		return secret, err
	}

	return secret, nil

}

func (n *NATS) Delete(id string) error {
	return n.kv.Delete(id)
}

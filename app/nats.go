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

package app

import (
	"encoding/json"
	"errors"

	"github.com/nats-io/nats.go"
)

type NATS struct {
	urls   string
	bucket string
	Conn   *nats.Conn
	opts   []nats.Option
	js     nats.JetStreamContext
	kv     nats.KeyValue
}

func NewNatsBackend(urls string, opts ...nats.Option) *NATS {
	return &NATS{
		bucket: "secrets",
		urls:   urls,
		opts:   opts,
	}
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect(n.urls, n.opts...)
	if err != nil {
		return err
	}

	js, err := conn.JetStream()
	if err != nil {
		return err
	}

	kv, err := js.KeyValue(n.bucket)
	if err != nil {
		return err
	}

	n.Conn = conn
	n.js = js
	n.kv = kv

	return nil
}

func (n *NATS) Write(s Secret) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
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
		return secret, errSecretNotFound
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

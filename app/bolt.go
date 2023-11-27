/*
Copyright © 2023 John Hooks john@hooks.technology

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

	bc "github.com/hooksie1/bclient"
)

const (
	secretBucket = "secrets"
)

type BoltDB struct {
	Client *bc.BoltClient
}

func (b *BoltDB) Init() error {
	bucket := bc.NewBucket("secrets")
	return b.Client.Write(bucket)
}

// BuildRecord does the initial record creation by generating a password,
// a bucket ID, encrypting the text, and then storing that data in the record.
func (b BoltDB) Write(s Secret) error {

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	bucket := bc.NewBucket(secretBucket)
	kv := bc.NewKV().SetBucket(bucket).SetKey(s.ID).SetValue(string(data))

	return b.Client.Write(kv)

}

// Read is just a wrapper to check the number of views left on a
// record. If the value is 0, then DeleteRecord is called.
func (b BoltDB) Read(id string) (Secret, error) {
	var secret Secret
	bucket := bc.NewBucket(secretBucket)
	kv := bc.NewKV().SetBucket(bucket).SetKey(id)

	if err := b.Client.Read(kv); err != nil {
		return secret, err
	}

	if kv.Value == "" {
		return Secret{}, errSecretNotFound
	}

	if err := json.Unmarshal([]byte(kv.Value), &secret); err != nil {
		return secret, err
	}

	return secret, nil
}

// Delete deletes a bucket based on the r.Bucket value.
func (b BoltDB) Delete(id string) error {
	bucket := bc.NewBucket(secretBucket)
	kv := bc.NewKV().SetBucket(bucket).SetKey(id)

	return b.Client.Delete(kv)
}

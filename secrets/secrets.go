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

package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/segmentio/ksuid"
)

var (
	errSecretNotFound = fmt.Errorf("secret not found")
	errBadAuth        = fmt.Errorf("bad password")
)

type Secret struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Password string `json:"password"`
	Views    int    `json:"views"`
}

// generateString takes an int and generates a random string based on the int size.
func generateString(size int) string {
	pass := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, pass)
	if err != nil {
		panic(err)
	}

	id := base64.RawURLEncoding.EncodeToString(pass)

	return string(id)
}

// itob takes an int and converts it into a byte slice
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// boti takes a byte slice and converts it into an int
func boti(bs []byte) int {
	var b [8]byte
	copy(b[8-len(bs):], bs)
	return int(binary.BigEndian.Uint64(b[:]))
}

func AddSecret(w Writer, s Secret) (Secret, error) {
	pass := generateString(24)
	s.ID = ksuid.New().String()

	if len(s.Text) > 100 {
		return Secret{}, NewSecretError(http.StatusBadRequest, "message must be less than 100 chars")
	}
	if s.Views < 1 {
		return Secret{}, NewSecretError(http.StatusBadRequest, "views must be greater than 0")
	}

	encryptedText, err := encrypt([]byte(s.Text), pass)
	if err != nil {
		return Secret{}, fmt.Errorf("Write: %w", err)
	}

	s.Text = string(toBase64(encryptedText))

	if err := w.Write(s); err != nil {
		return Secret{}, err
	}

	// don't set password until here so it's not written in the DB
	s.Password = pass
	s.Text = ""

	return s, nil
}

func GetSecret(s Secret, b Backend) (Secret, error) {

	if err := checkLength(s.Password); err != nil {
		return Secret{}, NewSecretError(http.StatusUnauthorized, errBadAuth.Error())
	}

	secret, err := b.Read(s.ID)
	if err != nil && errors.Is(err, errSecretNotFound) {
		return Secret{}, NewSecretError(http.StatusNotFound, errSecretNotFound.Error())
	}
	if err != nil {
		return Secret{}, err
	}

	decodedSecret, err := fromBase64(secret.Text)
	if err != nil {
		return Secret{}, fmt.Errorf("read: %w", err)
	}

	decryptedMessage, err := decrypt(decodedSecret, s.Password)
	if err != nil {
		return Secret{}, fmt.Errorf("read: %w", err)
	}

	secret.Views = secret.Views - 1

	if secret.Views < 1 {
		if err := b.Delete(secret.ID); err != nil {
			return Secret{}, err
		}
	} else {
		if err := b.Write(secret); err != nil {
			return Secret{}, err
		}
	}

	return Secret{
		Views: secret.Views,
		Text:  string(decryptedMessage),
	}, nil

}

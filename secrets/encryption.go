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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// generateKey generates a 32 byte key and returns a byte slice and an error
func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, fmt.Errorf("generateKey: error getting random data: %w", err)
	}

	return key, nil
}

// makeHash generates a sha256 hash based off of a string. Returns a 32 byte hash
// as a byte slice.
func makeHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// encrypt takes a plain text secret and a 32 byte key and encrypts
// the text using the key. It returns the encrypted text or an error.
func encrypt(plaintext []byte, pass string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(makeHash(pass)))
	if err != nil {
		return nil, fmt.Errorf("encrypt: error creating cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encrypt: error generating GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("encrypt: error reading nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// decrypt takes a byte slice and a 32 byte key and decrypts the text using the key.
// It returns the decrypted value or an error.
func decrypt(ciphertext []byte, pass string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(makeHash(pass)))
	if err != nil {
		return nil, fmt.Errorf("decrypt: error creating cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt: error generating GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("decrypt: malformed ciphertext: %w", err)
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// checkLength checks the length of the password and  returns an error if it is too short.
func checkLength(password string) error {
	if password == "" {
		return fmt.Errorf("checkLength: password is empty")
	}

	if len(password) < aes.BlockSize {
		return fmt.Errorf("checkLength: password is too short")
	}

	return nil
}

func toBase64(key []byte) string {
	return base64.RawStdEncoding.EncodeToString(key)

}

func fromBase64(key string) ([]byte, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

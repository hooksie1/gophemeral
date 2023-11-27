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

import "testing"

func TestDecrypt(t *testing.T) {
	key, err := generateKey()
	if err != nil {
		t.Errorf("error generating key")
	}

	pass := generateString(32)

	encrypted, err := encrypt(key, pass)
	if err != nil {
		t.Errorf("error encrypting data")
	}

	decrypted, err := decrypt(encrypted, pass)
	if err != nil {
		t.Errorf("error decrypting data")
	}

	if len(decrypted) != len(key) {
		t.Errorf("decrypted and key do not match")
	}

	for i := range decrypted {
		if decrypted[i] != key[i] {
			t.Errorf("decrypted and key do not match")
		}
	}

}

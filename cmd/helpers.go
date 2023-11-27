/*
Copyright © 2021 John Hooks john@hooks.technology

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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/spf13/viper"
)

type Printer interface {
	String() string
}

func PrintData(p Printer, r io.ReadCloser) error {
	jsonTrue := viper.GetBool("jsonTrue")
	if jsonTrue {
		data, err := GetJson(r)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
	}

	if !jsonTrue {
		GetData(r, p)
		fmt.Println(p.String())
	}

	return nil
}

func GetData(r io.ReadCloser, v interface{}) error {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}

	return nil

}

func GetJson(r io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

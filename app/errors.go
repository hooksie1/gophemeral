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

type RecordError struct {
	Status      int
	Description string
}

func (r RecordError) Error() string {
	return r.Description
}

func (r RecordError) Code() int {
	return r.Status
}

func (r RecordError) Body() string {
	return r.Description
}

func NewSecretError(status int, description string) RecordError {
	return RecordError{
		Status:      status,
		Description: description,
	}
}

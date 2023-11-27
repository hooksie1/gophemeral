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

package rest

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/hooksie1/gophemeral/app"
)

type AppHandlerFunc func(http.ResponseWriter, *http.Request, app.Backend) error

func getErrorDetails(err error) (int, string) {
	clientError, ok := err.(ClientError)
	if !ok {
		log.Printf("An error ocurred: %v", err)
		return 500, http.StatusText(http.StatusInternalServerError)
	}

	return clientError.Code(), clientError.Body()

}

func errHandlers(h AppHandlerFunc, b app.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r, b)
		if err == nil {
			return
		}

		status, body := getErrorDetails(err)

		apiErrDetails := fmt.Sprintf(`{"error": "%s"}`, body)

		w.WriteHeader(status)
		w.Write([]byte(apiErrDetails))
	}
}

type ClientError interface {
	Error() string
	Code() int
	Body() string
}

type HTTPError struct {
	Details string `json:"details"`
	Status  int    `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Details
}

func (e *HTTPError) Body() string {
	return fmt.Sprintf(`{"error": "%s"}`, e.Details)
}

func (e *HTTPError) Headers() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewHTTPError(status int, details string) error {
	return &HTTPError{
		Details: details,
		Status:  status,
	}
}

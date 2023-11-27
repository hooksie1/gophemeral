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

package rest

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/hooksie1/gophemeral/app"

	"github.com/gorilla/mux"
)

type Server struct {
	Backend app.Backend
	Router  *mux.Router
}

type IDPass struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

type TextViews struct {
	Text  string `json:"text,omitempty"`
	Views int    `json:"views"`
}

func (t *TextViews) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("TextViews was empty")
	}

	text, ok := data["text"].(string)
	if !ok {
		return fmt.Errorf("whoa")
	}

	t.Text = text

	views, ok := data["views"].(int)
	if ok {
		t.Views = views
		return nil
	}

	viewString, ok := data["views"].(string)
	if !ok {
		return fmt.Errorf("unacceptable value for views")
	}

	views, err := strconv.Atoi(viewString)
	if err != nil {
		return err
	}

	t.Views = views

	return nil

}

func (i *IDPass) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("IDPass was empty")
	}

	id, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("error getting id")
	}

	i.ID = id

	password, ok := data["password"].(string)
	if !ok {
		return fmt.Errorf("error getting password")
	}

	i.Password = password

	return nil

}

func NewServer(b app.Backend) Server {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
	router.Handle("/", http.FileServer(http.FS(sub)))

	hxRouter := router.PathPrefix("/hx").Subrouter().StrictSlash(true)
	hxRouter.Handle("/createSecret", http.HandlerFunc(errHandlers(addHxSecret, b))).Methods("POST")
	hxRouter.Handle("/lookupSecret", http.HandlerFunc(errHandlers(getHxSecret, b))).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)
	apiRouter.Handle("/secret", http.HandlerFunc(errHandlers(addSecret, b))).Methods("POST")
	apiRouter.Handle("/secret", http.HandlerFunc(errHandlers(getSecret, b))).Methods("GET")
	apiRouter.Handle("/health", http.HandlerFunc(getHealth)).Methods("GET")

	apiRouter.Use(logger)
	hxRouter.Use(logger)

	return Server{
		Router:  router,
		Backend: b,
	}
}

func (s *Server) Serve(port int) {
	address := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(address, s.Router); err != nil {
		log.Println(err)
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// logger logs the endpoint requested and times how long the request takes.
func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// AddRecord is a handler that creates a record
func addSecret(w http.ResponseWriter, r *http.Request, b app.Backend) error {
	var secret app.Secret

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		return err
	}

	record, err := app.AddSecret(b, secret)
	if err != nil {
		return err
	}

	resp := IDPass{
		ID:       record.ID,
		Password: record.Password,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("error encoding json data: %s", err)
	}

	return nil
}

// getRecord is a handler that retrieves a record
func getSecret(w http.ResponseWriter, r *http.Request, b app.Backend) error {
	password := r.Header.Get("X-Password")
	secret := app.Secret{
		ID:       r.URL.Query().Get("id"),
		Password: password,
	}

	record, err := app.GetSecret(secret, b)
	if err != nil {
		return err
	}

	resp := TextViews{
		Text:  record.Text,
		Views: record.Views,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("error encoding json data: %s", err)
	}

	return nil

}

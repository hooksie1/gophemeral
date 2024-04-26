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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CoverWhale/logr"
	"github.com/hooksie1/gophemeral/secrets"

	"github.com/gorilla/mux"
)

type Server struct {
	Hostname string
	Backend  secrets.Backend
	Router   *http.Server
	Logger   *logr.Logger
}

type IDPass struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
	Link     string `json:"link,omitempty"`
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

func NewServer(b secrets.Backend, l *logr.Logger, port int) Server {
	address := fmt.Sprintf(":%d", port)

	apiServer := &http.Server{
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	s := Server{
		Router:  apiServer,
		Backend: b,
		Logger:  l,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
	router.Handle("/", http.FileServer(http.FS(sub)))

	apiServer.Handler = router

	hxRouter := router.PathPrefix("/hx").Subrouter().StrictSlash(true)
	hxRouter.Handle("/createSecret", http.HandlerFunc(errHandlers(s.addHxSecret))).Methods("POST")
	hxRouter.Handle("/lookupSecret", http.HandlerFunc(errHandlers(s.getHxSecret))).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)
	apiRouter.Handle("/secret", http.HandlerFunc(errHandlers(s.addSecret))).Methods("POST")
	apiRouter.Handle("/secret", http.HandlerFunc(errHandlers(s.getSecret))).Methods("GET")
	apiRouter.Handle("/health", http.HandlerFunc(getHealth)).Methods("GET")

	apiRouter.Use(s.logger)
	hxRouter.Use(s.logger)

	apiServer.Handler = router

	return Server{
		Router:  apiServer,
		Backend: b,
		Logger:  l,
	}
}

func (s *Server) Serve(errChan chan<- error) {
	if err := s.Router.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// logger logs the endpoint requested and times how long the request takes.
func (s *Server) logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			ctx := s.Logger.WithContext(map[string]string{"method": r.Method, "path": r.RequestURI})
			ctx.Infof("processing time: %dms", time.Since(start).Milliseconds())

		}()

		inner.ServeHTTP(w, r)

	})
}

// AddRecord is a handler that creates a record
func (s *Server) addSecret(w http.ResponseWriter, r *http.Request) error {
	var secret secrets.Secret

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		return err
	}

	record, err := secrets.AddSecret(s.Backend, secret)
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
func (s *Server) getSecret(w http.ResponseWriter, r *http.Request) error {
	password := r.Header.Get("X-Password")
	secret := secrets.Secret{
		ID:       r.URL.Query().Get("id"),
		Password: password,
	}

	record, err := secrets.GetSecret(secret, s.Backend)
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

func (s *Server) AutoHandleErrors(ctx context.Context, errChan <-chan error) {
	go func() {
		serverErr := <-errChan
		if serverErr != nil {
			s.Logger.Errorf("error starting server: %v", serverErr)
			s.ShutdownServer(ctx)
		}
	}()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigTerm
	s.Logger.Infof("received signal: %s", sig)
	s.ShutdownServer(ctx)
}

func (s *Server) ShutdownServer(ctx context.Context) {
	s.Logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s.Logger.Info("server stopped")
	os.Exit(1)
}

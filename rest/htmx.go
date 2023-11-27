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
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"gitlab.com/hooksie1/gophemeral/app"
)

//go:embed static/*
var staticFS embed.FS

var lookupTemplate = `
<div id="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay">
		<div class="modal-content text-[#41454c] bg-[#fcfcfc] dark:text-[#ffffff] dark:bg-[#031022]">
				<h3 class="text-3xl text-white max-w-none">Secret Information</h3>
            {{ if .Text  }}
			    <div class="text-left items-left">
			    	<div><b>Secret Text</b>: {{ .Text }}</div>
			    	<div><b>Views</b>: {{ .Views }}
			    </div>
                {{ if eq .Views 0 }} 
                    <b class="text-red">This is the last time you can view this message</b>
                  </div>
                {{ end }}
            {{ else }}
			    <div class="text-left items-left">
			    	<div>{{ .Err }}</div>
			    </div>
            {{ end }}
			<div>
				<button class="mt-3 bg-transparent font-semibold hover:text-white py-2 px-4 border hover:border-transparent rounded" type="button" _="on click trigger closeModal">Close</button>
			</div>
		</div>
	</div>
</div>
`
var errorTemplate = `
<div id="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay">
		<div class="modal-content text-[#41454c] bg-[#fcfcfc] dark:text-[#ffffff] dark:bg-[#031022]">
			<h3 class="text-3xl text-white max-w-none">Secret Information</h3>
			<div class="text-left items-left">
				<div>{{ .Err }}</div>
			</div>
			<div>
				<button class="mt-3 bg-transparent font-semibold hover:text-white py-2 px-4 border hover:border-transparent rounded" type="button" _="on click trigger closeModal">Close</button>
			</div>
		</div>
	</div>
</div>
`

var createTemplate = `
<div id="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay">
		<div class="modal-content text-[#41454c] bg-[#fcfcfc] dark:text-[#ffffff] dark:bg-[#031022]">
				<h3 class="text-3xl text-white max-w-none">Secret Information</h3>
			<div class="text-left items-left">
				<div><b>Secret ID</b>: {{ .ID }}</div>
				<div><b>Password</b>: {{ .Password }}</div>
			<div>
				<button class="mt-3 bg-transparent font-semibold hover:text-white py-2 px-4 border hover:border-transparent rounded" type="button" _="on click trigger closeModal">Close</button>
			</div>
		</div>
	</div>
</div>
`

func addHxSecret(w http.ResponseWriter, r *http.Request, b app.Backend) error {
	var tv TextViews

	modal, err := template.New("modal").Parse(createTemplate)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(&tv); err != nil {
		return err
	}

	rec := app.Secret{
		Text:  tv.Text,
		Views: tv.Views,
	}

	resp, err := app.AddSecret(b, rec)
	if err != nil {
		return handleHTMXError(err, w)
	}

	return modal.Execute(w, resp)
}

func getHxSecret(w http.ResponseWriter, r *http.Request, b app.Backend) error {
	modal, err := template.New("modal").Parse(lookupTemplate)
	if err != nil {
		return err
	}

	r.ParseForm()

	secret := app.Secret{
		ID:       r.FormValue("id"),
		Password: r.FormValue("password"),
	}

	resp, err := app.GetSecret(secret, b)
	if err != nil {
		return handleHTMXError(err, w)
	}

	return modal.Execute(w, resp)
}

func handleHTMXError(err error, w io.Writer) error {
	code, errDetails := getErrorDetails(err)

	modal, err := template.New("modal").Parse(errorTemplate)
	if err != nil {
		return err
	}

	r := struct {
		Code int
		Err  string
	}{
		Code: code,
		Err:  errDetails,
	}

	return modal.Execute(w, r)
}

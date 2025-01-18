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
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/hooksie1/gophemeral/secrets"
)

//go:embed static/*
var staticFS embed.FS

var lookupTemplate = `
<div id="modal" _="on closeModal add .closing then wait for animationend then remove me">
	<div class="modal-underlay">
		<div class="modal-content text-[#41454c] bg-[#fcfcfc] dark:text-[#ffffff] dark:bg-[#031022]">
        <h3 class="text-3xl text-[#41454c] dark:text-[#ffffff] max-w-none">Secret Information</h3>
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
            <h3 class="text-3xl text-[#41454c] dark:text-[#ffffff] max-w-none">Secret Information</h3>
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
            <h3 class="text-3xl text-[#41454c] dark:text-[#ffffff] max-w-none">Secret Information</h3>
			<div class="text-left items-left">
				<div><b>Secret ID</b>: <a href={{ .Link }}>{{ .ID }}</a></div>
				
					<div><b>Password</b>: <span id="secretPassword" style="display:none">{{ .Password }}</span>
					<button class="px-4 text-white"
						_="on click show #secretPassword then hide">
						Show Password
					</button>
					<button class="px-4 text-white"
						_="on click call navigator.clipboard.writeText(#secretPassword.innerText) then put 'Secret Copied!' into #copyConfirmation then remove .hidden from #copyConfirmation">
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" style="width:20px">
							<path stroke-linecap="round" stroke-linejoin="round" d="M8.25 7.5V6.108c0-1.135.845-2.098 1.976-2.192.373-.03.748-.057 1.123-.08M15.75 18H18a2.25 2.25 0 0 0 2.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 0 0-1.123-.08M15.75 18.75v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5A3.375 3.375 0 0 0 6.375 7.5H5.25m11.9-3.664A2.251 2.251 0 0 0 15 2.25h-1.5a2.251 2.251 0 0 0-2.15 1.586m5.8 0c.065.21.1.433.1.664v.75h-6V4.5c0-.231.035-.454.1-.664M6.75 7.5H4.875c-.621 0-1.125.504-1.125 1.125v12c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V16.5a9 9 0 0 0-9-9Z" />
						</svg>
					</button>
				</div>
				<p id="copyConfirmation" class="hidden"></p>
			<div>
				<button class="mt-3 bg-transparent font-semibold hover:text-white py-2 px-4 border hover:border-transparent rounded" type="button" _="on click trigger closeModal">Close</button>
			</div>
		</div>
	</div>
</div>
`

func (s *Server) addHxSecret(w http.ResponseWriter, r *http.Request) error {
	var tv TextViews

	modal, err := template.New("modal").Parse(createTemplate)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(&tv); err != nil {
		return err
	}

	rec := secrets.Secret{
		Text:  tv.Text,
		Views: tv.Views,
	}

	resp, err := secrets.AddSecret(s.Backend, rec)
	if err != nil {
		return handleHTMXError(err, w)
	}

	url := r.Header.Get("origin")

	if r.Header.Get("x-forward-host") != "" {
		url = r.Header.Get("x-forward-host")
	}

	idPass := IDPass{
		ID:       resp.ID,
		Password: resp.Password,
		Link:     fmt.Sprintf(`%s?id=%s`, url, resp.ID),
	}

	return modal.Execute(w, idPass)
}

func (s *Server) getHxSecret(w http.ResponseWriter, r *http.Request) error {
	modal, err := template.New("modal").Parse(lookupTemplate)
	if err != nil {
		return err
	}

	r.ParseForm()

	secret := secrets.Secret{
		ID:       r.FormValue("id"),
		Password: r.FormValue("password"),
	}

	resp, err := secrets.GetSecret(secret, s.Backend)
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

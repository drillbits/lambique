//    Copyright 2017 drillbits
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package lambique

import "net/http"

// App is an application.
type App struct {
	Mux http.Handler
}

// WithMux creates a new application with mux.
func WithMux(mux http.Handler) *App {
	return &App{
		Mux: mux,
	}
}

// Server creates a new server.
func (app *App) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: app.Mux,
	}
}

// Start serves a HTTP server.
func (app *App) Start(addr string) error {
	s := app.Server(addr)
	return s.ListenAndServe()
}

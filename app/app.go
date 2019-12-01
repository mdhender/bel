// bel: a bel interpreter
// Copyright (C) 2019  Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package app implements the application server.
// It's still just an HTTP server.
package app

import (
	"fmt"
	"github.com/mdhender/bel/router"
	"log"
	"net/http"
	"time"
)

// Server is the application server.
// That just means that it has an http server embedded in it.
type Server struct{
	http.Server
}

// NewApp returns a default application server
func NewServer(port int) *Server {
	a := &Server{
		http.Server{
			Addr: fmt.Sprintf(":%d", port),
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
		},
	}
	a.Handler = a
	return a
}

// Handler supports
//   GET  /
//
func (srv *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	log.Printf("[app] about: %q", r.URL.Path)
	switch r.Method {
	case "GET":
		srv.respond(w, r, struct {Path string `json:"path"`}{Path: r.URL.Path}, http.StatusOK)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Handler supports
//   GET  /
//
func (srv *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		srv.respond(w, r, "index.html", http.StatusOK)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// ServeHTTP implements the http handler interface
//
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var route string
	route, r.URL.Path = router.Split(r.URL.Path)
	log.Printf("[app] route %q: rest %q", route, r.URL.Path)

	// route first
	switch route {
	case "":
		srv.handleIndex(w, r)
		return
	case "about":
		srv.handleAbout(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

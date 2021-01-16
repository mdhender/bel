/*
 * Bel - an implementation of Paul Graham's Bel
 *
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
type Server struct {
	http.Server
}

// NewApp returns a default application server
func NewServer(port int) *Server {
	a := &Server{
		http.Server{
			Addr:           fmt.Sprintf(":%d", port),
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
		srv.respond(w, r, struct {
			Path string `json:"path"`
		}{Path: r.URL.Path}, http.StatusOK)
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

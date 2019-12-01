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

package main

//go:generate stringer -type opcode bel/opcode.go

import (
	"fmt"
	"github.com/mdhender/bel/app"
	"github.com/mdhender/bel/bel"
	"io/ioutil"
	"log"
	"os"
)

// bel has four fundamental data types
type CHAR struct {}
type PAIR struct {
  left, right *PAIR
}
type STREAM struct {}
type SYMBOL struct {
  name string
}

func (c *CHAR) IsAtom() bool {
	return true
}

func (p *PAIR) IsAtom() bool {
	return false
}

func (s *STREAM) IsAtom() bool {
	return true
}

func (s *SYMBOL) IsAtom() bool {
	return true
}

var NIL *PAIR

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	fmt.Println("thank you")
}

func run() error {
	config := struct {
		root   string
		port   int
		public string
	}{
		root:   "/Volumes/ssda/mdhender/Software/bel",
		port:   3001,
		public: "public",
	}

	var err error
	if config.root == "" {
		if config.root, err = os.Getwd(); err != nil {
			return err
		}
	} else if err = os.Chdir(config.root); err != nil {
		return err
	}
	log.Printf("[app] root %q\n", config.root)

	src, err := ioutil.ReadFile("source.bel")
	if err != nil {
		return err
	}

	vm := bel.NewVM([]string{"source.bel"})
	if vm == nil {
		panic("missing vm")
	}
	c, err := vm.Execute(src)
	fmt.Printf("[app] result %d\n", len(c))
	if err != nil {
		return err
	}
	fmt.Println(c)
	if c != nil {
		return nil
	}

	a := app.NewServer(config.port)
	return a.ListenAndServe()
}

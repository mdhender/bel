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

// Package bel is the interpreter for paul graham's BEL
package bel

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// VM is the virtual machine for the interpreter
type VM struct {
	initFiles []string

	// nil is the ubiquitous empty list
	t     *cell  // t is the truthy cell.cell
	o     *cell  // something?
	apply *cell  // something?
	ins   stream // default input cellStream
	outs  stream // default output cellStream
	errs  stream // default error cellStream

	frames frameStack
	scanner *scanner

	// following are used during eval
	global struct {
		args *cell
		env *cell
		op   opcode
		saveRegisters *cell
		currentEnv *cell
		currentToken *token
		printFlag int
		isTopLevel bool
		tempOutputStream *cell
	}
}

type stream struct {
	r []io.Reader
	w []io.Writer
}

func NewVM(initFiles []string) *VM {
	vm := &VM{
		t: _TRUE,
		ins: stream{
			r: []io.Reader{os.Stdin},
		},
		outs: stream{
			w: []io.Writer{os.Stdout},
		},
		errs: stream{
			w: []io.Writer{os.Stderr},
		},
	}
	vm.global.env = mkpair(_NIL, _NIL)

	fmt.Println(_NIL)
	fmt.Println(_FALSE)
	fmt.Println(_TRUE)

	// we want to guarantee that the caller can release parameters
	for _, name := range initFiles {
		vm.initFiles = append(vm.initFiles, string([]byte(name)))
	}

	return vm
}

func (vm *VM) Run() {
	for _, name := range vm.initFiles {
		fmt.Printf("(load %q)\n", name)
		args := mkpair(mkstring(name), _NIL)
		vm.eval(opcLoad, args)
	}
}

func (vm *VM) Execute(b []byte) ([]*cell, error) {
	vm.ins.pushReader(bytes.NewReader(b))
	var result []*cell
	result = append(result, vm.eval(opcLoad, mkpair(mkstring("source.bel"), _NIL)))
	vm.ins.popReader()
	return result, nil
}

func (vm *VM) read(arg *cell) *cell {
	return nil
}

func (vm *VM) print(arg *cell) *cell {
	return nil
}

func (vm *VM) popReader() {
	vm.ins.popReader()
}

func (vm *VM) pushReader(r io.Reader) {
	vm.ins.pushReader(r)
}

func (vm *VM) popWriter() {
	vm.outs.popWriter()
}

func (vm *VM) pushWriter(w io.Writer) {
	vm.outs.pushWriter(w)
}

func (vm *VM) puts(s string) {
	vm.outs.puts(s)
}

func (s stream) popReader() {
	if !(len(s.r) > 1) {
		return
	}
	s.r[0] = nil
	s.r = s.r[:1]
}

func (s stream) pushReader(r io.Reader) {
	s.r = append(s.r, r)
}

func (s stream) popWriter() {
	if !(len(s.w) > 1) {
		return
	}
	s.w[0] = nil
	s.w = s.w[:1]
}

func (s stream) pushWriter(w io.Writer) {
	s.w = append(s.w, w)
}

func (s stream) puts(str string) {
	if len(s.w) == 0 {
		fmt.Printf("error: no writer: %q\n", str)
		return
	}
	fmt.Fprint(s.w[0], str)
}


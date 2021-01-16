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
	t     *cell       // t is the truthy cell.cell
	o     *cell       // something?
	apply *cell       // something?
	ins   []io.Reader // default input cellStream
	outs  []io.Writer // default output cellStream
	errs  []io.Writer // default error cellStream

	frames  frameStack
	scanner *scanner

	// following are used during eval
	global struct {
		args             *cell
		env              *cell
		op               opcode
		saveRegisters    *cell
		currentEnv       *cell
		currentToken     *token
		printFlag        int
		isTopLevel       bool
		tempOutputStream *cell
	}
}

type stream struct {
	name string
	r    io.Reader
	w    io.Writer
}

func NewVM(initFiles []string) *VM {
	vm := &VM{
		t:    _TRUE,
		ins:  []io.Reader{os.Stdin},
		outs: []io.Writer{os.Stdout},
		errs: []io.Writer{os.Stderr},
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
	var result []*cell
	result = append(result, vm.eval(opcLoad, mkpair(mkstreamr("source.bel", bytes.NewReader(b)), _NIL)))
	return result, nil
}

func (vm *VM) read(arg *cell) *cell {
	return nil
}

func (vm *VM) print(arg *cell) *cell {
	return nil
}

func (vm *VM) popReader() {
	if len(vm.ins) == 0 {
		panic("vm: ins: underflow")
	}
	vm.ins[0] = nil
	vm.ins = vm.ins[1:]
}

func (vm *VM) pushReader(r io.Reader) {
	vm.ins = append([]io.Reader{r}, vm.ins...)
}

func (vm *VM) popErrorWriter() {
	if len(vm.errs) == 0 {
		panic("vm: errs: underflow")
	}
	vm.errs[0] = nil
	vm.errs = vm.errs[1:]
}

func (vm *VM) pushErrorWriter(w io.Writer) {
	vm.errs = append([]io.Writer{w}, vm.errs...)
}

func (vm *VM) popWriter() {
	if len(vm.outs) == 0 {
		panic("vm: outs: underflow")
	}
	vm.outs[0] = nil
	vm.outs = vm.outs[1:]
}

func (vm *VM) pushWriter(w io.Writer) {
	vm.outs = append([]io.Writer{w}, vm.outs...)
}

func (vm *VM) puts(str string) {
	if len(vm.outs) == 0 {
		fmt.Printf("error: no writer: %q\n", str)
		return
	}
	fmt.Fprint(vm.outs[0], str)
}

func (s stream) puts(str string) {
	fmt.Fprint(s.w, str)
}

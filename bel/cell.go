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

package bel

import (
	"fmt"
	"io"
)

// cell implements a word of storage.
type cell struct {
	_flag   bitfield
	_object object
}

// flag is a bit-field for the type of cell
type bitfield uint

// enums for the flag bit-field
const (
	bfString       bitfield = 1     // 0000000000000001
	bfNumber       bitfield = 2     // 0000000000000010
	bfSymbol       bitfield = 4     // 0000000000000100
	bfSyntax       bitfield = 8     // 0000000000001000
	bfProc         bitfield = 16    // 0000000000010000
	bfPair         bitfield = 32    // 0000000000100000
	bfClosure      bitfield = 64    // 0000000001000000
	bfContinuation bitfield = 128   // 0000000010000000
	bfMacro        bitfield = 256   // 0000000100000000
	bfPromise      bitfield = 512   // 0000001000000000
	bfStream       bitfield = 1024  // 0000010000000000
	bfOpCode       bitfield = 2048  // 0000100000000000
	bfAtom         bitfield = 16384 // 0100000000000000
)

// object will eventually be the union of the atomic data types
type object struct {
	_string string
	_number number
	_pair   pair
	_stream stream
}

// number is just a plain integer for the moment
type number struct {
	_ivalue int
}

// pair is a cons cell
type pair struct {
	_car *cell
	_cdr *cell
}

func asstring(a *cell) string {
	return a._object._string
}

func car(a *cell) *cell {
	return a._object._pair._car
}

func cdr(a *cell) *cell {
	return a._object._pair._cdr
}

func ispair(a *cell) bool {
	return (a._flag & bfPair) != 0
}

func isstream(a *cell) bool {
	return (a._flag & bfStream) != 0
}

func isstring(a *cell) bool {
	return (a._flag & bfString) != 0
}

func mknumber(ivalue int) *cell {
	return &cell{
		_flag: bfNumber,
		_object: object{
			_number: number{
				_ivalue: ivalue,
			},
		},
	}
}

func mkpair(car, cdr *cell) *cell {
	fmt.Printf("(mkpair %s %s)\n", car, cdr)
	return &cell{
		_flag: bfPair,
		_object: object{
			_pair: pair{
				_car: car,
				_cdr: cdr,
			},
		},
	}
}

// mkself creates a new self-linking ATOM
func mkself() *cell {
	a := &cell{_flag: bfAtom}
	setcar(a, a)
	setcdr(a, a)
	return a
}

// mkstreamr creates a new STREAM reader cell
func mkstreamr(r io.Reader) *cell {
	return &cell{
		_flag: bfStream,
		_object: object{
			_stream: stream{
				r: []io.Reader{r},
			},
		},
	}
}

// mkstreamw creates a new STREAM writer cell
func mkstreamw(w io.Writer) *cell {
	return &cell{
		_flag: bfStream,
		_object: object{
			_stream: stream{
				w: []io.Writer{w},
			},
		},
	}
}

// mkstring creates a new STRING cell
func mkstring(s string) *cell {
	return &cell{
		_flag: bfString,
		_object: object{
			_string: s,
		},
	}
}

func setcar(a, b *cell) {
	a._object._pair._car = b
}

func setcdr(a, b *cell) {
	a._object._pair._cdr = b
}

func (a *cell) String() string {
	if a == _NIL {
		return "()"
	} else if a == _FALSE {
		return "#f"
	} else if a == _TRUE {
		return "#t"
	} else if ispair(a) {
		return fmt.Sprintf("(%s . %s)", car(a), cdr(a))
	} else if isstream(a) {
		if len(a._object._stream.r) != 0 {
			return "#reader"
		} else if len(a._object._stream.w) != 0 {
			return "#writer"
		}
		return "#?stream?"
	} else if isstring(a) {
		return fmt.Sprintf("%q", a._object._string)
	}
	return "#unknown"
}
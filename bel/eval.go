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
)

// framePop pops a frame
func (vm *VM) framePop() {
	f := vm.frames.pop()
	vm.global.op = f.op
	vm.global.currentEnv = f.env
}

// framePush pushes a frame
func (vm *VM) framePush(op opcode, args ...*cell) {
	vm.frames.push(op, vm.global.currentEnv, args...)
}

var loops int

func (vm *VM) eval(op opcode, args *cell) *cell {
	fmt.Printf("\n%03d$ (eval (%s %s) %s)\n", loops, op, args, vm.global.env)

	var err error
	vm.global.args = args
	vm.global.isTopLevel = true

	for {
		if err != nil {
			fmt.Printf("eval: %+v\n", err)
			return _NIL
		}

		loops = loops + 1
		if loops > 50 {
			panic("threshold exceeded")
		}

		fmt.Printf("\n%03d> (%s %s)\n", loops, op, vm.global.args)

		switch op {
		case opcError0:
			// args: (string ...)
			if !isstring(car(vm.global.args)) {
				vm.global.args = mkpair(mkstring("error: argument is not string"), _NIL)
				op = opcError0
				continue
			}
			vm.pushWriter(vm.errs.w[0])
			vm.puts("error: ")
			vm.puts(asstring(car(vm.global.args)))
			// output the remaining error strings
			vm.global.args = cdr(vm.global.args)
			op = opcError1
		case opcError1:
			// args: NIL or a list
			// if the args are NIL, print a newline and return to the top level.
			// otherwise, print the car and then the cdr of the list.
			if vm.global.args == _NIL {
				vm.puts("\n")
				vm.popWriter()
				vm.global.isTopLevel = true
				op = opcTopLevel0
				continue
			}
			vm.puts(" ")
			carArgs, cdrArgs := car(vm.global.args), cdr(vm.global.args)
			vm.framePush(opcError1, cdrArgs)
			vm.global.args = carArgs
			vm.global.printFlag = 1
			op = opcP0List
		case opcLoad:
			// args: (filename)
			name, args := car(vm.global.args), cdr(vm.global.args)
			if !isstring(name) {
				a := mkpair(mkstring("load: argument is not string"), _NIL)
				vm.global.args = a
				op = opcError0
				continue
			} else if args != _NIL {
				a := mkpair(mkstring("load: too many arguments"), _NIL)
				vm.global.args = a
				op = opcError0
				continue
			}
			// interactive, need better way to test this
			if asstring(name) == "*stdin*" {
				panic(fmt.Sprintf("%s: interactive: not implemented", op))
			}
			fmt.Printf("\nloading %s\n", asstring(name))
			op = opcTopLevel0
		case opcP0List:
			panic(fmt.Sprintf("%s: not implemented", op))
		case opcP1List:
			panic(fmt.Sprintf("%s: not implemented", op))
		case opcRead:
			if vm.global.currentToken, err = vm.scanner.nextToken(); err != nil {
				fmt.Printf("read: %+v\n", err)
				panic(fmt.Sprintf("op(%s): unhandled error", op))
			}
			op = opcReadSExpr
		case opcReadSExpr:
			if vm.global.currentToken.k == tkEOF {
				fmt.Printf("\n%s: found end-of-file\n", op)
				return _NIL
			}
			panic(fmt.Sprintf("%s: not implemented", op))
		case opcTopLevel0:
			// flush the output stream
			vm.puts("\n")
			// clear any existing frames
			vm.frames.reset()
			// reset the environment
			vm.global.currentEnv = vm.global.env
			// push two frames to run top level one and then to print the result
			vm.frames.push(opcValuePrint, _NIL)
			vm.frames.push(opcTopLevel1, _NIL)
			// display a prompt
			vm.puts("bel> ")
			// read in the next bit of input
			op = opcRead
		default:
			panic(fmt.Sprintf("%s: not implemented", op))
		}
	}
	return _NIL
}
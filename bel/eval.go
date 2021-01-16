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
			vm.pushWriter(vm.errs[0])
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
			// args: (filename) or (stream)
			input, args := car(vm.global.args), cdr(vm.global.args)
			if isstring(input) {
				if asstring(input) == "*stdin*" {
					panic(fmt.Sprintf("%s: interactive: not implemented", op))
				}
				a := mkpair(mkstring("load: filename: not implemented"), _NIL)
				vm.global.args = a
				op = opcError0
				continue
			} else if isstream(input) {
				if args != _NIL {
					a := mkpair(mkstring("load: too many arguments"), _NIL)
					vm.global.args = a
					op = opcError0
					continue
				}
				fmt.Printf("\nloading %s\n", input._object._stream.name)
				vm.pushReader(input._object._stream.r)
				op = opcTopLevel0
				continue
			}
			a := mkpair(mkstring("load: argument must be stream or string"), _NIL)
			vm.global.args = a
			op = opcError0
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

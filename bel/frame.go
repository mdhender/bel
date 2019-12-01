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

type frameStack struct {
	stack []*frame
}

type frame struct {
	op opcode
	env *cell
	args []*cell
	code *cell
}

func (fs *frameStack) pop() *frame {
	if len(fs.stack) == 0 {
		panic("frame: stack: underflow")
	}
	f := fs.stack[len(fs.stack)-1]
	fs.stack = fs.stack[:len(fs.stack)-1]
	return f
}

func (fs *frameStack) push(op opcode, env *cell, args ...*cell) {
	f := &frame{op: op, env: env}
	for _, arg := range args {
		f.args = append(f.args, arg)
	}

	fs.stack = append(fs.stack, f)
}

func (fs *frameStack) reset() {
	fs.stack = nil
}


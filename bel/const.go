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

// NIL is a special cell representing the empty list
var _NIL *cell

// FALSE is a special cell representing #f
var _FALSE *cell

// TRUE is a special cell representing #t
var _TRUE *cell

func init() {
	_NIL = mkself()	// special cell representing the empty list
	_FALSE = mkself() // special cell representing #f
	_TRUE = mkself() // special cell representing #t
}

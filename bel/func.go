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

// AsString returns the STRING of a cell
func (a *cell) AsString() string {
	if !IsString(a) {
		panic("cell::AsString")
	}
	return asstring(a)
}

// Car returns the CAR of a PAIR.
func (a *cell) Car() *cell {
	if !IsPair(a) {
		panic("cell:Car")
	}
	return car(a)
}

// Cdr returns the CDR of a PAIR.
func (a *cell) Cdr() *cell {
	return cdr(a)
}

// Car returns the CAR of a PAIR.
func Car(a *cell) *cell {
	return a._object._pair._car
}

// Cdr returns the CDR of a PAIR.
func Cdr(a *cell) *cell {
	return a._object._pair._cdr
}

// IsPair is a predicate returning true if the cell is a pair
func IsPair(a *cell) bool {
	return ispair(a)
}

// IsString is a predicate returning true if the cell is a string
func IsString(a *cell) bool {
	return isstring(a)
}

// MkPair creates a new PAIR
func MkPair(a, b *cell) *cell {
	return mkpair(a, b)
}

// MkString creates a new STRING cell
func MkString(s string) *cell {
	return mkstring(s)
}


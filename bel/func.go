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

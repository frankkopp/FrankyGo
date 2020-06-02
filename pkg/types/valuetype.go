//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package types

// ValueType is a set of constants for value types used in search and TtTable
//  Vnone   ValueType = 0
//	EXACT   ValueType = 1
//	ALPHA   ValueType = 2 // upper bound
//	BETA    ValueType = 3 // lower bound
//	Vlength int       = 4
type ValueType int8

// ValueType is a set of constants for value types used in search and TtTable
//noinspection GoVarAndConstTypeMayBeOmitted
const (
	Vnone   ValueType = 0
	EXACT   ValueType = 1
	ALPHA   ValueType = 2 // upper bound
	BETA    ValueType = 3 // lower bound
	Vlength int       = 4
)

// IsValid check if pt is a valid value type
func (vt ValueType) IsValid() bool {
	return vt < 4
}

// array of string labels for piece types
var valueTypeToString = [Vlength]string{"NoneValue", "ExactValue", "AlphaValue", "BetaValue"}

// String returns a string representation of a piece type
func (vt ValueType) String() string {
	return valueTypeToString[vt]
}

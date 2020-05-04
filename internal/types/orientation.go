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

import "fmt"

// Orientation is a set of constants for directions from a squares
type Orientation uint8

// Orientation is a set of constants for moving squares within a Bb
//noinspection GoVarAndConstTypeMayBeOmitted
const (
	NW Orientation = 0
	N  Orientation = 1
	NE Orientation = 2
	E  Orientation = 3
	SE Orientation = 4
	S  Orientation = 5
	SW Orientation = 6
	W  Orientation = 7
)

// IsValid tests if o is a valid Orientation value
func (o Orientation) IsValid() bool {
	return o < 8
}

// String returns a string representation of a Orientation (e.g. N, E, ...,NW,...)
func (d Orientation) String() string {
	switch d {
	case N:
		return "N"
	case E:
		return "E"
	case S:
		return "S"
	case W:
		return "W"
	case NE:
		return "NE"
	case SE:
		return "SE"
	case SW:
		return "SW"
	case NW:
		return "NW"
	default:
		panic(fmt.Sprintf("Invalid orientation %d", d))
	}
}

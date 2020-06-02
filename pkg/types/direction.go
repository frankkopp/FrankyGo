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

// Direction is a set of constants for moving squares within a Bb
//  North     Direction = 8
//	East      Direction = 1
//	South     Direction = -North
//	West      Direction = -East
//	Northeast Direction = North + East
//	Southeast Direction = South + East
//	Southwest Direction = South + West
//	Northwest Direction = North + West
type Direction int8

// Direction is a set of constants for moving squares within a Bb
//noinspection GoVarAndConstTypeMayBeOmitted
const (
	North     Direction = 8
	East      Direction = 1
	South     Direction = -North
	West      Direction = -East
	Northeast Direction = North + East
	Southeast Direction = South + East
	Southwest Direction = South + West
	Northwest Direction = North + West
)

var Directions = [8]Direction{North, East, South, West, Northeast, Southeast, Southwest, Northwest}

// String returns a string representation of a Direction (e.g. N, E, ...,NW,...)
func (d Direction) String() string {
	switch d {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	case Northeast:
		return "NE"
	case Southeast:
		return "SE"
	case Southwest:
		return "SW"
	case Northwest:
		return "NW"
	default:
		panic(fmt.Sprintf("Invalid direction %d", d))
	}
}

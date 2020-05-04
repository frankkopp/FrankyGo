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

// Color represents constants for each chess color White and Black
type Color uint8

// Constants for each color
const (
	White       Color = 0
	Black       Color = 1
	ColorLength int   = 2
)

// Flip returns the opposite color
func (c Color) Flip() Color {
	return c ^ 1
}

// IsValid checks if f represents a valid color
func (c Color) IsValid() bool {
	return c < 2
}

// String returns a string representation of color as "w" or "b"
func (c Color) String() string {
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	default:
		panic(fmt.Sprintf("Invalid color %d", c))
	}
}

// Color direction factor
var moveDirectionFactor = [2]int{1, -1}

// Direction returns positive 1 for White and negative 1 (-1) for Black
func (c Color) Direction() int {
	return moveDirectionFactor[c]
}

// Color pawn move direction
var pawnDir = [2]Direction{North, South}

// MoveDirection returns the direction of a pawn move for the color
func (c Color) MoveDirection() Direction {
	return pawnDir[c]
}

var promRankBb = [2]Bitboard{Rank8_Bb, Rank1_Bb}

// PromotionRankBb returns the rank on which the given color promotes
func (c Color) PromotionRankBb() Bitboard {
	return promRankBb[c]
}

var pawnDoubleRankBb = [2]Bitboard{Rank3_Bb, Rank6_Bb}

// PawnDoubleRank returns the rank from which the given color
// might make a second pawn push to achieve a pawn double move
func (c Color) PawnDoubleRank() Bitboard {
	return pawnDoubleRankBb[c]
}

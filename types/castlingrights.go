/*
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
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

package types

// CastlingRights encodes the castling state e.g. available castling
// and defines functions to change this state
type CastlingRights uint8

// Constants for Castling
const (
	CastlingNone = 0 // 0000

	CastlingWhiteOO  CastlingRights = 1                                  // 0001
	CastlingWhiteOOO CastlingRights = CastlingWhiteOO << 1               // 0010
	CastlingWhite    CastlingRights = CastlingWhiteOO | CastlingWhiteOOO // 0011

	CastlingBlackOO  CastlingRights = CastlingWhiteOO << 2               // 0100
	CastlingBlackOOO CastlingRights = CastlingBlackOO << 1               // 1000
	CastlingBlack    CastlingRights = CastlingBlackOO | CastlingBlackOOO // 1100

	CastlingAny    CastlingRights = CastlingWhite | CastlingBlack // 1111
	CastlingLength CastlingRights = 16
)

// Has checks if the state has the bit for the Castling right set and
// therefore this castling is available
func (lhs CastlingRights) Has(rhs CastlingRights) bool {
	return lhs & rhs > 0
}

// Remove removes a castling right from the input state (deletes right)
func (lhs *CastlingRights) Remove(rhs CastlingRights) CastlingRights {
	*lhs = *lhs & ^rhs
	return	*lhs
}

// Add adds a castling right ti the state
func (lhs *CastlingRights) Add(rhs CastlingRights) CastlingRights {
	*lhs = *lhs | rhs
	return *lhs
}

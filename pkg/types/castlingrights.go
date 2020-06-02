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

import (
	"strings"
)

// CastlingRights encodes the castling state e.g. available castling
// and defines functions to change this state
//  CastlingNone         CastlingRights = 0                                  // 0000
//	CastlingWhiteOO      CastlingRights = 1                                  // 0001
//	CastlingWhiteOOO                    = CastlingWhiteOO << 1               // 0010
//	CastlingWhite                       = CastlingWhiteOO | CastlingWhiteOOO // 0011
//	CastlingBlackOO                     = CastlingWhiteOO << 2               // 0100
//	CastlingBlackOOO                    = CastlingBlackOO << 1               // 1000
//	CastlingBlack                       = CastlingBlackOO | CastlingBlackOOO // 1100
//	CastlingAny                         = CastlingWhite | CastlingBlack      // 1111
//	CastlingRightsLength CastlingRights = 16
type CastlingRights uint8

// Constants for Castling
const (
	CastlingNone         CastlingRights = 0                                  // 0000
	CastlingWhiteOO      CastlingRights = 1                                  // 0001
	CastlingWhiteOOO                    = CastlingWhiteOO << 1               // 0010
	CastlingWhite                       = CastlingWhiteOO | CastlingWhiteOOO // 0011
	CastlingBlackOO                     = CastlingWhiteOO << 2               // 0100
	CastlingBlackOOO                    = CastlingBlackOO << 1               // 1000
	CastlingBlack                       = CastlingBlackOO | CastlingBlackOOO // 1100
	CastlingAny                         = CastlingWhite | CastlingBlack      // 1111
	CastlingRightsLength CastlingRights = 16
)

// Has checks if the state has the bit for the Castling right set and
// therefore this castling is available
func (cr CastlingRights) Has(rhs CastlingRights) bool {
	return cr&rhs != 0
}

// Remove removes a castling right from the input state (deletes right)
func (cr *CastlingRights) Remove(rhs CastlingRights) CastlingRights {
	*cr = *cr &^ rhs
	return *cr
}

// Add adds a castling right ti the state
func (cr *CastlingRights) Add(rhs CastlingRights) CastlingRights {
	*cr = *cr | rhs
	return *cr
}

// String returns a string representation for the castling right instance
// which can be used directly in a fen (e.g. "KQkq")
func (cr CastlingRights) String() string {
	if cr == CastlingNone {
		return "-"
	}
	var os strings.Builder
	if cr.Has(CastlingWhiteOO) {
		os.WriteString("K")
	}
	if cr.Has(CastlingWhiteOOO) {
		os.WriteString("Q")
	}
	if cr.Has(CastlingBlackOO) {
		os.WriteString("k")
	}
	if cr.Has(CastlingBlackOOO) {
		os.WriteString("q")
	}
	return os.String()
}

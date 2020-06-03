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

// PieceType is a set of constants for piece types in chess
//  test for non sliding pt & 0b0100 == 0 (must also be none zero)
//  test for sliding pt & 0b0100 == 1 (must also be < 7)
//  PtNone   = 0b0000
//  King     = 0b0001 // non sliding
//  Pawn     = 0b0010 // non sliding
//  Knight   = 0b0011 // non sliding
//  Bishop   = 0b0100 // sliding
//  Rook     = 0b0101 // sliding
//  Queen    = 0b0110 // sliding
//  PtLength = 0b0111
type PieceType uint8

// PieceType is a set of constants for piece types in chess
const (
	PtNone   PieceType = 0b0000
	King     PieceType = 0b0001
	Pawn     PieceType = 0b0010
	Knight   PieceType = 0b0011
	Bishop   PieceType = 0b0100
	Rook     PieceType = 0b0101
	Queen    PieceType = 0b0110
	PtLength PieceType = 0b0111
)

// IsValid check if pt is a valid piece type
func (pt PieceType) IsValid() bool {
	return pt < 7
}

// array of values for each piece type when calculating game phase
var gamePhaseValue = [PtLength]int{0, 0, 0, 1, 1, 2, 4}

// GamePhaseValue returns a value for calculating game phase
// by adding the number of certain piece type times this value
func (pt PieceType) GamePhaseValue() int {
	return gamePhaseValue[pt]
}

// array of static values of each piece type
var pieceTypeValue = [PtLength]Value{0, 2000, 100, 320, 330, 500, 900}

// ValueOf returns a static value for the piece type
func (pt PieceType) ValueOf() Value {
	return pieceTypeValue[pt]
}

// array of string labels for piece types
var pieceTypeToString = [PtLength]string{"NOPIECE", "King", "Pawn", "Knight", "Bishop", "Rook", "Queen"}

// String returns a string representation of a piece type
func (pt PieceType) String() string {
	return pieceTypeToString[pt]
}

// array of string labels for piece types
var pieceTypeToChar = "-KPNBRQ"

// Char returns a single char string representation of a piece type
func (pt PieceType) Char() string {
	return string(pieceTypeToChar[pt])
}

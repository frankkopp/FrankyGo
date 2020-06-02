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

// Piece is a set of constants for pieces in chess
// Can be used with masks:
//  No Piece = 0
//  White Piece is a non zero value with piece & 0b1000 == 0
//  Black Piece is a non zero value with piece & 0b1000 == 1
//  PieceNone  = 0b0000
//  WhiteKing  = 0b0001
//  WhitePawn  = 0b0010
//  WhiteKnight= 0b0011
//  WhiteBishop= 0b0100
//  WhiteRook  = 0b0101
//  WhiteQueen = 0b0110
//  BlackKing  = 0b1001
//  BlackPawn  = 0b1010
//  BlackKnight= 0b1011
//  BlackBishop= 0b1100
//  BlackRook  = 0b1101
//  BlackQueen = 0b1110
//  PieceLength= 0b10000
type Piece int8

// Pieces are a set of constants to represent the different pieces
// of a chess game.
const (
	PieceNone   Piece = 0
	WhiteKing   Piece = 1
	WhitePawn   Piece = 2
	WhiteKnight Piece = 3
	WhiteBishop Piece = 4
	WhiteRook   Piece = 5
	WhiteQueen  Piece = 6
	BlackKing   Piece = 9
	BlackPawn   Piece = 10
	BlackKnight Piece = 11
	BlackBishop Piece = 12
	BlackRook   Piece = 13
	BlackQueen  Piece = 14
	PieceLength Piece = 16
)

// MakePiece creates the piece given by color and piece type
func MakePiece(c Color, pt PieceType) Piece {
	return Piece((int(c) << 3) + int(pt))
}

// ColorOf returns the color of the given piece */
func (p Piece) ColorOf() Color {
	return Color(p >> 3)
}

// TypeOf returns the piece type of the given piece */
func (p Piece) TypeOf() PieceType {
	return PieceType(p & 7)
}

// ValueOf returns a value for calculating game phase
// by adding the number of certain piece type times this value
func (p Piece) ValueOf() Value {
	return pieceTypeValue[p.TypeOf()]
}

// PieceFromChar returns the Piece corresponding to the given character.
// If s contains not exactly one character or if the character is invalid this
// will return PieceNone
func PieceFromChar(s string) Piece {
	if len(s) != 1 || s == "-" {
		return PieceNone
	}
	index := strings.Index(pieceToString, s)
	if index == -1 {
		return PieceNone
	}
	return Piece(index)
}

// array of string labels for piece types
var pieceToString = string(" KPNBRQ- kpnbrq-")

// String returns a string representation of a piece type
func (p Piece) String() string {
	return string(pieceToString[p])
}

// array of string labels for pieces
var pieceToChar = " KONBRQ- k*nbrq-"

// Char returns a string representation of a piece type
// where pawns are O and * for white and black
func (p Piece) Char() string {
	return string(pieceToChar[p])
}

// array of unicode string labels for pieces
var pieceToUnicode = []string{" ", "♔", "♙", "♘", "♗", "♖", "♕", "-", " ", "♚", "♟", "♞", "♝", "♜", "♛", "-"}

// UniChar returns a unicode string representation of the given pieces
func (p Piece) UniChar() string {
	return pieceToUnicode[p]
}

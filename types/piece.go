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

// Piece is a set of constants for pieces in chess
type Piece int8

// Orientation is a set of constants for moving squares within a Bitboard
//noinspection GoVarAndConstTypeMayBeOmitted
const (
	PieceNone   = 0  // 0b0000
	WhiteKing   = 1  // 0b0001
	WhitePawn   = 2  // 0b0010
	WhiteKnight = 3  // 0b0011
	WhiteBishop = 4  // 0b0100
	WhiteRook   = 5  // 0b0101
	WhiteQueen  = 6  // 0b0110
	BlackKing   = 9  // 0b1001
	BlackPawn   = 10 // 0b1010
	BlackKnight = 11 // 0b1011
	BlackBishop = 12 // 0b1100
	BlackRook   = 13 // 0b1101
	BlackQueen  = 14 // 0b1110
	PieceLength = 16 // 0b10000
)

// array of string labels for piece types
var pieceToString = string("-KPNBRQ--kpnbrq-")

// String returns a string representation of a piece type
func (p Piece) String() string {
	return string(pieceToString[p])
}

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
func (p Piece) ValueOf() int {
	return pieceTypeValue[p.TypeOf()]
}

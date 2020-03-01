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

import (
	"fmt"
	"github.com/frankkopp/FrankyGo/config"
	"log"
	"math/bits"
	"strings"
)

// 64 bit for each square on the board
type Bitboard uint64

// Internal square to bitboard array. Needs to be initialized
var sqBb [64]Bitboard

// Pre computes various bitboards to avoid runtime calculation
// Will only run once (checks an initialized flag)
func initBb() {
	for i := SqA1; i < SqNone; i++ {
		sqBb[i] = i.bitboard_()
	}
}

// Returns a Bitboard of the square by shifting the
// square onto an empty bitboards.
// Usually one would use Bitboard() after initializing with InitBb
func (sq Square) bitboard_() Bitboard {
	return Bitboard(uint64(1) << sq)
}

// Returns a Bitboard of the square by accessing the pre calculated
// square to bitboard array.
// Initialize with InitBb() before use. Throws panic otherwise.
func (sq Square) Bitboard() Bitboard {
	// assertion
	if config.DEBUG && !initialized {
		log.Printf("Warning: Bitboards not initialized. Using runtime calculation.\n")
		return sq.bitboard_()
	}
	return sqBb[sq]
}

// sets the corresponding bit of the bitboard_ for the square
func (b Bitboard) put(s Square) Bitboard {
	return b | s.Bitboard()
}

// removes the corresponding bit of the bitboard_ for the square
func (b Bitboard) remove(s Square) Bitboard {
	return (b | s.Bitboard()) ^ s.Bitboard()
}

// Returns a string representation of the 64 bits
func (b Bitboard) str() string {
	return fmt.Sprintf("%-0.64b", b)
}

// Returns a string representation of the Bitboard
// as a board off 8x8 squares
func (b Bitboard) strBoard() string {
	var os strings.Builder
	os.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := Rank8 + 1; r != Rank1; r-- {
		for f := FileA; f <= FileH; f++ {
			if (b & squareOf(f, r-1).Bitboard()) > 0 {
				os.WriteString("| X ")
			} else {
				os.WriteString("|   ")
			}
		}
		os.WriteString("|\n+---+---+---+---+---+---+---+---+\n")
	}
	return os.String()
}

// Returns a string representation of the 64 bits grouped in 8
// Order is LSB to msb ==> A1 B1 ... G8 H8
func (b Bitboard) strGrp() string {
	var os strings.Builder
	for i := 0; i < 64; i++ {
		if i > 0 && i%8 == 0 {
			os.WriteString(".")
		}
		if (b & (BbOne << i)) != 0 {
			os.WriteString("1")
		} else {
			os.WriteString("0")
		}
	}
	os.WriteString(fmt.Sprintf(" (%d)", b))
	return os.String()
}

// Returns the least significant bit of the 64-bit Bitboard.
// This translates directly to the Square which is returned.
// If the bitboard is empty SqNone will be returned.
// Lsb() indexes from 0-63 - 0 being the the lsb and
// equal to SqA1
func (b Bitboard) Lsb() Square {
	if b == BbZero {
		return SqNone
	}
	return Square(bits.TrailingZeros64(uint64(b)))
}

// Returns the most significant bit of the 64-bit Bitboard.
// This translates directly to the Square which is returned.
// If the bitboard is empty SqNone will be returned.
// Msb() indexes from 0-63 - 63 being the the msb and
// equal to SqH8
func (b Bitboard) Msb() Square {
	if b == BbZero {
		return SqNone
	}
	return Square(63 - bits.LeadingZeros64(uint64(b)))
}

// Returns the Lsb square and removes it from the bitboard.
// The given bitboard is changed directly.
func (b *Bitboard) PopLsb() Square {
	if *b == BbZero {
		return SqNone
	}
	lsb := b.Lsb()
	*b = *b & (*b - 1)
	return lsb
}

// various constant bitboards for convenience
//noinspection ALL
const (
	BbZero Bitboard = Bitboard(0)
	BbAll  Bitboard = ^BbZero
	BbOne  Bitboard = Bitboard(1)

	FileA_Bb Bitboard = 0x0101010101010101
	FileB_Bb Bitboard = FileA_Bb << 1
	FileC_Bb Bitboard = FileA_Bb << 2
	FileD_Bb Bitboard = FileA_Bb << 3
	FileE_Bb Bitboard = FileA_Bb << 4
	FileF_Bb Bitboard = FileA_Bb << 5
	FileG_Bb Bitboard = FileA_Bb << 6
	FileH_Bb Bitboard = FileA_Bb << 7

	Rank1_Bb Bitboard = 0xFF
	Rank2_Bb Bitboard = Rank1_Bb << (8 * 1)
	Rank3_Bb Bitboard = Rank1_Bb << (8 * 2)
	Rank4_Bb Bitboard = Rank1_Bb << (8 * 3)
	Rank5_Bb Bitboard = Rank1_Bb << (8 * 4)
	Rank6_Bb Bitboard = Rank1_Bb << (8 * 5)
	Rank7_Bb Bitboard = Rank1_Bb << (8 * 6)
	Rank8_Bb Bitboard = Rank1_Bb << (8 * 7)

	// Go does not overflow const values when shifting a bit over msb when

	MsbMask   Bitboard = ^(Bitboard(1) << 63)
	Rank8Mask Bitboard = ^Rank8_Bb
	FileAMask Bitboard = ^FileA_Bb
	FileHMask Bitboard = ^FileH_Bb

	DiagUpA1 Bitboard = 0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001
	DiagUpB1 Bitboard = (MsbMask & DiagUpA1) << 1 & FileAMask // shift EAST
	DiagUpC1 Bitboard = (MsbMask & DiagUpB1) << 1 & FileAMask
	DiagUpD1 Bitboard = (MsbMask & DiagUpC1) << 1 & FileAMask
	DiagUpE1 Bitboard = (MsbMask & DiagUpD1) << 1 & FileAMask
	DiagUpF1 Bitboard = (MsbMask & DiagUpE1) << 1 & FileAMask
	DiagUpG1 Bitboard = (MsbMask & DiagUpF1) << 1 & FileAMask
	DiagUpH1 Bitboard = (MsbMask & DiagUpG1) << 1 & FileAMask
	DiagUpA2 Bitboard = (Rank8Mask & DiagUpA1) << 8 // shift NORTH
	DiagUpA3 Bitboard = (Rank8Mask & DiagUpA2) << 8
	DiagUpA4 Bitboard = (Rank8Mask & DiagUpA3) << 8
	DiagUpA5 Bitboard = (Rank8Mask & DiagUpA4) << 8
	DiagUpA6 Bitboard = (Rank8Mask & DiagUpA5) << 8
	DiagUpA7 Bitboard = (Rank8Mask & DiagUpA6) << 8
	DiagUpA8 Bitboard = (Rank8Mask & DiagUpA7) << 8

	DiagDownH1 Bitboard = 0b0000000100000010000001000000100000010000001000000100000010000000
	DiagDownH2 Bitboard = (Rank8Mask & DiagDownH1) << 8 // shift NORTH
	DiagDownH3 Bitboard = (Rank8Mask & DiagDownH2) << 8
	DiagDownH4 Bitboard = (Rank8Mask & DiagDownH3) << 8
	DiagDownH5 Bitboard = (Rank8Mask & DiagDownH4) << 8
	DiagDownH6 Bitboard = (Rank8Mask & DiagDownH5) << 8
	DiagDownH7 Bitboard = (Rank8Mask & DiagDownH6) << 8
	DiagDownH8 Bitboard = (Rank8Mask & DiagDownH7) << 8
	DiagDownG1 Bitboard = (DiagDownH1 >> 1) & FileHMask // shift WEST
	DiagDownF1 Bitboard = (DiagDownG1 >> 1) & FileHMask
	DiagDownE1 Bitboard = (DiagDownF1 >> 1) & FileHMask
	DiagDownD1 Bitboard = (DiagDownE1 >> 1) & FileHMask
	DiagDownC1 Bitboard = (DiagDownD1 >> 1) & FileHMask
	DiagDownB1 Bitboard = (DiagDownC1 >> 1) & FileHMask
	DiagDownA1 Bitboard = (DiagDownB1 >> 1) & FileHMask
)

// ////////////////////
// Pre compute helpers

// Used to pre compute an indexMap for diagonals
func (sq Square) lengthDiagUpMask() Bitboard {
	return (BbOne << lengthDiagUp[sq]) - 1
}

// Used to pre compute an indexMap for diagonals
func (sq Square) lengthDiagDownMask() Bitboard {
	return (BbOne << lengthDiagDown[sq]) - 1
}

// helper arrays
var (
	// Used to pre compute an indexMap for rotated boards
	rotateMapR90 = [SqLength]int{
		7, 15, 23, 31, 39, 47, 55, 63,
		6, 14, 22, 30, 38, 46, 54, 62,
		5, 13, 21, 29, 37, 45, 53, 61,
		4, 12, 20, 28, 36, 44, 52, 60,
		3, 11, 19, 27, 35, 43, 51, 59,
		2, 10, 18, 26, 34, 42, 50, 58,
		1, 9, 17, 25, 33, 41, 49, 57,
		0, 8, 16, 24, 32, 40, 48, 56}

	// Used to pre compute an indexMap for rotated boards
	rotateMapL90 = [SqLength]int{
		56, 48, 40, 32, 24, 16, 8, 0,
		57, 49, 41, 33, 25, 17, 9, 1,
		58, 50, 42, 34, 26, 18, 10, 2,
		59, 51, 43, 35, 27, 19, 11, 3,
		60, 52, 44, 36, 28, 20, 12, 4,
		61, 53, 45, 37, 29, 21, 13, 5,
		62, 54, 46, 38, 30, 22, 14, 6,
		63, 55, 47, 39, 31, 23, 15, 7}

	// Used to pre compute an indexMap for rotated boards
	rotateMapR45 = [SqLength]int{
		7,
		6, 15,
		5, 14, 23,
		4, 13, 22, 31,
		3, 12, 21, 30, 39,
		2, 11, 20, 29, 38, 47,
		1, 10, 19, 28, 37, 46, 55,
		0, 9, 18, 27, 36, 45, 54, 63,
		8, 17, 26, 35, 44, 53, 62,
		16, 25, 34, 43, 52, 61,
		24, 33, 42, 51, 60,
		32, 41, 50, 59,
		40, 49, 58,
		48, 57,
		56}

	// Used to pre compute an indexMap for rotated boards
	rotateMapL45 = [SqLength]int{
		0,
		8, 1,
		16, 9, 2,
		24, 17, 10, 3,
		32, 25, 18, 11, 4,
		40, 33, 26, 19, 12, 5,
		48, 41, 34, 27, 20, 13, 6,
		56, 49, 42, 35, 28, 21, 14, 7,
		57, 50, 43, 36, 29, 22, 15,
		58, 51, 44, 37, 30, 23,
		59, 52, 45, 38, 31,
		60, 53, 46, 39,
		61, 54, 47,
		62, 55,
		63}

	// Used to pre compute an indexMap for diagonals
	lengthDiagUp = [SqLength]int{
		8, 7, 6, 5, 4, 3, 2, 1,
		7, 8, 7, 6, 5, 4, 3, 2,
		6, 7, 8, 7, 6, 5, 4, 3,
		5, 6, 7, 8, 7, 6, 5, 4,
		4, 5, 6, 7, 8, 7, 6, 5,
		3, 4, 5, 6, 7, 8, 7, 6,
		2, 3, 4, 5, 6, 7, 8, 7,
		1, 2, 3, 4, 5, 6, 7, 8}

	// Used to pre compute an indexMap for diagonals
	lengthDiagDown = [SqLength]int{
		1, 2, 3, 4, 5, 6, 7, 8,
		2, 3, 4, 5, 6, 7, 8, 7,
		3, 4, 5, 6, 7, 8, 7, 6,
		4, 5, 6, 7, 8, 7, 6, 5,
		5, 6, 7, 8, 7, 6, 5, 4,
		6, 7, 8, 7, 6, 5, 4, 3,
		7, 8, 7, 6, 5, 4, 3, 2,
		8, 7, 6, 5, 4, 3, 2, 1}

	shiftsDiagUp = [SqLength]int{
		28, 21, 15, 10, 6, 3, 1, 0,
		36, 28, 21, 15, 10, 6, 3, 1,
		43, 36, 28, 21, 15, 10, 6, 3,
		49, 43, 36, 28, 21, 15, 10, 6,
		54, 49, 43, 36, 28, 21, 15, 10,
		58, 54, 49, 43, 36, 28, 21, 15,
		61, 58, 54, 49, 43, 36, 28, 21,
		63, 61, 58, 54, 49, 43, 36, 28}

	shiftsDiagDown = [SqLength]int{
		0, 1, 3, 6, 10, 15, 21, 28,
		1, 3, 6, 10, 15, 21, 28, 36,
		3, 6, 10, 15, 21, 28, 36, 43,
		6, 10, 15, 21, 28, 36, 43, 49,
		10, 15, 21, 28, 36, 43, 49, 54,
		15, 21, 28, 36, 43, 49, 54, 58,
		21, 28, 36, 43, 49, 54, 58, 61,
		28, 36, 43, 49, 54, 58, 61, 63}
)

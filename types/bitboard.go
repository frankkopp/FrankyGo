/*
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, To any person obtaining a copy
 * of this software and associated documentation files (the "Software"), To deal
 * in the Software without restriction, including without limitation the rights
 * To use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and To permit persons To whom the Software is
 * furnished To do so, subject To the following conditions:
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
	"github.com/frankkopp/FrankyGo/util"
	"log"
	"math/bits"
	"strings"
)

// 64 bit for each square on the board
type Bitboard uint64

// Returns a Bitboard of the square by accessing the pre calculated
// square To bitboard array.
// Initialize with InitBb() before use. Throws panic otherwise.
func (sq Square) Bitboard() Bitboard {
	// assertion
	if config.DEBUG && !initialized {
		log.Printf("Warning: Bitboards not initialized. Using runtime calculation.\n")
		return sq.bitboard_()
	}
	return sqBb[sq]
}

// Sets the corresponding bit of the bitboard for the square
func PushSquare(b Bitboard, s Square) Bitboard {
	return b | s.Bitboard()
}

// Sets the corresponding bit of the bitboard for the square
func (b *Bitboard) PushSquare(s Square) {
	*b |= s.Bitboard()
}

// Removes the corresponding bit of the bitboard for the square
func PopSquare(b Bitboard, s Square) Bitboard {
	return (b | s.Bitboard()) ^ s.Bitboard()
}

// Removes the corresponding bit of the bitboard for the square
func (b *Bitboard) PopSquare(s Square) {
	*b = (*b | s.Bitboard()) ^ s.Bitboard()
}

// Shifting all bits of a bitboard in the given direction by 1 square
func ShiftBitboard(b Bitboard, d Direction) Bitboard {
	// move the bits and clear the left our right file
	// after the shift To erase bits jumping over
	switch d {
	case North:
		return (Rank8Mask & b) << 8
	case East:
		return (MsbMask & b) << 1 & FileAMask
	case South:
		return b >> 8
	case West:
		return (b >> 1) & FileHMask
	case Northeast:
		return (Rank8Mask & b) << 9 & FileAMask
	case Southeast:
		return (b >> 7) & FileAMask
	case Southwest:
		return (b >> 9) & FileHMask
	case Northwest:
		return (b << 7) & FileHMask
	}
	return b
}

// Returns the least significant bit of the 64-bit Bitboard.
// This translates directly To the Square which is returned.
// If the bitboard is empty SqNone will be returned.
// Lsb() indexes from 0-63 - 0 being the the lsb and
// equal To SqA1
func (b Bitboard) Lsb() Square {
	return Square(bits.TrailingZeros64(uint64(b)))
}

// Returns the most significant bit of the 64-bit Bitboard.
// This translates directly To the Square which is returned.
// If the bitboard is empty SqNone will be returned.
// Msb() indexes from 0-63 - 63 being the the msb and
// equal To SqH8
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

// Returns a string representation of the 64 bits
func (b Bitboard) Str() string {
	return fmt.Sprintf("%-0.64b", b)
}

// Returns a string representation of the Bitboard
// as a board off 8x8 squares
func (b Bitboard) StrBoard() string {
	var os strings.Builder
	os.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := Rank8 + 1; r != Rank1; r-- {
		for f := FileA; f <= FileH; f++ {
			if (b & SquareOf(f, r-1).Bitboard()) > 0 {
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
// Order is LSB To msb ==> A1 B1 ... G8 H8
func (b Bitboard) StrGrp() string {
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

// Returns the absolute distance in squares between two files
func FileDistance(f1 File, f2 File) int {
	return util.Abs(int(f2) - int(f1))
}

// Returns the absolute distance in squares between two ranks
func RankDistance(r1 Rank, r2 Rank) int {
	return util.Abs(int(r2) - int(r1))
}

// Returns the absolute distance in squares between two squares
func SquareDistance(s1 Square, s2 Square) int {
	return squareDistance[s1][s2]
}


// Rotates a Bitboard by 90 degrees clockwise
func RotateR90(b Bitboard) Bitboard {
	return rotate(b, &rotateMapR90)
}

// Rotates a Bitboard by 90 degrees counter clockwise
func rotateL90(b Bitboard) Bitboard {
	return rotate(b, &rotateMapL90)
}

// Rotates a Bitboard by 45 degrees clockwise
// To get all upward diagonals in compact block of bits
// This is used To create a mask To find moves for
// quuen and bishop on the upward diagonal
func rotateR45(b Bitboard) Bitboard {
	return rotate(b, &rotateMapR45)
}

// Rotates a Bitboard by 45 degrees counter clockwise
// To get all downward diagonals in compact block of bits
// This is used To create a mask To find moves for
// quuen and bishop on the downward diagonal
func rotateL45(b Bitboard) Bitboard {
	return rotate(b, &rotateMapL45)
}

// Maps squares To the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareR90(sq Square) Square { return indexMapR90[sq] }

// Maps squares To the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareL90(sq Square) Square { return indexMapL90[sq] }

// Maps squares To the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareR45(sq Square) Square { return indexMapR45[sq] }

// Maps squares To the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareL45(sq Square) Square { return indexMapL45[sq] }

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
// Private

// Rotates a Bitboard using a mapping array which holds the position of
// the square in the rotated board indexed by the square.
// Basically the array tells bit x To move To bit y
func rotate(b Bitboard, rotationMap *[SqLength]int) Bitboard {
	rotated := BbZero
	for sq := SqA1; sq < SqNone; sq++ {
		if (b & Square(rotationMap[sq]).bitboard_()) != 0 {
			rotated |= sq.bitboard_()
		}
	}
	return rotated
}

// ////////////////////
// Pre compute helpers

// Returns a Bitboard of the square by shifting the
// square onto an empty bitboards.
// Usually one would use Bitboard() after initializing with InitBb
func (sq Square) bitboard_() Bitboard {
	return Bitboard(uint64(1) << sq)
}

// Used To pre compute an indexMap for diagonals
func (sq Square) lengthDiagUpMask() Bitboard {
	return (BbOne << lengthDiagUp[sq]) - 1
}

// Used To pre compute an indexMap for diagonals
func (sq Square) lengthDiagDownMask() Bitboard {
	return (BbOne << lengthDiagDown[sq]) - 1
}

// helper arrays
var (
	// Used To pre compute an indexMap for rotated boards
	rotateMapR90 = [SqLength]int{
		7, 15, 23, 31, 39, 47, 55, 63,
		6, 14, 22, 30, 38, 46, 54, 62,
		5, 13, 21, 29, 37, 45, 53, 61,
		4, 12, 20, 28, 36, 44, 52, 60,
		3, 11, 19, 27, 35, 43, 51, 59,
		2, 10, 18, 26, 34, 42, 50, 58,
		1, 9, 17, 25, 33, 41, 49, 57,
		0, 8, 16, 24, 32, 40, 48, 56}

	// Used To pre compute an indexMap for rotated boards
	rotateMapL90 = [SqLength]int{
		56, 48, 40, 32, 24, 16, 8, 0,
		57, 49, 41, 33, 25, 17, 9, 1,
		58, 50, 42, 34, 26, 18, 10, 2,
		59, 51, 43, 35, 27, 19, 11, 3,
		60, 52, 44, 36, 28, 20, 12, 4,
		61, 53, 45, 37, 29, 21, 13, 5,
		62, 54, 46, 38, 30, 22, 14, 6,
		63, 55, 47, 39, 31, 23, 15, 7}

	// Used To pre compute an indexMap for rotated boards
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

	// Used To pre compute an indexMap for rotated boards
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

	// Used To pre compute an indexMap for diagonals
	lengthDiagUp = [SqLength]int{
		8, 7, 6, 5, 4, 3, 2, 1,
		7, 8, 7, 6, 5, 4, 3, 2,
		6, 7, 8, 7, 6, 5, 4, 3,
		5, 6, 7, 8, 7, 6, 5, 4,
		4, 5, 6, 7, 8, 7, 6, 5,
		3, 4, 5, 6, 7, 8, 7, 6,
		2, 3, 4, 5, 6, 7, 8, 7,
		1, 2, 3, 4, 5, 6, 7, 8}

	// Used To pre compute an indexMap for diagonals
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

	// Reverse index To quickly calculate the index of a square in the rotated board
	indexMapR90 = [SqLength]Square{}
	// Reverse index To quickly calculate the index of a square in the rotated board
	indexMapL90 = [SqLength]Square{}
	// Reverse index To quickly calculate the index of a square in the rotated board
	indexMapR45 = [SqLength]Square{}
	// Reverse index To quickly calculate the index of a square in the rotated board
	indexMapL45 = [SqLength]Square{}
)

// Internal pre computed square To square bitboard array.
// Needs To be initialized with initBb()
var sqBb [SqLength]Bitboard

// Internal pre computed square To file bitboard array.
// Needs To be initialized with initBb()
var sqToFileBb [SqLength]Bitboard

// Internal pre computed square To rank bitboard array.
// Needs To be initialized with initBb()
var sqToRankBb [SqLength]Bitboard

// Internal pre computed square To diag up bitboard array.
// Needs To be initialized with initBb()
var sqDiagUpBb [SqLength]Bitboard

// Internal pre computed square To diag down bitboard array.
// Needs To be initialized with initBb()
var sqDiagDownBb [SqLength]Bitboard

// Internal pre computed index for quick square distance lookup
var squareDistance [SqLength][SqLength]int

// Internal pre computed index To map possible moves on a rank
// for each square and board occupation of this rank
var movesRank [SqLength][256]Bitboard

// Internal pre computed index To map possible moves on a file
// for each square and board occupation of this file
// (needs rotating and masking the index)
var movesFile [SqLength][256]Bitboard

// Internal pre computed index To map possible moves on a up diagonal
// for each square and board occupation of this up diagonal
// (needs rotating and masking the index)
var movesDiagUp [SqLength][256]Bitboard

// Internal pre computed index To map possible moves on a down diagonal
// for each square and board occupation of this down diagonal
// (needs rotating and masking the index)
var movesDiagDown [SqLength][256]Bitboard

// Pre computes various bitboards To avoid runtime calculation
func initBb() {
	for sq := SqA1; sq < SqNone; sq++ {
		sqBb[sq] = sq.bitboard_()

		// file and rank bitboards
		sqToFileBb[sq] = FileA_Bb << sq.FileOf()
		sqToRankBb[sq] = Rank1_Bb << (8 * sq.RankOf())

		// square diagonals // @formatter:off
		if        DiagUpA8&sq.bitboard_() > 0 { sqDiagUpBb[sq] = DiagUpA8
		} else if DiagUpA7&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA7
		} else if DiagUpA6&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA6
		} else if DiagUpA5&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA5
		} else if DiagUpA4&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA4
		} else if DiagUpA3&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA3
		} else if DiagUpA2&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA2
		} else if DiagUpA1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpA1
		} else if DiagUpB1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpB1
		} else if DiagUpC1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpC1
		} else if DiagUpD1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpD1
		} else if DiagUpE1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpE1
		} else if DiagUpF1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpF1
		} else if DiagUpG1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpG1
		} else if DiagUpH1&sq.bitboard_() > 0 {	sqDiagUpBb[sq] = DiagUpH1
		}

		if        DiagDownH8&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH8
		} else if DiagDownH7&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH7
		} else if DiagDownH6&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH6
		} else if DiagDownH5&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH5
		} else if DiagDownH4&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH4
		} else if DiagDownH3&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH3
		} else if DiagDownH2&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH2
		} else if DiagDownH1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownH1
		} else if DiagDownG1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownG1
		} else if DiagDownF1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownF1
		} else if DiagDownE1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownE1
		} else if DiagDownD1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownD1
		} else if DiagDownC1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownC1
		} else if DiagDownB1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownB1
		} else if DiagDownA1&sq.bitboard_() > 0 { sqDiagDownBb[sq] = DiagDownA1
		}
		// @formatter:on

		// Reverse index To quickly calculate the index of a square in the rotated board
		indexMapR90[rotateMapR90[sq]] = sq
		indexMapL90[rotateMapL90[sq]] = sq
		indexMapR45[rotateMapR45[sq]] = sq
		indexMapL45[rotateMapL45[sq]] = sq
	}

	// Distance between squares index
	for sq1 := SqA1; sq1 <= SqH8; sq1++ {
		for sq2 := SqA1; sq2 <= SqH8; sq2++ {
			if sq1 != sq2 {
				squareDistance[sq1][sq2] =
					util.Max(FileDistance(sq1.FileOf(), sq2.FileOf()), RankDistance(sq1.RankOf(), sq2.RankOf()))
			}
		}
	}

	// Pre-compute attacks and moves on a an empty board (pseudo attacks)

	// All sliding attacks with blockers - horizontal
	// Shamefully copied from Beowulf :)
	for file := int(FileA); file <= int(FileH); file++ {
		for j := 0; j < 256; j++ {
			mask := BbZero
			for x := file - 1; x >= 0; x-- {
				mask += BbOne << x
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for x := file + 1; x < 8; x++ {
				mask += BbOne << x
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for rank := int(Rank1); rank <= int(Rank8); rank++ {
				movesRank[(rank*8)+file][j] = mask << (rank * 8)
			}
		}
	}

	// All sliding attacks with blocker - vertical
	// Shamefully copied from Beowulf :)
	for rank := int(Rank1); rank <= int(Rank8); rank++ {
		for j := 0; j < 256; j++ {
			mask := BbZero;
			for x := 6 - rank; x >= 0; x-- {
				mask += BbOne << (8 * (7 - x))
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for x := 8 - rank; x < 8; x++ {
				mask += BbOne << (8 * (7 - x))
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for file := int(FileA); file <= int(FileH); file++ {
				movesFile[(rank * 8) + file][j] = mask << file
			}
		}
	}

	// All sliding attacks with blocker - up diag sliders
	// Shamefully copied from Beowulf :)
	for square := SqA1; square <= SqH8; square++ {
		file := square.FileOf()
		rank := square.RankOf()
		// Get the far left hand square on this diagonal
		diagstart := square - Square(9 * util.Min(int(file), int(rank)))
		dsfile := diagstart.FileOf()
		dl := lengthDiagUp[square]
		// Loop through all possible occupations of this diagonal line
		for sq := 0; sq < (1 << dl); sq++ {
			var mask, mask2 Bitboard
			/* Calculate possible target squares */
			for b1 := int(file) - int(dsfile) - 1; b1 >= 0; b1-- {
				mask += BbOne << b1
				if (sq & (1 << b1)) != 0 {
					break
				}
			}
			for b2 := int(file) - int(dsfile) + 1; b2 < dl; b2++ {
				mask += BbOne << b2
				if (sq & (1 << b2)) != 0 {
					break
				}
			}
			/* Rotate target squares back */
			for x := 0; x < dl; x++ {
				mask2 += ((mask >> x) & 1) << (int(diagstart) + (9 * x))
			}
			movesDiagUp[square][sq] = mask2
		}
	}

	// All sliding attacks with blocker - down diag sliders
	// Shamefully copied from Beowulf :)
	for square := SqA1; square <= SqH8; square++ {
		file := square.FileOf()
		rank := square.RankOf()
		// Get the far left hand square on this diagonal
		diagstart := Square(7 * (util.Min(int(file), 7 - int(rank))) + int(square))
		dsfile := diagstart.FileOf()
		dl := lengthDiagDown[square]
		// Loop through all possible occupations of this diagonal line
		for j := 0; j < (1 << dl); j++ {
			var mask, mask2 Bitboard
			// Calculate possible target squares
			for x := int(file) - int(dsfile) - 1; x >= 0; x-- {
				mask += BbOne << x;
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for x := int(file) - int(dsfile) + 1; x < dl; x++ {
				mask += BbOne << x;
				if (j & (1 << x)) != 0 {
					break
				}
			}
			/* Rotate the target line back onto the required diagonal */
			for x := 0; x < dl; x++  {
				mask2 += ((mask >> x) & 1) << (int(diagstart) - (7 * x))
			}
			movesDiagDown[square][j] = mask2
		}
	}

}

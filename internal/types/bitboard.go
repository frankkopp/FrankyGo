/*
 * FrankyGo - UCI chess engine in GO for learning purposes
 *
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
	"math/bits"
	"strings"

	"github.com/frankkopp/FrankyGo/internal/util"
)

// Bitboard is a 64 bit unsigned int with 1 bit for each square on the board
type Bitboard uint64

// Bb returns a Bitboard of the square by accessing the pre calculated
// square to bitboard array.
func (sq Square) Bb() Bitboard {
	return sqBb[sq]
}

// PushSquare sets the corresponding bit of the bitboard for the square
func PushSquare(b Bitboard, s Square) Bitboard {
	return b | s.Bb()
}

// PushSquare sets the corresponding bit of the bitboard for the square
func (b *Bitboard) PushSquare(s Square) Bitboard {
	*b |= s.Bb()
	return *b
}

// PopSquare removes the corresponding bit of the bitboard for the square
func PopSquare(b Bitboard, s Square) Bitboard {
	return b &^ s.Bb()
}

// PopSquare removes the corresponding bit of the bitboard for the square
func (b *Bitboard) PopSquare(s Square) Bitboard {
	*b = *b &^ s.Bb()
	return *b
}

// Has tests if a square (bit) is set
func (b Bitboard) Has(s Square) bool {
	return b&sqBb[s] != 0
}

// ShiftBitboard shifting all bits of a bitboard in the given direction by 1 square
func ShiftBitboard(b Bitboard, d Direction) Bitboard {
	// move the bits and clear the left our right file
	// after the shift to erase bits jumping over
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

// GetMovesOnRank returns a Bb for all possible horizontal moves
// on the rank of the square with the rank content (blocking pieces)
// determined from the given pieces bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesOnRank(sq Square, content Bitboard) Bitboard {
	// content = the pieces currently on the board and maybe blocking the moves
	// no rotation necessary for ranks - their squares are already in a row
	// shift to the least significant bit
	contentIdx := content >> (8 * int(sq.RankOf()))
	// retrieve all possible moves for this square with the current content
	// and mask with the first row to erase any other pieces
	return movesRank[sq][contentIdx&255]
}

// GetMovesOnFileRotated Bb for all possible horizontal moves on the
// rank of the square with the rank content (blocking pieces) determined
// from the given L90 rotated bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesOnFileRotated(sq Square, rotated Bitboard) Bitboard {
	// shift to the lsb
	contentIdx := rotated >> (int(sq.FileOf()) * 8)
	// retrieve all possible moves for this square with the current content
	// and mask with the first row to erase any other pieces not erased by shift
	return movesFile[sq][contentIdx&255]
}

// GetMovesOnFile Bb for all possible horizontal moves on the rank of
// the square with the rank content (blocking pieces) determined from the
// given bitboard (not rotated - use GetMovesOnFileRotated for already rotated
// bitboards)
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesOnFile(sq Square, content Bitboard) Bitboard {
	// content = the pieces currently on the board and maybe blocking the moves
	// rotate the content of the board to get all file squares in a row
	return GetMovesOnFileRotated(sq, RotateL90(content))
}

// GetMovesDiagUpRotated  Bb for all possible diagonal up moves of
// the square with the content (blocking pieces) determined from the
// given R45 rotated bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesDiagUpRotated(sq Square, rotated Bitboard) Bitboard {
	// shift the correct row to the lsb
	shifted := rotated >> shiftsDiagUp[sq]
	// mask the content with the length of the diagonal to erase any other
	// pieces which have not been erased by the shift
	contentMasked := shifted & ((BbOne << lengthDiagUp[sq]) - 1)
	// retrieve all possible moves for this square with the current content
	return movesDiagUp[sq][contentMasked]
}

// GetMovesDiagUp Bb for all possible diagonal up moves of the square with
// the content (blocking pieces) determined from the given non rotated
// bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesDiagUp(sq Square, content Bitboard) Bitboard {
	// content = the pieces currently on the board and maybe blocking the moves
	// rotate the content of the board to get all diagonals in a row
	return GetMovesDiagUpRotated(sq, RotateR45(content))
}

// GetMovesDiagDownRotated Bb for all possible diagonal up moves of the square with
// the content (blocking pieces) determined from the given L45 rotated
// bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesDiagDownRotated(sq Square, rotated Bitboard) Bitboard {
	// shift the correct row to the lsb
	shifted := rotated >> shiftsDiagDown[sq]
	// mask the content with the length of the diagonal to erase any other
	// pieces which have not been erased by the shift
	contentMasked := shifted & ((BbOne << lengthDiagDown[sq]) - 1)
	// retrieve all possible moves for this square with the current content
	return movesDiagDown[sq][contentMasked]
}

// GetMovesDiagDown Bb for all possible diagonal up moves of the square with
// the content (blocking pieces) determined from the given non rotated
// bitboard.
//
// Deprecated
// use GetAttacksBb(pt PieceType, sq Square, occupied Bitboard)
func GetMovesDiagDown(square Square, content Bitboard) Bitboard {
	// content = the pieces currently on the board and maybe blocking the moves
	// rotate the content of the board to get all diagonals in a row
	return GetMovesDiagDownRotated(square, RotateL45(content))
}

// Lsb returns the least significant bit of the 64-bit Bb.
// This translates directly to the Square which is returned.
// If the bitboard is empty SqNone will be returned.
// Lsb() indexes from 0-63 - 0 being the the lsb and
// equal to SqA1
func (b Bitboard) Lsb() Square {
	return Square(bits.TrailingZeros64(uint64(b)))
}

// Msb returns the most significant bit of the 64-bit Bb.
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

// PopLsb returns the Lsb square and removes it from the bitboard.
// The given bitboard is changed directly.
func (b *Bitboard) PopLsb() Square {
	if *b == BbZero {
		return SqNone
	}
	lsb := b.Lsb()
	*b = *b & (*b - 1)
	return lsb
}

// PopCount returns the number of one bits ("population count") in b.
// This equals the number of squares set in a Bitboard
func (b Bitboard) PopCount() int {
	return bits.OnesCount64(uint64(b))
}

// String returns a string representation of the 64 bits
func (b Bitboard) String() string {
	return fmt.Sprintf("%-0.64b", b)
}

// StringBoard returns a string representation of the Bb
// as a board off 8x8 squares
func (b Bitboard) StringBoard() string {
	var os strings.Builder
	os.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := Rank1; r <= Rank8; r++ {
		for f := FileA; f <= FileH; f++ {
			if (b & SquareOf(f, Rank8-r).Bb()) > 0 {
				os.WriteString("| X ")
			} else {
				os.WriteString("|   ")
			}
		}
		os.WriteString("|\n+---+---+---+---+---+---+---+---+\n")
	}
	return os.String()
}

// StringGrouped returns a string representation of the 64 bits grouped in 8.
// Order is LSB to msb ==> A1 B1 ... G8 H8
func (b Bitboard) StringGrouped() string {
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

// FileDistance returns the absolute distance in squares between two files
func FileDistance(f1 File, f2 File) int {
	return util.Abs(int(f2) - int(f1))
}

// RankDistance returns the absolute distance in squares between two ranks
func RankDistance(r1 Rank, r2 Rank) int {
	return util.Abs(int(r2) - int(r1))
}

// SquareDistance returns the absolute distance in squares between two squares
func SquareDistance(s1 Square, s2 Square) int {
	if !s1.IsValid() || !s2.IsValid() || s1 == s2 {
		return 0
	}
	return squareDistance[s1][s2]
}

// CenterDistance returns the distance to the nearest center square
func (sq Square) CenterDistance() int {
	return centerDistance[sq]
}

// RotateR90 rotates a Bb by 90 degrees clockwise
func RotateR90(b Bitboard) Bitboard {
	return rotate(b, &rotateMapR90)
}

// RotateL90 rotates a Bb by 90 degrees counter clockwise
func RotateL90(b Bitboard) Bitboard {
	return rotate(b, &rotateMapL90)
}

// RotateR45 rotates a Bb by 45 degrees clockwise
// to get all upward diagonals in compact block of bits
// This is used to create a mask to find moves for
// queen and bishop on the upward diagonal
func RotateR45(b Bitboard) Bitboard {
	return rotate(b, &rotateMapR45)
}

// RotateL45 rotates a Bb by 45 degrees counter clockwise
// to get all downward diagonals in compact block of bits
// This is used to create a mask to find moves for
// queen and bishop on the downward diagonal
func RotateL45(b Bitboard) Bitboard {
	return rotate(b, &rotateMapL45)
}

// RotateSquareR90 maps squares to the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareR90(sq Square) Square {
	return indexMapR90[sq]
}

// RotateSquareL90 maps squares to the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareL90(sq Square) Square {
	return indexMapL90[sq]
}

// RotateSquareR45 maps squares to the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareR45(sq Square) Square {
	return indexMapR45[sq]
}

// RotateSquareL45 maps squares to the sq of the rotated board. E.g. when rotating
// clockwise by 90 degree A1 becomes A8, A8 becomes H8, etc.
func RotateSquareL45(sq Square) Square {
	return indexMapL45[sq]
}

// GetAttacksBb returns a bitboard representing all the squares attacked by a
// piece of the given type pt (not pawn) placed on 'sq'.
// For sliding pieces this uses the pre-computed Magic Bitboard Attack arrays.
// For Knight and King this the occupied Bitboard is ignored (can be BbZero)
// as for these non sliders the pre-computed pseudo attacks are used
func GetAttacksBb(pt PieceType, sq Square, occupied Bitboard) Bitboard {
	if pt == Pawn {
		msg := fmt.Sprint("GetAttackBb called with piece type Pawn is not supported")
		panic(msg)
	}
	switch pt {
	case Bishop:
		return bishopMagics[sq].Attacks[bishopMagics[sq].index(occupied)]
	case Rook:
		return rookMagics[sq].Attacks[rookMagics[sq].index(occupied)]
	case Queen:
		return bishopMagics[sq].Attacks[bishopMagics[sq].index(occupied)] | rookMagics[sq].Attacks[rookMagics[sq].index(occupied)]
	default:
		return pseudoAttacks[pt][sq]
	}
}

// GetPseudoAttacks returns a Bb of possible attacks of a piece
// as if on an empty board
func GetPseudoAttacks(pt PieceType, sq Square) Bitboard {
	return pseudoAttacks[pt][sq]
}

// GetPawnAttacks returns a Bb of possible attacks of a pawn
func GetPawnAttacks(c Color, sq Square) Bitboard {
	return pawnAttacks[c][sq]
}

// FilesWestMask returns a Bb of the files west of the square
func (sq Square) FilesWestMask() Bitboard {
	return filesWestMask[sq]
}

// FilesEastMask returns a Bb of the files east of the square
func (sq Square) FilesEastMask() Bitboard {
	return filesEastMask[sq]
}

// FileWestMask returns a Bb of the file west of the square
func (sq Square) FileWestMask() Bitboard {
	return fileWestMask[sq]
}

// FileEastMask returns a Bb of the file east of the square
func (sq Square) FileEastMask() Bitboard {
	return fileEastMask[sq]
}

// RanksNorthMask returns a Bb of the ranks north of the square
func (sq Square) RanksNorthMask() Bitboard {
	return ranksNorthMask[sq]
}

// RanksSouthMask returns a Bb of the ranks south of the square
func (sq Square) RanksSouthMask() Bitboard {
	return ranksSouthMask[sq]
}

// NeighbourFilesMask returns a Bb of the file east and west of the square
func (sq Square) NeighbourFilesMask() Bitboard {
	return neighbourFilesMask[sq]
}

// Ray returns a Bb of squares outgoing from the
// square in direction of the orientation
func (sq Square) Ray(o Orientation) Bitboard {
	return rays[o][sq]
}

// Intermediate returns a Bb of squares between
// the given two squares
func Intermediate(sq1 Square, sq2 Square) Bitboard {
	return intermediate[sq1][sq2]
}

// Intermediate returns a Bb of squares between
// the given two squares
func (sq Square) Intermediate(sqTo Square) Bitboard {
	return intermediate[sq][sqTo]
}

// PassedPawnMask returns a Bitboards with all possible squares
// which have an opponents pawn which could stop this pawn.
// Use this mask and AND it with the opponents pawns bitboards
// to see if a pawn has passed.
func (sq Square) PassedPawnMask(c Color) Bitboard {
	return passedPawnMask[c][sq]
}

// KingSideCastleMask returns a Bb with the kings side
// squares used in castling without the king square
func KingSideCastleMask(c Color) Bitboard {
	return kingSideCastleMask[c]
}

// QueenSideCastMask returns a Bb with the queen side
// squares used in castling without the king square
func QueenSideCastMask(c Color) Bitboard {
	return queenSideCastleMask[c]
}

// GetCastlingRights returns the CastlingRights for
// changes on this square.
func GetCastlingRights(sq Square) CastlingRights {
	return castlingRights[sq]
}

// SquaresBb returns a Bb of all squares of the given color.
// E.g. can be used to find bishops of the same "color" for draw detection.
func SquaresBb(c Color) Bitboard {
	return squaresBb[c]
}

// Various constant bitboards
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

	CenterFiles   Bitboard = FileD_Bb | FileE_Bb
	CenterRanks   Bitboard = Rank4_Bb | Rank5_Bb
	CenterSquares Bitboard = CenterFiles & CenterRanks
)

// ////////////////////
// Private
// ////////////////////

// Rotates a Bb using a mapping array which holds the position of
// the square in the rotated board indexed by the square.
// Basically the array tells bit x to move to bit y
func rotate(b Bitboard, rotationMap *[SqLength]int) Bitboard {
	rotated := BbZero
	for sq := SqA1; sq < SqNone; sq++ {
		if (b & sqBb[Square(rotationMap[sq])]) != 0 {
			rotated |= sqBb[sq]
		}
	}
	return rotated
}

// ////////////////////
// Pre compute helpers

// Returns a Bb of the square by shifting the
// square onto an empty bitboards.
// Usually one would use Bb() after initializing with InitBb
func (sq Square) bitboard() Bitboard {
	return Bitboard(uint64(1) << sq)
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

	// Reverse index to quickly calculate the index of a square in the rotated board
	indexMapR90 = [SqLength]Square{}
	// Reverse index to quickly calculate the index of a square in the rotated board
	indexMapL90 = [SqLength]Square{}
	// Reverse index to quickly calculate the index of a square in the rotated board
	indexMapR45 = [SqLength]Square{}
	// Reverse index to quickly calculate the index of a square in the rotated board
	indexMapL45 = [SqLength]Square{}

	// Internal pre computed square to square bitboard array.
	// Needs to be initialized with initBb()
	sqBb [SqLength]Bitboard

	// Internal pre computed square to file bitboard array.
	// Needs to be initialized with initBb()
	sqToFileBb [SqLength]Bitboard

	// Internal pre computed square to rank bitboard array.
	// Needs to be initialized with initBb()
	sqToRankBb [SqLength]Bitboard

	// Internal pre computed square to diag up bitboard array.
	// Needs to be initialized with initBb()
	sqDiagUpBb [SqLength]Bitboard

	// Internal pre computed square to diag down bitboard array.
	// Needs to be initialized with initBb()
	sqDiagDownBb [SqLength]Bitboard

	// Internal pre computed rank bitboard array.
	// Needs to be initialized with initBb()
	rankBb [8]Bitboard

	// Internal pre computed file bitboard array.
	// Needs to be initialized with initBb()
	fileBb [8]Bitboard

	// Internal pre computed index for quick square distance lookup
	squareDistance [SqLength][SqLength]int

	// Internal pre computed index to map possible moves on a rank
	// for each square and board occupation of this rank
	movesRank [SqLength][256]Bitboard

	// Internal pre computed index to map possible moves on a file
	// for each square and board occupation of this file
	// (needs rotating and masking the index)
	movesFile [SqLength][256]Bitboard

	// Internal pre computed index to map possible moves on a up diagonal
	// for each square and board occupation of this up diagonal
	// (needs rotating and masking the index)
	movesDiagUp [SqLength][256]Bitboard

	// Internal pre computed index to map possible moves on a down diagonal
	// for each square and board occupation of this down diagonal
	// (needs rotating and masking the index)
	movesDiagDown [SqLength][256]Bitboard

	// Internal Bb for pawn attacks for each color for each square
	pawnAttacks [2][SqLength]Bitboard

	// Internal Bb for attacks for each piece for each square
	pseudoAttacks [PtLength][SqLength]Bitboard

	// magic bitboards - rook attacks
	rookTable  []Bitboard
	rookMagics [SqLength]Magic

	// magic bitboards - bishop attacks
	bishopTable  []Bitboard
	bishopMagics [SqLength]Magic

	// Internal pre computed bitboards
	filesWestMask      [SqLength]Bitboard
	filesEastMask      [SqLength]Bitboard
	ranksNorthMask     [SqLength]Bitboard
	ranksSouthMask     [SqLength]Bitboard
	fileWestMask       [SqLength]Bitboard
	fileEastMask       [SqLength]Bitboard
	neighbourFilesMask [SqLength]Bitboard

	// Internal pre computed arrays of rays which
	// have a bitboard per orientation and square
	rays [8][SqLength]Bitboard

	// intermediate holds bitboards for the squares between
	// to squares
	intermediate [SqLength][SqLength]Bitboard

	// mask to determine of pawn is passed e.g. has no
	// opponent pawns on the same file or the neighbour
	// files
	passedPawnMask [2][SqLength]Bitboard

	// helper mask for supporting castling moves
	kingSideCastleMask [2]Bitboard
	// helper mask for supporting castling moves
	queenSideCastleMask [2]Bitboard

	// array to store all possible CastlingRights for squares which impact castlings
	castlingRights [SqLength]CastlingRights

	// mask for all white  and black squares
	squaresBb [2]Bitboard

	// array with distance of a square to the center
	centerDistance [SqLength]int
)

// ///////////////////////////////////////
// Initialization
// ///////////////////////////////////////

// Pre computes various bitboards to avoid runtime calculation
func initBb() {
	squareBitboardsPreCompute()
	rankFileBbPreCompute()
	castleMasksPreCompute()
	squareDistancePreCompute()
	movesRankPreCompute()
	movesFilePreCompute()
	movesDiagUpPreCompute()
	movesDiagDownPreCompute()
	pseudoAttacksPreCompute()
	neighbourMasksPreCompute()
	raysPreCompute()
	intermediatePreCompute()
	maskPassedPawnsPreCompute()
	squareColorsPreCompute()
	centerDistancePreCompute()
	initMagicBitboards()
}

// start calculating the magic bitboards
// Taken from Stockfish and
// from  https://www.chessprogramming.org/Magic_Bitboards
func initMagicBitboards() {
	rookDirections := [4]Direction{North, East, South, West}
	bishopDirections := [4]Direction{Northeast, Southeast, Southwest, Northwest}

	rookTable = make([]Bitboard, 0x19000, 0x19000)
	bishopTable = make([]Bitboard, 0x1480, 0x1480)

	initMagics(&rookTable, &rookMagics, &rookDirections)
	initMagics(&bishopTable, &bishopMagics, &bishopDirections)
}

func rankFileBbPreCompute() {
	for i := Rank1; i <= Rank8; i++ {
		rankBb[i] = Rank1_Bb << (8 * i)
	}
	for i := FileA; i <= FileH; i++ {
		fileBb[i] = FileA_Bb << i
	}
}

func castleMasksPreCompute() {
	kingSideCastleMask[White] = sqBb[SqF1] | sqBb[SqG1] | sqBb[SqH1]
	kingSideCastleMask[Black] = sqBb[SqF8] | sqBb[SqG8] | sqBb[SqH8]
	queenSideCastleMask[White] = sqBb[SqD1] | sqBb[SqC1] | sqBb[SqB1] | sqBb[SqA1]
	queenSideCastleMask[Black] = sqBb[SqD8] | sqBb[SqC8] | sqBb[SqB8] | sqBb[SqA8]
	castlingRights[SqE1] = CastlingWhite
	castlingRights[SqA1] = CastlingWhiteOOO
	castlingRights[SqH1] = CastlingWhiteOO
	castlingRights[SqE8] = CastlingBlack
	castlingRights[SqA8] = CastlingBlackOOO
	castlingRights[SqH8] = CastlingBlackOO
}

func squareBitboardsPreCompute() {
	for sq := SqA1; sq < SqNone; sq++ {
		// pre compute bitboard for a single sq
		sqBb[sq] = sq.bitboard()

		// file and rank bitboards
		sqToFileBb[sq] = FileA_Bb << sq.FileOf()
		sqToRankBb[sq] = Rank1_Bb << (8 * sq.RankOf())

		// sq diagonals // @formatter:off
		//noinspection GoLinterLocal
		if        DiagUpA8&sq.bitboard() > 0 { sqDiagUpBb[sq] = DiagUpA8
		} else if DiagUpA7&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA7
		} else if DiagUpA6&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA6
		} else if DiagUpA5&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA5
		} else if DiagUpA4&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA4
		} else if DiagUpA3&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA3
		} else if DiagUpA2&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA2
		} else if DiagUpA1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpA1
		} else if DiagUpB1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpB1
		} else if DiagUpC1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpC1
		} else if DiagUpD1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpD1
		} else if DiagUpE1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpE1
		} else if DiagUpF1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpF1
		} else if DiagUpG1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpG1
		} else if DiagUpH1&sq.bitboard() > 0 {	sqDiagUpBb[sq] = DiagUpH1
		}

		if        DiagDownH8&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH8
		} else if DiagDownH7&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH7
		} else if DiagDownH6&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH6
		} else if DiagDownH5&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH5
		} else if DiagDownH4&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH4
		} else if DiagDownH3&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH3
		} else if DiagDownH2&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH2
		} else if DiagDownH1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownH1
		} else if DiagDownG1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownG1
		} else if DiagDownF1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownF1
		} else if DiagDownE1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownE1
		} else if DiagDownD1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownD1
		} else if DiagDownC1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownC1
		} else if DiagDownB1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownB1
		} else if DiagDownA1&sq.bitboard() > 0 { sqDiagDownBb[sq] = DiagDownA1
		}
		// @formatter:on

		// Reverse index to quickly calculate the index of a sq in the rotated board
		indexMapR90[rotateMapR90[sq]] = sq
		indexMapL90[rotateMapL90[sq]] = sq
		indexMapR45[rotateMapR45[sq]] = sq
		indexMapL45[rotateMapL45[sq]] = sq
	}
}

// pre computes distances to center squares by quadrant
func centerDistancePreCompute() {
	for square := SqA1; square <= SqH8; square++ {
		// left upper quadrant
		if (sqBb[square] & ranksNorthMask[27] & filesWestMask[36]) != 0 {
			centerDistance[square] = squareDistance[square][SqD5]
			// right upper quadrant
		} else if (sqBb[square] & ranksNorthMask[28] & filesEastMask[35]) != 0 {
			centerDistance[square] = squareDistance[square][SqE5]
			// left lower quadrant
		} else if (sqBb[square] & ranksSouthMask[35] & filesWestMask[28]) != 0 {
			centerDistance[square] = squareDistance[square][SqD4]
			// right lower quadrant
		} else if (sqBb[square] & ranksSouthMask[36] & filesEastMask[27]) != 0 {
			centerDistance[square] = squareDistance[square][SqE4]
		}
	}
}

// masks for each square color (good for bishops vs bishops or pawns)
func squareColorsPreCompute() {
	for square := SqA1; square <= SqH8; square++ {
		f := square.FileOf()
		r := square.RankOf()
		if (int(f)+int(r))%2 == 0 {
			squaresBb[Black] |= BbOne << square
		} else {
			squaresBb[White] |= BbOne << square
		}
	}
}

// pre computes passed pawn masks
func maskPassedPawnsPreCompute() {
	for square := SqA1; square <= SqH8; square++ {
		f := square.FileOf()
		r := square.RankOf()
		// white pawn - ignore that pawns can'*t be on all squares
		passedPawnMask[White][square] |= rays[N][square]
		if f < 7 && r < 7 {
			passedPawnMask[White][square] |= rays[N][square.To(East)]
		}
		if f > 0 && r < 7 {
			passedPawnMask[White][square] |= rays[N][square.To(West)]
		}
		// black pawn - ignore that pawns can'*t be on all squares
		passedPawnMask[Black][square] |= rays[S][square]
		if f < 7 && r > 0 {
			passedPawnMask[Black][square] |= rays[S][square.To(East)]
		}
		if f > 0 && r > 0 {
			passedPawnMask[Black][square] |= rays[S][square.To(West)]
		}
	}
}

// mask for intermediate squares in between two squares
func intermediatePreCompute() {
	for from := SqA1; from <= SqH8; from++ {
		for to := SqA1; to <= SqH8; to++ {
			toBB := sqBb[to]
			for o := 0; o < 8; o++ {
				if rays[Orientation(o)][from]&toBB != BbZero {
					intermediate[from][to] |=
						rays[Orientation(o)][from] & ^rays[Orientation(o)][to] & ^toBB
				}
			}
		}
	}
}

func raysPreCompute() {
	for sq := SqA1; sq <= SqH8; sq++ {
		rays[N][sq] = pseudoAttacks[Rook][sq] & ranksNorthMask[sq]
		rays[E][sq] = pseudoAttacks[Rook][sq] & filesEastMask[sq]
		rays[S][sq] = pseudoAttacks[Rook][sq] & ranksSouthMask[sq]
		rays[W][sq] = pseudoAttacks[Rook][sq] & filesWestMask[sq]

		rays[NW][sq] = pseudoAttacks[Bishop][sq] & filesWestMask[sq] & ranksNorthMask[sq]
		rays[NE][sq] = pseudoAttacks[Bishop][sq] & filesEastMask[sq] & ranksNorthMask[sq]
		rays[SE][sq] = pseudoAttacks[Bishop][sq] & filesEastMask[sq] & ranksSouthMask[sq]
		rays[SW][sq] = pseudoAttacks[Bishop][sq] & filesWestMask[sq] & ranksSouthMask[sq]
	}
}

// masks for files and ranks left, right, up and down from sq
func neighbourMasksPreCompute() {
	for square := SqA1; square <= SqH8; square++ {
		f := int(square.FileOf())
		r := int(square.RankOf())
		for j := 0; j <= 7; j++ {
			// file masks
			if j < f {
				filesWestMask[square] |= FileA_Bb << j
			}
			if 7-j > f {
				filesEastMask[square] |= FileA_Bb << (7 - j)
			}
			// rank masks
			if 7-j > r {
				ranksNorthMask[square] |= Rank1_Bb << (8 * (7 - j))
			}
			if j < r {
				ranksSouthMask[square] |= Rank1_Bb << (8 * j)
			}
		}
		if f > 0 {
			fileWestMask[square] = FileA_Bb << (f - 1)
		}
		if f < 7 {
			fileEastMask[square] = FileA_Bb << (f + 1)
		}
		neighbourFilesMask[square] = fileEastMask[square] | fileWestMask[square]
	}
}

// Distance between squares index
func squareDistancePreCompute() {
	for sq1 := SqA1; sq1 <= SqH8; sq1++ {
		for sq2 := SqA1; sq2 <= SqH8; sq2++ {
			if sq1 != sq2 {
				squareDistance[sq1][sq2] =
					util.Max(FileDistance(sq1.FileOf(), sq2.FileOf()), RankDistance(sq1.RankOf(), sq2.RankOf()))
			}
		}
	}
}

// pre compute all possible attacked sq per color, piece and sq
func pseudoAttacksPreCompute() {
	// steps for kings, pawns, knight for WHITE - negate to get BLACK
	var steps = [][]Direction{
		{},
		{Northwest, North, Northeast, East}, // king
		{Northwest, Northeast},              // pawn
		{West + Northwest, East + Northeast, North + Northwest, North + Northeast}} // knight

	// non-sliding attacks
	for c := White; c <= Black; c++ {
		for _, pt := range []PieceType{King, Pawn, Knight} {
			for s := SqA1; s <= SqH8; s++ {
				for i := 0; i < len(steps[pt]); i++ {
					to := Square(int(s) + c.Direction()*int(steps[pt][i]))
					if to.IsValid() && squareDistance[s][to] < 3 { // no wrap around board edges
						if pt == Pawn {
							pawnAttacks[c][s] |= sqBb[to]
						} else {
							pseudoAttacks[pt][s] |= sqBb[to]
						}
					}
				}
			}
		}
	}

	// sliding pieces pseudo attacks
	for square := SqA1; square <= SqH8; square++ {
		pseudoAttacks[Bishop][square] |= movesDiagUp[square][0]
		pseudoAttacks[Bishop][square] |= movesDiagDown[square][0]
		pseudoAttacks[Rook][square] |= movesFile[square][0]
		pseudoAttacks[Rook][square] |= movesRank[square][0]
		pseudoAttacks[Queen][square] |= pseudoAttacks[Bishop][square] | pseudoAttacks[Rook][square]
	}
}

// Pre-compute attacks and moves on a an empty board (pseudo attacks)
func movesDiagDownPreCompute() {
	// All sliding attacks with blocker - down diag sliders
	// Shamefully copied from Beowulf :)
	for square := SqA1; square <= SqH8; square++ {
		file := square.FileOf()
		rank := square.RankOf()
		// Get the far left hand square on this diagonal
		diagstart := Square(7*(util.Min(int(file), 7-int(rank))) + int(square))
		dsfile := diagstart.FileOf()
		dl := lengthDiagDown[square]
		// Loop through all possible occupations of this diagonal line
		for j := 0; j < (1 << dl); j++ {
			var mask, mask2 Bitboard
			// Calculate possible target squares
			for x := int(file) - int(dsfile) - 1; x >= 0; x-- {
				mask += BbOne << x
				if (j & (1 << x)) != 0 {
					break
				}
			}
			for x := int(file) - int(dsfile) + 1; x < dl; x++ {
				mask += BbOne << x
				if (j & (1 << x)) != 0 {
					break
				}
			}
			/* Rotate the target line back onto the required diagonal */
			for x := 0; x < dl; x++ {
				mask2 += ((mask >> x) & 1) << (int(diagstart) - (7 * x))
			}
			movesDiagDown[square][j] = mask2
		}
	}
}

// Pre-compute attacks and moves on a an empty board (pseudo attacks)
func movesDiagUpPreCompute() {
	// All sliding attacks with blocker - up diag sliders
	// Shamefully copied from Beowulf :)
	for square := SqA1; square <= SqH8; square++ {
		file := square.FileOf()
		rank := square.RankOf()
		// Get the far left hand square on this diagonal
		diagstart := square - Square(9*util.Min(int(file), int(rank)))
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
}

// Pre-compute attacks and moves on a an empty board (pseudo attacks)
func movesFilePreCompute() {
	// All sliding attacks with blocker - vertical
	// Shamefully copied from Beowulf :)
	for rank := int(Rank1); rank <= int(Rank8); rank++ {
		for j := 0; j < 256; j++ {
			mask := BbZero
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
				movesFile[(rank*8)+file][j] = mask << file
			}
		}
	}
}

// Pre-compute attacks and moves on a an empty board (pseudo attacks)
func movesRankPreCompute() {
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
}

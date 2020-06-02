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
	"fmt"
	"math/bits"
	"strings"

	"github.com/frankkopp/FrankyGo/internal/util"
)

// Bitboard is a 64 bit unsigned int with 1 bit for each square on the board
type Bitboard uint64

// Bb returns a Bitboard of the given file
func (f File) Bb() Bitboard {
	return fileBb[f]
}

// Bb returns a Bitboard of the given rank
func (r Rank) Bb() Bitboard {
	return rankBb[r]
}

// Bb returns a Bitboard of the square
func (sq Square) Bb() Bitboard {
	return sqBb[sq]
}

// FileBb returns a Bitboard of the file of this square
func (sq Square) FileBb() Bitboard {
	return sqToFileBb[sq]
}

// RankBb returns a Bitboard of the rank of this square
func (sq Square) RankBb() Bitboard {
	return sqToRankBb[sq]
}

// PushSquare sets the corresponding bit of the bitboard for the square
func PushSquare(b Bitboard, sq Square) Bitboard {
	return b | sqBb[sq]
}

// PushSquare sets the corresponding bit of the bitboard for the square
func (b *Bitboard) PushSquare(sq Square) Bitboard {
	*b |= sqBb[sq]
	return *b
}

// PopSquare removes the corresponding bit of the bitboard for the square
func PopSquare(b Bitboard, sq Square) Bitboard {
	return b &^ sqBb[sq]
}

// PopSquare removes the corresponding bit of the bitboard for the square
func (b *Bitboard) PopSquare(sq Square) Bitboard {
	*b = *b &^ sqBb[sq]
	return *b
}

// Has tests if a square (bit) is set
func (b Bitboard) Has(sq Square) bool {
	return b&sqBb[sq] != 0
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

// GetAttacksBb returns a bitboard representing all the squares attacked by a
// piece of the given type pt (not pawn) placed on 'sq'.
// For sliding pieces this uses the pre-computed Magic Bitboard Attack arrays.
// For Knight and King this the occupied Bitboard is ignored (can be BbZero)
// as for these non sliders the pre-computed pseudo attacks are used
func GetAttacksBb(pt PieceType, sq Square, occupied Bitboard) Bitboard {
	switch pt {
	case Bishop:
		return bishopMagics[sq].Attacks[bishopMagics[sq].index(occupied)]
	case Rook:
		return rookMagics[sq].Attacks[rookMagics[sq].index(occupied)]
	case Queen:
		return bishopMagics[sq].Attacks[bishopMagics[sq].index(occupied)] | rookMagics[sq].Attacks[rookMagics[sq].index(occupied)]
	case Knight:
		return nonSliderAttacks[pt][sq]
	case King:
		return nonSliderAttacks[pt][sq]
	default:
		msg := fmt.Sprintf("GetAttackBb called with piece type %d is not supported", pt)
		panic(msg)
	}
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

// Various constant bitboards
const (
	BbZero = Bitboard(0)
	BbAll  = ^BbZero
	BbOne  = Bitboard(1)

	FileA_Bb Bitboard = 0x0101010101010101
	FileB_Bb          = FileA_Bb << 1
	FileC_Bb          = FileA_Bb << 2
	FileD_Bb          = FileA_Bb << 3
	FileE_Bb          = FileA_Bb << 4
	FileF_Bb          = FileA_Bb << 5
	FileG_Bb          = FileA_Bb << 6
	FileH_Bb          = FileA_Bb << 7

	Rank1_Bb Bitboard = 0xFF
	Rank2_Bb          = Rank1_Bb << (8 * 1)
	Rank3_Bb          = Rank1_Bb << (8 * 2)
	Rank4_Bb          = Rank1_Bb << (8 * 3)
	Rank5_Bb          = Rank1_Bb << (8 * 4)
	Rank6_Bb          = Rank1_Bb << (8 * 5)
	Rank7_Bb          = Rank1_Bb << (8 * 6)
	Rank8_Bb          = Rank1_Bb << (8 * 7)

	MsbMask   = ^(Bitboard(1) << 63)
	Rank8Mask = ^Rank8_Bb
	FileAMask = ^FileA_Bb
	FileHMask = ^FileH_Bb

	DiagUpA1 Bitboard = 0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001
	DiagUpB1          = (MsbMask & DiagUpA1) << 1 & FileAMask // shift EAST
	DiagUpC1          = (MsbMask & DiagUpB1) << 1 & FileAMask
	DiagUpD1          = (MsbMask & DiagUpC1) << 1 & FileAMask
	DiagUpE1          = (MsbMask & DiagUpD1) << 1 & FileAMask
	DiagUpF1          = (MsbMask & DiagUpE1) << 1 & FileAMask
	DiagUpG1          = (MsbMask & DiagUpF1) << 1 & FileAMask
	DiagUpH1          = (MsbMask & DiagUpG1) << 1 & FileAMask
	DiagUpA2          = (Rank8Mask & DiagUpA1) << 8 // shift NORTH
	DiagUpA3          = (Rank8Mask & DiagUpA2) << 8
	DiagUpA4          = (Rank8Mask & DiagUpA3) << 8
	DiagUpA5          = (Rank8Mask & DiagUpA4) << 8
	DiagUpA6          = (Rank8Mask & DiagUpA5) << 8
	DiagUpA7          = (Rank8Mask & DiagUpA6) << 8
	DiagUpA8          = (Rank8Mask & DiagUpA7) << 8

	DiagDownH1 Bitboard = 0b0000000100000010000001000000100000010000001000000100000010000000
	DiagDownH2          = (Rank8Mask & DiagDownH1) << 8 // shift NORTH
	DiagDownH3          = (Rank8Mask & DiagDownH2) << 8
	DiagDownH4          = (Rank8Mask & DiagDownH3) << 8
	DiagDownH5          = (Rank8Mask & DiagDownH4) << 8
	DiagDownH6          = (Rank8Mask & DiagDownH5) << 8
	DiagDownH7          = (Rank8Mask & DiagDownH6) << 8
	DiagDownH8          = (Rank8Mask & DiagDownH7) << 8
	DiagDownG1          = (DiagDownH1 >> 1) & FileHMask // shift WEST
	DiagDownF1          = (DiagDownG1 >> 1) & FileHMask
	DiagDownE1          = (DiagDownF1 >> 1) & FileHMask
	DiagDownD1          = (DiagDownE1 >> 1) & FileHMask
	DiagDownC1          = (DiagDownD1 >> 1) & FileHMask
	DiagDownB1          = (DiagDownC1 >> 1) & FileHMask
	DiagDownA1          = (DiagDownB1 >> 1) & FileHMask

	CenterFiles   = FileD_Bb | FileE_Bb
	CenterRanks   = Rank4_Bb | Rank5_Bb
	CenterSquares = CenterFiles & CenterRanks
)

// ////////////////////
// Private
// ////////////////////

// ////////////////////
// Pre compute helpers

// Returns a Bb of the square by shifting the
// square onto an empty bitboards.
// Usually one would use sq.Bb()
func (sq Square) bitboard() Bitboard {
	return Bitboard(uint64(1) << sq)
}

// helper arrays
var (
	// Internal pre computed square to square bitboard array.
	sqBb [SqLength]Bitboard

	// Internal pre computed rank bitboard array.
	rankBb [8]Bitboard

	// Internal pre computed file bitboard array.
	fileBb [8]Bitboard

	// Internal pre computed square to file bitboard array.
	sqToFileBb [SqLength]Bitboard

	// Internal pre computed square to rank bitboard array.
	sqToRankBb [SqLength]Bitboard

	// Internal pre computed square to diag up bitboard array.
	sqDiagUpBb [SqLength]Bitboard

	// Internal pre computed square to diag down bitboard array.
	sqDiagDownBb [SqLength]Bitboard

	// Internal pre computed index for quick square distance lookup
	squareDistance [SqLength][SqLength]int

	// Internal Bb for pawn attacks for each color for each square
	pawnAttacks [2][SqLength]Bitboard

	// Internal Bb for attacks for each piece for each square
	nonSliderAttacks [PtLength][SqLength]Bitboard

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
// OBS: The order is important as some initializations depend on others
func initBb() {
	rankFileBbPreCompute()
	squareBitboardsPreCompute()
	squareDistancePreCompute()
	nonSlidingAttacksPreCompute()
	initMagicBitboards()
	neighbourMasksPreCompute()
	raysPreCompute()
	intermediatePreCompute()
	maskPassedPawnsPreCompute()
	centerDistancePreCompute()
	castleMasksPreCompute()
	squareColorsPreCompute()
}

func rankFileBbPreCompute() {
	for i := Rank1; i <= Rank8; i++ {
		rankBb[i] = Rank1_Bb << (8 * i)
	}
	for i := FileA; i <= FileH; i++ {
		fileBb[i] = FileA_Bb << i
	}
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
	}
}

// Distance between squares index
func squareDistancePreCompute() {
	for sq1 := SqA1; sq1 <= SqH8; sq1++ {
		for sq2 := SqA1; sq2 <= SqH8; sq2++ {
			if sq1 != sq2 {
				squareDistance[sq1][sq2] = util.Max(FileDistance(sq1.FileOf(), sq2.FileOf()), RankDistance(sq1.RankOf(), sq2.RankOf()))
			}
		}
	}
}

// pre compute all possible attacked sq per color, piece and sq
func nonSlidingAttacksPreCompute() {
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
							nonSliderAttacks[pt][s] |= sqBb[to]
						}
					}
				}
			}
		}
	}
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

func raysPreCompute() {
	for sq := SqA1; sq <= SqH8; sq++ {
		rays[N][sq] = GetAttacksBb(Rook, sq, BbZero) & ranksNorthMask[sq]
		rays[E][sq] = GetAttacksBb(Rook, sq, BbZero) & filesEastMask[sq]
		rays[S][sq] = GetAttacksBb(Rook, sq, BbZero) & ranksSouthMask[sq]
		rays[W][sq] = GetAttacksBb(Rook, sq, BbZero) & filesWestMask[sq]
		rays[NW][sq] = GetAttacksBb(Bishop, sq, BbZero) & filesWestMask[sq] & ranksNorthMask[sq]
		rays[NE][sq] = GetAttacksBb(Bishop, sq, BbZero) & filesEastMask[sq] & ranksNorthMask[sq]
		rays[SE][sq] = GetAttacksBb(Bishop, sq, BbZero) & filesEastMask[sq] & ranksSouthMask[sq]
		rays[SW][sq] = GetAttacksBb(Bishop, sq, BbZero) & filesWestMask[sq] & ranksSouthMask[sq]
	}
}

// mask for intermediate squares in between two squares
func intermediatePreCompute() {
	for from := SqA1; from <= SqH8; from++ {
		for to := SqA1; to <= SqH8; to++ {
			toBB := sqBb[to]
			for o := 0; o < 8; o++ {
				if rays[Orientation(o)][from]&toBB != BbZero {
					intermediate[from][to] |= rays[Orientation(o)][from] &^ rays[Orientation(o)][to] &^ toBB
				}
			}
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

// masks for each square color (good for bishops vs bishops or pawns)
func squareColorsPreCompute() {
	for square := SqA1; square <= SqH8; square++ {
		if (int(square.FileOf())+int(square.RankOf()))%2 == 0 {
			squaresBb[Black] |= BbOne << square
		} else {
			squaresBb[White] |= BbOne << square
		}
	}
}

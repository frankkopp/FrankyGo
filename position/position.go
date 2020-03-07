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

package position

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frankkopp/FrankyGo/assert"
	. "github.com/frankkopp/FrankyGo/types"
)

// Key is used for zobrist keys in chess positions.
// Zobrist keys need all 64 bits for distribution
type Key uint64

// Position
// This struct represents the chess board and its position.
// It uses a 8x8 piece board and bitboards, a stack for undo moves, zobrist keys
// for transposition tables, piece lists, material and positional value counter.
//
// Needs to be created with New() or New(fen string)
type Position struct {

	// The zobrist key to use as a hash key in transposition tables
	// The zobrist key will be updated incrementally every time one of the the
	// state variables change.
	zobristKey Key

	// Board State
	// unique chess position (exception is 3-fold repetition
	// which is also not represented in a FEN string)
	board           [SqLength]Piece
	castlingRights  CastlingRights
	enPassantSquare Square
	halfMoveClock   int
	nextPlayer      Color

	// Extended Board State
	// not necessary for a unique position
	// special for king squares
	kingSquare [ColorLength]Square
	// half move number - the actual half move number to determine the full move number
	nextHalfMoveNumber int
	// piece bitboards
	piecesBb [ColorLength][PtLength]Bitboard
	// occupied bitboards with rotations
	occupiedBb    [ColorLength]Bitboard
	occupiedBbR90 [ColorLength]Bitboard
	occupiedBbL90 [ColorLength]Bitboard
	occupiedBbR45 [ColorLength]Bitboard
	occupiedBbL45 [ColorLength]Bitboard
	// history information for undo and repetition detection
	historyCounter int
	history        [maxHistory]historyState

	// Calculated by doMove/undoMove

	// Material value will always be up to date
	material        [ColorLength]int
	materialNonPawn [ColorLength]int
	// Positional value will always be up to date
	psqMidValue [ColorLength]int
	psqEndValue [ColorLength]int
	// Game phase value
	gamePhase int

	// caches a hasCheck and hasMate Flag for the current position. Will be set
	// after a call to hasCheck() and reset to TBD every time a move is made or
	// unmade.
	hasCheckFlag int
}

type historyState struct {
	zobristKeyHistory      Key
	moveHistory            Move
	fromPieceHistory       Piece
	capturedPieceHistory   Piece
	castlingRightsHistory  CastlingRights
	enpassantSquareHistory Square
	halfMoveClockHistory   int
	hasCheckFlagHistory    int
}

const maxHistory int = MaxMoves

// state flag for cached values
const (
	flagTBD   int = 0
	flagFalse int = 1
	flagTrue  int = 2
)

var initialized = false

// New creates a new position with Start Fen as default
func New() Position {
	return NewFen(StartFen)
}

// NewFen creates a new position with the given fen string
// as board position
func NewFen(fen string) Position {
	if !initialized {
		initZobrist()
		initialized = true
	}
	p := Position{}
	if e := p.setupBoard(fen); e != nil {
		panic(fmt.Sprintf("fen for position setup not valid and position can't be created: %s", e))
	}
	return p
}

// String returns a string representing the board instance. This
// includes the fen, a board matrix, game phase, material and pos values.
func (p Position) String() string {
	var os strings.Builder
	os.WriteString(p.StringFen())
	os.WriteString("\n")
	os.WriteString(p.StringBoard())
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Game Phase     : %d", p.gamePhase))
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Material White : %d", p.material[White]))
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Material Black : %d", p.material[Black]))
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Pos Value White: %d/%d", p.psqMidValue[White], p.psqEndValue[White]))
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Pos Value Black: %d/%d", p.psqMidValue[Black], p.psqEndValue[Black]))
	os.WriteString("\n")
	return os.String()
}

// StringFen returns a string with the FEN of the current position
func (p Position) StringFen() string {
	return p.fen()
}

// StringBoard returns a visual matrix of the board and pieces
func (p Position) StringBoard() string {
	var os strings.Builder
	os.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := Rank1; r <= Rank8; r++ {
		for f := FileA; f <= FileH; f++ {
			os.WriteString("| ")
			os.WriteString(p.board[SquareOf(f, Rank8-r)].Char())
			os.WriteString(" ")
		}
		os.WriteString("|\n+---+---+---+---+---+---+---+---+\n")
	}
	return os.String()
}

func (p *Position) putPiece(piece Piece, square Square) {
	color := piece.ColorOf()
	pieceType := piece.TypeOf()

	assert.Assert(p.board[square] == PieceNone, "tried to put piece on an occupied square: %s", square.String())
	assert.Assert(!p.piecesBb[color][pieceType].Has(square), "tried to set bit on pieceBb which is already set: %s", square.String())
	assert.Assert(!p.occupiedBb[color].Has(square), "tried to set bit on occupiedBb which is already set: %s", square.String())

	// update board
	p.board[square] = piece
	if pieceType == King {
		p.kingSquare[color] = square
	}
	// update bitboards
	p.piecesBb[color][pieceType].PushSquare(square)
	p.occupiedBb[color].PushSquare(square)
	p.occupiedBbR90[color].PushSquare(RotateSquareR90(square))
	p.occupiedBbL90[color].PushSquare(RotateSquareL90(square))
	p.occupiedBbR45[color].PushSquare(RotateSquareR45(square))
	p.occupiedBbL45[color].PushSquare(RotateSquareL45(square))
	// zobrist
	p.zobristKey ^= zobristBase.pieces[piece][square]
	// game phase
	p.gamePhase += pieceType.GamePhaseValue()
	// material
	p.material[color] += pieceType.ValueOf()
	if pieceType > Pawn {
		p.materialNonPawn[color] += pieceType.ValueOf()
	}
	// position value
	p.psqMidValue[color] += PosMidValue(piece, square)
	p.psqEndValue[color] += PosEndValue(piece, square)
}

func (p *Position) removePiece(square Square) Piece {
	removed := p.board[square]
	color := removed.ColorOf()
	pieceType := removed.TypeOf()

	assert.Assert(p.board[square] != PieceNone, "tried to remove piece from an empty square: %s", square.String())
	assert.Assert(p.piecesBb[color][pieceType].Has(square), "tried to clear bit from pieceBb which is not set: %s", square.String())
	assert.Assert(p.occupiedBb[color].Has(square), "tried to clear bit from occupiedBb which is not set: %s", square.String())

	// update board
	p.board[square] = PieceNone
	// update bitboards
	p.piecesBb[color][pieceType].PopSquare(square)
	p.occupiedBb[color].PopSquare(square)
	p.occupiedBbR90[color].PopSquare(RotateSquareR90(square))
	p.occupiedBbL90[color].PopSquare(RotateSquareL90(square))
	p.occupiedBbR45[color].PopSquare(RotateSquareR45(square))
	p.occupiedBbL45[color].PopSquare(RotateSquareL45(square))
	// zobrist
	p.zobristKey ^= zobristBase.pieces[removed][square]
	// game phase
	p.gamePhase -= pieceType.GamePhaseValue()
	// material
	p.material[color] -= pieceType.ValueOf()
	if pieceType > Pawn {
		p.materialNonPawn[color] -= pieceType.ValueOf()
	}
	// position value
	p.psqMidValue[color] -= PosMidValue(removed, square)
	p.psqEndValue[color] -= PosEndValue(removed, square)
	return removed
}

func (p Position) fen() string {
	var fen strings.Builder
	// pieces
	for r := Rank1; r <= Rank8; r++ {
		emptySquares := 0
		for f := FileA;	f <= FileH; f++ {
			pc := p.board[SquareOf(f, Rank8-r)];
			if pc == PieceNone {
				emptySquares++;
			} else {
				if emptySquares > 0 {
					fen.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}
				fen.WriteString(pc.String());
			}
		}
		if emptySquares > 0 {
			fen.WriteString(strconv.Itoa(emptySquares))
		}
		if r < Rank8 {
			fen.WriteString("/")
		}
	}
	// next player
	fen.WriteString(" ")
	fen.WriteString(p.nextPlayer.String())
	// castling
	fen.WriteString(" ")
	fen.WriteString(p.castlingRights.String())
	// en passant
	fen.WriteString(" ")
	fen.WriteString(p.enPassantSquare.String())
	// half move clock
	fen.WriteString(" ")
	fen.WriteString(strconv.Itoa(p.halfMoveClock))
	// full move number
	fen.WriteString(" ")
	fen.WriteString(strconv.Itoa((p.nextHalfMoveNumber + 1) / 2))

	return fen.String()
}

// setupBoard sets up a board based on a fen. This is basically
// the only way to get a valid Position instance. Internal state
// will be setup as well as all struct data is initialized to 0.
func (p *Position) setupBoard(fen string) error {

	// we will analyse the fen and only require the initial board layout part
	// All other parts will have defaults. E.g. next player is white, no castling, etc.
	fen = strings.TrimSpace(fen)
	fenParts := strings.Split(fen, " ")

	if len(fenParts) == 0 {
		err := errors.New("fen must not be empty")
		return err
	}

	// make sure only valid chars are used
	match, _ := regexp.MatchString("[0-8pPnNbBrRqQkK/]+", fenParts[0])
	if !match {
		err := errors.New("fen position contains invalid characters")
		return err
	}

	// fen string starts at a8 and runs to h8
	// with / jumping to file A of next lower rank
	currentSquare := SqA8

	// loop over fen and check an execute information
	for _, c := range fenParts[0] {
		if number, e := strconv.Atoi(string(c)); e == nil { // is number
			currentSquare = Square(int(currentSquare) + (number * int(East)))
		} else if string(c) == "/" { // find rank separator
			currentSquare = currentSquare.To(South).To(South)
		} else { // find piece type
			piece := PieceFromChar(string(c))
			if piece == PieceNone {
				err := errors.New(fmt.Sprintf("invalid piece character: %s", string(c)))
				return err
			}
			p.putPiece(piece, currentSquare)
			currentSquare++
		}
	}
	if currentSquare != SqA2 { // after h1++ we reach a2 - a2 needs to be last current square
		err := errors.New("not reached last square (h1) after reading fen")
		return err
	}

	// set defaults
	p.nextHalfMoveNumber = 1
	p.enPassantSquare = SqNone

	// everything below is optional as we can apply defaults

	// next player
	if len(fenParts) >= 2 {
		match, _ = regexp.MatchString("^[w|b]$", fenParts[1])
		if !match {
			err := errors.New("fen next player contains invalid characters")
			return err
		}
		switch fenParts[1] {
		case "w":
			p.nextPlayer = White
		case "b":
			{
				p.nextPlayer = Black
				p.zobristKey ^= zobristBase.nextPlayer
				p.nextHalfMoveNumber++
			}
		}
	}

	// castling rights
	if len(fenParts) >= 3 {
		match, _ = regexp.MatchString("^(K?Q?k?q?|-)$", fenParts[2])
		if !match {
			err := errors.New("fen castling rights contains invalid characters")
			return err
		}
		// are there  rights to be encoded?
		if fenParts[2] != "-" {
			for _, c := range fenParts[2] {
				switch string(c) {
				case "K":
					p.castlingRights.Add(CastlingWhiteOO)
				case "Q":
					p.castlingRights.Add(CastlingWhiteOOO)
				case "k":
					p.castlingRights.Add(CastlingBlackOO)
				case "q":
					p.castlingRights.Add(CastlingBlackOOO)
				}
			}
		}
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights]
	}

	// en passant
	if len(fenParts) >= 4 {
		match, _ = regexp.MatchString("^([a-h][1-8]|-)$", fenParts[3])
		if !match {
			err := errors.New("fen castling rights contains invalid characters")
			return err
		}
		if fenParts[3] != "-" {
			p.enPassantSquare = MakeSquare(fenParts[3])
		}
	}

	// half move clock (50 moves rule)
	if len(fenParts) >= 5 {
		if number, e := strconv.Atoi(fenParts[4]); e == nil { // is number
			p.halfMoveClock = number
		} else {
			return e
		}
	}

	// move number
	if len(fenParts) >= 6 {
		// game move number - to be converted into next half move number (ply)
		if moveNumber, e := strconv.Atoi(fenParts[5]); e == nil { // is number
			if moveNumber == 0 {
				moveNumber = 1
			}
			p.nextHalfMoveNumber = 2*moveNumber - (1 - int(p.nextPlayer))
		} else {
			return e
		}
	}

	// return without error
	return nil
}

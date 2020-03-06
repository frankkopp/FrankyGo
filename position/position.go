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

	// **********************************************************
	// Board State START ----------------------------------------
	// unique chess position (exception is 3-fold repetition
	// which is also not represented in a FEN string)

	// piece Board
	board [SqLength]Piece

	// Castling rights
	castlingRights CastlingRights

	// enpassant field
	enPassantSquare Square

	// half move clock - number of half moves since last capture
	halfMoveClock int

	// next player color
	nextPlayer Color

	// Board State END ------------------------------------------
	// **********************************************************

	// **********************************************************
	// Extended Board State -------------------------------------
	// not necessary for a unique position

	// special for king squares
	kingSquare [ColorLength]Square

	// half move number - the actual half move number to determine the full move
	// number
	nextHalfMoveNumber int

	// piece bitboards
	piecesBB [ColorLength][PtLength]Bitboard

	// occupied bitboards with rotations
	occupiedBB    [ColorLength]Bitboard
	occupiedBBR90 [ColorLength]Bitboard
	occupiedBBL90 [ColorLength]Bitboard
	occupiedBBR45 [ColorLength]Bitboard
	occupiedBBL45 [ColorLength]Bitboard

	// Extended Board State END ---------------------------------
	// **********************************************************

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
	p.setupBoard(fen)
	return p
}

func (p Position) String() string {
	return p.StringFen()
}

func (p Position) StringFen() string {
	return ""
}

func (p Position) StringBoard() string {
	return ""
}

func (p Position) setupBoard(fen string) {

}

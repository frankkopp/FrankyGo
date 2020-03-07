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
	zobristKey      Key
	move            Move
	fromPiece       Piece
	capturedPiece   Piece
	castlingRights  CastlingRights
	enpassantSquare Square
	halfMoveClock   int
	hasCheckFlag    int
}

const maxHistory int = MaxMoves

// state flag for cached values
const (
	flagTBD   int = 0
	flagFalse int = 1
	flagTrue  int = 2
)

var initialized = false

// //////////////////////////////////////////////////////
// // Public functions
// //////////////////////////////////////////////////////

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

// DoMove commits a move to the board. Due to performance there is no check if this
// move is legal on the current position. Legal check needs to be done
// beforehand. Usually the move will be generated by a MoveGenerator and
// therefore the move will be assumed legal anyway.
func (p *Position) DoMove(m Move) {
	assert.Assert(m.IsValid(), "Position DoMove: Invalid move %s", m.String())
	moveType := m.MoveType()
	fromSq := m.From()
	fromPc := p.board[fromSq]
	assert.Assert(fromPc != PieceNone, "Position DoMove: No piece on %s", fromPc.String())
	fromPt := fromPc.TypeOf()
	myColor := fromPc.ColorOf()
	assert.Assert(myColor == p.nextPlayer, "Position DoMove: Piece to move does not belong to next player %s", fromPc.String())
	toSq := m.To()
	targetPc := p.board[toSq]
	promPt := m.PromotionType()

	p.history[p.historyCounter] = historyState{
		p.zobristKey,
		m,
		fromPc,
		targetPc,
		p.castlingRights,
		p.enPassantSquare,
		p.halfMoveClock,
		p.hasCheckFlag}
	p.historyCounter++
	assert.Assert(p.historyCounter < MaxMoves, "Position DoMove: Can't have more moves than MaxMoves=%d", MaxMoves)

	switch moveType {
	case Normal:
		if p.castlingRights != CastlingNone && (CastlingMask.Has(fromSq) || CastlingMask.Has(toSq)) {
			p.invalidateCastlingRights(fromSq, toSq)
		}
		p.clearEnPassant()
		if targetPc != PieceNone { // capture
			p.removePiece(toSq)
			p.halfMoveClock = 0 // reset half move clock because of capture
		} else if fromPt == Pawn {
			p.halfMoveClock = 0                    // reset half move clock because of pawn move
			if SquareDistance(fromSq, toSq) == 2 { // pawn double - set en passant
				// set new en passant target field - always one "behind" the toSquare
				p.enPassantSquare = toSq.To(Direction(myColor.Flip().MoveDirection()) * North)
				p.zobristKey ^= zobristBase.enPassantFile[p.enPassantSquare.FileOf()] // in
			}
		} else {
			p.halfMoveClock++
		}
		p.movePiece(fromSq, toSq)
	case Promotion:
		assert.Assert(fromPc == MakePiece(myColor, King), "Position DoMove: Move type promotion but From piece not king")
		assert.Assert(toSq.RankOf() == myColor.PromotionRank(), "Position DoMove: Promotion move but wrong Rank")
		if targetPc != PieceNone { // capture
			p.removePiece(toSq)
		}
		if p.castlingRights != CastlingNone && (CastlingMask.Has(fromSq) || CastlingMask.Has(toSq)) {
			p.invalidateCastlingRights(fromSq, toSq)
		}
		p.removePiece(fromSq)
		p.putPiece(MakePiece(myColor, promPt), toSq)
		p.clearEnPassant()
		p.halfMoveClock = 0 // reset half move clock because of pawn move
	case EnPassant:
		assert.Assert(fromPc == MakePiece(myColor, Pawn), "Position DoMove: Move type en passant but from piece not pawn")
		assert.Assert(p.enPassantSquare != SqNone, "Position DoMove: EnPassant move type without en passant")
		capSq := toSq.To(Direction(myColor.Flip().MoveDirection()) * North)
		assert.Assert(p.board[capSq] == MakePiece(myColor.Flip(), Pawn), "Position DoMove: Captured en passant piece invalid")
		p.removePiece(capSq)
		p.movePiece(fromSq, toSq)
		p.clearEnPassant()
		p.halfMoveClock = 0 // reset half move clock because of pawn move
	case Castling:
		assert.Assert(fromPc == MakePiece(myColor, King), "Position DoMove: Move type castling but from piece not king")
		switch toSq {
		case SqG1:
			assert.Assert(p.castlingRights.Has(CastlingWhiteOO), "Position DoMove: White king side castling not available")
			assert.Assert(fromSq == SqE1, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE1] != WhiteKing, "Position DoMove: SqE1 has no king for castling")
			assert.Assert(p.board[SqH1] != WhiteRook, "Position DoMove: SqH1 has no rook for castling")
			assert.Assert(p.getOccupied()&Intermediate(SqE1, SqH1) == 0, "Position DoMove: Castling king side blocked")

			p.movePiece(fromSq, toSq)                                    // King
			p.movePiece(SqH1, SqF1)                                      // Rook
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingWhite)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in;
		case SqC1:
			assert.Assert(p.castlingRights.Has(CastlingWhiteOOO), "Position DoMove: White queen side castling not available")
			assert.Assert(fromSq == SqE1, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE1] != WhiteKing, "Position DoMove: SqE1 has no king for castling")
			assert.Assert(p.board[SqA1] != WhiteRook, "Position DoMove: SqA1 has no rook for castling")
			assert.Assert(p.getOccupied()&Intermediate(SqE1, SqA1) == 0, "Position DoMove: Castling queen side blocked")

			p.movePiece(fromSq, toSq)                                    // King
			p.movePiece(SqA1, SqD1)                                      // Rook
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingWhite)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		case SqG8:
			assert.Assert(p.castlingRights.Has(CastlingBlackOO), "Position DoMove: Black king side castling not available")
			assert.Assert(fromSq == SqE8, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE8] != BlackKing, "Position DoMove: SqE8 has no king for castling")
			assert.Assert(p.board[SqH8] != BlackRook, "Position DoMove: SqH8 has no rook for castling")
			assert.Assert(p.getOccupied()&Intermediate(SqE8, SqH8) == 0, "Position DoMove: Castling king side blocked")

			p.movePiece(fromSq, toSq)                                    // King
			p.movePiece(SqH8, SqF8)                                      // Rook
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingBlack)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		case SqC8:
			assert.Assert(p.castlingRights.Has(CastlingBlackOOO), "Position DoMove: Black queen side castling not available")
			assert.Assert(fromSq == SqE8, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE8] != BlackKing, "Position DoMove: SqE8 has no king for castling")
			assert.Assert(p.board[SqA8] != BlackRook, "Position DoMove: SqA8 has no rook for castling")
			assert.Assert(p.getOccupied()&Intermediate(SqE8, SqA8) == 0, "Position DoMove: Castling queen side blocked")

			p.movePiece(fromSq, toSq)                                    // King
			p.movePiece(SqA8, SqD8)                                      // Rook
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingBlack)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
			break
		default:
			panic("Invalid castle move!")
		}
		p.clearEnPassant()
		p.halfMoveClock++
	}

	p.hasCheckFlag = flagTBD
	p.nextHalfMoveNumber++
	p.nextPlayer = p.nextPlayer.Flip()
	p.zobristKey ^= zobristBase.nextPlayer
}

// UndoMove resets the position to a state before the last move has been made
func (p *Position) UndoMove() {
	assert.Assert(p.historyCounter > 0, "Position UndoMove: Cannot undo initial position")

	// Restore state part 1
	p.historyCounter--
	p.nextHalfMoveNumber--
	p.nextPlayer = p.nextPlayer.Flip()
	move := p.history[p.historyCounter].move

	// undo piece move / restore board
	switch move.MoveType() {

	case Normal:
		p.movePiece(move.To(), move.From())
		if p.history[p.historyCounter].capturedPiece != PieceNone {
			p.putPiece(p.history[p.historyCounter].capturedPiece, move.To())
		}
	case Promotion:
		p.removePiece(move.To())
		p.putPiece(MakePiece(p.nextPlayer, Pawn), move.From())
		if p.history[p.historyCounter].capturedPiece != PieceNone {
			p.putPiece(p.history[p.historyCounter].capturedPiece, move.To())
		}
	case EnPassant:
		// ignore Zobrist Key as it will be restored via history
		p.movePiece(move.To(), move.From())
		p.putPiece(MakePiece(p.nextPlayer.Flip(), Pawn), move.To().To(Direction(p.nextPlayer.Flip().MoveDirection())*North))
		break

	case Castling:
		// ignore Zobrist Key as it will be restored via history
		// castling rights are restored via history
		p.movePiece(move.To(), move.From()) // King
		switch move.To() {
		case SqG1:
			p.movePiece(SqF1, SqH1) // Rook
		case SqC1:
			p.movePiece(SqD1, SqA1) // Rook
		case SqG8:
			p.movePiece(SqF8, SqH8) // Rook
		case SqC8:
			p.movePiece(SqD8, SqA8) // Rook
		default:
			panic("Invalid castle move!")
		}
	}

	// restore state part 2
	p.castlingRights = p.history[p.historyCounter].castlingRights
	p.enPassantSquare = p.history[p.historyCounter].enpassantSquare
	p.halfMoveClock = p.history[p.historyCounter].halfMoveClock
	p.hasCheckFlag = p.history[p.historyCounter].hasCheckFlag
	p.zobristKey = p.history[p.historyCounter].zobristKey
}

// String returns a string representing the board instance. This
// includes the fen, a board matrix, game phase, material and pos values.
func (p *Position) String() string {
	var os strings.Builder
	os.WriteString(p.StringFen())
	os.WriteString("\n")
	os.WriteString(p.StringBoard())
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Next Player    : %s", p.nextPlayer.String()))
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
func (p *Position) StringFen() string {
	return p.fen()
}

// StringBoard returns a visual matrix of the board and pieces
func (p *Position) StringBoard() string {
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

// //////////////////////////////////////////////////////
// // Private functions
// //////////////////////////////////////////////////////

func (p *Position) movePiece(fromSq Square, toSq Square) {
	p.putPiece(p.removePiece(fromSq), toSq)
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

func (p Position) invalidateCastlingRights(from Square, to Square) {
	// check for castling rights invalidation
	if p.castlingRights&CastlingWhite != 0 {
		if from == SqE1 || to == SqE1 {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingWhite)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
		if p.castlingRights == CastlingWhiteOO && (from == SqH1 || to == SqH1) {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingWhiteOO)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
		if p.castlingRights == CastlingWhiteOOO && (from == SqA1 || to == SqA1) {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingWhiteOOO)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
	}
	if p.castlingRights&CastlingBlack != 0 {
		if from == SqE8 || to == SqE8 {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingBlack)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
		if p.castlingRights == CastlingBlackOOO && (from == SqA8 || to == SqA8) {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingBlackOOO)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
		if p.castlingRights == CastlingBlackOO && (from == SqH8 || to == SqH8) {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(CastlingBlackOO)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
	}
}

func (p *Position) clearEnPassant() {
	if p.enPassantSquare != SqNone {
		p.zobristKey = p.zobristKey ^ zobristBase.enPassantFile[p.enPassantSquare.FileOf()] // out
		p.enPassantSquare = SqNone
	}
}

func (p *Position) getOccupied() Bitboard {
	return p.occupiedBb[White] | p.occupiedBb[Black]
}

func (p *Position) fen() string {
	var fen strings.Builder
	// pieces
	for r := Rank1; r <= Rank8; r++ {
		emptySquares := 0
		for f := FileA; f <= FileH; f++ {
			pc := p.board[SquareOf(f, Rank8-r)]
			if pc == PieceNone {
				emptySquares++
			} else {
				if emptySquares > 0 {
					fen.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}
				fen.WriteString(pc.String())
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

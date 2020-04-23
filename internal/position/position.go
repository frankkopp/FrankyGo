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

// Package position represents data structures and functions for a chess board
// and its position.
// It uses a 8x8 piece board and bitboards, a stack for undo moves, zobrist keys
// for transposition tables, piece lists, material and positional value counter.
//
// Create a new instance with NewPosition(...) with no parameters to get the
// chess start position.
package position

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/internal/assert"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

var log *logging.Logger

var initialized = false

// initialize package
func init() {
	if !initialized {
		initZobrist()
		initialized = true
	}
}

const (
	// StartFen is a string with the fen position for a standard chess game
	StartFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

// Key is used for zobrist keys in chess positions.
// Zobrist keys need all 64 bits for distribution
type Key uint64

// Position
// This struct represents the chess board and its position.
// It uses a 8x8 piece board and bitboards, a stack for undo moves, zobrist keys
// for transposition tables, piece lists, material and positional value counter.
//
// Needs to be created with NewPosition() or NewPosition(fen string)
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
	occupiedBb [ColorLength]Bitboard
	// occupiedBbR90 [ColorLength]Bitboard
	// occupiedBbL90 [ColorLength]Bitboard
	// occupiedBbR45 [ColorLength]Bitboard
	// occupiedBbL45 [ColorLength]Bitboard
	// history information for undo and repetition detection
	historyCounter int
	history        [maxHistory]historyState

	// Calculated by doMove/undoMove

	// Material value will always be up to date
	material        [ColorLength]Value
	materialNonPawn [ColorLength]Value
	// Positional value will always be up to date
	psqMidValue [ColorLength]Value
	psqEndValue [ColorLength]Value
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

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// NewPosition creates a new position.
// When called without an argument the position will have the start position
// When a fen string is given it will create a position with based on this fen.
// Additional fens/strings are ignored
func NewPosition(fen ...string) *Position {
	if len(fen) == 0 {
		f, _ := NewPositionFen(StartFen)
		return f
	}
	f, _ := NewPositionFen(fen[0])
	return f
}

// NewPositionFen creates a new position with the given fen string
// as board position
// It returns nil and an error if the fen was invalid.
func NewPositionFen(fen string) (*Position, error) {
	if log == nil {
		log = myLogging.GetLog()
	}
	p := &Position{}
	if e := p.setupBoard(fen); e != nil {
		log.Errorf("fen for position setup not valid and position can't be created: %s", e)
		return nil, e
	}
	return p, nil
}

// DoMove commits a move to the board. Due to performance there is no check if this
// move is legal on the current position. Legal check needs to be done
// beforehand or after in case of pseudo legal moves. Usually the move will be
// generated by a MoveGenerator and therefore the move will be assumed legal anyway.
func (p *Position) DoMove(m Move) {
	fromSq := m.From()
	fromPc := p.board[fromSq]
	myColor := fromPc.ColorOf()
	toSq := m.To()
	targetPc := p.board[toSq]

	if false { // DEBUG
		switch {
		case !m.IsValid():
			msg := fmt.Sprintf("Position DoMove: Invalid move %s", m.String())
			log.Criticalf(msg)
			panic(msg)
		case fromPc == PieceNone:
			msg := fmt.Sprintf("Position DoMove: No piece on %s for move %s", fromPc.String(), m.StringUci())
			log.Criticalf(msg)
			panic(msg)
		case myColor != p.nextPlayer:
			msg := fmt.Sprintf("Position DoMove: Piece to move does not belong to next player %s", fromPc.String())
			log.Criticalf(msg)
			panic(msg)
		case targetPc.TypeOf() == King:
			msg := fmt.Sprintf("Position DoMove: King cannot be captured!")
			log.Criticalf(msg)
			panic(msg)
		}
	} // DEBUG

	if assert.DEBUG {
		assert.Assert(m.IsValid(), "Position DoMove: Invalid move %s", m.String())
		assert.Assert(fromPc != PieceNone, "Position DoMove: No piece on %s for move %s", fromPc.String(), m.StringUci())
		assert.Assert(myColor == p.nextPlayer, "Position DoMove: Piece to move does not belong to next player %s", fromPc.String())
		assert.Assert(targetPc.TypeOf() != King, "Position DoMove: King cannot be captured yet target piece is %s", targetPc.String())
	}

	// Save state of board for undo
	// this helps the compiler to prove that it is in bounds for the several updates we do after
	tmpHistoryCounter := p.historyCounter
	// update existing history entry to not create and allocate a new one
	p.history[tmpHistoryCounter].zobristKey = p.zobristKey
	p.history[tmpHistoryCounter].move = m
	p.history[tmpHistoryCounter].fromPiece = fromPc
	p.history[tmpHistoryCounter].capturedPiece = targetPc
	p.history[tmpHistoryCounter].castlingRights = p.castlingRights
	p.history[tmpHistoryCounter].enpassantSquare = p.enPassantSquare
	p.history[tmpHistoryCounter].halfMoveClock = p.halfMoveClock
	p.history[tmpHistoryCounter].hasCheckFlag = p.hasCheckFlag
	// update counter
	p.historyCounter++

	// do move according to MoveType
	switch m.MoveType() {
	case Normal:
		p.doNormalMove(fromSq, toSq, targetPc, fromPc, myColor)
	case Promotion:
		p.doPromotionMove(m, fromPc, myColor, toSq, targetPc, fromSq)
	case EnPassant:
		p.doEnPassantMove(toSq, myColor, fromPc, fromSq)
	case Castling:
		p.doCastlingMove(fromPc, myColor, toSq, fromSq)
	}

	// update additional state info
	p.hasCheckFlag = flagTBD
	p.nextHalfMoveNumber++
	p.nextPlayer = p.nextPlayer.Flip()
	p.zobristKey ^= zobristBase.nextPlayer
}

// UndoMove resets the position to a state before the last move has been made
func (p *Position) UndoMove() {
	if assert.DEBUG {
		assert.Assert(p.historyCounter > 0, "Position UndoMove: Cannot undo initial position")
	}

	// Restore state part 1
	p.historyCounter--
	p.nextHalfMoveNumber--
	p.nextPlayer = p.nextPlayer.Flip()
	// this helps the compiler to prove that it is in bounds
	// for the several updates we do after
	tmpHistoryCounter := p.historyCounter
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
		p.putPiece(MakePiece(p.nextPlayer.Flip(), Pawn), move.To().To(p.nextPlayer.Flip().MoveDirection()))
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
	p.castlingRights = p.history[tmpHistoryCounter].castlingRights
	p.enPassantSquare = p.history[tmpHistoryCounter].enpassantSquare
	p.halfMoveClock = p.history[tmpHistoryCounter].halfMoveClock
	p.hasCheckFlag = p.history[tmpHistoryCounter].hasCheckFlag
	p.zobristKey = p.history[tmpHistoryCounter].zobristKey
}

// DoNullMove is used in Null Move Pruning. The position is basically unchanged but
// the next player changes. The state before the null move will be stored to
// history.
// The history entry will be changed. So in effect after an UndoNullMove()
// the external view on the position is unchanged (e.g. fenBeforeNull == fenAfterNull
// and zobristBeforeNull == zobristAfterNull but positionBeforeNull != positionAfterNull.
func (p *Position) DoNullMove() {
	// Save state of board for undo
	// this helps the compiler to prove that it is in bounds for the several updates we do after
	tmpHistoryCounter := p.historyCounter
	// update existing history entry to not create and allocate a new one
	p.history[tmpHistoryCounter].zobristKey = p.zobristKey
	p.history[tmpHistoryCounter].move = MoveNone
	p.history[tmpHistoryCounter].fromPiece = PieceNone
	p.history[tmpHistoryCounter].capturedPiece = PieceNone
	p.history[tmpHistoryCounter].castlingRights = p.castlingRights
	p.history[tmpHistoryCounter].enpassantSquare = p.enPassantSquare
	p.history[tmpHistoryCounter].halfMoveClock = p.halfMoveClock
	p.history[tmpHistoryCounter].hasCheckFlag = p.hasCheckFlag
	// update counter
	p.historyCounter++
	// update state for null move
	p.hasCheckFlag = flagTBD
	p.clearEnPassant()
	p.nextHalfMoveNumber++
	p.nextPlayer = p.nextPlayer.Flip()
	p.zobristKey ^= zobristBase.nextPlayer
}

// UndoNullMove restores the state of the position to before the DoNullMove() call.
// The history entry will be changed but the history counter reset. So in effect
// the external view on the position is unchanged (e.g. fenBeforeNull == fenAfterNull
// and zobristBeforeNull == zobristAfterNull but positionBeforeNull != positionAfterNull
// If positionBeforeNull != positionAfterNull would be required this function would have
// to be changed to reset the history entry as well. Currently this is not necessary
// and therefore we spare the time to do this.
func (p *Position) UndoNullMove() {
	// Restore state
	p.historyCounter--
	p.nextHalfMoveNumber--
	p.nextPlayer = p.nextPlayer.Flip()
	// this helps the compiler to prove that it is in bounds
	// for the several updates we do after
	tmpHistoryCounter := p.historyCounter
	p.castlingRights = p.history[tmpHistoryCounter].castlingRights
	p.enPassantSquare = p.history[tmpHistoryCounter].enpassantSquare
	p.halfMoveClock = p.history[tmpHistoryCounter].halfMoveClock
	p.hasCheckFlag = p.history[tmpHistoryCounter].hasCheckFlag
	p.zobristKey = p.history[tmpHistoryCounter].zobristKey
}

// IsAttacked checks if the given square is attacked by a piece
// of the given color.
func (p *Position) IsAttacked(sq Square, by Color) bool {

	// to test if a position is attacked we do a reverse attack from the
	// target square to see if we hit a piece of the same or similar type

	// non sliding
	if (GetPawnAttacks(by.Flip(), sq)&p.piecesBb[by][Pawn] != 0) || // check pawns
		(GetPseudoAttacks(Knight, sq)&p.piecesBb[by][Knight] != 0) || // check knights
		(GetPseudoAttacks(King, sq)&p.piecesBb[by][King] != 0) { // check king
		return true
	}

	// New code - using GetAttacksBb from Magics - slower see tests
	// we do check a reverse attack with a queen to see if we can hit any other sliders. If yes
	// they also could hit us which means the square is attacked.
	// TODO: Look at this again
	if GetAttacksBb(Bishop, sq, p.OccupiedAll())&p.piecesBb[by][Bishop] > 0 ||
		GetAttacksBb(Rook, sq, p.OccupiedAll())&p.piecesBb[by][Rook] > 0 ||
		GetAttacksBb(Queen, sq, p.OccupiedAll())&p.piecesBb[by][Queen] > 0 {
		return true
	}

	// // sliding rooks and queens
	// if (GetPseudoAttacks(Rook, sq)&p.piecesBb[by][Rook] != 0 || (GetPseudoAttacks(Rook, sq)&p.piecesBb[by][Queen] != 0)) &&
	// 	(((GetMovesOnRank(sq, p.OccupiedAll()) |
	// 		GetMovesOnFileRotated(sq, p.occupiedBbL90[White]|p.occupiedBbL90[Black])) &
	// 		(p.piecesBb[by][Rook] | p.piecesBb[by][Queen])) != 0) {
	// 	return true
	// }
	//
	// // sliding bishop and queens
	// if (GetPseudoAttacks(Bishop, sq)&p.piecesBb[by][Bishop] != 0 || (GetPseudoAttacks(Bishop, sq)&p.piecesBb[by][Queen] != 0)) &&
	// 	(((GetMovesDiagUpRotated(sq, p.occupiedBbR45[White]|p.occupiedBbR45[Black]) |
	// 		GetMovesDiagDownRotated(sq, p.occupiedBbL45[White]|p.occupiedBbL45[Black])) &
	// 		(p.piecesBb[by][Bishop] | p.piecesBb[by][Queen])) != 0) {
	// 	return true
	// }

	// check en passant
	if p.enPassantSquare != SqNone {
		switch by {
		case White: // white is attacker
			// black is target
			if p.board[p.enPassantSquare.To(South)] == BlackPawn &&
				// this is indeed the en passant attacked square
				p.enPassantSquare.To(South) == sq {
				// left
				square := sq.To(West)
				if p.board[square] == WhitePawn {
					return true
				}
				// right
				square = sq.To(East)
				return p.board[square] == WhitePawn
			}
		case Black: // black is attacker
			// white is target
			if p.board[p.enPassantSquare.To(North)] == WhitePawn &&
				// this is indeed the en passant attacked square
				p.enPassantSquare.To(North) == sq {
				// attack from left
				square := sq.To(West)
				if p.board[square] == BlackPawn {
					return true
				}
				// right
				square = sq.To(East)
				return p.board[square] == BlackPawn
			}
		}
	}
	return false
}

// IsLegalMove tests a move if it is legal on the current position.
// Basically tests if the king would be left in check after the move
// or if the king crosses an attacked square during castling.
func (p *Position) IsLegalMove(move Move) bool {
	// king is not allowed to pass a square which is attacked by opponent
	if move.MoveType() == Castling {
		// castling not allowed when in check
		// we can simply check the from square of the castling move
		// and check if the current opponent attacks it. Castling would not
		// be possible if the attack would be influenced by the castling
		// itself.
		if p.IsAttacked(move.From(), p.nextPlayer.Flip()) {
			return false
		}
		// castling crossing attacked square?
		switch move.To() {
		case SqG1:
			if p.IsAttacked(SqF1, p.nextPlayer.Flip()) {
				return false
			}
		case SqC1:
			if p.IsAttacked(SqD1, p.nextPlayer.Flip()) {
				return false
			}
		case SqG8:
			if p.IsAttacked(SqF8, p.nextPlayer.Flip()) {
				return false
			}
		case SqC8:
			if p.IsAttacked(SqD8, p.nextPlayer.Flip()) {
				return false
			}
		default:
			break
		}
	}
	// make the move on the position
	// then check if the move leaves the king in check
	p.DoMove(move)
	legal := !p.IsAttacked(p.kingSquare[p.nextPlayer.Flip()], p.nextPlayer)
	p.UndoMove()
	return legal
}

// WasLegalMove tests if the last move was legal. Basically tests if
// the king is now in check or if the king crossed an attacked square
// during castling or of there was a castling although in check.
// If the position does not have a last move (history empty) this
// will only check if the king of the opponent is attacked e.g. could
// now be captured by the next player.
func (p *Position) WasLegalMove() bool {
	// king attacked?
	if p.IsAttacked(p.kingSquare[p.nextPlayer.Flip()], p.nextPlayer) {
		return false
	}
	// look back and check if castling was legal
	if p.historyCounter > 0 {
		move := p.history[p.historyCounter-1].move
		if move.MoveType() == Castling {
			// castling not allowed when in check
			// we can simply check the from square of the last castling move
			// and check if the current player attacks it. Castling would not
			// be possible if the attack would be influenced by the castling
			// itself.
			if p.IsAttacked(move.From(), p.nextPlayer) {
				return false
			}
			// castling crossing attacked square?
			switch move.To() {
			case SqG1:
				if p.IsAttacked(SqF1, p.nextPlayer) {
					return false
				}
			case SqC1:
				if p.IsAttacked(SqD1, p.nextPlayer) {
					return false
				}
			case SqG8:
				if p.IsAttacked(SqF8, p.nextPlayer) {
					return false
				}
			case SqC8:
				if p.IsAttacked(SqD8, p.nextPlayer) {
					return false
				}
			default:
				break
			}
		}
	}
	return true
}

// HasCheck returns true if the next player is threatened by a check
// (king is attacked).
// This is cached for the current position. Multiple calls to this
// on the same position are therefore very efficient.
func (p *Position) HasCheck() bool {
	if p.hasCheckFlag != flagTBD {
		return p.hasCheckFlag == flagTrue
	}
	check := p.IsAttacked(p.kingSquare[p.nextPlayer], p.nextPlayer.Flip())
	if check {
		p.hasCheckFlag = flagTrue
	} else {
		p.hasCheckFlag = flagFalse
	}
	return check
}

// IsCapturingMove determines if a move on this position is a capturing move
// incl. en passant
func (p *Position) IsCapturingMove(move Move) bool {
	return p.occupiedBb[p.nextPlayer.Flip()].Has(move.To()) || move.MoveType() == EnPassant
}

// CheckRepetitions
// Repetition of a position:.
// To detect a 3-fold repetition the given position must occur at least 2
// times before:<br/> <code>position.checkRepetitions(2)</code> checks for 3
// fold-repetition <p> 3-fold repetition: This most commonly occurs when
// neither side is able to avoid repeating moves without incurring a
// disadvantage. The three occurrences of the position need not occur on
// consecutive moves for a claim to be valid. FIDE rules make no mention of
// perpetual check; this is merely a specific type of draw by threefold
// repetition.
//
// Return true if this position has been played reps times before
func (p *Position) CheckRepetitions(reps int) bool {
	/*
	   [0]     3185849660387886977 << 1st
	   [1]     447745478729458041
	   [2]     3230145143131659788
	   [3]     491763876012767476
	   [4]     3185849660387886977 << 2nd
	   [5]     447745478729458041
	   [6]     3230145143131659788
	   [7]     491763876012767476  <<< history
	   [8]     3185849660387886977 <<< 3rd REPETITION from current zobrist
	*/
	counter := 0
	i := p.historyCounter - 2
	lastHalfMove := p.halfMoveClock
	for i >= 0 {
		// every time the half move clock gets reset (non reversible position) there
		// can't be any more repetition of positions before this position
		if p.history[i].halfMoveClock >= lastHalfMove {
			break
		} else {
			lastHalfMove = p.history[i].halfMoveClock
		}
		if p.zobristKey == p.history[i].zobristKey {
			counter++
		}
		if counter >= reps {
			return true
		}
		i -= 2
	}
	return false
}

// HasInsufficientMaterial returns true if no side has enough material to
// force a mate (does not exclude combination where a helpmate would be
// possible, e.g. the opponent needs to support a mate by mistake)
func (p *Position) HasInsufficientMaterial() bool {

	// we use material value as minor pieces knights and bishops
	// have different values and it is assumed that this is faster
	// then a pop count on a bitboard - not empirically tested

	// no material
	// both sides have a bare king
	if p.material[White]+p.material[Black] == 0 {
		return true
	}

	// no more pawns
	if p.piecesBb[White][Pawn].PopCount() == 0 && p.piecesBb[Black][Pawn].PopCount() == 0 {
		// one side has a king and a minor piece against a bare king
		// both sides have a king and a minor piece each
		if p.materialNonPawn[White] < 400 && p.materialNonPawn[Black] < 400 {
			return true
		}
		// the weaker side has a minor piece against two knights
		if (p.materialNonPawn[White] == 2*Knight.ValueOf() && p.materialNonPawn[Black] <= Bishop.ValueOf()) ||
			(p.materialNonPawn[Black] == 2*Knight.ValueOf() && p.materialNonPawn[White] <= Bishop.ValueOf()) {
			return true
		}
		// two bishops draw against a bishop
		if (p.materialNonPawn[White] == 2*Bishop.ValueOf() && p.materialNonPawn[Black] == Bishop.ValueOf()) ||
			(p.materialNonPawn[Black] == 2*Bishop.ValueOf() && p.materialNonPawn[White] == Bishop.ValueOf()) {
			return true
		}
		// one side has two bishops a mate can be forced
		if p.materialNonPawn[White] == 2*Bishop.ValueOf() || p.materialNonPawn[Black] == 2*Bishop.ValueOf() {
			return false
		}
		// two minor pieces against one draw, except when the stronger side has a bishop pair
		if (p.materialNonPawn[White] < 2*Bishop.ValueOf() && p.materialNonPawn[Black] <= Bishop.ValueOf()) ||
			(p.materialNonPawn[White] <= Bishop.ValueOf() && p.materialNonPawn[Black] < 2*Bishop.ValueOf()) {
			return true
		}
	}
	return false
}

// GivesCheck determines if the given move will give check to the opponent
// of p.NextPlayer() and returns true if so.
func (p *Position) GivesCheck(move Move) bool {

	us := p.nextPlayer
	them := us.Flip()

	// opponents king square
	kingSq := p.kingSquare[them]

	// move details
	fromSq := move.From()
	toSq := move.To()
	fromPc := p.board[fromSq]
	fromPt := fromPc.TypeOf()
	epTargetSq := SqNone
	moveType := move.MoveType()

	switch moveType {
	case Promotion:
		// promotion moves - use new piece type
		fromPt = move.PromotionType()
	case Castling:
		// set the target square to the rook square and
		// piece type to ROOK. King can't give check
		// also no revealed check possible in castling
		fromPt = Rook
		switch toSq {
		case SqG1: // white king side castle
			toSq = SqF1
		case SqC1: // white queen side castle
			toSq = SqD1
		case SqG8: // black king side castle
			toSq = SqF8
		case SqC8: // black queen side castle
			toSq = SqD8
		}
	case EnPassant:
		// set en passant capture square
		epTargetSq = toSq.To(them.MoveDirection())
	}

	// get all pieces to check occupied intermediate squares
	boardAfterMove := p.OccupiedAll()

	// adapt board by moving the piece on the bitboard
	boardAfterMove.PopSquare(fromSq)
	boardAfterMove.PushSquare(toSq)
	if moveType == EnPassant {
		boardAfterMove.PopSquare(epTargetSq)
	}

	// Find direct checks
	switch fromPt {
	case Pawn:
		if GetPawnAttacks(us, toSq).Has(kingSq) {
			return true
		}
	case King:
	// ignore - can't give check
	default:
		if GetAttacksBb(fromPt, toSq, boardAfterMove).Has(kingSq) {
			return true
		}
	}

	// revealed checks
	// we only need to check for rook, bishop and queens
	// knight and pawn attacks can't be revealed
	// exception is en passant where the captured piece can reveal check
	switch {
	case GetAttacksBb(Bishop, kingSq, boardAfterMove)&p.piecesBb[us][Bishop] > 0:
		return true
	case GetAttacksBb(Rook, kingSq, boardAfterMove)&p.piecesBb[us][Rook] > 0:
		return true
	case GetAttacksBb(Queen, kingSq, boardAfterMove)&p.piecesBb[us][Queen] > 0:
		return true
	}

	// we did not find a check
	return false
}

// String returns a string representing the board instance. This
// includes the fen, a board matrix, game phase, material and pos values.
func (p *Position) String() string {
	var os strings.Builder
	os.WriteString(p.StringFen())
	os.WriteString("\n")
	os.WriteString(p.StringBoard())
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Next Player    : %s\n", p.nextPlayer.String()))
	os.WriteString(fmt.Sprintf("Game Phase     : %d\n", p.gamePhase))
	os.WriteString(fmt.Sprintf("Material White : %d\n", p.material[White]))
	os.WriteString(fmt.Sprintf("Material Black : %d\n", p.material[Black]))
	os.WriteString(fmt.Sprintf("Pos value White: %d/%d\n", p.psqMidValue[White], p.psqEndValue[White]))
	os.WriteString(fmt.Sprintf("Pos value Black: %d/%d\n", p.psqMidValue[Black], p.psqEndValue[Black]))
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

// //////////////////////////////////////////////////////////
// Private
// //////////////////////////////////////////////////////////

func (p *Position) doNormalMove(fromSq Square, toSq Square, targetPc Piece, fromPc Piece, myColor Color) {
	// If we still have castling rights and the move touches castling squares then invalidate
	// the corresponding castling right
	if p.castlingRights != CastlingNone {
		cr := GetCastlingRights(fromSq) | GetCastlingRights(toSq)
		if cr != CastlingNone {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(cr)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
	}
	p.clearEnPassant()
	if targetPc != PieceNone { // capture
		p.removePiece(toSq)
		p.halfMoveClock = 0 // reset half move clock because of capture
	} else if fromPc.TypeOf() == Pawn {
		p.halfMoveClock = 0                    // reset half move clock because of pawn move
		if SquareDistance(fromSq, toSq) == 2 { // pawn double - set en passant
			// set new en passant target field - always one "behind" the toSquare
			p.enPassantSquare = toSq.To(myColor.Flip().MoveDirection())
			p.zobristKey ^= zobristBase.enPassantFile[p.enPassantSquare.FileOf()] // in
		}
	} else {
		p.halfMoveClock++
	}
	p.movePiece(fromSq, toSq)
}

func (p *Position) doCastlingMove(fromPc Piece, myColor Color, toSq Square, fromSq Square) {
	if assert.DEBUG {
		assert.Assert(fromPc == MakePiece(myColor, King), "Position DoMove: Move type castling but from piece not king")
	}
	switch toSq {
	case SqG1:
		if assert.DEBUG {
			assert.Assert(p.castlingRights.Has(CastlingWhiteOO), "Position DoMove: White king side castling not available")
			assert.Assert(fromSq == SqE1, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE1] == WhiteKing, "Position DoMove: SqE1 has no king for castling")
			assert.Assert(p.board[SqH1] == WhiteRook, "Position DoMove: SqH1 has no rook for castling")
			assert.Assert(p.OccupiedAll()&Intermediate(SqE1, SqH1) == 0, "Position DoMove: Castling king side blocked")
		}
		p.movePiece(fromSq, toSq)                                    // King
		p.movePiece(SqH1, SqF1)                                      // Rook
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
		p.castlingRights.Remove(CastlingWhite)
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in;
	case SqC1:
		if assert.DEBUG {
			assert.Assert(p.castlingRights.Has(CastlingWhiteOOO), "Position DoMove: White queen side castling not available")
			assert.Assert(fromSq == SqE1, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE1] == WhiteKing, "Position DoMove: SqE1 has no king for castling")
			assert.Assert(p.board[SqA1] == WhiteRook, "Position DoMove: SqA1 has no rook for castling")
			assert.Assert(p.OccupiedAll()&Intermediate(SqE1, SqA1) == 0, "Position DoMove: Castling queen side blocked")
		}
		p.movePiece(fromSq, toSq)                                    // King
		p.movePiece(SqA1, SqD1)                                      // Rook
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
		p.castlingRights.Remove(CastlingWhite)
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
	case SqG8:
		if assert.DEBUG {
			assert.Assert(p.castlingRights.Has(CastlingBlackOO), "Position DoMove: Black king side castling not available")
			assert.Assert(fromSq == SqE8, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE8] == BlackKing, "Position DoMove: SqE8 has no king for castling")
			assert.Assert(p.board[SqH8] == BlackRook, "Position DoMove: SqH8 has no rook for castling")
			assert.Assert(p.OccupiedAll()&Intermediate(SqE8, SqH8) == 0, "Position DoMove: Castling king side blocked")
		}
		p.movePiece(fromSq, toSq)                                    // King
		p.movePiece(SqH8, SqF8)                                      // Rook
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
		p.castlingRights.Remove(CastlingBlack)
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
	case SqC8:
		if assert.DEBUG {
			assert.Assert(p.castlingRights.Has(CastlingBlackOOO), "Position DoMove: Black queen side castling not available")
			assert.Assert(fromSq == SqE8, "Position DoMove: Castling from square not correct")
			assert.Assert(p.board[SqE8] == BlackKing, "Position DoMove: SqE8 has no king for castling")
			assert.Assert(p.board[SqA8] == BlackRook, "Position DoMove: SqA8 has no rook for castling")
			assert.Assert(p.OccupiedAll()&Intermediate(SqE8, SqA8) == 0, "Position DoMove: Castling queen side blocked")
		}
		p.movePiece(fromSq, toSq)                                    // King
		p.movePiece(SqA8, SqD8)                                      // Rook
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
		p.castlingRights.Remove(CastlingBlack)
		p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
	default:
		panic("Invalid castle move!")
	}
	p.clearEnPassant()
	p.halfMoveClock++
}

func (p *Position) doEnPassantMove(toSq Square, myColor Color, fromPc Piece, fromSq Square) {
	capSq := toSq.To(myColor.Flip().MoveDirection())
	if assert.DEBUG {
		assert.Assert(fromPc == MakePiece(myColor, Pawn), "Position DoMove: Move type en passant but from piece not pawn")
		assert.Assert(p.enPassantSquare != SqNone, "Position DoMove: EnPassant move type without en passant")
		assert.Assert(p.board[capSq] == MakePiece(myColor.Flip(), Pawn), "Position DoMove: Captured en passant piece invalid")
	}
	p.removePiece(capSq)
	p.movePiece(fromSq, toSq)
	p.clearEnPassant()
	// reset half move clock because of pawn move
	p.halfMoveClock = 0
}

func (p *Position) doPromotionMove(m Move, fromPc Piece, myColor Color, toSq Square, targetPc Piece, fromSq Square) {
	if assert.DEBUG {
		assert.Assert(fromPc == MakePiece(myColor, Pawn), "Position DoMove: Move type promotion but From piece not Pawn")
		assert.Assert(myColor.PromotionRankBb().Has(toSq), "Position DoMove: Promotion move but wrong Rank")
	}
	if targetPc != PieceNone { // capture
		p.removePiece(toSq)
	}
	if p.castlingRights != CastlingNone {
		cr := GetCastlingRights(fromSq) | GetCastlingRights(toSq)
		if cr != CastlingNone {
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // out
			p.castlingRights.Remove(cr)
			p.zobristKey ^= zobristBase.castlingRights[p.castlingRights] // in
		}
	}
	p.removePiece(fromSq)
	p.putPiece(MakePiece(myColor, m.PromotionType()), toSq)
	p.clearEnPassant()
	p.halfMoveClock = 0 // reset half move clock because of pawn move
}

func (p *Position) movePiece(fromSq Square, toSq Square) {
	p.putPiece(p.removePiece(fromSq), toSq)
}

func (p *Position) putPiece(piece Piece, square Square) {
	color := piece.ColorOf()
	pieceType := piece.TypeOf()

	if assert.DEBUG {
		assert.Assert(p.board[square] == PieceNone, "tried to put piece on an occupied square: %s", square.String())
		assert.Assert(!p.piecesBb[color][pieceType].Has(square), "tried to set bit on pieceBb which is already set: %s", square.String())
		assert.Assert(!p.occupiedBb[color].Has(square), "tried to set bit on occupiedBb which is already set: %s", square.String())
	}

	// update board
	p.board[square] = piece
	if pieceType == King {
		p.kingSquare[color] = square
	}
	// update bitboards
	p.piecesBb[color][pieceType].PushSquare(square)
	p.occupiedBb[color].PushSquare(square)
	// p.occupiedBbR90[color].PushSquare(RotateSquareR90(square))
	// p.occupiedBbL90[color].PushSquare(RotateSquareL90(square))
	// p.occupiedBbR45[color].PushSquare(RotateSquareR45(square))
	// p.occupiedBbL45[color].PushSquare(RotateSquareL45(square))
	// zobrist
	p.zobristKey ^= zobristBase.pieces[piece][square]
	// game phase
	p.gamePhase += pieceType.GamePhaseValue()
	if p.gamePhase > GamePhaseMax {
		p.gamePhase = GamePhaseMax
	}
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

	if assert.DEBUG {
		assert.Assert(p.board[square] != PieceNone, "tried to remove piece from an empty square: %s", square.String())
		assert.Assert(p.piecesBb[color][pieceType].Has(square), "tried to clear bit from pieceBb which is not set: %s", square.String())
		assert.Assert(p.occupiedBb[color].Has(square), "tried to clear bit from occupiedBb which is not set: %s", square.String())
	}

	// update board
	p.board[square] = PieceNone
	// update bitboards
	p.piecesBb[color][pieceType].PopSquare(square)
	p.occupiedBb[color].PopSquare(square)
	// p.occupiedBbR90[color].PopSquare(RotateSquareR90(square))
	// p.occupiedBbL90[color].PopSquare(RotateSquareL90(square))
	// p.occupiedBbR45[color].PopSquare(RotateSquareR45(square))
	// p.occupiedBbL45[color].PopSquare(RotateSquareL45(square))
	// zobrist
	p.zobristKey ^= zobristBase.pieces[removed][square]
	// game phase
	p.gamePhase -= pieceType.GamePhaseValue()
	if p.gamePhase < 0 {
		p.gamePhase = 0
	}
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

func (p *Position) clearEnPassant() {
	if p.enPassantSquare != SqNone {
		p.zobristKey ^= zobristBase.enPassantFile[p.enPassantSquare.FileOf()] // out
		p.enPassantSquare = SqNone
	}
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

// regex for first part of fen (position of pieces)
var regexFenPos = regexp.MustCompile("[0-8pPnNbBrRqQkK/]+")

// regex for next player color in fen
var regexWorB = regexp.MustCompile("^[w|b]$")

// regex for castling rights in fen
var regexCastlingRights = regexp.MustCompile("^(K?Q?k?q?|-)$")

// regex for en passent square in fen
var regexEnPassant = regexp.MustCompile("^([a-h][1-8]|-)$")

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
	match := regexFenPos.MatchString(fenParts[0])
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
		match = regexWorB.MatchString(fenParts[1])
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
		match = regexCastlingRights.MatchString(fenParts[2])
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
		match = regexEnPassant.MatchString(fenParts[3])
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

// //////////////////////////////////////////////////////
// // Getter and Setter functions
// //////////////////////////////////////////////////////

// ZobristKey returns the current zobrist key for this position
func (p *Position) ZobristKey() Key {
	return p.zobristKey
}

// NextPlayer returns the next player as Color for the position
func (p *Position) NextPlayer() Color {
	return p.nextPlayer
}

// GetPiece returns the piece on the given square. Empty
// squares are initialized with PieceNone and return the same.
func (p *Position) GetPiece(sq Square) Piece {
	return p.board[sq]
}

// PiecesBb returns the Bitboard for the given piece type of the given color
func (p *Position) PiecesBb(c Color, pt PieceType) Bitboard {
	return p.piecesBb[c][pt]
}

// OccupiedAll returns a Bitboard of all pieces currently on the board
func (p *Position) OccupiedAll() Bitboard {
	return p.occupiedBb[White] | p.occupiedBb[Black]
}

// OccupiedBb returns a Bitboard f all pieces of Color c
func (p *Position) OccupiedBb(c Color) Bitboard {
	return p.occupiedBb[c]
}

// GamePhase returns the current game phase value of the position.
// GamePhase is 24 at the start of the game (24 is also the max).
// End games when no officers are left have a GamePhase value of 0.
func (p *Position) GamePhase() int {
	return p.gamePhase
}

// GamePhaseFactor returns a factor between 0 and 1 which reflects
// the ratio between the actual game phase and the max game phase
func (p *Position) GamePhaseFactor() float64 {
	return float64(p.gamePhase) / GamePhaseMax
}

// GetEnPassantSquare returns the en passant square or SqNone if not set
func (p *Position) GetEnPassantSquare() Square {
	return p.enPassantSquare
}

// CastlingRights returns the castling rights instance of the position
func (p *Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

// KingSquare returns the current square of the king of color c
func (p *Position) KingSquare(c Color) Square {
	return p.kingSquare[c]
}

// HalfMoveClock returns the positions half move clock
func (p *Position) HalfMoveClock() int {
	return p.halfMoveClock
}

// Material returns the material value for the given color
// on this position
func (p *Position) Material(c Color) Value {
	return p.material[c]
}

// MaterialNonPawn returns the non pawn material value for
// given color
func (p *Position) MaterialNonPawn(c Color) Value {
	return p.materialNonPawn[c]
}

// PsqMidValue returns the positional value for the given color
// for early game phases. Best used together with a game phase
// factor
func (p *Position) PsqMidValue(c Color) Value {
	return p.psqMidValue[c]
}

// PsqEndValue returns the positional value for the given color
// for later game phases. Best used together with a game phase
// factor
func (p *Position) PsqEndValue(c Color) Value {
	return p.psqEndValue[c]
}

// LastMove returns the last move made on the position or
// MoveNone if the position has no history of earlier moves.
func (p *Position) LastMove() Move {
	if p.historyCounter <= 0 {
		return MoveNone
	}
	return p.history[p.historyCounter-1].move
}

// LastCapturedPiece returns the captured piece of the the last
// move made on the position or MoveNone if the move was
// non-capturing or the position has no history of earlier moves.
func (p *Position) LastCapturedPiece() Piece {
	if p.historyCounter <= 0 {
		return PieceNone
	}
	return p.history[p.historyCounter-1].capturedPiece
}

// WasCapturingMove returns true if the last move was
// a capturing move.
func (p *Position) WasCapturingMove() bool {
	return p.LastCapturedPiece() != PieceNone
}

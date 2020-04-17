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

// Package movegen contains functionality to create moves on a
// chess position. It implements several variants like
// generate pseudo legal moves, legal moves or on demand
// generation of pseudo legal moves.
package movegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/op/go-logging"

	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

var log *logging.Logger

// Movegen data structure. Create new move generator via
//  movegen.NewMoveGen()
// Creating this directly will not work.
type Movegen struct {
	pseudoLegalMoves   *moveslice.MoveSlice
	legalMoves         *moveslice.MoveSlice
	onDemandMoves      *moveslice.MoveSlice
	killerMoves        [2]Move
	currentIteratorKey position.Key
	takeIndex          int
	pvMove             Move
	currentODStage     int8
	pvMovePushed       bool
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// GenMode generation modes for on demand move generation
type GenMode int

// GenMode generation modes for on demand move generation
const (
	GenZero   GenMode = 0b00
	GenCap    GenMode = 0b01
	GenNonCap GenMode = 0b10
	GenAll    GenMode = 0b11
)

// NewMoveGen creates a new instance of a move generator
func NewMoveGen() *Movegen {
	if log == nil {
		log = myLogging.GetLog()
	}
	tmpMg := &Movegen{
		pseudoLegalMoves:   moveslice.NewMoveSlice(MaxMoves),
		legalMoves:         moveslice.NewMoveSlice(MaxMoves),
		onDemandMoves:      moveslice.NewMoveSlice(MaxMoves),
		killerMoves:        [2]Move{MoveNone, MoveNone},
		pvMove:             MoveNone,
		currentODStage:     odNew,
		currentIteratorKey: 0,
		pvMovePushed:       false,
		takeIndex:          0,
	}
	return tmpMg
}

// GeneratePseudoLegalMoves generates pseudo moves for the next player. Does not check if
// king is left in check or if it passes an attacked square when castling or has been in check
// before castling.
func (mg *Movegen) GeneratePseudoLegalMoves(position *position.Position, mode GenMode) *moveslice.MoveSlice {
	mg.pseudoLegalMoves.Clear()
	if mode&GenCap != 0 {
		mg.generatePawnMoves(position, GenCap, mg.pseudoLegalMoves)
		mg.generateCastling(position, GenCap, mg.pseudoLegalMoves)
		mg.generateKingMoves(position, GenCap, mg.pseudoLegalMoves)
		mg.generateMoves(position, GenCap, mg.pseudoLegalMoves)
	}
	if mode&GenNonCap != 0 {
		mg.generatePawnMoves(position, GenNonCap, mg.pseudoLegalMoves)
		mg.generateCastling(position, GenNonCap, mg.pseudoLegalMoves)
		mg.generateKingMoves(position, GenNonCap, mg.pseudoLegalMoves)
		mg.generateMoves(position, GenNonCap, mg.pseudoLegalMoves)
	}
	// PV and Killer handling
	mg.pseudoLegalMoves.ForEach(func(i int) {
		at := mg.pseudoLegalMoves.At(i)
		switch {
		case at.MoveOf() == mg.pvMove:
			mg.pseudoLegalMoves.Set(i, at.SetValue(ValueMax))
		case at.MoveOf() == mg.killerMoves[0]:
			mg.pseudoLegalMoves.Set(i, at.SetValue(-4000))
		case at.MoveOf() == mg.killerMoves[1]:
			mg.pseudoLegalMoves.Set(i, at.SetValue(-4001))
		}
	})
	mg.pseudoLegalMoves.Sort()
	// remove internal sort value
	mg.pseudoLegalMoves.ForEach(func(i int) {
		mg.pseudoLegalMoves.Set(i, mg.pseudoLegalMoves.At(i).MoveOf())
	})
	return mg.pseudoLegalMoves
}

// GenerateLegalMoves generates legal moves for the next player.
// Uses GeneratePseudoLegalMoves and filters out illegal moves.
//
func (mg *Movegen) GenerateLegalMoves(position *position.Position, mode GenMode) *moveslice.MoveSlice {
	mg.legalMoves.Clear()
	mg.GeneratePseudoLegalMoves(position, mode)
	mg.pseudoLegalMoves.FilterCopy(mg.legalMoves, func(i int) bool {
		return position.IsLegalMove(mg.pseudoLegalMoves.At(i))
	})
	return mg.legalMoves
}

// GetNextMove returns the next move for the given position. Usually this would be used in a loop
// during search.
//
// If a PV move is set with setPV(Move pv) this will be returned first
// and will not be returned at its normal place.
// Killer moves will be played as soon as possible. As Killer moves are stored for
// the whole ply a Killer move might not be valid for the current position. Therefore
// we need to wait until they are generated by the phased move generation. Killers will
// then be pushed to the top of the list of the generation stage.
//
// To reuse this on the sames position a call to ResetOnDemand() is necessary. This
// is not necessary when a different position is called as this func will reset it self
// in this case.
func (mg *Movegen) GetNextMove(position *position.Position, mode GenMode) Move {

	// if the position changes during iteration the iteration
	// will be reset and generation will be restart with the
	// new position.
	if position.ZobristKey() != mg.currentIteratorKey {
		mg.onDemandMoves.Clear()
		mg.currentODStage = odNew
		mg.pvMovePushed = false
		mg.takeIndex = 0
		mg.currentIteratorKey = position.ZobristKey()
	}

	// ad takeIndex
	// With the takeIndex we can take from the front of the vector
	// without removing the element from the vector which would
	// be expensive as all elements would have to be shifted.
	// (although our Moveslice class can handle this efficiently
	// through a similar mechanism)

	// If the list is currently empty and we have not generated all moves yet
	// generate the next batch until we have new moves or there are no more
	// moves to generate
	if mg.onDemandMoves.Len() == 0 {
		mg.fillOnDemandMoveList(position, mode)
	}

	// If we have generated moves we will return the first move and
	// increase the takeIndex to the next move. If the list is empty
	// even after all stages of generating we have no more moves
	// and return MOVE_NONE
	// If we have pushed a pvMove into the list we will need to
	// skip this pvMove for each subsequent phases.
	if mg.onDemandMoves.Len() != 0 {

		// Handle PvMove
		// if we pushed a pv move and the list is not empty we
		// check if the pv is the next move in list and skip it.
		if mg.currentODStage != od1 &&
			mg.pvMovePushed &&
			(*mg.onDemandMoves)[mg.takeIndex].MoveOf() == mg.pvMove.MoveOf() {

			// skip pv move
			mg.takeIndex++

			// We found the pv move and skipped it.
			// No need to check this for this generation cycle
			mg.pvMovePushed = false

			// PV move last in move list
			if mg.takeIndex >= mg.onDemandMoves.Len() {
				// The pv move was the last move in this iterations list.
				// We will try to generate more moves. If no more moves
				// can be generated we will return MOVE_NONE.
				// Otherwise we return the move below.
				mg.takeIndex = 0
				mg.onDemandMoves.Clear()
				mg.fillOnDemandMoveList(position, mode)
				// no more moves - return MOVE_NONE
				if mg.onDemandMoves.Len() == 0 {
					return MoveNone
				}
			}
		}

		// we have at least one move in the list and
		// it is not the pvMove. Increase the takeIndex
		// and return the move
		move := (*mg.onDemandMoves)[mg.takeIndex].MoveOf()
		mg.takeIndex++
		if mg.takeIndex >= mg.onDemandMoves.Len() {
			mg.takeIndex = 0
			mg.onDemandMoves.Clear()
		}
		return move // remove internal sort value
	}

	// no more moves to be generated
	mg.takeIndex = 0
	mg.pvMovePushed = false
	return MoveNone
}

// ResetOnDemand resets the move on demand generator to start fresh.
// Also deletes Killer and PV moves
func (mg *Movegen) ResetOnDemand() {
	mg.onDemandMoves.Clear()
	mg.currentODStage = odNew
	mg.currentIteratorKey = 0
	mg.pvMove = MoveNone
	mg.pvMovePushed = false
	mg.takeIndex = 0
}

// SetPvMove sets a PV move which should be returned first by
// the OnDemand MoveGenerator.
func (mg *Movegen) SetPvMove(move Move) {
	mg.pvMove = move.MoveOf()
}

// StoreKiller provides the on demand move generator with a new killer move
// which should be returned as soon as possible when generating moves with
// the on demand generator.
func (mg *Movegen) StoreKiller(move Move) {
	// check if already stored in first slot - if so return
	moveOf := move.MoveOf()
	if mg.killerMoves[0] == moveOf {
		return
	} else if mg.killerMoves[1] == moveOf { // if in second slot move it to first
		mg.killerMoves[1] = mg.killerMoves[0]
		mg.killerMoves[0] = moveOf
	} else {
		// add it to first slot und move first to second
		mg.killerMoves[1] = mg.killerMoves[0]
		mg.killerMoves[0] = moveOf
	}
}

// HasLegalMove determines if we have at least one legal move. We only have to find
// one legal move. We search for any KING, PAWN, KNIGHT, BISHOP, ROOK, QUEEN move
// and return immediately if we found one.
// The order of our search is approx from the most likely to the least likely
func (mg *Movegen) HasLegalMove(position *position.Position) bool {

	nextPlayer := position.NextPlayer()
	nextPlayerBb := position.OccupiedBb(nextPlayer)

	// KING
	// We do not need to check castling as possible castling implies King or Rook moves
	kingSquare := position.KingSquare(nextPlayer)
	tmpMoves := GetPseudoAttacks(King, kingSquare) &^ nextPlayerBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		if position.IsLegalMove(CreateMove(kingSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	myPawns := position.PiecesBb(nextPlayer, Pawn)
	opponentBb := position.OccupiedBb(nextPlayer.Flip())

	// PAWN
	// normal pawn captures to the west (includes promotions)
	tmpMoves = ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North+West) & opponentBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North + East)
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	// normal pawn captures to the east - promotions first
	tmpMoves = ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North+East) & opponentBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North + West)
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	occupiedBb := position.OccupiedAll()

	// pawn pushes - check step one to unoccupied squares
	// don't have to test double steps as they would be redundant to single steps
	// for the purpose of finding at least one legal move
	tmpMoves = ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North) &^ occupiedBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North)
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	// OFFICERS
	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(nextPlayer, pt)
		for pieces != 0 {
			fromSquare := pieces.PopLsb()
			moves := GetPseudoAttacks(pt, fromSquare) &^ nextPlayerBb
			for moves != 0 {
				toSquare := moves.PopLsb()
				if pt > Knight { // sliding pieces
					if Intermediate(fromSquare, toSquare)&occupiedBb == 0 {
						if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
							return true
						}
					}
				} else { // knight cannot be blocked
					if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
						return true
					}
				}
			}
		}
	}

	// en passant captures
	enPassantSquare := position.GetEnPassantSquare()
	if enPassantSquare != SqNone {
		// left
		tmpMoves = ShiftBitboard(enPassantSquare.Bb(), Direction(nextPlayer.Flip().MoveDirection())*North+West) & myPawns
		if tmpMoves != 0 {
			fromSquare := tmpMoves.PopLsb()
			if position.IsLegalMove(CreateMove(fromSquare, fromSquare.To(Direction(nextPlayer.MoveDirection())*North+East), EnPassant, PtNone)) {
				return true
			}
		}
		// right
		tmpMoves = ShiftBitboard(enPassantSquare.Bb(), Direction(nextPlayer.Flip().MoveDirection())*North+East) & myPawns
		if tmpMoves != 0 {
			fromSquare := tmpMoves.PopLsb()
			if position.IsLegalMove(CreateMove(fromSquare, fromSquare.To(Direction(nextPlayer.MoveDirection())*North+West), EnPassant, PtNone)) {
				return true
			}
		}
	}

	// no move found
	return false
}

// Regex for UCI notation (UCI)
var regexUciMove = regexp.MustCompile("([a-h][1-8][a-h][1-8])([NBRQnbrq])?")

// GetMoveFromUci Generates all legal moves and matches the given UCI
// move string against them. If there is a match the actual move is returned.
// Otherwise MoveNone is returned.
//
// As this uses string creation and comparison this is not very efficient.
// Use only when performance is not critical.
func (mg *Movegen) GetMoveFromUci(posPtr *position.Position, uciMove string) Move {
	matches := regexUciMove.FindStringSubmatch(uciMove)
	if matches == nil {
		return MoveNone
	}

	// get the parts from the pattern match
	movePart := matches[1]
	promotionPart := ""
	if len(matches) == 3 {
		// we allow lower case promotion letters
		// not really UCI but many input files have this wrong
		promotionPart = strings.ToUpper(matches[2])
	}

	// check against all legal moves on position
	mg.GenerateLegalMoves(posPtr, GenAll)
	for _, m := range *mg.legalMoves {
		if m.StringUci() == movePart+promotionPart {
			// move found
			return m
		}
	}
	// move not found
	return MoveNone
}

var regexSanMove = regexp.MustCompile("([NBRQK])?([a-h])?([1-8])?x?([a-h][1-8]|O-O-O|O-O)(=?([NBRQ]))?([!?+#]*)?")

// GetMoveFromSan Generates all legal moves and matches the given SAN
// move string against them. If there is a match the actual move is returned.
// Otherwise MoveNone is returned.
//
// As this uses string creation and comparison this is not very efficient.
// Use only when performance is not critical.
func (mg *Movegen) GetMoveFromSan(posPtr *position.Position, sanMove string) Move {
	matches := regexSanMove.FindStringSubmatch(sanMove)
	if matches == nil {
		return MoveNone
	}

	// get parts
	pieceType := matches[1]
	disambFile := matches[2]
	disambRank := matches[3]
	toSquare := matches[4]
	promotion := matches[6]
	// checkSign := matches[7]

	movesFound := 0
	moveFromSAN := MoveNone

	// check against all legal moves on position
	mg.GenerateLegalMoves(posPtr, GenAll)
	for _, genMove := range *mg.legalMoves {

		// castling moves
		if genMove.MoveType() == Castling {
			kingToSquare := genMove.To()
			var castlingString string
			switch kingToSquare {
			case SqG1: // white king side
				fallthrough
			case SqG8: // black king side
				castlingString = "O-O"
				break
			case SqC1: // white queen side
				fallthrough
			case SqC8: // black queen side
				castlingString = "O-O-O"
				break
			default:
				log.Error("Move type CASTLING but wrong to square: %s %s", castlingString, kingToSquare.String())
				continue
			}
			if castlingString == toSquare {
				moveFromSAN = genMove
				movesFound++
				continue
			}
		}

		// normal moves
		moveTarget := genMove.To().String()
		if moveTarget == toSquare {

			// determine if piece types match - if not skip
			legalPt := posPtr.GetPiece(genMove.From()).TypeOf()
			legalPtChar := legalPt.Char()
			if (len(pieceType) == 0 || legalPtChar != pieceType) &&
				(len(pieceType) != 0 || legalPt != Pawn) {
				continue
			}

			// Disambiguation File
			if len(disambFile) != 0 && genMove.From().FileOf().String() != disambFile {
				continue
			}

			// Disambiguation Rank
			if len(disambRank) != 0 && genMove.From().RankOf().String() != disambRank {
				continue
			}

			// promotion
			if (len(promotion) != 0 && genMove.PromotionType().Char() != promotion) ||
				(len(promotion) == 0 && genMove.MoveType() == Promotion) {
				continue
			}

			// we should have our move if we end up here
			moveFromSAN = genMove
			movesFound++
		}
	}

	// we should only have one move here
	if movesFound > 1 {
		log.Warningf("SAN move %s is ambiguous (%d matches) on %s!", sanMove, movesFound, posPtr.StringFen())
	} else if movesFound == 0 || !moveFromSAN.IsValid() {
		log.Warningf("SAN move not valid! SAN move %s not found on position: %s", sanMove, posPtr.StringFen())
	} else {
		return moveFromSAN
	}
	// no move found
	return MoveNone
}

// ValidateMove validates if a move is a valid move on the given position
func (mg *Movegen) ValidateMove(p *position.Position, move Move) bool {
	if move == MoveNone {
		return false
	}
	ml := mg.GenerateLegalMoves(p, GenAll)
	for _, m := range *ml {
		if move.MoveOf() == m {
			return true
		}
	}
	return false
}

// PvMove returns the current PV move
func (mg *Movegen) PvMove() Move {
	return mg.pvMove
}

// KillerMoves returns a pointer to the killer moves array
func (mg *Movegen) KillerMoves() *[2]Move {
	return &mg.killerMoves
}

// String returns a string representation of a MoveGen instance
func (mg *Movegen) String() string {
	return fmt.Sprintf("MoveGen: { OnDemand Stage: { %d }, PV Move: %s Killer Move 1: %s Killer Move 2: %s }",
		mg.currentODStage, mg.pvMove.String(), mg.killerMoves[0].String(), mg.killerMoves[1].String())
}

// //////////////////////////////////////////////////////
// // Private
// //////////////////////////////////////////////////////

// States for the on demand move generator
const (
	odNew = iota
	odPv  = iota
	od1   = iota
	od2   = iota
	od3   = iota
	od4   = iota
	od5   = iota
	od6   = iota
	od7   = iota
	od8   = iota
	odEnd = iota
)

// This calls the actual generation of moves in phases. The phases match roughly
// the order of most promising moves first.
func (mg *Movegen) fillOnDemandMoveList(p *position.Position, mode GenMode) {
	for mg.onDemandMoves.Len() == 0 && mg.currentODStage < odEnd {
		switch mg.currentODStage {
		case odNew:
			mg.currentODStage = odPv
			fallthrough
		case odPv:
			// If a pvMove is set we return it first and filter it out before
			// returning a move
			if mg.pvMove != MoveNone {
				switch mode {
				case GenAll:
					mg.pvMovePushed = true
					mg.onDemandMoves.PushBack(mg.pvMove)
				case GenCap:
					if p.IsCapturingMove(mg.pvMove) {
						mg.pvMovePushed = true
						mg.onDemandMoves.PushBack(mg.pvMove)
					}
				case GenNonCap:
					if !p.IsCapturingMove(mg.pvMove) {
						mg.pvMovePushed = true
						mg.onDemandMoves.PushBack(mg.pvMove)
					}
				}
			}
			// decide which state we should continue with
			// captures or non captures or both
			if mode&GenCap != 0 {
				mg.currentODStage = od1
			} else {
				mg.currentODStage = od4
			}
		case od1: // capture
			mg.generatePawnMoves(p, GenCap, mg.onDemandMoves)
			mg.currentODStage = od2
		case od2:
			mg.generateMoves(p, GenCap, mg.onDemandMoves)
			mg.currentODStage = od3
		case od3:
			mg.generateKingMoves(p, GenCap, mg.onDemandMoves)
			mg.currentODStage = od4
		case od4:
			if mode&GenNonCap != 0 {
				mg.currentODStage = od5
			} else {
				mg.currentODStage = odEnd
			}
		case od5: // non capture
			mg.generatePawnMoves(p, GenNonCap, mg.onDemandMoves)
			mg.pushKiller(mg.onDemandMoves)
			mg.currentODStage = od6
		case od6:
			mg.generateCastling(p, GenNonCap, mg.onDemandMoves)
			mg.pushKiller(mg.onDemandMoves)
			mg.currentODStage = od7
		case od7:
			mg.generateMoves(p, GenNonCap, mg.onDemandMoves)
			mg.pushKiller(mg.onDemandMoves)
			mg.currentODStage = od8
		case od8:
			mg.generateKingMoves(p, GenNonCap, mg.onDemandMoves)
			mg.pushKiller(mg.onDemandMoves)
			mg.currentODStage = odEnd
		case odEnd:
			break
		}
		// sort the list according to sort values encoded in the move
		if mg.onDemandMoves.Len() > 0 {
			mg.onDemandMoves.Sort()
		}
	} // while onDemandMoves.empty()
}

func (mg *Movegen) pushKiller(m *moveslice.MoveSlice) {
	// Killer may only be returned if they actually are valid moves
	// in this position which we can't know as Killers are stored
	// for the whole ply. Obviously checking if the killer move is valid
	// is expensive (part of a whole move generation) so we only re-sort
	// them to the top once they are actually generated

	// Find the move in the list. If move not found ignore killer.
	// Otherwise move element to the front.
	for i := 0; i < len(*m); i++ {
		move := &(*m)[i]
		if mg.killerMoves[1] == move.MoveOf() {
			(*move).SetValue(Value(-4001))
		}
		if mg.killerMoves[0] == move.MoveOf() {
			(*move).SetValue(Value(-4000))
		}
	}
}

func (mg *Movegen) generatePawnMoves(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {

	nextPlayer := position.NextPlayer()
	myPawns := position.PiecesBb(nextPlayer, Pawn)
	oppPieces := position.OccupiedBb(nextPlayer.Flip())
	gamePhase := position.GamePhase()
	piece := MakePiece(nextPlayer, Pawn)

	// captures
	if mode&GenCap != 0 {

		// This algorithm shifts the own pawn bitboard in the direction of pawn captures
		// and ANDs it with the opponents pieces. With this we get all possible captures
		// and can easily create the moves by using a loop over all captures and using
		// the backward shift for the from-Square.
		// All moves get sort values so that sort order should be:
		//   captures: most value victim least value attacker - promotion piece value
		//   non captures: killer (TBD), promotions, castling, normal moves (position value)
		// Values for sorting are descending - the most valuable move has the highest value.
		// Values are not compatible to position evaluation values outside of the move
		// generator.

		var tmpCaptures, promCaptures Bitboard

		for _, dir := range []Direction{West, East} {
			// normal pawn captures - promotions first
			tmpCaptures = ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North+dir) & oppPieces
			promCaptures = tmpCaptures & nextPlayer.PromotionRankBb()
			// promotion captures
			for promCaptures != 0 {
				toSquare := promCaptures.PopLsb()
				fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North - dir)
				// value is the delta of values from the two pieces involved plus the positional value
				value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
					PosValue(piece, toSquare, gamePhase)
				// add the possible promotion moves to the move list and also add value of the promoted piece type
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, value+Queen.ValueOf()))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, value+Knight.ValueOf()))
				// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
				// therefore we give them lower sort order
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, value+Rook.ValueOf()-Value(2000)))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, value+Bishop.ValueOf()-Value(2000)))
			}
			// non promotion pawn captures
			tmpCaptures &= ^nextPlayer.PromotionRankBb()
			for tmpCaptures != 0 {
				toSquare := tmpCaptures.PopLsb()
				fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North - dir)
				// value is the delta of values from the two pieces involved plus the positional value
				value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
					PosValue(piece, toSquare, gamePhase)
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
			}
		}

		// en passant captures
		enPassantSquare := position.GetEnPassantSquare()
		if enPassantSquare != SqNone {
			for _, dir := range []Direction{West, East} {
				tmpCaptures = ShiftBitboard(enPassantSquare.Bb(),
					Direction(nextPlayer.Flip().MoveDirection())*North+dir) & myPawns
				if tmpCaptures != 0 {
					fromSquare := tmpCaptures.PopLsb()
					toSquare := fromSquare.To(Direction(nextPlayer.MoveDirection())*North - dir)
					// value is the positional value of the piece at this game phase
					value := PosValue(piece, toSquare, gamePhase)
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, EnPassant, PtNone, value))
				}
			}
		}
	}

	// non captures
	if mode&GenNonCap != 0 {

		//  Move my pawns forward one step and keep all on not occupied squares
		//  Move pawns now on rank 3 (rank 6) another square forward to check for pawn doubles.
		//  Loop over pawns remaining on unoccupied squares and add moves.

		// pawns - check step one to unoccupied squares
		tmpMoves := ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North) & ^position.OccupiedAll()
		// pawns double - check step two to unoccupied squares
		tmpMovesDouble := ShiftBitboard(tmpMoves&nextPlayer.PawnDoubleRank(), Direction(nextPlayer.MoveDirection())*North) & ^position.OccupiedAll()

		// single pawn steps - promotions first
		promMoves := tmpMoves & nextPlayer.PromotionRankBb()
		for promMoves != 0 {
			toSquare := promMoves.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			// value for non captures is lowered by 10k
			value := Value(-10_000)
			// add the possible promotion moves to the move list and also add value of the promoted piece type
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, value+Queen.ValueOf()))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, value+Knight.ValueOf()))
			// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
			// therefore we give them lower sort order
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, value+Rook.ValueOf()-Value(2000)))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, value+Bishop.ValueOf()-Value(2000)))
		}
		// double pawn steps
		for tmpMovesDouble != 0 {
			toSquare := tmpMovesDouble.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North).
				To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
		// normal single pawn steps
		tmpMoves &= ^nextPlayer.PromotionRankBb()
		for tmpMoves != 0 {
			toSquare := tmpMoves.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}
}

func (mg *Movegen) generateCastling(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	occupiedBB := position.OccupiedAll()

	// castling - pseudo castling - we will not check if we are in check after the move
	// or if we have passed an attacked square with the king or if the king has been in check

	if mode&GenNonCap != 0 && position.CastlingRights() != CastlingNone {
		cr := position.CastlingRights()
		if nextPlayer == White { // white
			if cr.Has(CastlingWhiteOO) && Intermediate(SqE1, SqH1)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE1, SqG1, Castling, PtNone, Value(-5000)))
			}
			if cr.Has(CastlingWhiteOOO) && Intermediate(SqE1, SqA1)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE1, SqC1, Castling, PtNone, Value(-5000)))
			}
		} else { // black
			if cr.Has(CastlingBlackOO) && Intermediate(SqE8, SqH8)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE8, SqG8, Castling, PtNone, Value(-5000)))
			}
			if cr.Has(CastlingBlackOOO) && Intermediate(SqE8, SqA8)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE8, SqC8, Castling, PtNone, Value(-5000)))
			}
		}
	}
}

func (mg *Movegen) generateKingMoves(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	piece := MakePiece(nextPlayer, King)
	gamePhase := position.GamePhase()
	kingSquareBb := position.PiecesBb(nextPlayer, King)
	fromSquare := kingSquareBb.PopLsb()

	// pseudo attacks include all moves no matter if the king would be in check
	pseudoMoves := GetPseudoAttacks(King, fromSquare)

	// captures
	if mode&GenCap != 0 {
		captures := pseudoMoves & position.OccupiedBb(nextPlayer.Flip())
		for captures != 0 {
			toSquare := captures.PopLsb()
			value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
				PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}

	// non captures
	if mode&GenNonCap != 0 {
		nonCaptures := pseudoMoves &^ position.OccupiedAll()
		for nonCaptures != 0 {
			toSquare := nonCaptures.PopLsb()
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}
}

// generates officers moves using the attacks pre-computed with magic bitboards
// Performance improvement to the previous loop based version:
// Old version:
// Test took 2.0049508s for 10.000.000 iterations
// Test took 200 ns per iteration
// Iterations per sec 4.987.653
// This version:
// Test took 1.516326s for 10.000.000 iterations
// Test took 151 ns per iteration
// Iterations per sec 6.594.887
// Improvement: +32%
func (mg *Movegen) generateMoves(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	gamePhase := position.GamePhase()
	occupiedBb := position.OccupiedAll()

	// loop through all piece types, get pseudo attacks for the piece and
	// AND it with the opponents pieces.
	// For sliding pieces check if there are other pieces in between the
	// piece and the target square. If free this is a valid move (or
	// capture)

	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(nextPlayer, pt)
		piece := MakePiece(nextPlayer, pt)

		for pieces != 0 {
			fromSquare := pieces.PopLsb()

			moves := GetAttacksBb(pt, fromSquare, occupiedBb)

			// captures
			if mode&GenCap != 0 {
				captures := moves & position.OccupiedBb(nextPlayer.Flip())
				for captures != 0 {
					toSquare := captures.PopLsb()
					value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() + PosValue(piece, toSquare, gamePhase)
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
				}
			}

			// non captures
			if mode&GenNonCap != 0 {
				nonCaptures := moves &^ occupiedBb
				for nonCaptures != 0 {
					toSquare := nonCaptures.PopLsb()
					value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
				}
			}
		}
	}
}

// generates officers moves
// DEBUG - kept for testing
func (mg *Movegen) generateMovesOld(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	gamePhase := position.GamePhase()
	occupiedBb := position.OccupiedAll()

	// loop through all piece types, get pseudo attacks for the piece and
	// AND it with the opponents pieces.
	// For sliding pieces check if there are other pieces in between the
	// piece and the target square. If free this is a valid move (or
	// capture)

	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(nextPlayer, pt)
		piece := MakePiece(nextPlayer, pt)

		for pieces != 0 {
			fromSquare := pieces.PopLsb()
			pseudoMoves := GetPseudoAttacks(pt, fromSquare)

			// captures
			if mode&GenCap != 0 {
				captures := pseudoMoves & position.OccupiedBb(nextPlayer.Flip())
				for captures != 0 {
					toSquare := captures.PopLsb()
					if pt > Knight { // sliding pieces
						if Intermediate(fromSquare, toSquare)&occupiedBb == 0 {
							value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
								PosValue(piece, toSquare, gamePhase)
							ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
						}
					} else { // king and knight cannot be blocked
						value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
							PosValue(piece, toSquare, gamePhase)
						ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
					}
				}
			}

			// non captures
			if mode&GenNonCap != 0 {
				nonCaptures := pseudoMoves &^ occupiedBb
				for nonCaptures != 0 {
					toSquare := nonCaptures.PopLsb()
					if pt > Knight { // sliding pieces
						if Intermediate(fromSquare, toSquare)&occupiedBb == 0 {
							value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
							ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
						}
					} else { // king and knight cannot be blocked
						value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
						ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
					}
				}
			}
		}
	}
}

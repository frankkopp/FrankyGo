//
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

// Package movegen contains functionality to create moves on a
// chess position. It implements several variants like
// generate pseudo legal moves, legal moves or on demand
// generation of pseudo legal moves.
package movegen

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/frankkopp/FrankyGo/internal/attacks"
	"github.com/frankkopp/FrankyGo/internal/history"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/pkg/position"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

const removeSortValue = true

// Movegen data structure. Create new move generator via
//  movegen.NewMoveGen()
// Creating this directly will not work.
type Movegen struct {
	pseudoLegalMoves *moveslice.MoveSlice
	legalMoves       *moveslice.MoveSlice

	onDemandMoves          *moveslice.MoveSlice
	currentODZobrist       Key
	onDemandEvasionTargets Bitboard
	currentODStage         int8
	takeIndex              int

	killerMoves  [2]Move
	pvMove       Move
	pvMovePushed bool
	historyData  *history.History
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// GenMode generation modes for on demand move generation.
//  GenZero     GenMode = 0b00
//	GenNonQuiet GenMode = 0b01
//	GenQuiet    GenMode = 0b10
//	GenAll      GenMode = 0b11
type GenMode int

// GenMode generation modes for on demand move generation.
const (
	GenZero     GenMode = 0b00
	GenNonQuiet GenMode = 0b01
	GenQuiet    GenMode = 0b10
	GenAll      GenMode = 0b11
)

// NewMoveGen creates a new instance of a move generator
// This is the only time when we allocate new memory. The instance
// will not create any move lists during normal move generation
// as it will reuse pre-created internal lists which will
// be returned via pointer to a caller.
// OBS: Be careful when trying to store the list of generated
// moves as the underlying list will be changed when move gen
// is called again. A deep copy is necessary if you need a
// copy of the move list.
// For a deep copy use:
//  moveslice.MoveSlice.Clone()
func NewMoveGen() *Movegen {
	tmpMg := &Movegen{
		pseudoLegalMoves: moveslice.NewMoveSlice(MaxMoves),
		legalMoves:       moveslice.NewMoveSlice(MaxMoves),

		onDemandMoves:          moveslice.NewMoveSlice(MaxMoves),
		currentODZobrist:       0,
		onDemandEvasionTargets: BbZero,
		currentODStage:         odNew,
		takeIndex:              0,

		killerMoves:  [2]Move{MoveNone, MoveNone},
		pvMove:       MoveNone,
		pvMovePushed: false,
		historyData:  nil,
	}
	return tmpMg
}

// GeneratePseudoLegalMoves generates pseudo moves for the next player. Does not check if
// king is left in check or if it passes an attacked square when castling or has been in check
// before castling.
//
// If a PV move is set with setPV(Move pv) this move will be returned first and will
// not be returned at its normal place.
//
// Killer moves will be played as soon as possible after non quiet moves. As Killer moves
// are stored for the whole ply a Killer move might not be valid for the current position.
// Therefore we need to wait until they are generated. Killer moves will then be pushed
// to the top of the the quiet moves.
//
// Evasion is a parameter given when the position is in check and only evasion moves should
// be generated. For testing purposes this is a parameter but obviously we could determine
// checks very quickly internally in this function.
// The idea of evasion is to avoid generating moves which are obviously not getting the
// king out of check. This may reduce the total number of generated moves but there might
// still be a few non legal moves. This is the case if considering and calculating all
// possible scenarios is more expensive than to just generate the move and dismiss it later.
// Because of beta cuts off we quite often will never have to check the full legality
// of these moves anyway.
func (mg *Movegen) GeneratePseudoLegalMoves(p *position.Position, mode GenMode, evasion bool) *moveslice.MoveSlice {
	// re-use move list
	mg.pseudoLegalMoves.Clear()

	// when in check only generate moves either blocking or capturing the attacker
	if evasion {
		mg.onDemandEvasionTargets = mg.getEvasionTargets(p)
	}

	// first generate all non quiet moves
	if mode&GenNonQuiet != 0 {
		mg.generatePawnMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
		// castling never captures
		mg.generateKingMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
		mg.generateMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
	}
	// second generate all other moves
	if mode&GenQuiet != 0 {
		mg.generatePawnMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
		if !evasion { // no castling when in check
			mg.generateCastling(p, GenQuiet, mg.pseudoLegalMoves)
		}
		mg.generateKingMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
		mg.generateMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.pseudoLegalMoves)
	}

	// PV, Killer and history handling
	mg.updateSortValues(p, mg.pseudoLegalMoves)

	// sort moves
	mg.pseudoLegalMoves.Sort()

	// remove internal sort value
	if removeSortValue {
		mg.pseudoLegalMoves.ForEach(func(i int) {
			mg.pseudoLegalMoves.Set(i, mg.pseudoLegalMoves.At(i).MoveOf())
		})
	}

	return mg.pseudoLegalMoves
}

// GenerateLegalMoves generates legal moves for the next player.
// Uses GeneratePseudoLegalMoves and filters out illegal moves.
// Usually only used for root moves generation as this is expensive. During
// the AlphaBeta search we will only use pseudo legal move generation.
func (mg *Movegen) GenerateLegalMoves(position *position.Position, mode GenMode) *moveslice.MoveSlice {
	mg.legalMoves.Clear()
	mg.GeneratePseudoLegalMoves(position, mode, false)
	mg.pseudoLegalMoves.FilterCopy(mg.legalMoves, func(i int) bool {
		return position.IsLegalMove(mg.pseudoLegalMoves.At(i))
	})
	return mg.legalMoves
}

// GetNextMove is the main function for phased generation of pseudo legal moves.
// It returns the next move for the given position and will usually be called in a
// loop during search. As we hope for an early beta cut this will save time as not
// all moves will have been generated.
//
// To reuse this on the same position a call to ResetOnDemand() is necessary. This
// is not necessary when a different position is called as this func will reset it self
// in this case.
//
// If a PV move is set with setPV(Move pv) this will be returned first
// and will not be returned at its normal place.
//
// Killer moves will be played as soon as possible. As Killer moves are stored for
// the whole ply a Killer move might not be valid for the current position. Therefore
// we need to wait until they are generated by the phased move generation. Killers will
// then be pushed to the top of the list of the generation stage.
//
// Evasion is a parameter given when the position is in check and only evasion moves should
// be generated. For testing purposes this is a parameter for now but obviously we could
// determine checks very quickly internally in this function.
// The idea of evasion is to avoid generating moves which are obviously not getting the
// king out of check. This may reduce the total number of generated moves but there might
// still be a few non legal moves. This is the case if considering and calculating all
// possible scenarios is more expensive than to just generate the move and dismiss it later.
// Because of beta cuts off we quite often will never have to check the full legality
// of these moves anyway.
func (mg *Movegen) GetNextMove(p *position.Position, mode GenMode, evasion bool) Move {

	// if the position changes during iteration the iteration
	// will be reset and generation will be restart with the
	// new position.
	if p.ZobristKey() != mg.currentODZobrist {
		mg.onDemandMoves.Clear()
		mg.onDemandEvasionTargets = BbZero
		mg.currentODStage = odNew
		mg.pvMovePushed = false
		mg.takeIndex = 0
		mg.currentODZobrist = p.ZobristKey()
	}

	// when in check only generate moves either blocking or capturing the attacker
	if evasion && mg.onDemandEvasionTargets == BbZero {
		mg.onDemandEvasionTargets = mg.getEvasionTargets(p)
	}

	// ad takeIndex
	// With the takeIndex we can take from the front of the vector
	// without removing the element from the vector which would
	// be expensive as all elements would have to be shifted.
	// (although our Moveslice class can handle this efficiently
	// through a similar mechanism of Go slices)

	// If the list is currently empty and we have not generated all moves yet
	// generate the next batch until we have new moves or there are no more
	// moves to generate
	if mg.onDemandMoves.Len() == 0 {
		mg.fillOnDemandMoveList(p, mode, evasion)
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
			// No need to check this again for this generation cycle
			mg.pvMovePushed = false

			// PV move last in move list
			if mg.takeIndex >= mg.onDemandMoves.Len() {
				// The pv move was the last move in this iterations list.
				// We will try to generate more moves. If no more moves
				// can be generated we will return MOVE_NONE.
				// Otherwise we return the move below.
				mg.takeIndex = 0
				mg.onDemandMoves.Clear()
				mg.fillOnDemandMoveList(p, mode, evasion)
				// no more moves - return MOVE_NONE
				if mg.onDemandMoves.Len() == 0 {
					return MoveNone
				}
			}
		}

		// we have at least one move in the list and it is not the
		// pvMove. Increase the takeIndex and return the move.
		// Also remove internal sort value before returning the move.
		var move Move
		if removeSortValue {
			move = (*mg.onDemandMoves)[mg.takeIndex].MoveOf()
		} else {
			move = (*mg.onDemandMoves)[mg.takeIndex]
		}
		mg.takeIndex++
		if mg.takeIndex >= mg.onDemandMoves.Len() {
			mg.takeIndex = 0
			mg.onDemandMoves.Clear()
		}
		return move
	}

	// no more moves to be generated
	mg.takeIndex = 0
	mg.pvMovePushed = false
	return MoveNone
}

// ResetOnDemand resets the move on demand generator to start fresh.
// Also deletes PV moves.
func (mg *Movegen) ResetOnDemand() {
	mg.onDemandMoves.Clear()
	mg.onDemandEvasionTargets = BbZero
	mg.currentODStage = odNew
	mg.currentODZobrist = 0
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
	}
	// if in second slot or not there at all move it to first
	// add it to first slot und move first to second
	mg.killerMoves[1] = mg.killerMoves[0]
	mg.killerMoves[0] = moveOf
}

// SetHistoryData provides a pointer to the search's history data
// for the move generator so it can optimize sorting.
func (mg *Movegen) SetHistoryData(historyData *history.History) {
	mg.historyData = historyData
}

// HasLegalMove determines if we have at least one legal move. We only have to find
// one legal move. We search for any KING, PAWN, KNIGHT, BISHOP, ROOK, QUEEN move
// and return immediately if we found one.
// The order of our search is approx from the most likely to the least likely.
func (mg *Movegen) HasLegalMove(position *position.Position) bool {

	us := position.NextPlayer()
	usBb := position.OccupiedBb(us)

	// KING
	// We do not need to check castling as possible castling implies King or Rook moves
	kingSquare := position.KingSquare(us)
	tmpMoves := GetAttacksBb(King, kingSquare, BbZero) &^ usBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		if position.IsLegalMove(CreateMove(kingSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	myPawns := position.PiecesBb(us, Pawn)
	occupiedBb := position.OccupiedAll()
	opponentBb := position.OccupiedBb(us.Flip())

	// PAWN
	// pawns - check step one to unoccupied squares
	tmpMoves = ShiftBitboard(myPawns, us.MoveDirection()) & ^position.OccupiedAll()
	// pawns double - check step two to unoccupied squares
	tmpMovesDouble := ShiftBitboard(tmpMoves&us.PawnDoubleRank(), us.MoveDirection()) & ^position.OccupiedAll()
	// double pawn steps
	for tmpMovesDouble != 0 {
		toSquare := tmpMovesDouble.PopLsb()
		fromSquare := toSquare.To(us.Flip().MoveDirection()).To(us.Flip().MoveDirection())
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}
	// normal single pawn steps
	tmpMoves &= ^us.PromotionRankBb()
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(us.Flip().MoveDirection())
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	// normal pawn captures to the west (includes promotions)
	tmpMoves = ShiftBitboard(myPawns, us.MoveDirection()+West) & opponentBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(us.Flip().MoveDirection() + East)
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	// normal pawn captures to the east - promotions first
	tmpMoves = ShiftBitboard(myPawns, us.MoveDirection()+East) & opponentBb
	for tmpMoves != 0 {
		toSquare := tmpMoves.PopLsb()
		fromSquare := toSquare.To(us.Flip().MoveDirection() + West)
		if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
			return true
		}
	}

	// OFFICERS
	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(us, pt)
		for pieces != 0 {
			fromSquare := pieces.PopLsb()
			moves := GetAttacksBb(pt, fromSquare, occupiedBb) &^ usBb
			for moves != 0 {
				toSquare := moves.PopLsb()
				if position.IsLegalMove(CreateMove(fromSquare, toSquare, Normal, PtNone)) {
					return true
				}
			}
		}
	}

	// en passant captures
	enPassantSquare := position.GetEnPassantSquare()
	if enPassantSquare != SqNone {
		// left
		tmpMoves = ShiftBitboard(enPassantSquare.Bb(), us.Flip().MoveDirection()+West) & myPawns
		if tmpMoves != 0 {
			fromSquare := tmpMoves.PopLsb()
			if position.IsLegalMove(CreateMove(fromSquare, fromSquare.To(us.MoveDirection()+East), EnPassant, PtNone)) {
				return true
			}
		}
		// right
		tmpMoves = ShiftBitboard(enPassantSquare.Bb(), us.Flip().MoveDirection()+East) & myPawns
		if tmpMoves != 0 {
			fromSquare := tmpMoves.PopLsb()
			if position.IsLegalMove(CreateMove(fromSquare, fromSquare.To(us.MoveDirection()+West), EnPassant, PtNone)) {
				return true
			}
		}
	}

	// no move found
	return false
}

// Regex for UCI notation (UCI).
var regexUciMove = regexp.MustCompile("([a-h][1-8][a-h][1-8])([NBRQnbrq])?")

// GetMoveFromUci Generates all legal moves and matches the given UCI
// move string against them. If there is a match the actual move is returned.
// Otherwise MoveNone is returned.
//
// As this uses string creation and comparison this is not very efficient.
// Use only when performance is not critical.
func (mg *Movegen) GetMoveFromUci(posPtr *position.Position, uciMove string) (Move, error) {
	matches := regexUciMove.FindStringSubmatch(uciMove)
	if matches == nil {
		return MoveNone, fmt.Errorf("provided uci move string does not match uci move pattern")
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
			return m, nil
		}
	}
	// move not found
	return MoveNone, fmt.Errorf("UCI move not valid! UCI move %s not found on position: %s", uciMove, posPtr.StringFen())

}

var regexSanMove = regexp.MustCompile("([NBRQK])?([a-h])?([1-8])?x?([a-h][1-8]|O-O-O|O-O)(=?([NBRQ]))?([!?+#]*)?")

// GetMoveFromSan Generates all legal moves and matches the given SAN
// move string against them. If there is a match the actual move is returned.
// Otherwise MoveNone is returned.
//
// As this uses string creation and comparison this is not very efficient.
// Use only when performance is not critical.
func (mg *Movegen) GetMoveFromSan(posPtr *position.Position, sanMove string) (Move, error) {
	matches := regexSanMove.FindStringSubmatch(sanMove)
	if matches == nil {
		return MoveNone, fmt.Errorf("provided san move string does not match san move pattern")
	}

	// get parts
	pieceType := matches[1]
	disambFile := matches[2]
	disambRank := matches[3]
	toSquare := matches[4]
	promotion := matches[6]
	// checkSign := matches[7] - ignore

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
			case SqC1: // white queen side
				fallthrough
			case SqC8: // black queen side
				castlingString = "O-O-O"
			default:
				log.Panicf("Move type CASTLING but wrong to square: %s %s", castlingString, kingToSquare.String())
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

	// we should only have one move here or none was found to match the SAN string
	if movesFound > 1 {
		return MoveNone, fmt.Errorf("SAN move %s is ambiguous (%d matches) on %s", sanMove, movesFound, posPtr.StringFen())
	} else if movesFound == 0 || !moveFromSAN.IsValid() {
		return MoveNone, fmt.Errorf("SAN move not valid! SAN move %s not found on position: %s", sanMove, posPtr.StringFen())
	}
	return moveFromSAN, nil
}

// ValidateMove validates if a move is a valid legal move on the given position
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
func (mg *Movegen) fillOnDemandMoveList(p *position.Position, mode GenMode, evasion bool) {
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
				case GenNonQuiet:
					if p.IsCapturingMove(mg.pvMove) {
						mg.pvMovePushed = true
						mg.onDemandMoves.PushBack(mg.pvMove)
					}
				case GenQuiet:
					if !p.IsCapturingMove(mg.pvMove) {
						mg.pvMovePushed = true
						mg.onDemandMoves.PushBack(mg.pvMove)
					}
				}
			}
			// decide which state we should continue with
			// captures or non captures or both
			if mode&GenNonQuiet != 0 {
				mg.currentODStage = od1
			} else {
				mg.currentODStage = od4
			}
		case od1: // pawns: capture and high value promotion
			mg.generatePawnMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
			mg.currentODStage = od2
		case od2: // officer capture
			mg.generateMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
			mg.currentODStage = od3
		case od3: // king captures
			mg.generateKingMoves(p, GenNonQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
			mg.currentODStage = od4
		case od4:
			if mode&GenQuiet != 0 {
				mg.currentODStage = od5
			} else {
				mg.currentODStage = odEnd
			}
		case od5: // pawn: non capture
			mg.generatePawnMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
			mg.currentODStage = od6
		case od6: // castling
			if !evasion { // no castlings when in check
				mg.generateCastling(p, GenQuiet, mg.onDemandMoves)
				mg.updateSortValues(p, mg.onDemandMoves)
			}
			mg.currentODStage = od7
		case od7: // officer non capture
			mg.generateMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
			mg.currentODStage = od8
		case od8: // king non capture
			mg.generateKingMoves(p, GenQuiet, evasion, mg.onDemandEvasionTargets, mg.onDemandMoves)
			mg.updateSortValues(p, mg.onDemandMoves)
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

// Move order heuristics based on history data.
func (mg *Movegen) updateSortValues(p *position.Position, moveList *moveslice.MoveSlice) {
	us := p.NextPlayer()
	// iterate over all available moves and update the
	// sort value if the move is the PV or a Killer move.
	// Also update the sort value for history and counter
	// move significance.
	for i := 0; i < len(*moveList); i++ {
		move := &(*moveList)[i]
		switch {
		case move.MoveOf() == mg.pvMove: // PV move
			(*move).SetValue(ValueMax)
		case move.MoveOf() == mg.killerMoves[1]: // Killer 2
			(*move).SetValue(1000)
		case move.MoveOf() == mg.killerMoves[0]: // Killer 1
			(*move).SetValue(1001)
		case mg.historyData != nil: // historical search data

			// History Count
			// Moves that cause a beta cut in the search get an increasing value
			// which favors many repetitions and deep searches.
			// We use the history count to improve the sort value of a move
			// If and how much a sort value has to be improved for a move is
			// difficult to predict - this needs testing and experimentation.
			// The current way is a hard cut for values <1000 and then 1 point
			// per 1000 count points.
			// It is also yet unclear if the history count table should be
			// reused for several consecutive searches or just for one search.
			// TODO: Testing
			count := mg.historyData.HistoryCount[us][move.From()][move.To()]
			value := Value(count / 100)

			// Counter Move History
			// When we have a counter move which caused a beta cut off before we
			// bump up its sort value
			// TODO: Testing
			if mg.historyData.CounterMoves[p.LastMove().From()][p.LastMove().To()] == move.MoveOf() {
				value += 500
			}

			// update move sort value
			if value > 0 { // only touch the value if it would be improved
				(*move).SetValue(move.ValueOf() + value)
				// out.Printf("HistoryCount: %s = %d / %d ==> %d \n", move.StringUci(), count, preValue, preValue+value)
			}
		}
	}
}

// getEvasionTargets returns the number of attackers and a Bitboard with target
// squares for generated moves when the position has check against the next
// player. Most of the moves will not even be generated as they will not
// have these target squares. These target squares cover the attacking
// (checker) piece and any squares in between the attacker and the king
// in case of the attacker being a slider.
// If we have more than one attacker we can skip everything apart from
// king moves.
func (mg *Movegen) getEvasionTargets(p *position.Position) Bitboard {
	us := p.NextPlayer()
	ourKing := p.KingSquare(us)
	// find all target squares which either capture or block the attacker
	evasionTargets := attacks.AttacksTo(p, ourKing, us.Flip())
	// we can only block attacks of sliders of there is not more
	// than one attacker
	popCount := evasionTargets.PopCount()
	if popCount == 1 {
		atck := evasionTargets.Lsb()
		// sliding pieces
		if p.GetPiece(atck).TypeOf() > Knight {
			evasionTargets |= Intermediate(atck, ourKing)
			return evasionTargets
		}
	}
	if popCount > 1 {
		return BbZero
	}
	return evasionTargets
}

func (mg *Movegen) generatePawnMoves(position *position.Position, mode GenMode, evasion bool, evasionTargets Bitboard, ml *moveslice.MoveSlice) {

	nextPlayer := position.NextPlayer()
	myPawns := position.PiecesBb(nextPlayer, Pawn)
	oppPieces := position.OccupiedBb(nextPlayer.Flip())
	gamePhase := position.GamePhase()
	piece := MakePiece(nextPlayer, Pawn)

	// captures
	if mode&GenNonQuiet != 0 {

		// This algorithm shifts the own pawn bitboard in the direction of pawn captures
		// and ANDs it with the opponents pieces. With this we get all possible captures
		// and can easily create the moves by using a loop over all captures and using
		// the backward shift for the from-Square.
		// All moves get sort values so that sort order should be:
		//   captures: most value victim least value attacker - promotion piece value
		//   non captures: promotions, castling, normal moves (position value)

		// When we are in check only evasion moves are generated. E.g. all moves need to
		// target these evasion squares. That is either capturing the attacker or blocking
		// a sliding attacker.

		var tmpCaptures, promCaptures Bitboard

		for _, dir := range []Direction{West, East} {
			// all pawn captures
			tmpCaptures = ShiftBitboard(myPawns, nextPlayer.MoveDirection()+dir) & oppPieces

			// filter evasion targets if in check
			if evasion {
				tmpCaptures &= evasionTargets
			}

			// normal pawn captures - promotions first
			promCaptures = tmpCaptures & nextPlayer.PromotionRankBb()
			// promotion captures
			for promCaptures != 0 {
				toSquare := promCaptures.PopLsb()
				fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection() - dir)
				// value is the delta of values from the two pieces involved minus the promoted pawn
				value := position.GetPiece(toSquare).ValueOf() - (2 * Pawn.ValueOf())
				// add the possible promotion moves to the move list and also add value of the promoted piece type
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, value+Queen.ValueOf()+5000))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, value+Knight.ValueOf()+1500))
				// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
				// therefore we give them a lower sort order
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, value+Rook.ValueOf()-Value(5000)))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, value+Bishop.ValueOf()-Value(5000)))
			}

			// non promotion pawn captures
			tmpCaptures &= ^nextPlayer.PromotionRankBb()
			for tmpCaptures != 0 {
				toSquare := tmpCaptures.PopLsb()
				fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection() - dir)
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
				tmpCaptures = ShiftBitboard(enPassantSquare.Bb(), nextPlayer.Flip().MoveDirection()+dir) & myPawns
				if tmpCaptures != 0 {
					fromSquare := tmpCaptures.PopLsb()
					toSquare := fromSquare.To(nextPlayer.MoveDirection() - dir)
					// value is the positional value of the piece at this game phase
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, EnPassant, PtNone, PosValue(piece, toSquare, gamePhase)))
				}
			}
		}

		// we treat Queen and Knight promotions as non quiet moves
		promMoves := ShiftBitboard(myPawns, nextPlayer.MoveDirection()) &^ position.OccupiedAll() & nextPlayer.PromotionRankBb()
		if evasion {
			promMoves &= evasionTargets
		}
		for promMoves != 0 {
			toSquare := promMoves.PopLsb()
			fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection())
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, 2000-Pawn.ValueOf()+Queen.ValueOf()))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, 1500-Pawn.ValueOf()+Knight.ValueOf()))
		}
	}

	// non captures
	if mode&GenQuiet != 0 {

		//  Move my pawns forward one step and keep all on not occupied squares
		//  Move pawns now on rank 3 (rank 6) another square forward to check for pawn doubles.
		//  Loop over pawns remaining on unoccupied squares and add moves.

		// When we are in check only evasion moves are generated. E.g. all moves need to
		// target these evasion squares. That is either capturing the attacker or blocking
		// a sliding attacker.

		// pawns - check step one to unoccupied squares
		tmpMoves := ShiftBitboard(myPawns, nextPlayer.MoveDirection()) & ^position.OccupiedAll()
		// pawns double - check step two to unoccupied squares
		tmpMovesDouble := ShiftBitboard(tmpMoves&nextPlayer.PawnDoubleRank(), nextPlayer.MoveDirection()) & ^position.OccupiedAll()

		// filter evasion targets if in check
		if evasion {
			tmpMoves &= evasionTargets
			tmpMovesDouble &= evasionTargets
		}

		// single pawn steps - promotions first
		promMoves := tmpMoves & nextPlayer.PromotionRankBb()
		for promMoves != 0 {
			toSquare := promMoves.PopLsb()
			fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection())
			// value for non captures is lowered by 10k
			// we treat Queen and Knight promotions as non quiet moves and they are generated above
			// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
			// therefore we give them lower sort order
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, Rook.ValueOf()-Value(6000)))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, Bishop.ValueOf()-Value(6000)))
		}
		// double pawn steps
		for tmpMovesDouble != 0 {
			toSquare := tmpMovesDouble.PopLsb()
			fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection()).To(nextPlayer.Flip().MoveDirection())
			value := PosValue(piece, toSquare, gamePhase) - 2000
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
		// normal single pawn steps
		tmpMoves &= ^nextPlayer.PromotionRankBb()
		for tmpMoves != 0 {
			toSquare := tmpMoves.PopLsb()
			fromSquare := toSquare.To(nextPlayer.Flip().MoveDirection())
			value := PosValue(piece, toSquare, gamePhase) - 2000
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}
}

func (mg *Movegen) generateCastling(position *position.Position, mode GenMode, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	occupiedBB := position.OccupiedAll()

	// castling - pseudo castling - we will not check if we are in check after the move
	// or if we have passed an attacked square with the king or if the king has been in check

	if mode&GenQuiet != 0 && position.CastlingRights() != CastlingNone {
		cr := position.CastlingRights()
		if nextPlayer == White { // white
			if cr.Has(CastlingWhiteOO) && Intermediate(SqE1, SqH1)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE1, SqG1, Castling, PtNone, Value(0)))
			}
			if cr.Has(CastlingWhiteOOO) && Intermediate(SqE1, SqA1)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE1, SqC1, Castling, PtNone, Value(0)))
			}
		} else { // black
			if cr.Has(CastlingBlackOO) && Intermediate(SqE8, SqH8)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE8, SqG8, Castling, PtNone, Value(0)))
			}
			if cr.Has(CastlingBlackOOO) && Intermediate(SqE8, SqA8)&occupiedBB == 0 {
				ml.PushBack(CreateMoveValue(SqE8, SqC8, Castling, PtNone, Value(0)))
			}
		}
	}
}

func (mg *Movegen) generateKingMoves(p *position.Position, mode GenMode, evasion bool, evasionTargets Bitboard, ml *moveslice.MoveSlice) {
	us := p.NextPlayer()
	them := us.Flip()
	piece := MakePiece(us, King)
	gamePhase := p.GamePhase()
	kingSquareBb := p.PiecesBb(us, King)
	fromSquare := kingSquareBb.PopLsb()

	// attacks include all moves no matter if the king would be in check
	pseudoMoves := GetAttacksBb(King, fromSquare, BbZero)

	// captures
	if mode&GenNonQuiet != 0 {
		captures := pseudoMoves & p.OccupiedBb(them)
		for captures != 0 {
			toSquare := captures.PopLsb()
			// in case we are in check we only generate king moves to target squares which
			// are not attacked by the opponent
			if !evasion || attacks.AttacksTo(p, toSquare, them).PopCount() == 0 {
				value := 2000 + p.GetPiece(toSquare).ValueOf() - p.GetPiece(fromSquare).ValueOf() + PosValue(piece, toSquare, gamePhase)
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
			}
		}
	}

	// non captures
	if mode&GenQuiet != 0 {
		nonCaptures := pseudoMoves &^ p.OccupiedAll()
		for nonCaptures != 0 {
			toSquare := nonCaptures.PopLsb()
			// in case we are in check we only generate king moves to target squares which
			// are not attacked by the opponent
			if !evasion || attacks.AttacksTo(p, toSquare, them).PopCount() == 0 {
				value := PosValue(piece, toSquare, gamePhase) - 2000
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
			}
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
func (mg *Movegen) generateMoves(position *position.Position, mode GenMode, evasion bool, evasionTargets Bitboard, ml *moveslice.MoveSlice) {
	nextPlayer := position.NextPlayer()
	gamePhase := position.GamePhase()
	occupiedBb := position.OccupiedAll()

	// Loop through all piece types, get attacks for the piece.
	// When we are in check (evasion=true) only evasion moves are generated. E.g. all
	// moves need to target these evasion squares. That is either capturing the
	// attacker or blocking a sliding attacker.

	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(nextPlayer, pt)
		piece := MakePiece(nextPlayer, pt)

		for pieces != 0 {
			fromSquare := pieces.PopLsb()

			moves := GetAttacksBb(pt, fromSquare, occupiedBb)

			// captures
			if mode&GenNonQuiet != 0 {
				captures := moves & position.OccupiedBb(nextPlayer.Flip())
				if evasion {
					captures &= evasionTargets
				}
				for captures != 0 {
					toSquare := captures.PopLsb()
					value := 2000 + position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() + PosValue(piece, toSquare, gamePhase)
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
				}
			}

			// non captures
			if mode&GenQuiet != 0 {
				nonCaptures := moves &^ occupiedBb
				if evasion {
					nonCaptures &= evasionTargets
				}
				for nonCaptures != 0 {
					toSquare := nonCaptures.PopLsb()
					value := PosValue(piece, toSquare, gamePhase) - 2000
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
				}
			}
		}
	}
}

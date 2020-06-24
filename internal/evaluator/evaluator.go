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

// Package evaluator contains structures and functions to calculate
// the value of a chess position to be used in a chess engine search
package evaluator

import (
	"strings"

	"github.com/op/go-logging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	attacks2 "github.com/frankkopp/FrankyGo/internal/attacks"
	. "github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/util"
	"github.com/frankkopp/FrankyGo/pkg/position"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

var out = message.NewPrinter(language.German)

// Evaluator  represents a data structure and functionality for
// evaluating chess positions by using various evaluation
// heuristics like material, positional values, pawn structure, etc.
//  Create a new instance with NewEvaluator()
type Evaluator struct {
	log *logging.Logger

	position        *position.Position
	gamePhaseFactor float64
	us              Color
	them            Color
	ourKing         Square
	theirKing       Square
	kingRing        [ColorLength]Bitboard
	allPieces       Bitboard
	ourPieces       Bitboard

	score Score

	attacks *attacks2.Attacks

	pawnCache *pawnCache
}

// to avoid object creation and memory allocation
// during evaluation we reuse this tmp Score.
var tmpScore = Score{}

// pre-computed list.
var threshold [GamePhaseMax + 1]int16

// initialize pre-computed values.
func init() {
	for i := 0; i <= GamePhaseMax; i++ {
		gamePhaseFactor := float64(i) / GamePhaseMax
		threshold[i] = Settings.Eval.LazyEvalThreshold + int16(float64(Settings.Eval.LazyEvalThreshold)*gamePhaseFactor)
	}
}

// NewEvaluator creates a new instance of an Evaluator.
func NewEvaluator() *Evaluator {
	e := &Evaluator{
		log:       myLogging.GetLog(),
		attacks:   attacks2.NewAttacks(),
		pawnCache: nil,
	}
	if Settings.Eval.UsePawnCache {
		e.pawnCache = newPawnCache()
	} else {
		e.log.Info("Pawn Cache is disabled in configuration")
	}
	return e
}

// InitEval initializes data structures and values which are used several times
// Is called at the beginning of Evaluate() but can be called separately to be able
// to run single evaluations in unit tests.
func (e *Evaluator) InitEval(p *position.Position) {
	// set some value which we need regularly
	e.position = p
	e.gamePhaseFactor = p.GamePhaseFactor()
	e.us = p.NextPlayer()
	e.them = e.us.Flip()
	e.ourKing = e.position.KingSquare(e.us)
	e.theirKing = e.position.KingSquare(e.them)
	e.kingRing[e.us] = GetAttacksBb(King, e.ourKing, BbZero)
	e.kingRing[e.them] = GetAttacksBb(King, e.theirKing, BbZero)
	e.allPieces = e.position.OccupiedAll()
	e.ourPieces = e.position.OccupiedBb(e.us)

	// reset all values
	e.score.MidGameValue = 0
	e.score.EndGameValue = 0

	// reset attacks
	if Settings.Eval.UseAttacksInEval {
		e.attacks.Clear()
	}
}

// Evaluate calculates a value for a chess positions by
// using various evaluation heuristics like material,
// positional values, pawn structure, etc.
// It calls InitEval and then the internal evaluation function
// which calculates the value for the position of the given
// position for the current game phase and from the
// view of the next player.
func (e *Evaluator) Evaluate(position *position.Position) Value {
	e.InitEval(position)
	return e.evaluate()
}

// value adds up the mid and end games scores after multiplying
// them with the game phase factor.
func (e *Evaluator) value() Value {
	return e.score.ValueFromScore(e.gamePhaseFactor)
}

// internal evaluation to sum up all partial evaluations.
// This assumes that InitEval() has been called beforehand.
func (e *Evaluator) evaluate() Value {
	// if not enough material on the board to achieve a mate it is a draw
	if e.position.HasInsufficientMaterial() {
		return ValueDraw
	}

	// Each position is evaluated from the view of the white
	// player. Before returning the value this will be adjusted
	// to the next player's color.
	// All heuristic should return a value in centi pawns or
	// have a dedicated configurable weight to adjust and test

	// Material
	if Settings.Eval.UseMaterialEval {
		e.score.MidGameValue = int16(e.position.Material(White) - e.position.Material(Black))
		e.score.EndGameValue = e.score.MidGameValue
	}

	// Positional values
	if Settings.Eval.UsePositionalEval {
		e.score.MidGameValue += int16(e.position.PsqMidValue(White) - e.position.PsqMidValue(Black))
		e.score.EndGameValue += int16(e.position.PsqEndValue(White) - e.position.PsqEndValue(Black))
	}

	// early exit
	// arbitrary threshold - in early phases (game phase = 1.0) this is doubled
	// in late phases it stands as it is
	if Settings.Eval.UseLazyEval {
		valueFromScore := e.value()
		th := threshold[e.position.GamePhase()]
		value := util.Abs16(int16(valueFromScore))
		if value > th {
			return e.finalEval(valueFromScore)
		}
	}

	// Get all attacks
	// find out where this should be done to be most effective
	// This is expensive and we should use this investment as often as
	// possible. If we could use it in search as well we could move
	// creating this to an earlier point in time in the search
	if Settings.Eval.UseAttacksInEval {
		e.attacks.Compute(e.position)
	}

	// evaluate pawns
	if Settings.Eval.UsePawnEval {
		// white and black are handled in evaluatePawns()
		e.score.Add(e.evaluatePawns())
	}

	// evaluate pieces - builds attacks and mobility
	if Settings.Eval.UseAdvancedPieceEval {
		e.score.Add(e.evalPiece(White, Knight))
		e.score.Sub(e.evalPiece(Black, Knight))
		e.score.Add(e.evalPiece(White, Bishop))
		e.score.Sub(e.evalPiece(Black, Bishop))
		e.score.Add(e.evalPiece(White, Rook))
		e.score.Sub(e.evalPiece(Black, Rook))
		e.score.Add(e.evalPiece(White, Queen))
		e.score.Sub(e.evalPiece(Black, Queen))
	}

	// mobility
	if Settings.Eval.UseAttacksInEval && Settings.Eval.UseMobility {
		e.score.MidGameValue += (e.attacks.Mobility[White] - e.attacks.Mobility[Black]) * Settings.Eval.MobilityBonus
		e.score.EndGameValue += e.score.MidGameValue
	}

	// evaluate king
	if Settings.Eval.UseKingEval {
		e.score.Add(e.evalKing(White))
		e.score.Sub(e.evalKing(Black))
	}

	// TEMPO Bonus for the side to move (helps with evaluation alternation -
	// less difference between side which makes aspiration search faster
	// (not empirically tested)
	e.score.MidGameValue += Settings.Eval.Tempo
	// e.score.EndGameValue += Value(0) // can be ignored

	// value is always from the view of the next player
	valueFromScore := e.value()

	return e.finalEval(valueFromScore)
}

// finalEval returns the value which is calculated always from the view of
// white from the view of the next player of the position.
func (e *Evaluator) finalEval(value Value) Value {
	// we can use the Direction factor to avoid an if statement
	// Direction returns positive 1 for White and negative 1 (-1) for Black
	return value * Value(e.position.NextPlayer().Direction())
}

func (e *Evaluator) evalKing(c Color) *Score {
	tmpScore.MidGameValue = 0
	tmpScore.EndGameValue = 0
	us := c
	them := us.Flip()

	// pawn shield - pawns in front of a castled king get a bonus
	// Higher bonus for middle game, lower or none in end game
	if KingSideCastleMask(us).Has(e.position.KingSquare(us)) {
		count := int16((ShiftBitboard(KingSideCastleMask(us), us.MoveDirection()) & e.position.PiecesBb(us, Pawn)).PopCount())
		tmpScore.MidGameValue += count * Settings.Eval.KingCastlePawnShieldBonus
	} else if QueenSideCastMask(us).Has(e.position.KingSquare(us)) {
		count := int16((ShiftBitboard(QueenSideCastMask(us), us.MoveDirection()) & e.position.PiecesBb(us, Pawn)).PopCount())
		tmpScore.MidGameValue += count * Settings.Eval.KingCastlePawnShieldBonus
	}

	// king safety / attacks to the king and king ring
	if Settings.Eval.UseAttacksInEval {
		enemyAttacks := e.kingRing[us] & e.attacks.All[them]
		ourDefence := e.kingRing[us] & e.attacks.All[us]
		// malus for difference between attacker and defender
		if enemyAttacks > ourDefence {
			tmpScore.MidGameValue -= int16(enemyAttacks.PopCount()-ourDefence.PopCount()) * Settings.Eval.KingDangerMalus
			tmpScore.EndGameValue -= tmpScore.MidGameValue
		} else {
			tmpScore.MidGameValue += int16(ourDefence.PopCount()-enemyAttacks.PopCount()) * Settings.Eval.KingDefenderBonus
			tmpScore.EndGameValue += tmpScore.MidGameValue
		}

		// king ring attacks
		if a := e.attacks.All[us] & e.kingRing[them]; a > 0 {
			tmpScore.MidGameValue += Settings.Eval.KingRingAttacksBonus
			tmpScore.EndGameValue += Settings.Eval.KingRingAttacksBonus
		}
	}
	return &tmpScore
}

// evalPiece is the evaluation function for all pieces except pawns and kings.
func (e *Evaluator) evalPiece(c Color, pieceType PieceType) *Score {
	tmpScore.MidGameValue = 0
	tmpScore.EndGameValue = 0
	us := c
	them := us.Flip()

	// get bitboard with all pieces of this color and type
	pieceBb := e.position.PiecesBb(us, pieceType)

	// piece type specific evaluation which are done once
	// for all pieces of one type
	switch pieceType {
	case Bishop:
		// bonus for pair
		if pieceBb.PopCount() > 1 {
			tmpScore.MidGameValue += Settings.Eval.BishopPairBonus
			tmpScore.EndGameValue += Settings.Eval.BishopPairBonus
		}
	}

	// loop over all pieces of the given piece type
	for pieceBb != BbZero {
		sq := pieceBb.PopLsb()

		// piece type specific evaluations per piece of piece type
		switch pieceType {
		case Knight:
			e.knightEval(us, them, sq)
		case Bishop:
			e.bishopEval(us, them, sq)
		case Rook:
			e.rookEval(sq, us)
		case Queen:
			// none yet
		}
	}

	return &tmpScore
}

func (e *Evaluator) rookEval(sq Square, us Color) {
	// same file as queen
	if sq.FileOf().Bb()&e.position.PiecesBb(us, Queen) > 0 {
		tmpScore.MidGameValue += Settings.Eval.RookOnQueenFileBonus
		tmpScore.EndGameValue += Settings.Eval.RookOnQueenFileBonus
	}

	// open file / semi open file (no own pawns on the file)
	if sq.FileOf().Bb()&e.position.PiecesBb(us, Pawn) == 0 {
		tmpScore.MidGameValue += Settings.Eval.RookOnOpenFileBonus
		// s.EndGameValue += 0
	}

	// trapped by king
	// on same row as king but on the outside from king
	kingSquare := e.position.KingSquare(us)
	if KingSideCastleMask(us).Has(kingSquare) {
		if sq.RankOf() == kingSquare.RankOf() && sq > kingSquare { // east of king
			tmpScore.MidGameValue -= Settings.Eval.RookTrappedMalus
		}
	} else if QueenSideCastMask(us).Has(kingSquare) {
		if sq.RankOf() == kingSquare.RankOf() && sq < kingSquare { // west of king
			tmpScore.MidGameValue -= Settings.Eval.RookTrappedMalus
		}
	}
}

func (e *Evaluator) bishopEval(us Color, them Color, sq Square) {
	// behind a pawn
	down := them.MoveDirection()
	if ShiftBitboard(e.position.PiecesBb(us, Pawn), down)&sq.Bb() > 0 {
		tmpScore.MidGameValue += Settings.Eval.MinorBehindPawnBonus
		// s.EndGameValue += 0
	}

	// malus for pawns on same color - worse in end game
	if SquaresBb(White).Has(sq) { // on white square
		popCount := int16((e.position.PiecesBb(us, Pawn) & SquaresBb(White)).PopCount())
		// s.MidGameValue -= 0
		tmpScore.EndGameValue -= Settings.Eval.BishopPawnMalus * popCount
	} else { // on black square
		popCount := int16((e.position.PiecesBb(us, Pawn) & SquaresBb(Black)).PopCount())
		// s.MidGameValue -= 0
		tmpScore.EndGameValue -= Settings.Eval.BishopPawnMalus * popCount
	}

	// long diagonal / seeing center
	popCount := int16((GetAttacksBb(Bishop, sq, BbZero) & CenterSquares).PopCount())
	tmpScore.MidGameValue += Settings.Eval.BishopCenterAimBonus * popCount
	// s.EndGameValue += 0

	// bishop blocked / mobility
	if (us == White && sq.RankOf() == Rank1) || (us == Black && sq.RankOf() == Rank8) {
		if GetAttacksBb(Bishop, sq, e.allPieces)&^e.position.OccupiedBb(us) == BbZero {
			tmpScore.MidGameValue -= Settings.Eval.BishopBlockedMalus
			tmpScore.EndGameValue -= Settings.Eval.BishopBlockedMalus
		}
	}
}

func (e *Evaluator) knightEval(us Color, them Color, sq Square) {
	// Knight behind pawn
	down := them.MoveDirection()
	if ShiftBitboard(e.position.PiecesBb(us, Pawn), down)&sq.Bb() > 0 {
		tmpScore.MidGameValue += Settings.Eval.MinorBehindPawnBonus
		// s.EndGameValue += 0
	}
}

// Report prints a report about the evaluations done. Used in debugging.
func (e *Evaluator) Report() string {
	var report strings.Builder

	report.WriteString("Evaluation Report\n")
	report.WriteString("=============================================\n")
	report.WriteString(out.Sprintf("Position: %s\n", e.position.StringFen()))
	report.WriteString(out.Sprintf("%s\n", e.position.StringBoard()))
	report.WriteString(out.Sprintf("GamePhase Factor: %f\n", e.position.GamePhaseFactor()))
	report.WriteString(out.Sprintf("(evals from the view of white player)\n", e.Evaluate(e.position)))
	// report.WriteString(out.Sprintf("Material    : %d\n", e.material()))
	// report.WriteString(out.Sprintf("Positional  : %d\n", e.positional()))
	// report.WriteString(out.Sprintf("Tempo       : %d\n", e.tempo()))
	report.WriteString(out.Sprintf("-------------------------\n", e.Evaluate(e.position)))
	report.WriteString(out.Sprintf("Eval value  : %d \n(from the view of next player = %s)\n", e.Evaluate(e.position), e.position.NextPlayer().String()))

	return report.String()
}

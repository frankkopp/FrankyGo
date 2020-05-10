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

package config

type evalConfiguration struct {

	// evaluation values
	UseLazyEval       bool
	LazyEvalThreshold int16

	Tempo int16

	UseAttacksInEval bool

	UsePawnEval   bool
	UsePawnCache  bool
	PawnCacheSize int

	UseMobility   bool
	MobilityBonus int16

	UseAdvancedPieceEval bool
	BishopPairBonus      int16
	MinorBehindPawnBonus int16
	BishopPawnMalus      int16
	BishopCenterAimBonus int16
	BishopBlockedMalus   int16
	RookOnQueenFileBonus int16
	RookOnOpenFileBonus  int16
	RookTrappedMalus     int16
	KingRingAttacksBonus int16

	UseKingEval       bool
	KingDangerMalus   int16
	KingDefenderBonus int16
}

// sets defaults which might be overwritten by config file.
func init() {
	Settings.Eval.UseLazyEval = false
	Settings.Eval.LazyEvalThreshold = 700

	// evaluation value
	Settings.Eval.Tempo = 34

	Settings.Eval.UsePawnEval = true
	Settings.Eval.UsePawnCache = true // not implemented yet
	Settings.Eval.PawnCacheSize = 64  // not implemented yet

	Settings.Eval.UseAttacksInEval = false

	Settings.Eval.UseMobility = false
	Settings.Eval.MobilityBonus = 5 // per piece and attacked square

	Settings.Eval.UseAdvancedPieceEval = false
	Settings.Eval.KingRingAttacksBonus = 10 // per piece and attacked king ring square
	Settings.Eval.MinorBehindPawnBonus = 15 // per piece and times game phase
	Settings.Eval.BishopPairBonus = 20      // once
	Settings.Eval.BishopPawnMalus = 5       // per pawn and times ~game phase
	Settings.Eval.BishopCenterAimBonus = 20 // per bishop and times game phase
	Settings.Eval.BishopBlockedMalus = 40   // per bishop
	Settings.Eval.RookOnQueenFileBonus = 6  // per rook
	Settings.Eval.RookOnOpenFileBonus = 25  // per rook and time game phase
	Settings.Eval.RookTrappedMalus = 40     // per rook and time game phase

	Settings.Eval.UseKingEval = false
	Settings.Eval.KingDangerMalus = 50   // number of number of attacker - defender times malus if attacker > defender
	Settings.Eval.KingDefenderBonus = 10 // number of number of defender - attacker times bonus if attacker <= defender

}

// set defaults for configurations here in case a configuration
// is not available from the config file.
func setupEval() {

}

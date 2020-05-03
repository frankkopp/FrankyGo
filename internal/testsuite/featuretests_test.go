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

package testsuite

import (
	"testing"
	"time"

	"github.com/frankkopp/FrankyGo/internal/config"
)

func TestFeatureTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup tests
	searchTime := 200 * time.Millisecond
	searchDepth := 0

	// Feature Settings
	{
		config.Settings.Search.UseQuiescence = true
		config.Settings.Search.UseQSStandpat = true
		config.Settings.Search.UseSEE = true
		config.Settings.Search.UsePromNonQuiet = true

		config.Settings.Search.UseTT = true
		config.Settings.Search.TTSize = 256
		config.Settings.Search.UseTTValue = true
		config.Settings.Search.UseQSTT = true

		config.Settings.Search.UsePVS = true
		config.Settings.Search.UseAspiration = false
		config.Settings.Search.UseMTDf = false

		config.Settings.Search.UseTTMove = true
		config.Settings.Search.UseIID = true
		config.Settings.Search.IIDDepth = 6
		config.Settings.Search.IIDReduction = 2
		config.Settings.Search.UseKiller = true
		config.Settings.Search.UseHistoryCounter = true
		config.Settings.Search.UseCounterMoves = true

		config.Settings.Search.UseMDP = true
		config.Settings.Search.UseRazoring = true
		config.Settings.Search.RazorMargin = 531
		config.Settings.Search.UseNullMove = true
		config.Settings.Search.NmpDepth = 3
		config.Settings.Search.NmpReduction = 2

		config.Settings.Search.UseExt = true
		config.Settings.Search.UseExtAddDepth = true
		config.Settings.Search.UseCheckExt = true
		config.Settings.Search.UseThreatExt = false

		config.Settings.Search.UseRFP = true
		config.Settings.Search.UseFP = true
		config.Settings.Search.UseQFP = true
		config.Settings.Search.UseLmr = true
		config.Settings.Search.LmrDepth = 3
		config.Settings.Search.LmrMovesSearched = 3
		config.Settings.Search.UseLmp = true

		config.Settings.Eval.Tempo = 34
		config.Settings.Eval.UseLazyEval = true
		config.Settings.Eval.LazyEvalThreshold = 700

		config.Settings.Eval.UsePawnCache = false
		config.Settings.Eval.PawnCacheSize = 64
		config.Settings.Eval.UseAttacksInEval = false
		config.Settings.Eval.UseMobility = false
		config.Settings.Eval.MobilityBonus = 5
		config.Settings.Eval.UseAdvancedPieceEval = false
		config.Settings.Eval.BishopPairBonus = 20
		config.Settings.Eval.MinorBehindPawnBonus = 15
		config.Settings.Eval.BishopPawnMalus = 5
		config.Settings.Eval.BishopCenterAimBonus = 20
		config.Settings.Eval.BishopBlockedMalus = 40
		config.Settings.Eval.RookOnQueenFileBonus = 6
		config.Settings.Eval.RookOnOpenFileBonus = 25
		config.Settings.Eval.RookTrappedMalus = 40
		config.Settings.Eval.KingRingAttacksBonus = 10
		config.Settings.Eval.UseKingEval = false
		config.Settings.Eval.KingDangerMalus = 50
		config.Settings.Eval.KingDefenderBonus = 10
	}

	folder := "test/testdata/featuretests/"

	out.Println(FeatureTests(folder, searchTime, searchDepth))
}

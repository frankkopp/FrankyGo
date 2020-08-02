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

package search

import (
	"testing"
	"time"

	"github.com/pkg/profile"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/util"
	"github.com/frankkopp/FrankyGo/pkg/position"
	. "github.com/frankkopp/FrankyGo/pkg/types"
	"github.com/frankkopp/FrankyGo/test/testdata"
)

func Test_savePV(t *testing.T) {
	src := moveslice.NewMoveSlice(10)
	dest := moveslice.NewMoveSlice(10)

	src.PushBack(Move(1234))
	src.PushBack(Move(2345))
	src.PushBack(Move(3456))
	src.PushBack(Move(4567))

	savePV(Move(9999), src, dest)

	// logTest.Debug(dest.String())
	assert.EqualValues(t, 5, dest.Len())
	assert.EqualValues(t, 9999, dest.At(0))
	assert.EqualValues(t, 4567, dest.At(4))
}

func TestMate(t *testing.T) {
	config.Settings.Search.UseBook = false
	s := NewSearch()
	p, _ := position.NewPositionFen("8/8/8/8/8/3K4/R7/5k2 w - -")
	sl := NewSearchLimits()
	sl.Depth = 9
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	assert.EqualValues(t, 9993, s.lastSearchResult.BestValue)
}

func TestDevelopAndTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("./bin")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof

	config.Settings.Search.UseBook = false
	config.Settings.Search.UsePonder = false

	config.Settings.Search.UseQuiescence = true
	config.Settings.Search.UseQSStandpat = false
	config.Settings.Search.UseSEE = true

	config.Settings.Search.UseTT = false
	config.Settings.Search.UseTTMove = false
	config.Settings.Search.UseTTValue = false
	config.Settings.Search.UseQSTT = false
	config.Settings.Search.UseEvalTT = false

	config.Settings.Search.UseAlphaBeta = false
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false

	config.Settings.Search.UseIID = false
	config.Settings.Search.UseKiller = false
	config.Settings.Search.UseHistoryCounter = false
	config.Settings.Search.UseCounterMoves = false

	config.Settings.Search.UseMDP = false
	config.Settings.Search.UseRazoring = false
	config.Settings.Search.UseNullMove = false
	config.Settings.Search.UseExt = false
	config.Settings.Search.UseExtAddDepth = false
	config.Settings.Search.UseCheckExt = false
	config.Settings.Search.UseThreatExt = false
	config.Settings.Search.UseRFP = false
	config.Settings.Search.UseFP = false
	config.Settings.Search.UseQFP = false
	config.Settings.Search.UseLmr = false
	config.Settings.Search.UseLmp = false

	config.Settings.Eval.UseMaterialEval = false
	config.Settings.Eval.UsePositionalEval = false
	config.Settings.Eval.UseLazyEval = false
	config.Settings.Eval.UsePawnEval = false
	config.Settings.Eval.UsePawnCache = false
	config.Settings.Eval.PawnCacheSize = 32
	config.Settings.Eval.UseAttacksInEval = false
	config.Settings.Eval.UseMobility = false
	config.Settings.Eval.UseAdvancedPieceEval = false
	config.Settings.Eval.UseKingEval = false
	config.Settings.Eval.Tempo = 34

	s := NewSearch()
	// "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	// rnbqkbnr/ppppp1pp/5p2/8/3PP3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 2
	// kiwipete
	// r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -
	// 4k2r/1q1p1pp1/p3p3/1pb1P3/2r3P1/P1N1P2p/1PP1Q2P/2R1R1K1 b k -
	// ASP Issue: r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 w kq - 0 1
	// 8/8/8/6k1/8/4K3/8/r7 b - - 23 94
	p := position.NewPosition("2rr2k1/1p2qp1p/1pn1pp2/1N6/3P4/P6P/1P2QPP1/2R2RK1 w - - 0 1")
	sl := NewSearchLimits()
	sl.Depth = 1
	// sl.TimeControl = true
	// sl.MoveTime = 10 * time.Second
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	out.Println("NPS : ", util.Nps(s.nodesVisited, s.lastSearchResult.SearchTime))
}

func TestTimingTTSize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	config.LogLevel = 2
	config.SearchLogLevel = 2
	config.Settings.Search.UseBook = false

	for _, fen := range testdata.Fens[10:30] {
		out.Println(fen)
		var results []string
		for ttSize := 1; ttSize < 10_000; ttSize = ttSize * 2 {
			out.Println("TT Size", ttSize)
			// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()
			// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof
			config.Settings.Search.TTSize = ttSize
			s := NewSearch()
			// p := position.NewPosition()
			p := position.NewPosition(fen)
			sl := NewSearchLimits()
			sl.Depth = 0
			sl.TimeControl = true
			sl.MoveTime = 5 * time.Second
			s.StartSearch(*p, *sl)
			s.WaitWhileSearching()
			nps := util.Nps(s.nodesVisited, s.lastSearchResult.SearchTime)
			results = append(results, out.Sprintf("tt size: %-6d time: %-12s nodes: %-12d depth: %2d/%-2d nps: %-12d stats: %s tt: %s",
				ttSize, s.lastSearchResult.SearchTime, s.nodesVisited, s.lastSearchResult.SearchDepth, s.lastSearchResult.ExtraDepth,
				nps, s.statistics.String(), s.tt.String()))
		}
		for _, r := range results {
			out.Println(r)
		}
		out.Println()
	}
}

// StartPos
// v0.8.0 24.4.2020
// NPS :  2.677.839 (evasion move gen)
// v0.8.0 27.4.2020
// NPS :  2.475.123.
// v1.1.0 7.5.2020
// NPS :  1.935.496
// v1.1.0 10.5. Asp
// NPS :  2.181.181
func TestTiming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./bin")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof

	config.Settings.Search.UseBook = false

	s := NewSearch()
	// "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	// rnbqkbnr/ppppp1pp/5p2/8/3PP3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 2
	// kiwipete
	// r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -
	p := position.NewPosition()
	sl := NewSearchLimits()
	// sl.Depth = 10
	sl.TimeControl = true
	sl.MoveTime = 30 * time.Second
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	out.Println("NPS : ", util.Nps(s.nodesVisited, s.lastSearchResult.SearchTime))
}

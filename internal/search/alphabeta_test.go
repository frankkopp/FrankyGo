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

package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
	"github.com/frankkopp/FrankyGo/internal/util"
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
	sl.Depth = 8
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	assert.EqualValues(t, 9993, s.lastSearchResult.BestValue)
}

func TestTimingTTSize(t *testing.T) {
	t.SkipNow()

	results := []string{}

	for ttSize := 1; ttSize < 10_000; ttSize = ttSize * 2 {

		out.Println("TT Size", ttSize)

		// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()
		// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof
		config.LogLevel = 2
		config.SearchLogLevel = 2
		config.Settings.Search.UseBook = false
		config.Settings.Search.TTSize = ttSize

		s := NewSearch()
		p := position.NewPosition()
		sl := NewSearchLimits()
		// sl.Depth = 8
		sl.TimeControl = true
		sl.MoveTime = 10 * time.Second
		s.StartSearch(*p, *sl)
		s.WaitWhileSearching()
		nps := util.Nps(s.nodesVisited, s.lastSearchResult.SearchTime)
		results = append(results, out.Sprintf("tt size: %-6d time: %-12s nodes: %-12d depth: %2d/%-2d nps: %-12d stats: %s tt: %s",
			ttSize, s.lastSearchResult.SearchTime, s.nodesVisited, s.lastSearchResult.SearchDepth, s.lastSearchResult.ExtraDepth,
			nps, s.statistics.String(), s.tt.String()))
	}

	out.Println()
	for _, r := range results {
		out.Println(r)
	}
}

func TestTiming(t *testing.T) {
	t.SkipNow()
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof
	config.Settings.Search.UseBook = false
	s := NewSearch()
	p := position.NewPosition("r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3")
	sl := NewSearchLimits()
	sl.Depth = 10
	// sl.TimeControl = true
	sl.MoveTime = 30 * time.Second
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	out.Println("TT  : ", s.tt.String())
	out.Println("NPS : ", util.Nps(s.nodesVisited, s.lastSearchResult.SearchTime))
}
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

// Package searchtreesize provides data structures and functionality to
// test the size of the search tree when certain heuristics and prunings
// are activated or deactivated.
package searchtreesize

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/search"
	"github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/util"
)

var out = message.NewPrinter(language.German)

type SingleTest struct {
	Name    string
	Nodes   uint64
	Nps     uint64
	Depth   int
	Extra   int
	Time    time.Duration
	Special uint64
	Move    types.Move
	Value   types.Value
	Pv      moveslice.MoveSlice
}

type Result struct {
	Fen   string
	Tests []SingleTest
}

type TestSums struct {
	SumCounter uint64
	SumNodes   uint64
	SumNps     uint64
	SumDepth   int
	SumExtra   int
	SumTime    time.Duration
	Special    uint64
}

var ptrToSpecial *uint64

func featureTest(depth int, movetime time.Duration, fen *string) Result {
	s := search.NewSearch()
	sl := search.NewSearchLimits()
	sl.Depth = depth
	sl.MoveTime = movetime
	if movetime > 0 {
		sl.TimeControl = true
	}
	result := Result{Fen: *fen}
	p := position.NewPositionFen(*fen)
	// turn off all options to turn them on later for each test
	turnOffFeatures()

	// /////////////////////////////////////////////////////////////////
	// TESTS

	// define which special data pointer to collect
	ptrToSpecial = &s.Statistics().BetaCuts

	// Base
	result.Tests = append(result.Tests, measure(s, sl, p, "00 Base"))

	// + Quiescence
	config.Settings.Search.UseQuiescence = true
	result.Tests = append(result.Tests, measure(s, sl, p, "10 QS"))

	// TESTS
	// /////////////////////////////////////////////////////////////////

	return result
}

func sizeTest(depth int, movetime time.Duration, startFen int, endFen int) {

	out.Printf("Start Search Tree Size Test for depth %d\n", depth)

	// prepare the slice for the tests
	if endFen > len(Fens) {
		endFen = len(Fens)
	}
	if startFen > endFen {
		startFen = endFen
	}
	testFens := Fens[startFen:endFen]

	// prepare slice of results to store them for the report
	results := make([]Result, 0, len(Fens))

	// execute tests and store results
	for _, fen := range testFens {
		results = append(results, featureTest(depth, movetime, &fen))
	}

	// Print result
	out.Printf("################## Results for depth %d ##########################\n\n", depth)
	out.Printf("%-15s | %-6s | %-8s | %-15s | %-12s | %-10s | %-7s | %-12s | %s | %s\n",
		"Test Name", "Move", "Value", "Nodes", "Nps", "Time", "Depth", "Special", "PV", "Fen")
	out.Println("----------------------------------------------------------------------------------------------------------------------------------------------")

	sums := make(map[string]TestSums, len(results))

	// loop through all results and each test within
	// sum up results to later print a summary
	for _, result := range results {
		for _, test := range result.Tests {
			// sum up result for total report
			sums[test.Name] = TestSums{
				SumCounter: sums[test.Name].SumCounter + 1,
				SumNodes:   sums[test.Name].SumNodes + test.Nodes,
				SumNps:     sums[test.Name].SumNps + test.Nps,
				SumDepth:   sums[test.Name].SumDepth + test.Depth,
				SumExtra:   sums[test.Name].SumExtra + test.Extra,
				SumTime:    sums[test.Name].SumTime + test.Time,
				Special:    sums[test.Name].Special + test.Special,
			}
			// print single test result
			out.Printf("%-15s | %-6s | %-8s | %-15d | %-12d | %-10d | %3d/%-3d | %-12d | %s | %s\n",
				test.Name, test.Move.StringUci(), test.Value.String(), test.Nodes, test.Nps,
				test.Time.Milliseconds(), test.Depth, test.Extra, test.Special, test.Pv.StringUci(), result.Fen)
		}
	}
	out.Println("----------------------------------------------------------------------------------------------------------------------------------------------")
	out.Println()

	// print Totals
	for _, test := range results[0].Tests {
		sum := sums[test.Name]
		out.Printf("Test: %-12s  Nodes: %-14d  Nps: %-14d  Time: %-10d Depth: %3d/%-3d Special: %-16d\n",
			test.Name,
			sum.SumNodes/sum.SumCounter,
			sum.SumNps/sum.SumCounter,
			uint64(sum.SumTime.Milliseconds())/sum.SumCounter,
			uint64(sum.SumDepth)/sum.SumCounter,
			uint64(sum.SumExtra)/sum.SumCounter,
			sum.Special/sum.SumCounter)
	}
	out.Println()
}

func measure(s *search.Search, sl *search.Limits, p *position.Position, name string) SingleTest {
	out.Printf("\nTesting  %s ###############################\n", name)
	out.Printf("Position %s \n", p.StringFen())
	s.ClearHash()
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()

	test := SingleTest{
		Name:    name,
		Nodes:   s.NodesVisited(),
		Nps:     util.Nps(s.NodesVisited(), s.LastSearchResult().SearchTime),
		Time:    s.LastSearchResult().SearchTime,
		Depth:   s.LastSearchResult().SearchDepth,
		Extra:   s.LastSearchResult().ExtraDepth,
		Special: 0,
		Move:    s.LastSearchResult().BestMove,
		Value:   s.LastSearchResult().BestValue,
		Pv:      s.LastSearchResult().Pv,
	}

	if ptrToSpecial != nil {
		test.Special = *ptrToSpecial
	}

	return test
}

func turnOffFeatures() {
	config.Settings.Search.UseBook = false
	config.Settings.Search.UsePonder = false
	config.Settings.Search.UseTT = false
	config.Settings.Search.UseQuiescence = false
}

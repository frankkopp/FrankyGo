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

	. "github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/position"
	"github.com/frankkopp/FrankyGo/internal/search"
	"github.com/frankkopp/FrankyGo/internal/types"
	"github.com/frankkopp/FrankyGo/internal/util"
	"github.com/frankkopp/FrankyGo/test/testdata"
)

var out = message.NewPrinter(language.German)

// singleTest holds the result data for a single test
// A single test is one fen with one set of feature executing
// one search according to the settings (depth odr time)
type singleTest struct {
	Name     string
	Nodes    uint64
	Nps      uint64
	Depth    int
	Extra    int
	Time     time.Duration
	Special  uint64
	Special2 uint64
	Move     types.Move
	Value    types.Value
	Pv       moveslice.MoveSlice
}

// result is representing a series of single tests for a single position (FEN)
type result struct {
	Fen   string
	Tests []singleTest
}

// testSums is a helper data structure to sum up all results from a list of
// single tests for a set of features to create a total reports at the end
type testSums struct {
	SumCounter uint64
	SumNodes   uint64
	SumNps     uint64
	SumDepth   int
	SumExtra   int
	SumTime    time.Duration
	Special    uint64
	Special2   uint64
}

var ptrToSpecial *uint64
var ptrToSpecial2 *uint64

// featureTest is called for each set of features for all configured test positions (fens).
// a feature test creates a result instance and stores all single tests into it.
// It sets up the search for the tests and configures the various features for each test.
// Define feature tests in this function.
func featureTest(depth int, movetime time.Duration, fen string) result {
	s := search.NewSearch()
	sl := search.NewSearchLimits()
	sl.Depth = depth
	sl.MoveTime = movetime
	if movetime > 0 {
		sl.TimeControl = true
	}
	r := result{Fen: fen}
	p, _ := position.NewPositionFen(fen)
	// turn off all options to turn them on later for each test
	turnOffFeatures()

	// /////////////////////////////////////////////////////////////////
	// TESTS

	// define which special data pointer to collect
	ptrToSpecial = &s.Statistics().RfpPrunings
	ptrToSpecial2 = &s.Statistics().FpPrunings

	// Base
	// r.Tests = append(r.Tests, measure(s, sl, p, "Base"))

	// + Quiescence
	Settings.Search.UseQuiescence = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "Base+QS"))

	// + QS Standpat
	Settings.Search.UseQSStandpat = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "Standpat"))

	// + TT
	Settings.Search.UseTT = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "TT"))

	// + TTValue
	Settings.Search.UseTTValue = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "TT"))

	// + QS TT
	Settings.Search.UseQSTT = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "QSTT"))

	// + MDP
	Settings.Search.UseMDP = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "MDP"))

	// PVS
	Settings.Search.UsePVS = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "PVS"))

	// PVS
	Settings.Search.UseKiller = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "Killer"))

	Settings.Search.UseTTMove = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "TTMove"))

	// IID
	Settings.Search.UseIID = true
	r.Tests = append(r.Tests, measure(s, sl, p, "BASE"))

	// SEE for qsearch
	Settings.Search.UseSEE = true
	r.Tests = append(r.Tests, measure(s, sl, p, "SEE"))

	// Reverse Futility
	Settings.Search.UseRFP = true
	r.Tests = append(r.Tests, measure(s, sl, p, "RFP"))

	// Null Move
	Settings.Search.UseNullMove = true
	r.Tests = append(r.Tests, measure(s, sl, p, "NMP"))

	// Extensions
	Settings.Search.UseExt = true
	Settings.Search.UseExtAddDepth = true

	Settings.Search.UseCheckExt = true
	r.Tests = append(r.Tests, measure(s, sl, p, "CHECK"))

	// Settings.Search.UseThreatExt = true
	// r.Tests = append(r.Tests, measure(s, sl, p, "THREAT"))

	// Futility
	Settings.Search.UseFP = true
	r.Tests = append(r.Tests, measure(s, sl, p, "FP"))

	// Late Move Reduction
	Settings.Search.UseLmr = true
	// Late Move Pruning
	Settings.Search.UseLmp = true
	r.Tests = append(r.Tests, measure(s, sl, p, "LMR & LMP"))

	Settings.Search.UseExtAddDepth = false
	r.Tests = append(r.Tests, measure(s, sl, p, "NOADD"))

	// TESTS
	// /////////////////////////////////////////////////////////////////

	return r
}

// SizeTest is the main function to call for testing a series of positions
// defined in fens.go with certain search limits (depth or time).
// Results are printed directly to Stdout.
func SizeTest(depth int, movetime time.Duration, startFen int, endFen int) {

	out.Printf("Start Search Tree Size Test for depth %d\n", depth)

	// prepare the slice for the tests
	if endFen > len(testdata.Fens) {
		endFen = len(testdata.Fens)
	}
	if startFen > endFen {
		startFen = endFen
	}
	testFens := testdata.Fens[startFen:endFen]

	// prepare slice of results to store them for the report
	results := make([]result, 0, len(testdata.Fens))

	// execute tests and store results
	for _, fen := range testFens {
		results = append(results, featureTest(depth, movetime, fen))
	}

	// Print result
	out.Printf("\n################## Results for depth %d ##########################\n\n", depth)

	out.Printf("%-15s | %-6s | %-8s | %-15s | %-12s | %-10s | %-7s | %-12s | %-12s |%s | %s\n",
		"Test Name", "Move", "value", "Nodes", "Nps", "Time", "Depth", "Special", "Special2", "PV", "fen")
	out.Println("----------------------------------------------------------------------------------------------------------------------------------------------")

	sums := make(map[string]testSums, len(results))

	// loop through all results and each test within.
	// sum up results to later print a summary
	for _, r := range results {
		for _, test := range r.Tests {
			// sum up result for total report
			sums[test.Name] = testSums{
				SumCounter: sums[test.Name].SumCounter + 1,
				SumNodes:   sums[test.Name].SumNodes + test.Nodes,
				SumNps:     sums[test.Name].SumNps + test.Nps,
				SumDepth:   sums[test.Name].SumDepth + test.Depth,
				SumExtra:   sums[test.Name].SumExtra + test.Extra,
				SumTime:    sums[test.Name].SumTime + test.Time,
				Special:    sums[test.Name].Special + test.Special,
				Special2:   sums[test.Name].Special2 + test.Special2,
			}
			// print single test result
			out.Printf("%-15s | %-6s | %-8s | %-15d | %-12d | %-10d | %3d/%-3d | %-12d | %-12d |%s | %s\n",
				test.Name, test.Move.StringUci(), test.Value.String(), test.Nodes, test.Nps,
				test.Time.Milliseconds(), test.Depth, test.Extra, test.Special, test.Special2, test.Pv.StringUci(), r.Fen)
		}
		out.Println()
	}
	out.Println("----------------------------------------------------------------------------------------------------------------------------------------------")
	out.Print("\n################## Totals/Avg results for each feature test ##################\n\n")

	out.Printf("Date                   : %s\n", time.Now().Local())
	out.Printf("SearchTime             : %s\n", movetime)
	out.Printf("MaxDepth               : %d\n", depth)
	out.Printf("Number of feature tests: %d\n", len(results[0].Tests))
	out.Printf("Number of fens         : %d\n", len(testFens))
	out.Printf("Total tests            : %d\n\n", len(results[0].Tests)*len(testFens))

	// print Totals
	// obs: GO does not order map entries. To get an order when iterating one must iterate over a
	// parallel data structure (e.g. array of map keys) which can be sorted.
	for _, test := range results[0].Tests {
		sum := sums[test.Name]
		out.Printf("Test: %-12s  Nodes: %-14d  Nps: %-14d  Time: %-10d Depth: %3d/%-3d Special: %-14d Special2: %-14d\n",
			test.Name,
			sum.SumNodes/sum.SumCounter,
			sum.SumNps/sum.SumCounter,
			uint64(sum.SumTime.Milliseconds())/sum.SumCounter,
			uint64(sum.SumDepth)/sum.SumCounter,
			uint64(sum.SumExtra)/sum.SumCounter,
			sum.Special/sum.SumCounter,
			sum.Special2/sum.SumCounter)
	}
	out.Println()
}

// measure starts a single search for a feature set on one position and returns one
// singleTest instance as a result
func measure(s *search.Search, sl *search.Limits, p *position.Position, name string) singleTest {
	out.Printf("\nTesting  %s ###############################\n", name)
	out.Printf("Position %s \n", p.StringFen())

	s.ClearHash()
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()

	test := singleTest{
		Name:     name,
		Nodes:    s.NodesVisited(),
		Nps:      util.Nps(s.NodesVisited(), s.LastSearchResult().SearchTime),
		Time:     s.LastSearchResult().SearchTime,
		Depth:    s.LastSearchResult().SearchDepth,
		Extra:    s.LastSearchResult().ExtraDepth,
		Special:  0,
		Special2: 0,
		Move:     s.LastSearchResult().BestMove,
		Value:    s.LastSearchResult().BestValue,
		Pv:       s.LastSearchResult().Pv,
	}

	if ptrToSpecial != nil {
		test.Special = *ptrToSpecial
	}
	if ptrToSpecial2 != nil {
		test.Special2 = *ptrToSpecial2
	}

	return test
}

func turnOffFeatures() {
	Settings.Search.UseBook = false
	Settings.Search.UsePonder = false
	Settings.Search.UseQuiescence = false
	Settings.Search.UseQSStandpat = false
	Settings.Search.UseSEE = false
	Settings.Search.UseTT = false
	Settings.Search.UseTTMove = false
	Settings.Search.UseTTValue = false
	Settings.Search.UseQSTT = false
	Settings.Search.UseIID = false
	Settings.Search.UsePVS = false
	Settings.Search.UseKiller = false
	Settings.Search.UseMDP = false
	Settings.Search.UseNullMove = false
	Settings.Search.UseExt = false
	Settings.Search.UseExtAddDepth = false
	Settings.Search.UseCheckExt = false
	Settings.Search.UseThreatExt = false
	Settings.Search.UseRFP = false
	Settings.Search.UseFP = false
	Settings.Search.UseLmr = false
	Settings.Search.UseLmp = false
}

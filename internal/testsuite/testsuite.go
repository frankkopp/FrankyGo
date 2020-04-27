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

// Package testsuite contains data structures and functionality to be able to
// run chess test suites which contain chess positions as EPD (Extended Position Description).
// EPD contain a standard FEN of a position but also meta data to describe a result for
// a successful test. This could be best move on the position, mate in x or avoid moves.
// https://www.chessprogramming.org/Extended_Position_Description
// For the purpose of testing our chess engine only the opcodes "bm" (best move), "am"
// (avoid move) and "dm" (direct mate) are implemented.
package testsuite

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/movegen"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/position"
	"github.com/frankkopp/FrankyGo/internal/search"
	. "github.com/frankkopp/FrankyGo/internal/types"
	"github.com/frankkopp/FrankyGo/internal/util"
)

var out = message.NewPrinter(language.German)
var log *logging.Logger

// testType defines the data type for the implemented opcode for EPD tests
// which are defined as constants below.
type testType uint8

// Implemented test types
const (
	None testType = iota
	DM   testType = iota
	BM   testType = iota
	AM   testType = iota
)

// resultType define possible results for a tests as a type and constants
type resultType uint8

// resultType define possible results for a tests as a type and constants
const (
	NotTested resultType = iota
	Skipped   resultType = iota
	Failed    resultType = iota
	Success   resultType = iota
)

// SuiteResult data structure to collect sum of the results of tests
type SuiteResult struct {
	Counter          int
	SuccessCounter   int
	FailedCounter    int
	SkippedCounter   int
	NotTestedCounter int
}

// Test defines the data structure for a test after reading in the
// test files. Each EPD from the read file will create an instance
// of this struct and when the tests are run the result will be
// stored back to this instance.
type Test struct {
	id          string
	fen         string
	tType       testType
	targetMoves moveslice.MoveSlice
	mateDepth   int
	target      Move
	actual      Move
	value       Value
	rType       resultType
	line        string
	nps         uint64
}

// TestSuite is the data structure for the running a file of EPD tests.
type TestSuite struct {
	Tests      []*Test
	Time       time.Duration
	Depth      int
	FilePath   string
	LastResult *SuiteResult
}

// NewTestSuite creates an instance of a TestSuite and reads in the given file
// to create test cases which can be run with RunTests()
func NewTestSuite(filePath string, searchTime time.Duration, depth int) (*TestSuite, error) {
	out.Println("Preparing Test Suite", filePath)

	if log == nil {
		log = myLogging.GetLog()
	}

	// Configure logging
	config.LogLevel = 2
	config.SearchLogLevel = 2

	// Configure search for running the tests of the test suite
	config.Settings.Search.UseBook = false

	// read complete file into array of strings
	lines, err := getTestLines(filePath)
	if err != nil {
		return nil, err
	}

	// create the TestSuite instance
	newTestSuite := &TestSuite{
		Tests:    make([]*Test, 0, len(*lines)),
		Time:     searchTime,
		Depth:    depth,
		FilePath: filePath,
	}

	// create tests from given input lines
	for _, line := range *lines {
		test := getTest(line)
		if test == nil {
			continue
		}
		newTestSuite.Tests = append(newTestSuite.Tests, test)
	}

	return newTestSuite, nil
}

// RunTests runs tests on a successfully created instance of a TestSuite
func (ts *TestSuite) RunTests() {

	if len(ts.Tests) == 0 {
		out.Printf("No tests to run\n")
		return
	}

	startTime := time.Now()

	// setup search
	s := search.NewSearch()
	sl := search.NewSearchLimits()
	sl.MoveTime = ts.Time
	sl.Depth = ts.Depth
	if sl.MoveTime > 0 {
		sl.TimeControl = true
	}

	out.Printf("Running Test Suite\n")
	out.Printf("==================================================================\n")
	out.Printf("EPD File:    %s\n", ts.FilePath)
	out.Printf("SearchTime:  %d ms\n", ts.Time.Milliseconds())
	out.Printf("MaxDepth:    %d\n", ts.Depth)
	out.Printf("Date:        %s\n", time.Now().Local())
	out.Printf("No of tests: %d\n", len(ts.Tests))
	out.Println()

	// execute all tests and store results in the
	// test instance
	for i, t := range ts.Tests {
		out.Printf("Test %d of %d\nTest: %s -- Target Result %s\n", i+1, len(ts.Tests), t.line, t.targetMoves.StringUci())
		startTime2 := time.Now()
		runSingleTest(s, sl, t)
		elapsedTime := time.Since(startTime2)
		t.nps = util.Nps(s.NodesVisited(), s.LastSearchResult().SearchTime)
		out.Printf("Test finished in %d ms with result %s (%s) - nps: %d\n\n",
			elapsedTime.Milliseconds(), t.rType.String(), t.actual.StringUci(), t.nps)
	}

	// sum up result for report
	tr := &SuiteResult{}
	for _, t := range ts.Tests {
		tr.Counter++
		switch t.rType {
		case NotTested:
			tr.NotTestedCounter++
		case Skipped:
			tr.SkippedCounter++
		case Failed:
			tr.FailedCounter++
		case Success:
			tr.SuccessCounter++
		}
	}
	ts.LastResult = tr

	elapsed := time.Since(startTime)

	// print report
	out.Printf("Results for Test Suite\n", ts.FilePath)
	out.Printf("------------------------------------------------------------------------------------------------------------------------------------\n")
	out.Printf("EPD File:   %s\n", ts.FilePath)
	out.Printf("SearchTime: %d ms\n", ts.Time.Milliseconds())
	out.Printf("MaxDepth:   %d\n", ts.Depth)
	out.Printf("Date:       %s\n", time.Now().Local())
	out.Printf("====================================================================================================================================\n")
	out.Printf(" %-4s | %-10s | %-8s | %-8s | %-15s | %s | %s\n", " Nr.", "Result", "Move", "Value", "Expected Result", "Fen", "Id")
	out.Printf("====================================================================================================================================\n")
	for i, t := range ts.Tests {
		if t.tType == DM {
			out.Printf(" %-4d | %-10s | %-8s | %-8s | %s%-15d | %s | %s\n",
				i+1, t.rType.String(), t.actual.StringUci(), t.value.String(), "dm ", t.mateDepth, t.fen, t.id)
		} else {
			out.Printf(" %-4d | %-10s | %-8s | %-8s | %s %-15s | %s | %s\n",
				i+1, t.rType.String(), t.actual.StringUci(), t.value.String(), t.tType.String(), t.targetMoves.StringUci(), t.fen, t.id)
		}
	}
	out.Printf("====================================================================================================================================\n")
	out.Printf("Summary:\n")
	out.Printf("EPD File:   %s\n", ts.FilePath)
	out.Printf("SearchTime: %d ms\n", ts.Time.Milliseconds())
	out.Printf("MaxDepth:   %d\n", ts.Depth)
	out.Printf("Date:       %s\n", time.Now().Local())
	out.Printf("Successful: %-3d (%d %%)\n", tr.SuccessCounter, 100*tr.SuccessCounter/tr.Counter)
	out.Printf("Failed:     %-3d (%d %%)\n", tr.FailedCounter, 100*tr.FailedCounter/tr.Counter)
	out.Printf("Skipped:    %-3d (%d %%)\n", tr.SkippedCounter, 100*tr.SkippedCounter/tr.Counter)
	out.Printf("Not tested: %-3d (%d %%)\n", tr.NotTestedCounter, 100*tr.NotTestedCounter/tr.Counter)
	out.Printf("Test time: %s\n", elapsed)
	out.Printf("Configuration: %s\n", config.Settings.String())
}

// determines which test type the test is and call the appropriate
// test function
func runSingleTest(s *search.Search, sl *search.Limits, t *Test) {
	// reset search and search limits
	s.NewGame()
	sl.Mate = 0
	// create position
	p, _ := position.NewPositionFen(t.fen)
	switch t.tType {
	case DM:
		directMateTest(s, sl, p, t)
	case BM:
		bestMoveTest(s, sl, p, t)
	case AM:
		avoidMoveMateTest(s, sl, p, t)
	default:
		log.Warningf("Unknown Test type: %d", t.tType)
	}
}

func directMateTest(s *search.Search, sl *search.Limits, p *position.Position, t *Test) {
	// prepare
	sl.Mate = t.mateDepth
	// start search
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	// check and store result
	if s.LastSearchResult().BestValue.String() == fmt.Sprintf("mate %d", t.mateDepth) {
		// success
		log.Infof("TestSet: id = '%s' SUCCESS", t.id)
		t.actual = s.LastSearchResult().BestMove
		t.value = s.LastSearchResult().BestValue
		t.rType = Success
		return
	}
	// Failed
	log.Infof("TestSet: id = '%s' FAILED", t.id)
	t.actual = s.LastSearchResult().BestMove
	t.value = s.LastSearchResult().BestValue
	t.rType = Failed
}

func bestMoveTest(s *search.Search, sl *search.Limits, p *position.Position, t *Test) {
	// start search
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	// check and store result
	for _, m := range t.targetMoves {
		if m == s.LastSearchResult().BestMove {
			// success
			log.Infof("TestSet: id = '%s' SUCCESS", t.id)
			t.actual = s.LastSearchResult().BestMove
			t.value = s.LastSearchResult().BestValue
			t.rType = Success
			return
		} else {
			continue
		}
	}
	// Failed
	log.Infof("TestSet: id = '%s' FAILED", t.id)
	t.actual = s.LastSearchResult().BestMove
	t.value = s.LastSearchResult().BestValue
	t.rType = Failed
}

func avoidMoveMateTest(s *search.Search, sl *search.Limits, p *position.Position, t *Test) {
	// start search
	s.StartSearch(*p, *sl)
	s.WaitWhileSearching()
	// check and store result
	for _, m := range t.targetMoves {
		if m == s.LastSearchResult().BestMove {
			// success
			log.Infof("TestSet: id = '%s' FAILED", t.id)
			t.actual = s.LastSearchResult().BestMove
			t.value = s.LastSearchResult().BestValue
			t.rType = Failed
			return
		} else {
			continue
		}
	}
	// Failed
	log.Infof("TestSet: id = '%s' SUCCESS", t.id)
	t.actual = s.LastSearchResult().BestMove
	t.value = s.LastSearchResult().BestValue
	t.rType = Success
}

var leadingComments = regexp.MustCompile("^\\s*#.*$")
var trailingComments = regexp.MustCompile("^(.*)#([^;]*)$")

// takes a line with an EPD and creates a test from it
func getTest(line string) *Test {
	// cleanup the line string
	line = strings.TrimSpace(line)
	line = leadingComments.ReplaceAllString(line, "")
	line = trailingComments.ReplaceAllString(line, "")

	if len(line) == 0 {
		return nil
	}

	// Find a EPD line
	var epdRegex = regexp.MustCompile("^\\s*(.*?) (bm|dm|am) (.*?);(.* id \"(.*?)\";)?.*$")
	if !epdRegex.MatchString(line) {
		log.Warningf("No EPD found in %s", line)
		return nil
	}

	// dissect the string
	parts := epdRegex.FindStringSubmatch(line)

	// part 1 is fen - test against position
	p, err := position.NewPositionFen(parts[1])
	if err != nil {
		log.Warningf("fen part of EPD is invalid. %s", parts[1])
		return nil
	}
	fen := parts[1]

	// part 2 opcode
	var ttype testType
	switch parts[2] {
	case "dm":
		ttype = DM
	case "bm":
		ttype = BM
	case "am":
		ttype = AM
	default:
		log.Warningf("Opcode from EPD is invalid or not implemented %s", parts[2])
		return nil
	}

	// part 3 target result
	resultMoves := moveslice.NewMoveSlice(4)
	dmDepth := 0
	if ttype == BM || ttype == AM {
		result := parts[3]
		strings.ReplaceAll(result, "!", "")
		strings.ReplaceAll(result, "?", "")

		// check if results are even valid on the position
		// and store the moves into the test
		mg := movegen.NewMoveGen()
		results := strings.Split(result, " ")
		for _, r := range results {
			r = strings.TrimSpace(r)
			m := mg.GetMoveFromSan(p, r)
			if m != MoveNone {
				resultMoves.PushBack(m)
			}
		}
		if resultMoves.Len() == 0 {
			log.Warningf("Result moves from EPD is/are invalid on this position %s", parts[3])
			return nil
		}
	} else if ttype == DM {
		dmDepth, err = strconv.Atoi(parts[3])
		if err != nil {
			log.Warningf("Direct mate depth from EPD is invalid %s", parts[3])
			return nil
		}
	}

	// create the test instance
	test := &Test{
		id:          parts[5],
		fen:         fen,
		tType:       ttype,
		targetMoves: *resultMoves,
		mateDepth:   dmDepth,
		target:      0,
		actual:      0,
		value:       0,
		rType:       0,
		line:        line,
	}

	return test
}

// reads a file a returns all lines as a slice of strings
func getTestLines(filePath string) (*[]string, error) {
	// get path to file
	if !filepath.IsAbs(filePath) {
		wd, _ := os.Getwd()
		filePath = wd + "/" + filePath
	}
	filePath = filepath.Clean(filePath)

	// check file path
	if _, err := os.Stat(filePath); err != nil {
		log.Errorf("File \"%s\" does not exist\n", filePath)
		return nil, err
	}

	// read tests from file
	log.Infof("Reading test suite tests from file: %s\n", filePath)
	startReading := time.Now()
	lines, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	elapsedReading := time.Since(startReading)
	log.Infof("Finished reading %d lines from file in: %d ms\n", len(*lines), elapsedReading.Milliseconds())
	return lines, nil
}

// reads a complete file into a slice of strings
func readFile(filePath string) (*[]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf("File \"%s\" could not be read; %s\n", filePath, err)
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Errorf("File \"%s\" could not be closed: %s\n", filePath, err)
		}
	}()
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	err = s.Err()
	if err != nil {
		log.Errorf("Error while reading file \"%s\": %s\n", filePath, err)
		return nil, err
	}
	return &lines, nil
}

func (rt *resultType) String() string {
	switch *rt {
	case NotTested:
		return "Not tested"
	case Skipped:
		return "Skipped"
	case Failed:
		return "Failed"
	case Success:
		return "Success"
	default:
		return "N/A"
	}
}

func (tt *testType) String() string {
	switch *tt {
	case BM:
		return "bm"
	case AM:
		return "am"
	case DM:
		return "dm"
	default:
		return "N/A"
	}
}

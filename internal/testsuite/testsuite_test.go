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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
)

var logTest *logging.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	config.Setup()
	log = myLogging.GetLog()
	logTest = myLogging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func TestGetTest(t *testing.T) {
	line := "2b4k/8/8/8/8/3N3N/P4p2/1K6 w - - bm Nhxf2 Ndxf2; id \"FRANKY-1 #7\";"
	test := getTest(line)
	assert.NotNil(t, test)
	assert.EqualValues(t, "2b4k/8/8/8/8/3N3N/P4p2/1K6 w - -", test.fen)
	assert.EqualValues(t, "h3f2 d3f2", test.targetMoves.StringUci())
	assert.EqualValues(t, "FRANKY-1 #7", test.id)
	assert.EqualValues(t, BM, test.tType)

	line = "6k1/P7/8/8/8/8/8/3K4 w - - bm a8=Q; id \"FRANKY-1 #4\";"
	test = getTest(line)
	assert.NotNil(t, test)
	assert.EqualValues(t, "6k1/P7/8/8/8/8/8/3K4 w - -", test.fen)
	assert.EqualValues(t, "a7a8Q", test.targetMoves.StringUci())
	assert.EqualValues(t, "FRANKY-1 #4", test.id)
	assert.EqualValues(t, BM, test.tType)

	// invalid epds
	line = "6k1/P7/8/9/8/8/8/3K4 w - - bm a8=Q; id \"FRANKY-1 #4\";"
	test = getTest(line)
	assert.Nil(t, test)
	line = "6k1/P7/8/8/8/8/8/3K4 w - - aa a8=Q; id \"FRANKY-1 #4\";"
	test = getTest(line)
	assert.Nil(t, test)
	line = "2b4k/8/8/8/8/3N3N/P4p2/1K6 w - - bm Nhxf2 Naxf2; id \"FRANKY-1 #7\";"
	test = getTest(line)
	assert.NotNil(t, test) // ok as only one result move is invalid
	line = "2b4k/8/8/8/8/3N3N/P4p2/1K6 w - - bm Nbxf2 Naxf2; id \"FRANKY-1 #7\";"
	test = getTest(line)
	assert.Nil(t, test)
}

func TestNewTestSuite(t *testing.T) {
	ts, err := NewTestSuite("test/testdata/testsets/franky_tests.epd", 2*time.Second, 0)
	assert.NotNil(t, ts)
	assert.Nil(t, err)
	assert.EqualValues(t, 13, len(ts.Tests))
}

// Summary:
// EPD File:   test/testdata/testsets/franky_tests.epd
// SearchTime: 3.000 ms
// MaxDepth:   0
// Date:       2020-05-01 16:33:20.1755465 +0200 CEST
// Successful: 13  (100 %)
// Failed:     0   (0 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 37.9227221s
// Configuration: Search Config:
// 0 : UseBook                bool   = false
// 1 : BookPath               string = ./assets/books
// 2 : BookFile               string = book.txt
// 3 : BookFormat             string = Simple
// 4 : UsePonder              bool   = true
// 5 : UseQuiescence          bool   = true
// 6 : UseQSStandpat          bool   = true
// 7 : UseSEE                 bool   = true
// 8 : UsePromNonQuiet        bool   = true
// 9 : UseAspiration          bool   = false
// 10: UseMTDf                bool   = false
// 11: UsePVS                 bool   = true
// 12: UseIID                 bool   = true
// 13: UseKiller              bool   = true
// 14: UseHistoryCounter      bool   = true
// 15: UseCounterMoves        bool   = true
// 16: IIDDepth               int    = 6
// 17: IIDReduction           int    = 2
// 18: UseTT                  bool   = true
// 19: TTSize                 int    = 256
// 20: UseTTMove              bool   = true
// 21: UseTTValue             bool   = true
// 22: UseQSTT                bool   = true
// 23: UseEvalTT              bool   = false
// 24: UseMDP                 bool   = true
// 25: UseRFP                 bool   = true
// 26: UseNullMove            bool   = true
// 27: NmpDepth               int    = 3
// 28: NmpReduction           int    = 2
// 29: UseExt                 bool   = true
// 30: UseExtAddDepth         bool   = true
// 31: UseCheckExt            bool   = true
// 32: UseThreatExt           bool   = false
// 33: UseFP                  bool   = true
// 34: UseLmp                 bool   = true
// 35: UseLmr                 bool   = true
// 36: LmrDepth               int    = 3
// 37: LmrMovesSearched       int    = 3
func TestRunTestSuiteTest(t *testing.T) {
	config.Settings.Search.UseRFP = true
	config.Settings.Search.UseFP = true
	ts, _ := NewTestSuite("test/testdata/testsets/franky_tests.epd", 3*time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 10.000 ms
// MaxDepth:   0
// Date:       2020-05-01 16:35:56.2143076 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result |Fen | Id
// ====================================================================================================================================
// 1    | Failed     | e1e6     | cp -222  | bm e1f1 | 8/8/p1p5/1p5p/1P5p/8/PPP2K1p/4R1rk w - - | zugzwang.001
// 2    | Success    | g5h6     | cp 42    | bm g5h6 | 1q1k4/2Rr4/8/2Q3K1/8/8/8/8 w - - | zugzwang.002
// 3    | Failed     | f7e7     | cp 22    | bm g4g5 | 7k/5K2/5P1p/3p4/6P1/3p4/8/8 w - - | zugzwang.003
// 4    | Failed     | g7h6     | cp 5     | bm h3h4 | 8/6B1/p5p1/Pp4kp/1P5r/5P1Q/4q1PK/8 w - - | zugzwang.004
// 5    | Failed     | d6d8     | cp -61   | bm f4d5 | 8/8/1p1r1k2/p1pPN1p1/P3KnP1/1P6/8/3R4 b - - | zugzwang.005
// ====================================================================================================================================
// Summary:
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 10.000 ms
// MaxDepth:   0
// Date:       2020-05-01 16:35:56.2143076 +0200 CEST
// Successful: 1   (20 %)
// Failed:     4   (80 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 50.0610856s
// Configuration: Search Config:
// 0 : UseBook                bool   = false
// 1 : BookPath               string = ./assets/books
// 2 : BookFile               string = book.txt
// 3 : BookFormat             string = Simple
// 4 : UsePonder              bool   = true
// 5 : UseQuiescence          bool   = true
// 6 : UseQSStandpat          bool   = true
// 7 : UseSEE                 bool   = true
// 8 : UsePromNonQuiet        bool   = true
// 9 : UseAspiration          bool   = false
// 10: UseMTDf                bool   = false
// 11: UsePVS                 bool   = true
// 12: UseIID                 bool   = true
// 13: UseKiller              bool   = true
// 14: UseHistoryCounter      bool   = true
// 15: UseCounterMoves        bool   = true
// 16: IIDDepth               int    = 6
// 17: IIDReduction           int    = 2
// 18: UseTT                  bool   = true
// 19: TTSize                 int    = 256
// 20: UseTTMove              bool   = true
// 21: UseTTValue             bool   = true
// 22: UseQSTT                bool   = true
// 23: UseEvalTT              bool   = false
// 24: UseMDP                 bool   = true
// 25: UseRFP                 bool   = true
// 26: UseNullMove            bool   = true
// 27: NmpDepth               int    = 3
// 28: NmpReduction           int    = 2
// 29: UseExt                 bool   = true
// 30: UseExtAddDepth         bool   = true
// 31: UseCheckExt            bool   = true
// 32: UseThreatExt           bool   = false
// 33: UseFP                  bool   = true
// 34: UseLmp                 bool   = true
// 35: UseLmr                 bool   = true
// 36: LmrDepth               int    = 3
// 37: LmrMovesSearched       int    = 3
func TestZugzwangTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseRazoring = true
	ts, _ := NewTestSuite("test/testdata/testsets/nullMoveZugZwangTest.epd", 10*time.Second, 0)
	ts.RunTests()
}

// Summary:
// config.Settings.Search.UseRazoring = true
// EPD File:   test/testdata/testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-05-03 11:10:54.9756125 +0200 CEST
// Successful: 15  (75 %)
// Failed:     5   (25 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 3m44.6555458s
func TestMateTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseRazoring = true
	ts, _ := NewTestSuite("test/testdata/testsets/mate_test_suite.epd", 15*time.Second, 0)
	ts.RunTests()
}

// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-01 17:30:59.8765561 +0200 CEST
// Successful: 191 (95 %)
// Failed:     10  (4 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 15m18.0870814s
func TestWACTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ts, _ := NewTestSuite("test/testdata/testsets/wac.epd", 5*time.Second, 0)
	ts.RunTests()
}

// Summary:
// EPD File:   test/testdata/testsets/crafty_test.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-24 18:09:23.5572576 +0200 CEST
// Successful: 169 (48 %)
// Failed:     176 (51 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 28m49.734886s
//
// RootAlpha
// Summary:
// EPD File:   test/testdata/testsets/crafty_test.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-01 18:01:13.2116238 +0200 CEST
// Successful: 161 (46 %)
// Failed:     184 (53 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 28m42.2265061s
func TestCraftyTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ts, _ := NewTestSuite("test/testdata/testsets/crafty_test.epd", 5*time.Second, 0)
	ts.RunTests()
}

// Summary:
// EPD File:   test/testdata/testsets/ecm98.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-01 20:26:51.2385188 +0200 CEST
// Successful: 491 (63 %)
// Failed:     278 (36 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1h4m12.0971456s
func TestECMTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	config.Settings.Search.UseRazoring = false
	// Summary:
	// EPD File:   test/testdata/testsets/ecm98.epd
	// SearchTime: 5.000 ms
	// Successful: 497 (64 %)
	// Failed:     272 (35 %)
	ts, _ := NewTestSuite("test/testdata/testsets/ecm98.epd", 5*time.Second, 0)
	ts.RunTests()

	config.Settings.Search.UseRazoring = true
	// Summary:
	// EPD File:   test/testdata/testsets/ecm98.epd
	// SearchTime: 5.000 ms
	// Date:       2020-05-03 03:21:06.0137971 +0200 CEST
	// Successful: 505 (65 %)
	// Failed:     264 (34 %)
	ts2, _ := NewTestSuite("test/testdata/testsets/ecm98.epd", 5*time.Second, 0)
	ts2.RunTests()
}

func TestStressTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	files, err := ioutil.ReadDir("test/testdata/testsets/")
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".epd" {
			list = append(list, f.Name())
		}
	}
	for {
		for _, t := range list {
			ts, _ := NewTestSuite("test/testdata/testsets/"+t, 5*time.Second, 0)
			ts.RunTests()
		}
	}
}

func TestFeatureTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup tests
	searchTime := 200 * time.Millisecond
	searchDepth := 0

	// Feature Settings
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
	config.Settings.Search.UseLmr = true
	config.Settings.Search.LmrDepth = 3
	config.Settings.Search.LmrMovesSearched = 3
	config.Settings.Search.UseLmp = true

	config.Settings.Eval.Tempo = 50
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

	folder := "test/testdata/featuretests/"

	out.Println(FeatureTests(folder, searchTime, searchDepth))
}

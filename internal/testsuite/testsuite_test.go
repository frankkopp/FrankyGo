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
//
// MTDf
// Summary:
// EPD File:   test/testdata/testsets/franky_tests.epd
// SearchTime: 3.000 ms
// MaxDepth:   0
// Date:       2020-05-02 02:59:53.803276 +0200 CEST
// Successful: 13  (100 %)
// Failed:     0   (0 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 35.4837035s
// Configuration: Search Config:
func TestRunTestSuiteTest(t *testing.T) {
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
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
//
// MTDf (TODO - check)
// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 10.000 ms
// MaxDepth:   0
// Date:       2020-05-02 03:01:36.8535347 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result |Fen | Id
// ====================================================================================================================================
// 1    | Failed     | e1e6     | cp -245  | bm e1f1 | 8/8/p1p5/1p5p/1P5p/8/PPP2K1p/4R1rk w - - | zugzwang.001
// 2    | Failed     | c7d7     | cp 0     | bm g5h6 | 1q1k4/2Rr4/8/2Q3K1/8/8/8/8 w - - | zugzwang.002
// 3    | Failed     | f7e7     | cp 11    | bm g4g5 | 7k/5K2/5P1p/3p4/6P1/3p4/8/8 w - - | zugzwang.003
// 4    | Failed     | g7h6     | cp 0     | bm h3h4 | 8/6B1/p5p1/Pp4kp/1P5r/5P1Q/4q1PK/8 w - - | zugzwang.004
// 5    | Failed     | d6d8     | cp -51   | bm f4d5 | 8/8/1p1r1k2/p1pPN1p1/P3KnP1/1P6/8/3R4 b - - | zugzwang.005
// ====================================================================================================================================
// Summary:
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 10.000 ms
// MaxDepth:   0
// Date:       2020-05-02 03:01:36.8535347 +0200 CEST
// Successful: 0   (0 %)
// Failed:     5   (100 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 50.0528239s
func TestZugzwangTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
	ts, _ := NewTestSuite("test/testdata/testsets/nullMoveZugZwangTest.epd", 10*time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   test/testdata/testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-05-01 16:40:53.049496 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result | Fen| Id
// ====================================================================================================================================
// 1    | Success    | f1e2     | mate 3   | dm 3               | 4r1b1/1p4B1/pN2pR2/RB2k3/1P2N2p/2p3b1/n2P1p1r/5K1n w - - |
// 2    | Success    | a6b7     | mate 3   | dm 3               | b7/p1BR2pK/B1p5/pNk3p1/Pp4Pp/1P3n2/4R2r/7n w - - |
// 3    | Success    | b3b6     | mate 3   | dm 3               | 6K1/n1P2N1p/6pr/b1pp3b/n2Bp1k1/1R2R1Pp/3p1P2/2qN1B2 w - - |
// 4    | Success    | e7e8B    | mate 3   | dm 3               | 8/2P1P1P1/3PkP2/8/4K3/8/8/8 w - - |
// 5    | Success    | d1h5     | mate 4   | dm 4               | r4rk1/pppqbp1p/3pp1p1/8/4P3/1P1P3R/PBP2PPP/R2Q2K1 w - - |
// 6    | Success    | h5h7     | mate 4   | dm 4               | r4qrk/ppp1b1pp/3p1p2/4pPPQ/4P2P/3PB3/PPP5/1K4RR w - - |
// 7    | Success    | e8e3     | mate 4   | dm 4               | r3r3/p1p2p1k/3p2pp/2p5/2P2n2/2N2B2/PPR1PP1q/3RQK2 b - - |
// 8    | Failed     | d4a4     | mate 5   | dm 4               | 3r4/1p1r4/1Pp5/3p4/p2R4/K1NN4/1P6/kqBB3R w - - |
// 9    | Success    | d4g1     | mate 4   | dm 4               | n7/3p1p2/NpkNp1p1/1p2P3/3Q4/6B1/b7/4K3 w - - |
// 10   | Failed     | d2d4     | mate 5   | dm 4               | 1b1R2B1/p1n1p3/p3P2K/N1k5/2N2P2/5P2/2PP4/R7 w - - |
// 11   | Success    | f1a6     | mate 4   | dm 4               | K6Q/1p6/pPq4P/P2p2P1/4pP1N/7k/n5R1/1n2BB2 w - - |
// 12   | Success    | g4h6     | mate 5   | dm 5               | 1r4k1/1b2K1pp/7b/2pp3P/6NB/2Q2pp1/4p3/5r2 w - - |
// 13   | Success    | d1g1     | mate 5   | dm 5               | r2r4/1p1R3p/5pk1/b1B1Pp2/p4P2/P7/1P5P/1K1R4 w - - |
// 14   | Success    | f2h1     | mate 5   | dm 5               | 5rk1/pp4p1/8/3N3p/2P4P/1P4K1/P2r1n2/R3R3 b - - |
// 15   | Failed     | d4f6     | mate 7   | dm 5               | 6b1/4Kpk1/5r2/8/3B2P1/7R/8/8 w - - |
// 16   | Failed     | c1b1     | mate 14  | dm 5               | 8/8/8/p7/8/8/R6p/2K2Rbk w - - |
// 17   | Failed     | h6f4     | cp -1480 | dm 5               | b7/8/7B/7p/8/2p3r1/2P1P1pp/4K1kq w - - |
// 18   | Success    | f8e8     | mate 5   | dm 5               | 5R2/6r1/3P4/1BBk4/8/3N4/8/K7 w - - |
// 19   | Failed     | c3d4     | cp 478   | dm 5               | 1r1r2k1/p4ppp/1qp5/4Pb2/3b1P2/1PP2N2/P2BQ1PP/2KR3R w - - |
// 20   | Success    | e3c5     | mate 6   | dm 6               | 4k3/8/4K3/8/4N3/4B3/3P1P2/8 w - - |
// ====================================================================================================================================
// Summary:
// EPD File:   test/testdata/testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-05-01 16:40:53.049496 +0200 CEST
// Successful: 14  (70 %)
// Failed:     6   (30 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 4m1.077031s
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
func TestMateTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
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
//
// MTDf
// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-02 03:23:18.4027014 +0200 CEST
// Successful: 192 (95 %)
// Failed:     9   (4 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 15m17.0509323s
func TestWACTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
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
//
// MTDf
// Summary:
// EPD File:   test/testdata/testsets/crafty_test.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-02 03:52:52.0807002 +0200 CEST
// Successful: 155 (44 %)
// Failed:     190 (55 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 28m41.2219044s
// Configuration: Search Config:
func TestCraftyTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
	ts, _ := NewTestSuite("test/testdata/testsets/crafty_test.epd", 5*time.Second, 0)
	ts.RunTests()
}

// Summary:
// EPD File:   test/testdata/testsets/ecm98.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-28 17:43:19.2814557 +0200 CEST
// Successful: 424 (55 %)
// Failed:     345 (44 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1h4m15.4163679s
//
// Summary: asp true
// EPD File:   test/testdata/testsets/ecm98.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-02 03:01:41.256517 +0200 CEST
// Successful: 485 (63 %)
// Failed:     284 (36 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1h4m11.360931673s
//
// MTDf
// Summary:
// EPD File:   test/testdata/testsets/ecm98.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-05-02 12:01:10.5886708 +0200 CEST
// Successful: 566 (73 %)
// Failed:     203 (26 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1h4m11.4441303s
func TestECMTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UsePVS = true
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = false
	// Summary:
	// EPD File:   test/testdata/testsets/ecm98.epd
	// SearchTime: 5.000 ms
	// MaxDepth:   0
	// Date:       2020-05-02 13:51:01.848948 +0200 CEST
	// Successful: 492 (63 %)
	// Failed:     277 (36 %)
	// Skipped:    0   (0 %)
	// Not tested: 0   (0 %)
	// Test time: 1h4m10.888877218s
	ts, _ := NewTestSuite("test/testdata/testsets/ecm98.epd", 5*time.Second, 0)
	ts.RunTests()
	config.Settings.Search.UsePVS = true
	config.Settings.Search.UseAspiration = true
	config.Settings.Search.UseMTDf = false
	// Summary:
	// EPD File:   test/testdata/testsets/ecm98.epd
	// SearchTime: 5.000 ms
	// MaxDepth:   0
	// Date:       2020-05-02 14:55:12.936006 +0200 CEST
	// Successful: 492 (63 %)
	// Failed:     277 (36 %)
	// Skipped:    0   (0 %)
	// Not tested: 0   (0 %)
	// Test time: 1h4m11.001542203s
	ts1, _ := NewTestSuite("test/testdata/testsets/ecm98.epd", 5*time.Second, 0)
	ts1.RunTests()
	config.Settings.Search.UsePVS = false
	config.Settings.Search.UseAspiration = false
	config.Settings.Search.UseMTDf = true
	// Summary:
	// EPD File:   test/testdata/testsets/ecm98.epd
	// SearchTime: 5.000 ms
	// MaxDepth:   0
	// Date:       2020-05-02 15:59:24.341412 +0200 CEST
	// Successful: 560 (72 %)
	// Failed:     209 (27 %)
	// Skipped:    0   (0 %)
	// Not tested: 0   (0 %)
	// Test time: 1h4m11.311020444s
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

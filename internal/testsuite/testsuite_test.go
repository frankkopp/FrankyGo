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
// Date:       2020-04-22 23:33:11.4600101 +0200 CEST
// Successful: 13  (100 %)
// Failed:     0   (0 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 39.1836873s
// Configuration: Search Config:
// 0 : UseBook                bool   = false
// 1 : BookPath               string = ./assets/books
// 2 : BookFile               string = book.txt
// 3 : BookFormat             string = Simple
// 4 : UsePonder              bool   = true
// 5 : UseQuiescence          bool   = true
// 6 : UseQSStandpat          bool   = true
// 7 : UseSEE                 bool   = true
// 8 : UsePVS                 bool   = true
// 9 : UseKiller              bool   = true
// 10: UseIID                 bool   = true
// 11: IIDDepth               int    = 6
// 12: IIDReduction           int    = 2
// 13: UseTT                  bool   = true
// 14: TTSize                 int    = 256
// 15: UseTTMove              bool   = true
// 16: UseTTValue             bool   = true
// 17: UseQSTT                bool   = true
// 18: UseEvalTT              bool   = false
// 19: UseMDP                 bool   = true
// 20: UseRFP                 bool   = true
// 21: UseNullMove            bool   = true
// 22: NmpDepth               int    = 3
// 23: NmpReduction           int    = 2
// 24: UseExt                 bool   = true
// 25: UseCheckExt            bool   = true
// 26: UseThreatExt           bool   = false
// 27: UseFP                  bool   = true
// 28: UseLmp                 bool   = true
// 29: UseLmr                 bool   = true
// 30: LmrDepth               int    = 3
// 31: LmrMovesSearched       int    = 3
func TestRunTestSuiteTest(t *testing.T) {
	config.Settings.Search.UseRFP = true
	config.Settings.Search.UseFP = true
	ts, _ := NewTestSuite("test/testdata/testsets/franky_tests.epd", 3*time.Second, 0)
	ts.RunTests()
	// assert.GreaterOrEqual(t, ts.LastResult.SuccessCounter, 12)
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-22 23:34:42.3348595 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result |Fen | Id
// ====================================================================================================================================
// 1    | Success    | e1f1     | cp -205  | bm e1f1 | 8/8/p1p5/1p5p/1P5p/8/PPP2K1p/4R1rk w - - | zugzwang.001
// 2    | Success    | g5h6     | cp 387   | bm g5h6 | 1q1k4/2Rr4/8/2Q3K1/8/8/8/8 w - - | zugzwang.002
// 3    | Failed     | f7e7     | cp 22    | bm g4g5 | 7k/5K2/5P1p/3p4/6P1/3p4/8/8 w - - | zugzwang.003
// 4    | Failed     | g7h6     | cp 15    | bm h3h4 | 8/6B1/p5p1/Pp4kp/1P5r/5P1Q/4q1PK/8 w - - | zugzwang.004
// 5    | Failed     | d6d8     | cp -42   | bm f4d5 | 8/8/1p1r1k2/p1pPN1p1/P3KnP1/1P6/8/3R4 b - - | zugzwang.005
// ====================================================================================================================================
// Summary:
// EPD File:   test/testdata/testsets/nullMoveZugZwangTest.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-22 23:34:42.3348595 +0200 CEST
// Successful: 2   (40 %)
// Failed:     3   (60 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 25.040373s
func TestZugzwangTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseRFP = true
	config.Settings.Search.UseFP = true
	ts, _ := NewTestSuite("test/testdata/testsets/nullMoveZugZwangTest.epd", 5*time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   test/testdata/testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-04-24 15:23:03.556285 +0200 CEST
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
// 10   | Failed     | d8c8     | mate 5   | dm 4               | 1b1R2B1/p1n1p3/p3P2K/N1k5/2N2P2/5P2/2PP4/R7 w - - |
// 11   | Success    | f1a6     | mate 4   | dm 4               | K6Q/1p6/pPq4P/P2p2P1/4pP1N/7k/n5R1/1n2BB2 w - - |
// 12   | Success    | g4h6     | mate 5   | dm 5               | 1r4k1/1b2K1pp/7b/2pp3P/6NB/2Q2pp1/4p3/5r2 w - - |
// 13   | Success    | d1g1     | mate 5   | dm 5               | r2r4/1p1R3p/5pk1/b1B1Pp2/p4P2/P7/1P5P/1K1R4 w - - |
// 14   | Success    | f2h1     | mate 5   | dm 5               | 5rk1/pp4p1/8/3N3p/2P4P/1P4K1/P2r1n2/R3R3 b - - |
// 15   | Success    | h3h6     | mate 5   | dm 5               | 6b1/4Kpk1/5r2/8/3B2P1/7R/8/8 w - - |
// 16   | Failed     | a2f2     | mate 6   | dm 5               | 8/8/8/p7/8/8/R6p/2K2Rbk w - - |
// 17   | Success    | h6c1     | mate 5   | dm 5               | b7/8/7B/7p/8/2p3r1/2P1P1pp/4K1kq w - - |
// 18   | Success    | f8e8     | mate 5   | dm 5               | 5R2/6r1/3P4/1BBk4/8/3N4/8/K7 w - - |
// 19   | Failed     | c3d4     | cp 476   | dm 5               | 1r1r2k1/p4ppp/1qp5/4Pb2/3b1P2/1PP2N2/P2BQ1PP/2KR3R w - - |
// 20   | Success    | e3c5     | mate 6   | dm 6               | 4k3/8/4K3/8/4N3/4B3/3P1P2/8 w - - |
// ====================================================================================================================================
// Summary:
// EPD File:   test/testdata/testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-04-24 15:23:03.556285 +0200 CEST
// Successful: 16  (80 %)
// Failed:     4   (20 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 5m0.3997113s
// Configuration: Search Config:
// 0 : UseBook                bool   = false
// 1 : BookPath               string = ./assets/books
// 2 : BookFile               string = book.txt
// 3 : BookFormat             string = Simple
// 4 : UsePonder              bool   = true
// 5 : UseQuiescence          bool   = true
// 6 : UseQSStandpat          bool   = true
// 7 : UseSEE                 bool   = true
// 8 : UsePVS                 bool   = true
// 9 : UseKiller              bool   = true
// 10: UseIID                 bool   = true
// 11: IIDDepth               int    = 6
// 12: IIDReduction           int    = 2
// 13: UseTT                  bool   = true
// 14: TTSize                 int    = 256
// 15: UseTTMove              bool   = true
// 16: UseTTValue             bool   = true
// 17: UseQSTT                bool   = true
// 18: UseEvalTT              bool   = false
// 19: UseMDP                 bool   = true
// 20: UseRFP                 bool   = true
// 21: UseNullMove            bool   = true
// 22: NmpDepth               int    = 3
// 23: NmpReduction           int    = 2
// 24: UseExt                 bool   = true
// 25: UseCheckExt            bool   = true
// 26: UseThreatExt           bool   = false
// 27: UseFP                  bool   = true
// 28: UseLmp                 bool   = true
// 29: UseLmr                 bool   = true
// 30: LmrDepth               int    = 3
// 31: LmrMovesSearched       int    = 3
//
// Evaluation Config:
// 0 : UsePawnCache           bool   = false
// 1 : PawnCacheSize          int    = 64
// 2 : UseLazyEval            bool   = true
// 3 : LazyEvalThreshold      int    = 700
// 4 : Tempo                  int    = 30
// 5 : UseAttacksInEval       bool   = false
// 6 : UseMobility            bool   = false
// 7 : MobilityBonus          int    = 5
// 8 : UseAdvancedPieceEval   bool   = false
// 9 : BishopPairBonus        int    = 20
// 10: MinorBehindPawnBonus   int    = 15
// 11: BishopPawnMalus        int    = 5
// 12: BishopCenterAimBonus   int    = 20
// 13: BishopBlockedMalus     int    = 40
// 14: RookOnQueenFileBonus   int    = 6
// 15: RookOnOpenFileBonus    int    = 25
// 16: RookTrappedMalus       int    = 40
// 17: KingRingAttacksBonus   int    = 10
// 18: UseKingEval            bool   = false
// 19: KingDangerMalus        int    = 50
// 20: KingDefenderBonus      int    = 10
func TestMateTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseRFP = true
	config.Settings.Search.UseFP = true
	ts, _ := NewTestSuite("test/testdata/testsets/mate_test_suite.epd", 15*time.Second, 0)
	ts.RunTests()
}

// Summary:
// EPD File:   testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-09 19:14:13.1463946 +0200 CEST
// Successful: 170 (84 %)
// Failed:     31  (15 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 16m47.1436222s

// +LMP
// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:    2020-04-21 12:21:06.2003336 +0200 CEST
// Successful: 182 (90 %)
// Failed:    19  (9 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 16m47.7678427s

// -LMP
// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:    2020-04-21 12:38:51.1514884 +0200 CEST
// Successful: 183 (91 %)
// Failed:    18  (8 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 16m47.683437s
//
// RFP/FP
// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-23 00:10:20.008491 +0200 CEST
// Successful: 183 (91 %)
// Failed:     18  (8 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 16m47.799935s
//
// HistoryCount
// Summary:
// EPD File:   test/testdata/testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-25 18:38:18.093394 +0200 CEST
// Successful: 187 (93 %)
// Failed:     14  (6 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 16m47.9795186s
func TestWACTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseHistoryCounter = true
	config.Settings.Search.UseCounterMoves = true
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
func TestCraftyTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseHistoryCounter = true
	config.Settings.Search.UseCounterMoves = true
	ts, _ := NewTestSuite("test/testdata/testsets/crafty_test.epd", 5*time.Second, 0)
	ts.RunTests()
}

// All but FP
// Summary:
// EPD File:   test/testdata/testsets/ecm98.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-22 20:38:43.7638852 +0200 CEST
// Successful: 460 (59 %)
// Failed:     309 (40 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1h4m14.3847089s
func TestECMTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.Settings.Search.UseRFP = true
	config.Settings.Search.UseFP = true
	ts, _ := NewTestSuite("test/testdata/testsets/ecm98.epd", 5*time.Second, 0)
	ts.RunTests()
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

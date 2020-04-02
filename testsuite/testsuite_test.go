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

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

var logTest *logging.Logger

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
	ts, err := NewTestSuite("testsets/franky_tests.epd", 2 * time.Second, 0)
	assert.NotNil(t, ts)
	assert.Nil(t, err)
	assert.EqualValues(t, 13, len(ts.Tests))
	for _,tt := range ts.Tests {
		out.Println(tt)
	}
}

// Test the testsuite itself
// Results for Test Suite testsets/franky_tests.epd
// ====================================================================================================================================
//  Nr. | rType      | Move     | value    | Expected rType  | fen | ID
// ====================================================================================================================================
// 13   | Success    | d3e3     | mate 4   | dm 4               | 8/8/8/8/8/3K4/R7/5k2 w - - | FRANKY-1 #1
// 13   | Success    | d3e3     | mate 4   | bm d3e3            | 8/8/8/8/8/3K4/R7/5k2 w - - | FRANKY-1 #2
// 13   | Success    | e4d3     | mate 5   | dm 5               | 8/8/8/8/4K3/8/R7/4k3 w - - | FRANKY-1 #3
// 13   | Success    | a7a8Q    | cp 963   | bm a7a8Q           | 6k1/P7/8/8/8/8/8/3K4 w - - | FRANKY-1 #4
// 13   | Success    | a4a5     | cp 961   | bm a4a5            | 5k2/8/8/8/P7/3p4/3K4/8 w - - | FRANKY-1 #5
// 13   | Success    | e5f3     | cp 801   | bm e5f3            | 7k/8/3p4/4N3/8/5p2/P7/1K2N3 w - - | FRANKY-1 #6
// 13   | Success    | h3f2     | cp 447   | bm h3f2 d3f2       | 2b4k/8/8/8/8/3N3N/P4p2/1K6 w - - | FRANKY-1 #7
// 13   | Success    | c3d4     | cp 18    | bm c3d4            | 8/3r1pk1/p1R2p2/1p5p/r2p4/PRP1K1P1/5P1P/8 w - - | Franky-1 #8
// 13   | Success    | e4d3     | cp 143   | bm e4d3            | 8/3r1pk1/p1R2p2/1p5p/r2Pp3/PRP3P1/4KP1P/8 b - d3 | FRANKY-1 #9
// 13   | Success    | e1g1     | cp -205  | bm e1g1            | r1bqk2r/pp3ppp/2pb4/3pp3/4n1n1/P4N2/1PPPBPPP/RNBQK2R w KQkq - | Franky-1 #10
// 13   | Success    | f3h1     | cp 280   | bm f3h1            | 8/2r1kpp1/1p6/pB1Pp1P1/Pbp1P3/2N2b1P/1PPK1P2/7R b - - | FRANKY-1 #11
// 13   | Success    | f1e2     | mate 3   | dm 3               | 4r1b1/1p4B1/pN2pR2/RB2k3/1P2N2p/2p3b1/n2P1p1r/5K1n w - - |
// 13   | Success    | d6f8     | cp 653   | bm d6f8 d4f5 g4f4  | 3r3k/1r3p1p/p1pB1p2/8/p1qNP1Q1/P6P/1P4P1/3R3K w - - | WAC.294
// ====================================================================================================================================
// Successful: 13  (100 %)
// Failed:     0   (0 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
//
// Test time: 26.101 ms
func TestRunTestSuiteTest(t *testing.T) {
	ts, _ := NewTestSuite("testsets/franky_tests.epd", 2 * time.Second, 0)
	ts.RunTests()
}

func TestBlunderTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/franky_blunders.epd", 2 * time.Second, 0)
	ts.RunTests()
}

// 3.4.2020 2sec
// Successful: 159 (79 %)
// Failed:     42  (20 %)
func TestWACTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/wac.epd", 2 * time.Second, 0)
	ts.RunTests()
}

// 3.4.2020 2sec
// Successful: 13  (65 %)
// Failed:     7   (35 %)
// 5sec
// Successful: 17  (85 %)
// Failed:     3   (15 %)
func TestMateTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/mate_test_suite.epd", 15 * time.Second, 0)
	ts.RunTests()
}
// 15sec
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result | Fen| Id
// ====================================================================================================================================
// 1    | Success    | f1e2     | mate 3   | dm 3               | 4r1b1/1p4B1/pN2pR2/RB2k3/1P2N2p/2p3b1/n2P1p1r/5K1n w - - |
// 2    | Success    | a6b7     | mate 3   | dm 3               | b7/p1BR2pK/B1p5/pNk3p1/Pp4Pp/1P3n2/4R2r/7n w - - |
// 3    | Success    | b3b6     | mate 3   | dm 3               | 6K1/n1P2N1p/6pr/b1pp3b/n2Bp1k1/1R2R1Pp/3p1P2/2qN1B2 w - - |
// 4    | Success    | g7g8Q    | mate 4   | dm 4               | 8/2P1P1P1/3PkP2/8/4K3/8/8/8 w - - |
// 5    | Success    | d1h5     | mate 4   | dm 4               | r4rk1/pppqbp1p/3pp1p1/8/4P3/1P1P3R/PBP2PPP/R2Q2K1 w - - |
// 6    | Success    | h5h7     | mate 4   | dm 4               | r4qrk/ppp1b1pp/3p1p2/4pPPQ/4P2P/3PB3/PPP5/1K4RR w - - |
// 7    | Success    | e8e3     | mate 4   | dm 4               | r3r3/p1p2p1k/3p2pp/2p5/2P2n2/2N2B2/PPR1PP1q/3RQK2 b - - |
// 8    | Success    | d4h4     | mate 4   | dm 4               | 3r4/1p1r4/1Pp5/3p4/p2R4/K1NN4/1P6/kqBB3R w - - |
// 9    | Success    | d4g1     | mate 4   | dm 4               | n7/3p1p2/NpkNp1p1/1p2P3/3Q4/6B1/b7/4K3 w - - |
// 10   | Success    | d8c8     | mate 4   | dm 4               | 1b1R2B1/p1n1p3/p3P2K/N1k5/2N2P2/5P2/2PP4/R7 w - - |
// 11   | Success    | f1a6     | mate 4   | dm 4               | K6Q/1p6/pPq4P/P2p2P1/4pP1N/7k/n5R1/1n2BB2 w - - |
// 12   | Success    | g4h6     | mate 5   | dm 5               | 1r4k1/1b2K1pp/7b/2pp3P/6NB/2Q2pp1/4p3/5r2 w - - |
// 13   | Success    | d1g1     | mate 5   | dm 5               | r2r4/1p1R3p/5pk1/b1B1Pp2/p4P2/P7/1P5P/1K1R4 w - - |
// 14   | Success    | f2h1     | mate 5   | dm 5               | 5rk1/pp4p1/8/3N3p/2P4P/1P4K1/P2r1n2/R3R3 b - - |
// 15   | Success    | h3h6     | mate 5   | dm 5               | 6b1/4Kpk1/5r2/8/3B2P1/7R/8/8 w - - |
// 16   | Success    | a2f2     | mate 5   | dm 5               | 8/8/8/p7/8/8/R6p/2K2Rbk w - - |
// 17   | Success    | h6c1     | mate 5   | dm 5               | b7/8/7B/7p/8/2p3r1/2P1P1pp/4K1kq w - - |
// 18   | Failed     | f8e8     | cp 2409  | dm 5               | 5R2/6r1/3P4/1BBk4/8/3N4/8/K7 w - - |
// 19   | Failed     | c3d4     | cp 422   | dm 5               | 1r1r2k1/p4ppp/1qp5/4Pb2/3b1P2/1PP2N2/P2BQ1PP/2KR3R w - - |
// 20   | Success    | e3c5     | mate 6   | dm 6               | 4k3/8/4K3/8/4N3/4B3/3P1P2/8 w - - |
// ====================================================================================================================================
// Successful: 18  (90 %)
// Failed:     2   (10 %)

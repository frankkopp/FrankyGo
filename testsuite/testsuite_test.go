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

func TestWACTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/wac.epd", 2 * time.Second, 0)
	ts.RunTests()
}

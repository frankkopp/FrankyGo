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
	assert.EqualValues(t, 13, ts.LastResult.SuccessCounter)
}

func TestBlunderTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/franky_blunders.epd", 2 * time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   testsets/mate_test_suite.epd
// SearchTime: 15.000 ms
// MaxDepth:   0
// Date:       2020-04-03 10:13:27.7114253 +0200 CEST
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
// 18   | Success    | f8e8     | mate 5   | dm 5               | 5R2/6r1/3P4/1BBk4/8/3N4/8/K7 w - - |
// 19   | Failed     | c3d4     | cp 447   | dm 5               | 1r1r2k1/p4ppp/1qp5/4Pb2/3b1P2/1PP2N2/P2BQ1PP/2KR3R w - - |
// 20   | Success    | e3c5     | mate 6   | dm 6               | 4k3/8/4K3/8/4N3/4B3/3P1P2/8 w - - |
// ====================================================================================================================================
// Successful: 19  (95 %)
// Failed:     1   (5 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 300.238 ms
func _TestMateTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/mate_test_suite.epd", 15 * time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   testsets/wac.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-03 10:33:59.5739813 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result | Fen | Id
// ====================================================================================================================================
// 1    | Success    | e3g3     | cp 233   | bm e3g3            | 5rk1/1ppb3p/p1pb4/6q1/3P1p1r/2P1R2P/PP1BQ1P1/5RKN w - - | WAC.003
// 2    | Success    | h6h7     | mate 2   | bm h6h7            | r1bq2rk/pp3pbp/2p1p1pQ/7P/3P4/2PB1N2/PP3PPR/2KR4 w - - | WAC.004
// 3    | Success    | c6c4     | mate 2   | bm c6c4            | 5k2/6pp/p1qN4/1p1p4/3P4/2PKP2Q/PP3r2/3R4 b - - | WAC.005
// 4    | Success    | b6b7     | cp 653   | bm b6b7            | 7k/p7/1R5K/6r1/6p1/6P1/8/8 w - - | WAC.006
// 5    | Success    | g4e3     | cp 460   | bm g4e3            | rnbqkb1r/pppp1ppp/8/4P3/6n1/7P/PPPNPPP1/R1BQKBNR b KQkq - | WAC.007
// 6    | Success    | e7f7     | cp 811   | bm e7f7            | r4q1k/p2bR1rp/2p2Q1N/5p2/5p2/2P5/PP3PPP/R5K1 w - - | WAC.008
// 7    | Success    | h4h7     | cp 258   | bm h4h7            | 2br2k1/2q3rn/p2NppQ1/2p1P3/Pp5R/4P3/1P3PPP/3R2K1 w - - | WAC.010
// 8    | Success    | f3c6     | cp 350   | bm f3c6            | r1b1kb1r/3q1ppp/pBp1pn2/8/Np3P2/5B2/PPP3PP/R2Q1RK1 w kq - | WAC.011
// 9    | Success    | g4f3     | mate 2   | bm g4f3            | 4k1r1/2p3r1/1pR1p3/3pP2p/3P2qP/P4N2/1PQ4P/5R1K b - - | WAC.012
// 10   | Success    | f1f8     | cp 109   | bm f1f8            | 5rk1/pp4p1/2n1p2p/2Npq3/2p5/6P1/P3P1BP/R4Q1K w - - | WAC.013
// 11   | Success    | h3h7     | cp 499   | bm h3h7            | r2rb1k1/pp1q1p1p/2n1p1p1/2bp4/5P2/PP1BPR1Q/1BPN2PP/R5K1 w - - | WAC.014
// 12   | Success    | e2c3     | cp 337   | bm e2c3            | r4rk1/ppp2ppp/2n5/2bqp3/8/P2PB3/1PP1NPPP/R2Q1RK1 w - - | WAC.016
// 13   | Success    | c4e5     | cp 112   | bm c4e5            | 1k5r/pppbn1pp/4q1r1/1P3p2/2NPp3/1QP5/P4PPP/R1B1R1K1 w - - | WAC.017
// 14   | Success    | a8h8     | cp 0     | bm a8h8            | R7/P4k2/8/8/8/8/r7/6K1 w - - | WAC.018
// 15   | Success    | c5c6     | cp 85    | bm c5c6            | r1b2rk1/ppbn1ppp/4p3/1QP4q/3P4/N4N2/5PPP/R1B2RK1 w - - | WAC.019
// 16   | Success    | d7b5     | cp 804   | bm d7b5            | r2qkb1r/1ppb1ppp/p7/4p3/P1Q1P3/2P5/5PPP/R1B2KNR b kq - | WAC.020
// 17   | Success    | d2h6     | cp 390   | bm d2h6            | 5rk1/1b3p1p/pp3p2/3n1N2/1P6/P1qB1PP1/3Q3P/4R1K1 w - - | WAC.021
// 18   | Failed     | c4a2     | cp -2    | bm g5f7            | r1bqk2r/ppp1nppp/4p3/n5N1/2BPp3/P1P5/2P2PPP/R1BQK2R w KQkq - | WAC.022
// 19   | Success    | g2g4     | cp 377   | bm g2g4            | r3nrk1/2p2p1p/p1p1b1p1/2NpPq2/3R4/P1N1Q3/1PP2PPP/4R1K1 w - - | WAC.023
// 20   | Success    | g7d4     | cp 581   | bm g7d4            | 6k1/1b1nqpbp/pp4p1/5P2/1PN5/4Q3/P5PP/1B2B1K1 b - - | WAC.024
// 21   | Success    | g4h4     | cp 1024  | bm g4h4            | 3R1rk1/8/5Qpp/2p5/2P1p1q1/P3P3/1P2PK2/8 b - - | WAC.025
// 22   | Success    | d7f5     | cp 75    | bm d7f5            | 3r2k1/1p1b1pp1/pq5p/8/3NR3/2PQ3P/PP3PP1/6K1 b - - | WAC.026
// 23   | Success    | a3f8     | mate 2   | bm a3f8            | 7k/pp4np/2p3p1/3pN1q1/3P4/Q7/1r3rPP/2R2RK1 w - - | WAC.027
// 24   | Success    | b4e1     | cp 284   | bm b4e1            | 1r1r2k1/4pp1p/2p1b1p1/p3R3/RqBP4/4P3/1PQ2PPP/6K1 b - - | WAC.028
// 25   | Success    | c5c6     | cp 186   | bm c5c6            | r2q2k1/pp1rbppp/4pn2/2P5/1P3B2/6P1/P3QPBP/1R3RK1 w - - | WAC.029
// 26   | Success    | e4d6     | cp 100   | bm e4d6            | 1r3r2/4q1kp/b1pp2p1/5p2/pPn1N3/6P1/P3PPBP/2QRR1K1 w - - | WAC.030
// 27   | Success    | d1d8     | cp 40    | bm d1d8            | 6k1/p4p1p/1p3np1/2q5/4p3/4P1N1/PP3PPP/3Q2K1 w - - | WAC.032
// 28   | Success    | d4g1     | cp 149   | bm d4g1            | 7k/1b1r2p1/p6p/1p2qN2/3bP3/3Q4/P5PP/1B1R3K b - - | WAC.034
// 29   | Success    | h4h7     | mate 4   | bm h4h7            | r3r2k/2R3pp/pp1q1p2/8/3P3R/7P/PP3PP1/3Q2K1 w - - | WAC.035
// 30   | Success    | e7e1     | cp 264   | bm e7e1            | 3r4/2p1rk2/1pQq1pp1/7p/1P1P4/P4P2/6PP/R1R3K1 b - - | WAC.036
// 31   | Success    | c6d4     | cp 89    | bm c6d4            | 2r5/2rk2pp/1pn1pb2/pN1p4/P2P4/1N2B3/nPR1KPPP/3R4 b - - | WAC.037
// 32   | Success    | c3a4     | cp 230   | bm c3a4            | r1br2k1/pp2bppp/2nppn2/8/2P1PB2/2N2P2/PqN1B1PP/R2Q1R1K w - - | WAC.039
// 33   | Success    | d8c8     | cp 374   | bm d8c8            | 3r1r1k/1p4pp/p4p2/8/1PQR4/6Pq/P3PP2/2R3K1 b - - | WAC.040
// 34   | Success    | b4a5     | cp 410   | bm b4a5            | r1b1r1k1/pp1n1pbp/1qp3p1/3p4/1B1P4/Q3PN2/PP2BPPP/R4RK1 w - - | WAC.042
// 35   | Success    | d5c4     | cp 207   | bm d5c4            | 3rb1k1/pq3pbp/4n1p1/3p4/2N5/2P2QB1/PP3PPP/1B1R2K1 b - - | WAC.044
// 36   | Success    | f1a1     | cp 482   | bm f1a1            | 7k/2p1b1pp/8/1p2P3/1P3r2/2P3Q1/1P5P/R4qBK b - - | WAC.045
// 37   | Success    | c3b5     | cp 90    | bm c3b5            | r1bqr1k1/pp1nb1p1/4p2p/3p1p2/3P4/P1N1PNP1/1PQ2PP1/3RKB1R w K - | WAC.046
// 38   | Success    | c6d4     | cp 46    | bm c6d4            | r1b2rk1/pp2bppp/2n1pn2/q5B1/2BP4/2N2N2/PP2QPPP/2R2RK1 b - - | WAC.047
// 39   | Success    | b8b4     | cp 170   | bm b8b4            | 1rbq1rk1/p1p1bppp/2p2n2/8/Q1BP4/2N5/PP3PPP/R1B2RK1 b - - | WAC.048
// 40   | Success    | b7b6     | mate 3   | bm b7b6            | k4r2/1R4pb/1pQp1n1p/3P4/5p1P/3P2P1/r1q1R2K/8 w - - | WAC.050
// 41   | Success    | f4g4     | cp 726   | bm f4g4            | r1bq1r2/pp4k1/4p2p/3pPp1Q/3N1R1P/2PB4/6P1/6K1 w - - | WAC.051
// 42   | Success    | e3e1     | cp 291   | bm e3e1            | 6k1/6p1/p7/3Pn3/5p2/4rBqP/P4RP1/5QK1 b - - | WAC.053
// 43   | Success    | h5h1     | mate 2   | bm h5h1            | r3kr2/1pp4p/1p1p4/7q/4P1n1/2PP2Q1/PP4P1/R1BB2K1 b q - | WAC.054
// 44   | Success    | d4g7     | mate 4   | bm d4g7            | r3r1k1/pp1q1pp1/4b1p1/3p2B1/3Q1R2/8/PPP3PP/4R1K1 w - - | WAC.055
// 45   | Success    | c5f2     | cp 549   | bm c5f2            | r1bqk2r/pppp1ppp/5n2/2b1n3/4P3/1BP3Q1/PP3PPP/RNB1K1NR b KQkq - | WAC.056
// 46   | Success    | f3f8     | mate 3   | bm f3f8            | r3q1kr/ppp5/3p2pQ/8/3PP1b1/5R2/PPP3P1/5RK1 w - - | WAC.057
// 47   | Success    | d3d1     | cp 478   | bm d3d1            | 8/8/2R5/1p2qp1k/1P2r3/2PQ2P1/5K2/8 w - - | WAC.058
// 48   | Success    | c3d5     | cp 482   | bm c3d5            | r1b2rk1/2p1qnbp/p1pp2p1/5p2/2PQP3/1PN2N1P/PB3PP1/3R1RK1 w - - | WAC.059
// 49   | Success    | g6g3     | cp 15    | bm g6g3            | 6r1/3Pn1qk/p1p1P1rp/2Q2p2/2P5/1P4P1/P3R2P/5RK1 b - - | WAC.062
// 50   | Success    | e5f7     | cp 203   | bm e5f7            | r1brnbk1/ppq2pp1/4p2p/4N3/3P4/P1PB1Q2/3B1PPP/R3R1K1 w - - | WAC.063
// 51   | Success    | g2g4     | mate 3   | bm g2g4            | 8/6pp/3q1p2/3n1k2/1P6/3NQ2P/5PP1/6K1 w - - | WAC.064
// 52   | Success    | d5e7     | cp 792   | bm d5e7            | 1r1r1qk1/p2n1p1p/bp1Pn1pQ/2pNp3/2P2P1N/1P5B/P6P/3R1RK1 w - - | WAC.065
// 53   | Success    | c7e5     | cp 295   | bm c7e5            | 1k1r2r1/ppq5/1bp4p/3pQ3/8/2P2N2/PP4P1/R4R1K b - - | WAC.066
// 54   | Success    | e5d5     | cp 327   | bm e5d5            | 3r2k1/p2q4/1p4p1/3rRp1p/5P1P/6PK/P3R3/3Q4 w - - | WAC.067
// 55   | Success    | e2e3     | cp 505   | bm e2e3            | 6k1/5ppp/1q6/2b5/8/2R1pPP1/1P2Q2P/7K w - - | WAC.068
// 56   | Success    | b4a2     | cp 585   | bm b4a2            | 2kr3r/pppq1ppp/3p1n2/bQ2p3/1n1PP3/1PN1BN1P/1PP2PP1/2KR3R b - - | WAC.070
// 57   | Failed     | h3b3     | cp -91   | bm b5a7            | 2kr3r/pp1q1ppp/5n2/1Nb5/2Pp1B2/7Q/P4PPP/1R3RK1 w - - | WAC.071
// 58   | Success    | e5e6     | cp 209   | bm e5e6            | r3r1k1/pp1n1ppp/2p5/4Pb2/2B2P2/B1P5/P5PP/R2R2K1 w - - | WAC.072
// 59   | Success    | e2d2     | cp 400   | bm e2d2            | r1q3rk/1ppbb1p1/4Np1p/p3pP2/P3P3/2N4R/1PP1Q1PP/3R2K1 w - - | WAC.073
// 60   | Success    | d7d6     | cp 383   | bm d7d6            | r3r1k1/pppq1ppp/8/8/1Q4n1/7P/PPP2PP1/RNB1R1K1 b - - | WAC.075
// 61   | Success    | g5f6     | cp 251   | bm g5f6            | r1b1qrk1/2p2ppp/pb1pnn2/1p2pNB1/3PP3/1BP5/PP2QPPP/RN1R2K1 w - - | WAC.076
// 62   | Success    | f5g3     | cp 823   | bm f5g3            | 3r2k1/ppp2ppp/6q1/b4n2/3nQB2/2p5/P4PPP/RN3RK1 b - - | WAC.077
// 63   | Success    | e4g5     | cp 201   | bm e4g5            | r2q3r/ppp2k2/4nbp1/5Q1p/2P1NB2/8/PP3P1P/3RR1K1 w - - | WAC.078
// 64   | Success    | d1a1     | cp 417   | bm d1a1            | r4rk1/p1B1bpp1/1p2pn1p/8/2PP4/3B1P2/qP2QP1P/3R1RK1 w - - | WAC.080
// 65   | Success    | e7d6     | cp 53    | bm e7d6            | r4rk1/1bR1bppp/4pn2/1p2N3/1P6/P3P3/4BPPP/3R2K1 b - - | WAC.081
// 66   | Success    | e4h7     | cp 355   | bm e4h7            | 3rr1k1/pp3pp1/4b3/8/2P1B2R/6QP/P3q1P1/5R1K w - - | WAC.082
// 67   | Success    | d5g8     | mate 2   | bm d5g8            | r2q1r1k/2p1b1pp/p1n5/1p1Q1bN1/4n3/1BP1B3/PP3PPP/R4RK1 w - - | WAC.084
// 68   | Success    | f6g4     | cp 30    | bm f6g4            | 8/p7/1ppk1n2/5ppp/P1PP4/2P1K1P1/5N1P/8 b - - | WAC.086
// 69   | Failed     | c5d4     | cp 178   | bm e6e5            | 8/p3k1p1/4r3/2ppNpp1/PP1P4/2P3KP/5P2/8 b - - | WAC.087
// 70   | Success    | g6g2     | mate 5   | bm g6g2            | r6k/p1Q4p/2p1b1rq/4p3/B3P3/4P3/PPP3P1/4RRK1 b - - | WAC.088
// 71   | Success    | f5g7     | cp 117   | bm f5g7            | 3qrrk1/1pp2pp1/1p2bn1p/5N2/2P5/P1P3B1/1P4PP/2Q1RRK1 w - - | WAC.090
// 72   | Failed     | c8e6     | cp 34    | bm b3e6            | 2qr2k1/4b1p1/2p2p1p/1pP1p3/p2nP3/PbQNB1PP/1P3PK1/4RB2 b - - | WAC.091
// 73   | Failed     | a8e8     | cp -44   | bm e6g4            | r4rk1/1p2ppbp/p2pbnp1/q7/3BPPP1/2N2B2/PPP4P/R2Q1RK1 b - - | WAC.092
// 74   | Success    | c1h6     | cp 159   | bm c1h6            | r1b1k1nr/pp3pQp/4pq2/3pn3/8/P1P5/2P2PPP/R1B1KBNR w KQkq - | WAC.093
// 75   | Success    | e5e4     | cp 729   | bm e5e4            | 8/k7/p7/3Qp2P/n1P5/3KP3/1q6/8 b - - | WAC.094
// 76   | Success    | f6g4     | cp 403   | bm f6g4            | 2r5/1r6/4pNpk/3pP1qp/8/2P1QP2/5PK1/R7 w - - | WAC.095
// 77   | Success    | g2a8     | mate 3   | bm g2a8            | 6k1/5p2/p5np/4B3/3P4/1PP1q3/P3r1QP/6RK w - - | WAC.097
// 78   | Success    | c5e4     | cp 619   | bm c5e4            | 1r3rk1/5pb1/p2p2p1/Q1n1q2p/1NP1P3/3p1P1B/PP1R3P/1K2R3 b - - | WAC.098
// 79   | Success    | e5h5     | mate 2   | bm e5h5            | r1bq1r1k/1pp1Np1p/p2p2pQ/4R3/n7/8/PPPP1PPP/R1B3K1 w - - | WAC.099
// 80   | Success    | d2e3     | cp 206   | bm d2e3 b5b6       | 8/k1b5/P4p2/1Pp2p1p/K1P2P1P/8/3B4/8 w - - | WAC.100
// 81   | Success    | d4c3     | cp 64    | bm d4c3            | 5rk1/p5pp/8/8/2Pbp3/1P4P1/7P/4RN1K b - - | WAC.101
// 82   | Success    | h6g6     | mate 4   | bm h6g6            | 6k1/2pb1r1p/3p1PpQ/p1nPp3/1q2P3/2N2P2/PrB5/2K3RR w - - | WAC.103
// 83   | Success    | d4b5     | cp 173   | bm d4b5            | 5n2/pRrk2p1/P4p1p/4p3/3N4/5P2/6PP/6K1 w - - | WAC.107
// 84   | Success    | c5e5     | cp 253   | bm c5e5            | r5k1/1q4pp/2p5/p1Q5/2P5/5R2/4RKPP/r7 w - - | WAC.108
// 85   | Success    | c4c3     | cp 97    | bm c4c3            | rn2k1nr/pbp2ppp/3q4/1p2N3/2p5/QP6/PB1PPPPP/R3KB1R b KQkq - | WAC.109
// 86   | Success    | a7e3     | cp 691   | bm a7e3            | 2kr4/bp3p2/p2p2b1/P7/2q5/1N4B1/1PPQ2P1/2KR4 b - - | WAC.110
// 87   | Success    | g1f1     | cp 567   | bm g1f1            | 6k1/p5p1/5p2/2P2Q2/3pN2p/3PbK1P/7P/6q1 b - - | WAC.111
// 88   | Success    | e1e6     | cp 361   | bm e1e6            | r4kr1/ppp5/4bq1b/7B/2PR1Q1p/2N3P1/PP3P1P/2K1R3 w - - | WAC.112
// 89   | Success    | d8f6     | cp 87    | bm d8f6            | rnbqkb1r/1p3ppp/5N2/1p2p1B1/2P5/8/PP2PPPP/R2QKB1R b KQkq - | WAC.113
// 90   | Success    | d3h7     | cp 268   | bm d3h7            | r1b1rnk1/1p4pp/p1p2p2/3pN2n/3P1PPq/2NBPR1P/PPQ5/2R3K1 w - - | WAC.114
// 91   | Success    | e8d6     | cp 210   | bm e8d6            | 4N2k/5rpp/1Q6/p3q3/8/P5P1/1P3P1P/5K2 w - - | WAC.115
// 92   | Success    | d8d2     | cp 109   | bm d8d2            | r2r2k1/2p2ppp/p7/1p2P1n1/P6q/5P2/1PB1QP1P/R5RK b - - | WAC.116
// 93   | Success    | d6e4     | cp 679   | bm d6e4            | 3r1rk1/q4ppp/p1Rnp3/8/1p6/1N3P2/PP3QPP/3R2K1 b - - | WAC.117
// 94   | Success    | f4h4     | cp 210   | bm f4h4            | r5k1/pb2rpp1/1p6/2p4q/5R2/2PB2Q1/P1P3PP/5R1K w - - | WAC.118
// 95   | Success    | d8d3     | cp 370   | bm d8d3            | r2qr1k1/p1p2ppp/2p5/2b5/4nPQ1/3B4/PPP3PP/R1B2R1K b - - | WAC.119
// 96   | Success    | c6f3     | cp 890   | bm c6f3            | 6k1/5p1p/2bP2pb/4p3/2P5/1p1pNPPP/1P1Q1BK1/1q6 b - - | WAC.121
// 97   | Success    | d1f1     | cp 859   | bm d1f1            | 1k6/ppp4p/1n2pq2/1N2Rb2/2P2Q2/8/P4KPP/3r1B2 b - - | WAC.122
// 98   | Success    | g4g3     | cp 167   | bm g4g3            | 6k1/3r4/2R5/P5P1/1P4p1/8/4rB2/6K1 b - - | WAC.124
// 99   | Success    | b6d4     | cp 300   | bm b6d4            | r1bqr1k1/pp3ppp/1bp5/3n4/3B4/2N2P1P/PPP1B1P1/R2Q1RK1 b - - | WAC.125
// 100  | Success    | f6c6     | cp 476   | bm f6c6            | r5r1/pQ5p/1qp2R2/2k1p3/4P3/2PP4/P1P3PP/6K1 w - - | WAC.126
// 101  | Success    | b2b7     | cp 338   | bm b2b7            | 2k4r/1pr1n3/p1p1q2p/5pp1/3P1P2/P1P1P3/1R2Q1PP/1RB3K1 w - - | WAC.127
// 102  | Success    | f7g6     | cp 94    | bm f7g6            | 6rk/1pp2Qrp/3p1B2/1pb1p2R/3n1q2/3P4/PPP3PP/R6K w - - | WAC.128
// 103  | Success    | b7f3     | cp 224   | bm b7f3            | 3r1r1k/1b2b1p1/1p5p/2p1Pp2/q1B2P2/4P2P/1BR1Q2K/6R1 b - - | WAC.129
// 104  | Success    | g7h6     | cp 15    | bm g7h6            | 6k1/1pp3q1/5r2/1PPp4/3P1pP1/3Qn2P/3B4/4R1K1 b - - | WAC.130
// 105  | Success    | g3h4     | cp 311   | bm g3h4            | r1b1k2r/1pp1q2p/p1n3p1/3QPp2/8/1BP3B1/P5PP/3R1RK1 w kq - | WAC.133
// 106  | Success    | d8d1     | mate 4   | bm d8d1            | 3r2k1/p6p/2Q3p1/4q3/2P1p3/P3Pb2/1P3P1P/2K2BR1 b - - | WAC.134
// 107  | Success    | e6d4     | cp 123   | bm e6d4            | 3r1r1k/N2qn1pp/1p2np2/2p5/2Q1P2N/3P4/PP4PP/3R1RK1 b - - | WAC.135
// 108  | Success    | d1d7     | cp 23    | bm d1d7            | 3b1rk1/1bq3pp/5pn1/1p2rN2/2p1p3/2P1B2Q/1PB2PPP/R2R2K1 w - - | WAC.137
// 109  | Success    | h4h5     | cp 1091  | bm h4h5            | r1bq3r/ppppR1p1/5n1k/3P4/6pP/3Q4/PP1N1PP1/5K1R w - - | WAC.138
// 110  | Success    | e4f6     | mate 4   | bm e4f6            | rnb3kr/ppp2ppp/1b6/3q4/3pN3/Q4N2/PPP2KPP/R1B1R3 w - - | WAC.139
// 111  | Success    | c3c7     | cp 920   | bm e5c7 c3c7       | r2b1rk1/pq4p1/4ppQP/3pB1p1/3P4/2R5/PP3PP1/5RK1 w - - | WAC.140
// 112  | Failed     | g2f1     | cp -121  | bm c1f4            | 4r1k1/p1qr1p2/2pb1Bp1/1p5p/3P1n1R/1B3P2/PP3PK1/2Q4R w - - | WAC.141
// 113  | Success    | g6h6     | mate 3   | bm g6h6            | 5b2/pp2r1pk/2pp1pRp/4rP1N/2P1P3/1P4QP/P3q1P1/5R1K w - - | WAC.143
// 114  | Success    | d4d3     | cp 165   | bm d4d3            | r2q1rk1/pp3ppp/2p2b2/8/B2pPPb1/7P/PPP1N1P1/R2Q1RK1 b - - | WAC.144
// 115  | Success    | f6g4     | cp 198   | bm f6g4            | r2r2k1/ppqbppbp/2n2np1/2pp4/6P1/1P1PPNNP/PBP2PB1/R2QK2R b KQ - | WAC.147
// 116  | Success    | g1g7     | cp 746   | bm g1g7            | 2r1k3/6pr/p1nBP3/1p3p1p/2q5/2P5/P1R4P/K2Q2R1 w - - | WAC.148
// 117  | Success    | e4g2     | cp 105   | bm e4g2            | 6k1/6p1/2p4p/4Pp2/4b1qP/2Br4/1P2RQPK/8 b - - | WAC.149
// 118  | Success    | a4c3     | cp 73    | bm a4c3            | 8/3b2kp/4p1p1/pr1n4/N1N4P/1P4P1/1K3P2/3R4 w - - | WAC.151
// 119  | Success    | c3e4     | cp 122   | bm c3e4            | 1br2rk1/1pqb1ppp/p3pn2/8/1P6/P1N1PN1P/1B3PP1/1QRR2K1 w - - | WAC.152
// 120  | Success    | f2f7     | mate 2   | bm f2f7            | r1b2rk1/2p2ppp/p7/1p6/3P3q/1BP3bP/PP3QP1/RNB1R1K1 w - - | WAC.154
// 121  | Failed     | c7c8     | cp 35    | bm d5d6            | 5bk1/1rQ4p/5pp1/2pP4/3n1PP1/7P/1q3BB1/4R1K1 w - - | WAC.155
// 122  | Success    | h3h6     | mate 2   | bm h3h6            | r1b1qN1k/1pp3p1/p2p3n/4p1B1/8/1BP4Q/PP3KPP/8 w - - | WAC.156
// 123  | Success    | d5e7     | cp -93   | bm d5e7            | 5rk1/p4ppp/2p1b3/3Nq3/4P1n1/1p1B2QP/1PPr2P1/1K2R2R w - - | WAC.157
// 124  | Success    | g5e6     | cp 640   | bm g5e6            | r1b2r2/5P1p/ppn3pk/2p1p1Nq/1bP1PQ2/3P4/PB4BP/1R3RK1 w - - | WAC.159
// 125  | Success    | c4d5     | cp 225   | bm c4d5            | r3kbnr/p4ppp/2p1p3/8/Q1B3b1/2N1B3/PP3PqP/R3K2R w KQkq - | WAC.162
// 126  | Failed     | c6d5     | cp 54    | bm f3g2            | 5rk1/2p4p/2p4r/3P4/4p1b1/1Q2NqPp/PP3P1K/R4R2 b - - | WAC.163
// 127  | Success    | c2c4     | cp 295   | bm c2c4            | 8/6pp/4p3/1p1n4/1NbkN1P1/P4P1P/1PR3K1/r7 w - - | WAC.164
// 128  | Success    | e3e2     | cp 246   | bm e3e2            | 1r5k/p1p3pp/8/8/4p3/P1P1R3/1P1Q1qr1/2KR4 w - - | WAC.165
// 129  | Success    | d5d4     | cp 263   | bm d5d4            | r3r1k1/5pp1/p1p4p/2Pp4/8/q1NQP1BP/5PP1/4K2R b K - | WAC.166
// 130  | Success    | d7d2     | cp 681   | bm d7d2            | r3k2r/pb1q1p2/8/2p1pP2/4p1p1/B1P1Q1P1/P1P3K1/R4R2 b kq - | WAC.168
// 131  | Success    | g7h6     | cp 344   | bm g7h6            | 5rk1/1pp3bp/3p2p1/2PPp3/1P2P3/2Q1B3/4q1PP/R5K1 b - - | WAC.169
// 132  | Success    | a2c4     | cp 344   | bm a2c4            | 5r1k/6Rp/1p2p3/p2pBp2/1qnP4/4P3/Q4PPP/6K1 w - - | WAC.170
// 133  | Success    | e3h6     | cp 182   | bm e3h6            | 2rq4/1b2b1kp/p3p1p1/1p1nNp2/7P/1B2B1Q1/PP3PP1/3R2K1 w - - | WAC.171
// 134  | Success    | e3h6     | mate 3   | bm e3h6            | 2r1b3/1pp1qrk1/p1n1P1p1/7R/2B1p3/4Q1P1/PP3PP1/3R2K1 w - - | WAC.173
// 135  | Success    | f4h5     | cp 195   | bm f4h5            | r5k1/pppb3p/2np1n2/8/3PqNpP/3Q2P1/PPP5/R4RK1 w - - | WAC.175
// 136  | Success    | f4e6     | cp 134   | bm f4e6            | 3r2k1/p1rn1p1p/1p2pp2/6q1/3PQNP1/5P2/P1P4R/R5K1 w - - | WAC.178
// 137  | Failed     | a7a5     | cp 49    | bm f6d5            | r1q2rk1/p3bppb/3p1n1p/2nPp3/1p2P1P1/6NP/PP2QPB1/R1BNK2R b KQ - | WAC.180
// 138  | Success    | f6g4     | cp 616   | bm f6g4            | r3k2r/2p2p2/p2p1n2/1p2p3/4P2p/1PPPPp1q/1P5P/R1N2QRK b kq - | WAC.181
// 139  | Success    | d1h5     | cp 609   | bm d1h5            | r1b2rk1/ppqn1p1p/2n1p1p1/2b3N1/2N5/PP1BP3/1B3PPP/R2QK2R w KQ - | WAC.182
// 140  | Success    | f4h3     | cp 513   | bm f4h3            | 6k1/5p2/p3p3/1p3qp1/2p1Qn2/2P1R3/PP1r1PPP/4R1K1 b - - | WAC.187
// 141  | Success    | f6g7     | mate 2   | bm f6g7            | 3RNbk1/pp3p2/4rQpp/8/1qr5/7P/P4P2/3R2K1 w - - | WAC.188
// 142  | Success    | d7h3     | cp 410   | bm d7h3            | 8/p2b2kp/1q1p2p1/1P1Pp3/4P3/3B2P1/P2Q3P/2Nn3K b - - | WAC.190
// 143  | Success    | b4d3     | cp 657   | bm b4d3            | r3k3/ppp2Npp/4Bn2/2b5/1n1pp3/N4P2/PPP3qP/R2QKR2 b Qq - | WAC.192
// 144  | Success    | f5h6     | cp 517   | bm f5h6            | 5rk1/ppq2ppp/2p5/4bN2/4P3/6Q1/PPP2PPP/3R2K1 w - - | WAC.194
// 145  | Success    | g2g3     | cp 134   | bm g2g3            | 3r1rk1/1p3p2/p3pnnp/2p3p1/2P2q2/1P5P/PB2QPPN/3RR1K1 w - - | WAC.195
// 146  | Failed     | b8b6     | cp 5     | bm c6b4            | rr4k1/p1pq2pp/Q1n1pn2/2bpp3/4P3/2PP1NN1/PP3PPP/R1B1K2R b KQ - | WAC.196
// 147  | Success    | f2f1     | mate 3   | bm f2f1            | 7k/1p4p1/7p/3P1n2/4Q3/2P2P2/PP3qRP/7K b - - | WAC.197
// 148  | Success    | d8d3     | cp 281   | bm d8d3            | 2br2k1/ppp2p1p/4p1p1/4P2q/2P1Bn2/2Q5/PP3P1P/4R1RK b - - | WAC.198
// 149  | Success    | b2f6     | cp 238   | bm b2f6            | 2rqrn1k/pb4pp/1p2pp2/n2P4/2P3N1/P2B2Q1/1B3PPP/2R1R1K1 w - - | WAC.200
// 150  | Success    | a1a7     | cp 1235  | bm a1a7            | 2b2r1k/4q2p/3p2pQ/2pBp3/8/6P1/1PP2P1P/R5K1 w - - | WAC.201
// 151  | Success    | c2a2     | cp 34    | bm c2a2            | QR2rq1k/2p3p1/3p1pPp/8/4P3/8/P1r3PP/1R4K1 b - - | WAC.202
// 152  | Success    | g5h6     | mate 3   | bm g5h6            | r4rk1/5ppp/p3q1n1/2p2NQ1/4n3/P3P3/1B3PPP/1R3RK1 w - - | WAC.203
// 153  | Success    | e1e5     | cp -236  | bm e1e5            | r1b1qrk1/1p3ppp/p1p5/3Nb3/5N2/P7/1P4PQ/K1R1R3 w - - | WAC.204
// 154  | Success    | d2g5     | cp 202   | bm d2g5            | r3rnk1/1pq2bb1/p4p2/3p1Pp1/3B2P1/1NP4R/P1PQB3/2K4R w - - | WAC.205
// 155  | Success    | d6c6     | cp 120   | bm d6c6            | 1Qq5/2P1p1kp/3r1pp1/8/8/7P/p4PP1/2R3K1 b - - | WAC.206
// 156  | Failed     | f2f4     | cp -290  | bm g4g7            | r1bq2kr/p1pp1ppp/1pn1p3/4P3/2Pb2Q1/BR6/P4PPP/3K1BNR w - - | WAC.207
// 157  | Success    | h5f7     | cp 305   | bm h5f7            | 3r1bk1/ppq3pp/2p5/2P2Q1B/8/1P4P1/P6P/5RK1 w - - | WAC.208
// 158  | Success    | d1h1     | cp -56   | bm d1h1            | 3r1rk1/pp1q1ppp/3pn3/2pN4/5PP1/P5PQ/1PP1B3/1K1R4 w - - | WAC.210
// 159  | Failed     | h8h7     | cp -185  | bm h8g7            | rn1qr2Q/pbppk1p1/1p2pb2/4N3/3P4/2N5/PPP3PP/R4RK1 w - - | WAC.212
// 160  | Failed     | e2b2     | cp -505  | bm h5h7            | 3r1r1k/1b4pp/ppn1p3/4Pp1R/Pn5P/3P4/4QP2/1qB1NKR1 w - - | WAC.213
// 161  | Success    | d3h7     | mate 4   | bm d3h7            | 3r2k1/pb1q1pp1/1p2pb1p/8/3N4/P2QB3/1P3PPP/1Br1R1K1 w - - | WAC.215
// 162  | Success    | f7f1     | mate 3   | bm f7f1            | 7k/p4q1p/1pb5/2p5/4B2Q/2P1B3/P6P/7K b - - | WAC.219
// 163  | Failed     | e4f6     | cp 114   | bm e2f1            | 3rr1k1/ppp2ppp/8/5Q2/4n3/1B5R/PPP1qPP1/5RK1 b - - | WAC.220
// 164  | Failed     | b2a3     | cp 30    | bm h4f6            | 2r1r2k/1q3ppp/p2Rp3/2p1P3/6QB/p3P3/bP3PPP/3R2K1 w - - | WAC.222
// 165  | Success    | d6d5     | cp 221   | bm d6d5            | 2k1rb1r/ppp3pp/2np1q2/5b2/2B2P2/2P1BQ2/PP1N1P1P/2KR3R b - - | WAC.227
// 166  | Success    | d3e4     | cp 221   | bm d3e4            | r4rk1/1bq1bp1p/4p1p1/p2p4/3BnP2/1N1B3R/PPP3PP/R2Q2K1 w - - | WAC.228
// 167  | Success    | c1g5     | cp 35    | bm c1g5            | r4rk1/1b1nqp1p/p5p1/1p2PQ2/2p5/5N2/PP3PPP/R1BR2K1 w - - | WAC.231
// 168  | Success    | c3c1     | cp 94    | bm c3c1            | 1R6/p5pk/4p2p/4P3/8/2r3qP/P3R1b1/4Q1K1 b - - | WAC.236
// 169  | Failed     | c2c5     | cp -67   | bm c2c1            | r5k1/pQp2qpp/8/4pbN1/3P4/6P1/PPr4P/1K1R3R b - - | WAC.237
// 170  | Failed     | c5b5     | cp -19   | bm g2b7            | 1k1r4/pp1r1pp1/4n1p1/2R5/2Pp1qP1/3P2QP/P4PB1/1R4K1 w - - | WAC.238
// 171  | Success    | c2c6     | cp 275   | bm c2c6            | 2b4k/p1b2p2/2p2q2/3p1PNp/3P2R1/3B4/P1Q2PKP/4r3 w - - | WAC.240
// 172  | Success    | d1d7     | cp 310   | bm d1d7            | r1b1r1k1/pp1nqp2/2p1p1pp/8/4N3/P1Q1P3/1P3PPP/1BRR2K1 w - - | WAC.242
// 173  | Failed     | e3e1     | cp 131   | bm e3e8            | 1b5k/7P/p1p2np1/2P2p2/PP3P2/4RQ1R/q2r3P/6K1 w - - | WAC.250
// 174  | Success    | e3e8     | mate 4   | bm e3e8            | k5r1/p4b2/2P5/5p2/3P1P2/4QBrq/P5P1/4R1K1 w - - | WAC.253
// 175  | Success    | f4h3     | cp 275   | bm f4h3            | r6k/pp3p1p/2p1bp1q/b3p3/4Pnr1/2PP2NP/PP1Q1PPN/R2B2RK b - - | WAC.254
// 176  | Success    | f6g6     | cp 102   | bm f6g6            | 3r3r/p4pk1/5Rp1/3q4/1p1P2RQ/5N2/P1P4P/2b4K w - - | WAC.255
// 177  | Success    | d4f5     | cp 176   | bm d4f5            | 3r1rk1/1pb1qp1p/2p3p1/p7/P2Np2R/1P5P/1BP2PP1/3Q1BK1 w - - | WAC.256
// 178  | Failed     | d1d3     | cp 20    | bm d1d4            | 4r1k1/pq3p1p/2p1r1p1/2Q1p3/3nN1P1/1P6/P1P2P1P/3RR1K1 w - - | WAC.257
// 179  | Success    | h5g6     | cp -133  | bm h5g6            | r3brkn/1p5p/2p2Ppq/2Pp3B/3Pp2Q/4P1R1/6PP/5R1K w - - | WAC.258
// 180  | Failed     | c7e6     | cp -42   | bm d5e6            | 2r2b1r/p1Nk2pp/3p1p2/N2Qn3/4P3/q6P/P4PP1/1R3K1R w - - | WAC.260
// 181  | Success    | g6h6     | cp -35   | bm g6h6            | 6k1/p1B1b2p/2b3r1/2p5/4p3/1PP1N1Pq/P2R1P2/3Q2K1 b - - | WAC.262
// 182  | Success    | f7g8     | mate 4   | bm f7g8            | rnbqr2k/pppp1Qpp/8/b2NN3/2B1n3/8/PPPP1PPP/R1B1K2R w KQ - | WAC.263
// 183  | Failed     | c2e4     | cp 18    | bm e5f6            | 2r1k2r/2pn1pp1/1p3n1p/p3PP2/4q2B/P1P5/2Q1N1PP/R4RK1 w k - | WAC.265
// 184  | Failed     | f2g3     | cp -9    | bm h8h2            | r3q2r/2p1k1p1/p5p1/1p2Nb2/1P2nB2/P7/2PNQbPP/R2R3K b - - | WAC.266
// 185  | Success    | d5c7     | cp 1117  | bm d5c7            | 2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - | WAC.267
// 186  | Failed     | e1g1     | cp 81    | bm a3b4            | 2kr2nr/pp1n1ppp/2p1p3/q7/1b1P1B2/P1N2Q1P/1PP1BPP1/R3K2R w KQ - | WAC.269
// 187  | Success    | d4g4     | cp 355   | bm d4g4            | 2r1r1k1/pp1q1ppp/3p1b2/3P4/3Q4/5N2/PP2RPPP/4R1K1 w - - | WAC.270
// 188  | Success    | e6d6     | cp 37    | bm e6d6            | 2kr4/ppp3Pp/4RP1B/2r5/5P2/1P6/P2p4/3K4 w - - | WAC.271
// 189  | Success    | e3c5     | cp 254   | bm e3c5            | nrq4r/2k1p3/1p1pPnp1/pRpP1p2/P1P2P2/2P1BB2/1R2Q1P1/6K1 w - - | WAC.272
// 190  | Failed     | d4e3     | cp 83    | bm c4f7            | r2qkb1r/pppb2pp/2np1n2/5pN1/2BQP3/2N5/PPP2PPP/R1B1K2R w KQkq - | WAC.278
// 191  | Success    | c7h7     | cp 563   | bm c7h7            | 2R5/2R4p/5p1k/6n1/8/1P2QPPq/r7/6K1 w - - | WAC.281
// 192  | Success    | h1h8     | mate 4   | bm h1h8            | 6k1/2p3p1/1p1p1nN1/1B1P4/4PK2/8/2r3b1/7R w - - | WAC.282
// 193  | Success    | h3g5     | mate 4   | bm h3g5            | 3q1rk1/4bp1p/1n2P2Q/3p1p2/6r1/Pp2R2N/1B4PP/7K w - - | WAC.283
// 194  | Success    | d8d5     | cp 153   | bm d8d5            | 3r1k2/1p6/p4P2/2pP2Qb/8/1P1KB3/P6r/8 b - - | WAC.286
// 195  | Success    | h5f6     | cp 348   | bm h5f6            | r1b2rk1/p4ppp/1p1Qp3/4P2N/1P6/8/P3qPPP/3R1RK1 w - - | WAC.288
// 196  | Success    | e6e5     | cp 349   | bm e6e5            | 2r3k1/5p1p/p3q1p1/2n3P1/1p1QP2P/1P4N1/PK6/2R5 b - - | WAC.289
// 197  | Success    | d5d6     | cp 144   | bm d5d6            | 4r3/1Q1qk2p/p4pp1/3Pb3/P7/6PP/5P2/4R1K1 w - - | WAC.292
// 198  | Failed     | e2e3     | cp 0     | bm f3g5            | 1nbq1r1k/3rbp1p/p1p1pp1Q/1p6/P1pPN3/5NP1/1P2PPBP/R4RK1 w - - | WAC.293
// 199  | Success    | d1d5     | mate 3   | bm d1d5            | 4r3/p4r1p/R1p2pp1/1p1bk3/4pNPP/2P1K3/2P2P2/3R4 w - - | WAC.295
// 200  | Success    | d8h8     | mate 4   | bm d8h8            | 3Q4/p3b1k1/2p2rPp/2q5/4B3/P2P4/7P/6RK w - - | WAC.298
// 201  | Success    | g5g6     | cp 281   | bm g5g6            | b2b1r1k/3R1ppp/4qP2/4p1PQ/4P3/5B2/4N1K1/8 w - - | WAC.300
// ====================================================================================================================================
// Successful: 176 (87 %)
// Failed:     25  (12 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
//
// Test time: 1.006.697 ms
func _TestWACTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/wac.epd", 5 * time.Second, 0)
	ts.RunTests()
}

// Results for Test Suite
// ------------------------------------------------------------------------------------------------------------------------------------
// EPD File:   testsets/crafty_test.epd
// SearchTime: 5.000 ms
// MaxDepth:   0
// Date:       2020-04-03 11:40:24.3555133 +0200 CEST
// ====================================================================================================================================
//  Nr. | Result     | Move     | Value    | Expected Result | Fen | Id
// ====================================================================================================================================
// 1    | Success    | h2h3     | cp 48    | bm h2h3            | rn2kb1r/pp3ppp/4pn2/2pq4/3P2b1/2P2N2/PP2BPPP/RNBQK2R w KQkq - | Crafty Test Pos.1
// 2    | Success    | f3e5     | cp 67    | bm f3e5            | r3k2r/pp2qppp/2n1pn2/bN5b/3P4/P3BN1P/1P2BPP1/R2Q1RK1 w kq - | Crafty Test Pos.2
// 3    | Failed     | e2g4     | cp 58    | bm f1d1            | 2rr2k1/1p2qp1p/1pn1pp2/1N6/3P4/P6P/1P2QPP1/2R2RK1 w - - | Crafty Test Pos.3
// 4    | Success    | b5d6     | cp 86    | bm b5d6            | 6rk/1p3p1p/2n2q2/1NQ2p2/3p4/PP5P/5PP1/2R3K1 w - - | Crafty Test Pos.4
// 5    | Success    | d6f7     | cp 304   | bm d6f7            | 7k/2R2p1p/3N1q2/3Q4/3p4/PP3pPP/5n1K/4r3 w - - | Crafty Test Pos.5
// 6    | Failed     | e2e3     | cp -36   | bm f3e5            | r1bqkb1r/pp3ppp/2n1pn2/2p5/2pP4/5NP1/PP2PPBP/RNBQ1RK1 w kq - | Crafty Test Pos.6
// 7    | Failed     | c1f4     | cp -37   | bm e5c6            |1r1q1rk1/p2b1ppp/3bpn2/4N3/3p4/5QP1/PP2PPBP/R1B2RK1 w - - | Crafty Test Pos.7
// 8    | Failed     | e2e3     | cp -114  | bm a4d7            | 1q3rk1/p4p1p/1r3p2/4p3/Qb1p4/6P1/P3PPBP/1R3RK1 w - - | Crafty Test Pos.8
// 9    | Success    | e4h7     | cp 10    | bm e4h7            | 1r6/2q2pkp/5p2/4p3/3pB3/2bQ2P1/P3PP1P/5RK1 w - - | Crafty Test Pos.9
// 10   | Success    | g4c8     | cp 16    | bm g4c8            | 5k1q/5p2/5p2/4p3/3pB1QP/6P1/4PP2/b5K1 w - - | Crafty Test Pos.10
// 11   | Success    | f3d5     | cp 17    | bm f3d5            | 2q5/4kp2/5p2/4p3/2Bp3P/2b2QP1/4PP2/6K1 w - - | Crafty Test Pos.11
// 12   | Success    | c4f7     | cp 49    | bm c4f7            | 8/4kp2/3q4/4pp2/2Bp3P/2b3P1/Q3PP2/6K1 w - - | Crafty Test Pos.12
// 13   | Success    | b8c7     | cp 74    | bm b8c7            | 1Q6/4k3/5q2/1B3p2/3pp2P/6P1/3bPP2/6K1 w - - | Crafty Test Pos.13
// 14   | Failed     | c4d3     | cp 187   | bm c7c8            | 4k3/2Q5/5q2/5p2/2Bp1P1P/6P1/3b4/5K2 w - - | Crafty Test Pos.14
// 15   | Failed     | f5e4     | cp 201   | bm e7d8            | 1k6/4Q3/2q5/5B2/3p1P1P/4b1P1/8/5K2 w - - | Crafty Test Pos.15
// 16   | Failed     | h2h3     | cp 48    | bm e1g1            | rn2kb1r/pp3ppp/4pn2/2pq4/3P2b1/2P2N2/PP2BPPP/RNBQK2R w KQkq - | Crafty Test Pos.16
// 17   | Failed     | e5c4     | cp 24    | bm b2c3            |r3k2r/pp3ppp/2nqpn2/4N3/3P4/P1b1B3/1P2QPPP/R4RK1 w kq - | Crafty Test Pos.17
//
// 18   | Success    | f4b8     | cp 16    | bm f4b8            |4k2r/p4ppp/1p2pn2/8/2rP1B2/P1P5/5PPP/RR4K1 w k - | Crafty Test Pos.18
// 19   | Failed     | d6b4     | cp -7    | bm f1e2            |r5k1/5ppp/p1RBpn2/1p6/r2P4/P1P5/5PPP/1R3K2 w - - | Crafty Test Pos.19
// 20   | Success    | c3d4     | cp 11    | bm c3d4            |8/3r1pk1/p1R2p2/1p5p/r2p4/PRP1K1P1/5P1P/8 w - - | Crafty Test Pos.20
// 21   | Failed     | c4c5     | cp 39    | bm e3e4            |r1bqk2r/pp1n1ppp/2pbpn2/3p4/2PP4/3BPN2/PP1N1PPP/R1BQK2R w KQkq - | Crafty Test Pos.21
// 22   | Success    | d4c3     | cp 39    | bm d4c3            |r1bq1rk1/pp1n1pp1/2p4p/2b5/2PQ4/5N2/PPB2PPP/R1B1R1K1 w - - | Crafty Test Pos.22
// 23   | Failed     | e3f3     | cp 42    | bm h2h3            | r1q1r1k1/1p3pp1/2p1bn1p/p3N3/2P2P2/P1Q1R3/1PB3PP/4R1K1 w - - | Crafty Test Pos.23
// 24   | Failed     | e5d4     | cp 24    | bm c4g4            | r5k1/3q1pp1/2p4p/p2nQP2/2R5/P6P/1PB3P1/6K1 w - - | Crafty Test Pos.24
// 25   | Success    | e4c4     | cp 76    | bm e4c4            | 1q1r4/6pk/2B2p1p/p2n1P2/2p1R3/P6P/1P3QP1/7K w - - | Crafty Test Pos.25
// 26   | Failed     | g3f4     | cp 2     | bm a5e5            | 3r4/5qpk/5p1p/R3nP2/8/5BQP/6P1/7K w - - | Crafty Test Pos.26
// 27   | Failed     | g6f7     | cp -53   | bm g6h5            | r6k/2Q3p1/6Bp/5P2/3q4/7P/6PK/8 w - - | Crafty Test Pos.27
// 28   | Success    | f1d3     | cp 1     | bm f1d3            | r1bqk2r/p1pp1ppp/2p2n2/8/1b2P3/2N5/PPP2PPP/R1BQKB1R w KQkq - | Crafty Test Pos.28
// 29   | Success    | g5f4     | cp 21    | bm g5f4            | r1bqr1k1/p3bpp1/2p2n1p/3p2B1/8/3B1Q2/PPP1NPPP/4RRK1 w - - | Crafty Test Pos.29
// 30   | Success    | f4d2     | cp -5    | bm f4d2            | 4r1k1/p2b1pp1/1q3n1p/3p4/3N1Q2/3B4/PP3PPP/5RK1 w - - | Crafty Test Pos.30
// 31   | Failed     | h2h3     | cp -30   | bm d1d2            | 3r2k1/p4bp1/1q5p/8/3Npp2/1PQ5/P4PPP/3R2K1 w - - | Crafty Test Pos.31
// 32   | Success    | f4e5     | cp -88   | bm f4e5            | 8/p4bpk/7p/3rq3/3N1P2/PPQR1P2/6KP/4q3 w - - | Crafty Test Pos.32
// 33   | Failed     | f2g3     | cp -155  | bm f2e1            | 8/p6k/7p/4P1p1/1Pb5/P3RP2/3r1K1P/8 w - - | Crafty Test Pos.33
// 34   | Failed     | d1c2     | cp 16    | bm g2g3            | r2q1rk1/pp1b1ppp/2nbp3/3p4/2PP1n2/1P3N2/PB1N1PPP/1BRQR1K1 w - - | Crafty Test Pos.35
// 35   | Success    | b3b4     | cp 7     | bm b3b4            | 2rr2k1/pp1qnppp/2n1p3/b2p4/2PP3P/PP2RNP1/1B3P2/1BRQ2K1 w - - | Crafty Test Pos.36
// 36   | Failed     | e2e1     | cp 39    | bm a3a4            | 2r1r3/ppbqnpk1/4p1p1/1PPp1n1p/3P3P/P2Q1NP1/3BRP2/1BR3K1 w - - | Crafty Test Pos.37
// 37   | Success    | f3d4     | cp 94    | bm f3d4            | rb6/1p2rpk1/pP2p1p1/P1PpPn1p/q6P/2BQ1NP1/4RP2/2R3K1 w - - | Crafty Test Pos.38
// 38   | Failed     | c5e7     | cp 175   | bm e2b2            | rbr5/3q1p2/pPp1pBpk/P1QpP2p/7P/6P1/4RP2/2R3K1 w - - | Crafty Test Pos.39
// 39   | Success    | e2f3     | cp 30    | bm e2f3            | r2qkb1r/pppnpppp/8/3p4/4nP2/1P2Pb2/PBPPB1PP/RN1QK2R w KQkq - | Crafty Test Pos.40
// 40   | Success    | c3b5     | cp -23   | bm c3b5            | 2kr1b1r/ppp3pp/8/2n2p2/4pP1q/1PN1P3/PBPP3P/R2Q1R1K w - - | Crafty Test Pos.41
// 41   | Failed     | d2d4     | cp -113  | bm g7c7            | 2k4r/1pp3Rp/p7/2n2r2/5P2/1P2Pp2/P1PP3P/R6K w - - | Crafty Test Pos.42
// 42   | Failed     | h2h3     | cp -458  | bm c2c3            | 8/1pk4p/p7/4r3/5r2/1P1P1p2/P1PR1K1P/8 w - - | Crafty Test Pos.43
// 43   | Failed     | e1g1     | cp 132   | bm h2h3            | rn2kb1r/pp2pppp/2p2n2/q7/2BP2b1/2N2N2/PPP2PPP/R1BQK2R w KQkq - | Crafty Test Pos.44
// 44   | Failed     | g5e6     | cp -10   | bm c3a4            | 2kr1b1r/ppqn1Bpp/2p2nb1/6N1/3p2P1/2N4P/PPPBQP2/2KR3R w - - | Crafty Test Pos.45
// 45   | Success    | e6f8     | cp -20   | bm e6f8            | 2k1rb1r/p2n2pp/2pqN3/1pN1n1P1/3p2Q1/7P/PPPB1P2/2KRR3 w - - | Crafty Test Pos.46
// 46   | Failed     | e1g1     | cp 132   | bm h2h3            | rn2kb1r/pp2pppp/2p2n2/q7/2BP2b1/2N2N2/PPP2PPP/R1BQK2R w KQkq - | Crafty Test Pos.47
// 47   | Failed     | c4d5     | cp 58    | bm c3e4            | r3k2r/ppqn1ppp/2pbp3/3n2P1/2BP4/2N2Q1P/PPPB1P2/2KR3R w kq - | Crafty Test Pos.48
// 48   | Failed     | c4d5     | cp -10   | bm c4e2            | r4rk1/5ppp/1np1p3/p2n2P1/2BP2p1/6PP/PP1B4/2KR3R w - - | Crafty Test Pos.49
// 49   | Success    | f3d5     | cp -39   | bm f3d5            | 2r3k1/5ppp/4p3/1rnn2P1/p7/1P3BPP/P2B4/1KR4R w - - | Crafty Test Pos.50
// 50   | Failed     | b3b4     | cp 70    | bm c2c3            | 5k2/5ppp/8/3p2P1/8/pP4PP/P1K5/8 w - - | Crafty Test Pos.51
// 51   | Success    | c3d4     | cp 61    | bm c3d4            | r1b1kb1r/pp2pppp/2n2n2/3q4/3p4/2P2N2/PP2BPPP/RNBQK2R w KQkq - | Crafty Test Pos.52
// 52   | Failed     | e5c6     | cp 39    | bm c4b3            | 2rq1rk1/pb2bppp/1pn1pn2/4N3/P1BP4/2N1B3/1P3PPP/R2Q1RK1 w - - | Crafty Test Pos.53
// 53   | Failed     | h2h3     | cp 38    | bm e3d2            | q4rk1/pb2bppp/1p2p3/3nN3/P2P4/1B1QBP2/1P4PP/2R3K1 w - - | Crafty Test Pos.54
// 54   | Failed     | f5e6     | cp 85    | bm c1f1            | 2q3k1/1br1b1pp/pp2pp2/3n1P2/P1BP2N1/1P1Q2P1/3B3P/2R3K1 w - - | Crafty Test Pos.55
// 55   | Success    | f5b5     | cp 135   | bm f5b5            | q6k/2r1b1pp/5p2/1p3R2/3P4/1P1Q2P1/3B3P/6K1 w - - | Crafty Test Pos.56
// 56   | Failed     | d1b3     | cp 48    | bm c1e3            | rn1qkb1r/ppp2ppp/4pn2/8/2PP2b1/5N2/PP2BPPP/RNBQK2R w KQkq - | Crafty Test Pos.57
// 57   | Success    | h4g6     | cp 83    | bm h4g6            | r2q1rk1/pp1nbppp/2p1pnb1/8/2PP2PN/P3B2P/1P1NBP2/R2Q1RK1 w - - | Crafty Test Pos.58
// 58   | Success    | a3b4     | cp 70    | bm a3b4            | r2q1rk1/ppbn1pp1/4p1p1/2P3P1/1p1P1P2/P3B2P/4B3/R2Q1RK1 w - - | Crafty Test Pos.59
// 59   | Success    | f1f7     | cp 492   | bm f1f7            | 1rr3k1/1Qbnqpp1/6p1/2p1P1P1/1PBP4/4B2P/8/R4RK1 w - - | Crafty Test Pos.60
// 60   | Failed     | c1e3     | cp -12   | bm f3g5            | r1bqkb1r/ppp1p1pp/1nnpp3/8/2PP4/5N2/PP3PPP/RNBQKB1R w KQkq - | Crafty Test Pos.61
// 61   | Success    | d2f3     | cp -18   | bm d2f3            | r1b2k1r/pp2p3/1n1qp1p1/2pp1n2/2P5/1P1B4/P2N1PPP/R1BQ1RK1 w - - | Crafty Test Pos.62
// 62   | Failed     | g3h4     | cp 19    | bm e5h5            | r4k2/pp1bp2r/1q2p1p1/2p1R3/2Pp2Pn/1P1B2B1/P4P1P/R2Q2K1 w - - | Crafty Test Pos.63
// 63   | Success    | h7g6     | cp 279   | bm h7g6            | 2r2k2/pp2p2B/6q1/2p1p2p/2Pp3B/1P3P1b/P6P/4R1K1 w - - | Crafty Test Pos.64
// 64   | Success    | f4e3     | cp 319   | bm f4e3            | 5k2/1p2p3/1r6/p1p4B/P1P2B2/1P1p1P1b/7P/4R1K1 w - - | Crafty Test Pos.65
// 65   | Success    | e8b5     | cp 377   | bm e8b5            | 4Br2/1p4k1/8/p5B1/P1P5/7b/3p3P/3R2K1 w - - | Crafty Test Pos.66
// 66   | Success    | d2a5     | cp 446   | bm d2a5            | 8/8/6k1/p7/P2r4/5B1b/3B3P/3R2K1 w - - | Crafty Test Pos.67
// 67   | Success    | c6c4     | cp 451   | bm c6c4            | 8/8/2R5/3B4/5kbP/2B5/4r3/6K1 w - - | Crafty Test Pos.68
// 68   | Success    | d4c2     | cp -17   | bm d4c2            | rnb1kb1r/pp3ppp/2p2n2/3q4/3Np3/6P1/PP1PPPBP/RNBQK2R w KQkq - | Crafty Test Pos.69
// 69   | Success    | d2e3     | cp -89   | bm d2e3            | rn3rk1/p4ppp/2p2n2/1p3b1q/4p2P/2N1b1P1/PPQPPPB1/R1B1K2R w KQ - | Crafty Test Pos.70
// 70   | Failed     | e1g1     | cp -42   | bm e1c1            | 3r2k1/p2n1ppp/2p1r1b1/3n3q/PpN1p2P/1P2P1P1/1BQ1PPB1/R3K2R w KQ - | Crafty TestPos.71
// 71   | Failed     | h1f1     | cp -155  | bm h1g1            | 6k1/p1r2pp1/2p1rn2/3n1q2/PpNBp1pP/1P2PP2/1KQ1P3/3R3R w - - | Crafty Test Pos.72
// 72   | Success    | d7d1     | cp -343  | bm d7d1            | 6k1/p2Q1pB1/8/2p1Nq2/Pp2p1rP/1P6/1K2P3/3n4 w - - | Crafty Test Pos.73
// 73   | Success    | h5g6     | cp -777  | bm h5g6            | 7k/p7/6p1/2pq3P/Pp2p3/1P2r3/1K2P2Q/8 w - - | Crafty Test Pos.74
// 74   | Failed     | c4b5     | cp 18    | bm c4b3            | r1bqkb1r/pp1ppppp/1nn5/4P3/2Bp4/2P2N2/PP3PPP/RNBQK2R w KQkq - | Crafty Test Pos.75
// 75   | Failed     | c1d2     | cp 30    | bm a2a3            | r1bq1rk1/pp2bppp/2n1p3/3n4/3P4/1BN2NP1/PP3P1P/R1BQR1K1 w - - | Crafty Test Pos.76
// 76   | Success    | c2a4     | cp 1     | bm c2a4            | 2rqr1k1/pb2bp1p/1p2p1p1/n7/3P4/P1PQ1NP1/2BB1P1P/R3R1K1 w - - | Crafty Test Pos.77
// 77   | Failed     | g4f5     | cp -79   | bm g4h5            | 4r1k1/pb2b3/1p6/3qpppp/2nP2P1/P1P2N1P/5PK1/R1BQR3 w - - | Crafty Test Pos.78
// 78   | Failed     | d4e5     | cp -525  | bm a1e1            | 5r2/p3b2k/1p6/4n2P/3P2p1/P1P2b2/3B1P2/R5K1 w - - | Crafty Test Pos.79
// 79   | Failed     | c4c5     | cp 22    | bm c1e3            | r1bqkb1r/ppp1pppp/1nn5/4P3/2PP4/8/PP4PP/RNBQKBNR w KQkq - | Crafty Test Pos.80
// 80   | Success    | d1c1     | cp 61    | bm d1c1            | r3kb1r/pppn1ppp/2n1p3/2P1P3/3P1q2/2N2P2/PP2BB1P/R2QK2R w KQkq - | Crafty Test Pos.81
// 81   | Failed     | f2g3     | cp 65    | bm e2b2            | 2kr2r1/1pnnb2p/p1p1p3/2P1Pp2/PP1P4/2N2B2/4RB1P/3R3K w - - | Crafty Test Pos.83
// 82   | Success    | b5b6     | cp 120   | bm b5b6            | Rnk3r1/1p4rp/4p3/1PPpPpb1/3P4/5B2/5B1P/R6K w - - | Crafty Test Pos.84
// 83   | Success    | e3f4     | cp 225   | bm e3f4            | Rnkb2r1/1p4r1/1P2p3/2PpP2p/B2P1p2/4B3/7P/1R5K w - - | Crafty Test Pos.85
// 84   | Success    | d5c6     | cp 43    | bm d5c6            | rnbqk2r/pp2nppp/2pb4/3P4/3P1p2/2N2N2/PPP3PP/R1BQKB1R w KQkq - | Crafty Test Pos.86
// 85   | Failed     | e1g1     | cp 48    | bm b3a2            | r1b2rk1/p3nppp/3b4/q1nP4/5p2/PBN2N2/1PP3PP/R1BQK2R w KQ - | Crafty Test Pos.87
// 86   | Success    | e1d2     | cp 83    | bm e1d2            | r3r1k1/p4ppp/b7/2bn4/2PN1p2/P7/B3N1PP/R1B1K2R w KQ - | Crafty Test Pos.88
// 87   | Success    | g1g2     | cp -69   | bm g1g2            | 3r2k1/p4ppp/b7/2b5/2PN4/P4r2/BB1KN1nP/6R1 w - - | Crafty Test Pos.89
// 88   | Success    | c1d2     | cp 131   | bm c1d2            | pr4k1/p4p1p/6p1/2b5/P1PN1R2/8/BB1rN2P/2K5 w - - | Crafty Test Pos.90
// 89   | Failed     | f1h1     | cp -167  | bm h2h4            | 8/p5kp/6p1/2bP1p2/P7/1r6/2NK3P/5R2 w - - | Crafty Test Pos.91
// 90   | Failed     | f1b5     | cp 17    | bm f1e2            | r1bq1rk1/ppp1bppp/2n1pn2/3p4/3P1B2/4PN1P/PPPN1PP1/R2QKB1R w KQ - | Crafty TestPos.92
// 91   | Failed     | a1c1     | cp 67    | bm a2a3            | r1bq1rk1/pp5p/3bpnp1/2ppNp2/2PP4/4PN1P/PP2BPP1/R2Q1RK1 w - - | Crafty Test Pos.93
// 92   | Failed     | e2d3     | cp 31    | bm e2b5            | 3r1r1k/2q4p/1p1bpnp1/p1pbNp2/Q2P4/P3PN1P/1PR1BPP1/3R2K1 w - - | Crafty Test Pos.94
// 93   | Failed     | e1c1     | cp 80    | bm e1f1            | 3r3k/2qN2r1/1p2p1p1/pBpbNp1p/Q2Pn2b/P3P2P/1PR2PP1/4R1K1 w - - | Crafty Test Pos.95
// 94   | Success    | d2e2     | cp 15    | bm d2e2            | 3r3k/2qNb1r1/1p2p3/pBpnNp2/Q2PbPp1/P3P3/1P1R2P1/2R2K2 w - - | Crafty Test Pos.96
// 95   | Success    | b5f1     | cp -1654 | bm b5f1            | 3r4/3N2kr/1p6/pBpn1p2/Q2PR1p1/P7/1P4P1/2q3K1 w - - | Crafty Test Pos.97
// 96   | Success    | e1g1     | cp 106   | bm e1g1            | rn1qkb1r/pp3ppp/2p1pnb1/3p2B1/3P4/2NBPN2/PPP2PPP/R2QK2R w KQkq - | Crafty TestPos.98
// 97   | Success    | c3e4     | cp 49    | bm c3e4            | r3k2r/pp1n1pp1/2pqpnp1/8/3Pp3/2NQ1N2/PPP2PPP/R3R1K1 w kq - | Crafty Test Pos.99
// 98   | Failed     | d4e5     | cp 30    | bm f3e5            | 1k1r3r/pp3pp1/1np1p1p1/4q3/1P1P4/5N2/P1P1RPPP/4R1K1 w - - | Crafty Test Pos.100
// 99   | Success    | f2f4     | cp 0     | bm f2f4            | 3r1r2/1pk2pp1/1pp1p1p1/8/1P1P3P/2P3P1/P1R2P2/4R1K1 w - - | Crafty Test Pos.101
// 100  | Failed     | g3g4     | cp 5     | bm e2h2            | r7/1p4p1/2p1ppp1/1p1k4/1P1P1P1P/r1PK2P1/PR2R3/8 w - - | Crafty Test Pos.102
//
// 101  | Failed     | g3g4     | cp 11    | bm e2e4            | r7/4r1p1/1ppkppp1/1p6/1P1P1P1P/2PK2P1/P2RR3/8 w - - | Crafty Test Pos.103
// 102  | Success    | e1g1     | cp -21   | bm e1g1            | r2qkb1r/pp1n1ppp/2p1bn2/3pp1B1/8/2NP1NP1/PPP1PPBP/R2QK2R w KQkq - | Crafty Test Pos.104
// 103  | Success    | d3d4     | cp -5    | bm d3d4            | 2kr3r/pp1n1p2/2pbbq1p/4p1p1/4P3/2PP2P1/P2NNPBP/R2Q1RK1 w - - | Crafty Test Pos.105
// 104  | Success    | d2b3     | cp -93   | bm d2b3            | 1k1r3r/ppb2p2/1np1b2p/4p1pq/1Q1PP3/2P3P1/1R1NNPBP/R5K1 w - - | Crafty Test Pos.106
// 105  | Success    | b4a5     | cp 109   | bm b4a5            | 1k1rr3/ppb2p2/1np4p/2N3p1/1Q1PP3/2N2qP1/1R3P1P/R5K1 w - - | Crafty Test Pos.107
// 106  | Failed     | d7b6     | mate 3   | bm a5c5            | k7/1pRN1p2/p3r2p/Q2N2p1/3PPP2/8/5P1P/R5K1 w - - | Crafty Test Pos.108
// 107  | Failed     | g1f3     | cp 80    | bm a2a3            | r2qkb1r/pb1p1ppp/1pn1pn2/2p5/3PP3/2PBB3/PP1N1PPP/R2QK1NR w KQkq - | Crafty Test Pos.109
// 108  | Failed     | e3g5     | cp 88    | bm h2h3            | 2r2rk1/pb1q1pbp/1pnppnp1/8/1P1PP3/P2BBQ2/3NNPPP/2R2RK1 w - - | Crafty Test Pos.110
// 109  | Success    | b5d3     | cp 32    | bm b5d3            | 3q1rk1/pb2npb1/1p1np1pp/1B1p4/1P1PP3/P3BPQP/3NN1P1/2R3K1 w - - | Crafty Test Pos.111
// 110  | Success    | d6e5     | cp 5     | bm d6e5            | b1nqr1k1/5pb1/p2Bp1pp/1pR5/1PpPP3/P4PQP/3NN1P1/6K1 w - - | Crafty Test Pos.112
// 111  | Success    | a4b5     | cp 51    | bm a4b5            | b7/3q1pk1/p2nr1pp/1pR5/PPp1PQ1P/5P2/3NN1P1/6K1 w - - | Crafty Test Pos.113
// 112  | Failed     | c5d5     | cp 121   | bm c3d5            | 8/3q1pk1/6pp/1pRb4/1r5P/2N5/5QP1/6K1 w - - | Crafty Test Pos.114
// 113  | Failed     | c1g5     | cp 39    | bm g2g3            | r2qkb1r/pp1npppp/2p2n2/5b2/2QP4/2N2N2/PP2PPPP/R1B1KB1R w KQkq - | Crafty Test Pos.115
// 114  | Failed     | c2d3     | cp 60    | bm c2a4            | r3kb1r/pp2pppp/2p1b3/5q2/3Pn3/2P2NP1/PBQ1PPBP/R3K2R w KQkq - | Crafty Test Pos.116
// 115  | Success    | d4d5     | cp 95    | bm d4d5            | 2kr1b1r/1p2pb1p/pQp2p2/5qp1/2PPn3/5NP1/PB2PPBP/2R1R1K1 w - - | Crafty Test Pos.117
// 116  | Success    | b1b7     | cp 57    | bm b1b7            | 3k1b1r/1r5p/p1Q1bp2/2p2qp1/2P5/2n3P1/P2NPP1P/1R2R1K1 w - - | Crafty Test Pos.118
// 117  | Success    | e2e3     | cp 360   | bm e2e3            | 4kr2/4q2p/p7/2p3p1/2PbN3/3Q2P1/P3PP1P/4R1K1 w - - | Crafty Test Pos.119
// 118  | Failed     | f1e2     | cp 15    | bm b2b3            | rnbqk1nr/pp3ppp/3b4/3p4/2pP4/2P1BN2/PP3PPP/RN1QKB1R w KQkq - | Crafty Test Pos.120
// 119  | Failed     | e1g1     | cp 32    | bm b5c3            | rb3rk1/pp1qnppp/2n5/1N1p1b2/R1PP4/1P1BBN2/5PPP/3QK2R w K - | Crafty Test Pos.121
// 120  | Failed     | h4f3     | cp -8    | bm h2h3            | 3rr1k1/1p1qnppp/p1n5/b1Pp4/R2P3N/1PNQB3/5PPP/4R1K1 w - - | Crafty Test Pos.122
// 121  | Failed     | e2d2     | cp 42    | bm b3b4            | 4r1k1/1p4p1/p1n1rp2/2Pp1q1p/3P4/1PQ1B2P/4RPP1/4R1K1 w - - | Crafty Test Pos.123
// 122  | Failed     | e1d1     | cp 34    | bm c1d2            | 6k1/1p2r3/p3rp2/1nPp1q1p/1P1P2pP/4B3/4RPP1/2Q1R1K1 w - - | Crafty Test Pos.124
// 123  | Success    | c6b7     | cp -46   | bm c6b7            | 8/1p2rk2/2P5/pP1p1p1p/2nPr1pP/1Q2BqP1/4RP1K/4R3 w - - | Crafty Test Pos.125
// 124  | Failed     | d1b1     | cp -315  | bm f2g2            | 8/4rk2/1P6/p2p3p/2nP1P1P/4r3/5R2/3R2K1 w - - | Crafty Test Pos.126
// 125  | Failed     | f1e2     | cp -2    | bm c1f4            | r2qkbnr/pp2pppp/2p5/3Pn3/2p1P1b1/2N2N2/PP3PPP/R1BQKB1R w KQkq - | Crafty Test Pos.127
// 126  | Failed     | h1f1     | cp 25    | bm d4c5            | r3kbnr/pp2pppp/5q2/1N1P4/2BQ4/4Bb2/PP3P1P/R3K2R w KQkq - | Crafty Test Pos.128
// 127  | Success    | e1e2     | cp 114   | bm e1e2            | r3k1nr/p4ppp/2p1p3/3P4/1b6/5Q2/PP3P1P/R3K2R w KQkq - | Crafty Test Pos.129
// 128  | Failed     | h2h3     | cp 66    | bm c1c6            | r4k1r/1R3pp1/3bpn1p/p2p4/Q7/8/PP2KP1P/2R5 w - - | Crafty Test Pos.130
// 129  | Success    | b8d6     | cp 208   | bm b8d6            | 1Q5r/5ppk/3np2p/P2p4/8/8/P3KP2/8 w - - | Crafty Test Pos.131
// 130  | Failed     | e1g1     | cp 6     | bm c1f4            | rnbqk1nr/pp2bppp/3pp3/8/2B1P3/2N2N2/PP3PPP/R1BQK2R w KQkq - | Crafty Test Pos.132
// 131  | Failed     | f3d2     | cp -73   | bm f4d2            | r1b2rk1/pp2bppp/2n1p3/8/4BB2/5N2/PP1q1PPP/2R2RK1 w - - | Crafty Test Pos.133
// 132  | Failed     | f4d2     | cp -88   | bm f4e3            | 2rrb1k1/pp3p2/2n1pb1p/6p1/P3BB2/1P3N1P/5PP1/3RR1K1 w - - | Crafty Test Pos.134
// 133  | Success    | b3b4     | cp -66   | bm b3b4            | 2rb2k1/p7/1pb1p2p/n5p1/P4p2/1P1B1N1P/1B3PP1/3R2K1 w - - | Crafty Test Pos.135
// 134  | Success    | b7f3     | cp 114   | bm b7f3            | 8/pB3k2/1p1rp2p/2b3p1/P4p2/2B4P/5PP1/2R3K1 w - - | Crafty Test Pos.136
// 135  | Success    | h5g4     | cp 190   | bm h5g4            | 8/p3k1B1/1p1bp2p/6pB/P3Rp2/7P/5PPK/2r5 w - - | Crafty Test Pos.137
// 136  | Success    | e4d4     | cp 188   | bm e4d4            | 8/p7/1p1k3p/6p1/P2bR1P1/1B3p1P/5P1K/5r2 w - - | Crafty Test Pos.138
// 137  | Success    | e3e6     | cp 389   | bm e3e6            | 8/5B2/p6p/1p4p1/P2k2P1/4RK1P/5P2/1r6 w - - | Crafty Test Pos.139
// 138  | Success    | a6a5     | cp 457   | bm a6a5            | 1r6/8/R7/6p1/2B3P1/p3K2P/1k3P2/8 w - - | Crafty Test Pos.140
// 139  | Success    | c4d5     | mate 6   | bm c4d5            |8/8/8/R5p1/2B3P1/p6r/2K2P2/k7 w - - | Crafty Test Pos.141
// 140  | Success    | f3d2     | cp -26   | bm f3d2            |rnbqk2r/ppp1n1pp/3p4/5p2/2PPp3/2N2N2/PP2PPPP/R2QKB1R w KQkq - | Crafty Test Pos.142
// 141  | Failed     | d1c2     | cp 37    | bm e1g1            |r3k2r/pp2nbpp/1q1p1n2/2pP1p2/2P5/1NN1PP2/PP2B2P/R2QK2R w KQkq - | Crafty Test Pos.143
// 142  | Success    | f4d3     | cp 77    | bm f4d3            |5k1r/1p2n1pp/rq1p2b1/p1pP1p2/Q1P2N2/2NnPP2/PP5P/R4RK1 w - - | Crafty Test Pos.144
// 143  | Failed     | a1d1     | cp 50    | bm c7e6            |5qkr/rpNbn2p/3p2p1/p1pP1p2/2P2N2/4PP2/PPQ4P/R5RK w - - | Crafty Test Pos.145
// 144  | Failed     | g1e1     | cp 65    | bm e3e4            |r6r/4nk1p/3pNqp1/p1pP1p2/8/3QPP2/PP4RP/6RK w - - | Crafty Test Pos.146
// 145  | Failed     | d2a5     | cp 159   | bm c7e8            |6kr/1rN1n2p/3p1qp1/p1pP1p2/4P3/P4P2/1P1Q2RP/4R2K w - - | Crafty Test Pos.147
// 146  | Success    | g2g6     | cp 191   | bm g2g6            | 6kr/8/6p1/p2n4/2N5/P3pP2/1P4RP/7K w - - | Crafty Test Pos.148
// 147  | Failed     | c4d6     | cp -383  | bm c4b2            | 2r5/7k/8/8/2N2n1P/Pp3P2/4p3/4R2K w - - | Crafty Test Pos.149
// 148  | Failed     | g2g4     | cp 56    | bm b1c3            | rn1qkb1r/ppp2ppp/4pn2/7b/2BP4/4PN1P/PP3PP1/RNBQK2R w KQkq - | Crafty Test Pos.150
// 149  | Failed     | e1g1     | cp 80    | bm h1g1            | r3kb1r/pppnqpp1/3np3/7p/3P2P1/1B2PN1P/PP1B1P2/R2QK2R w KQkq - | Crafty Test Pos.151
// 150  | Success    | f3d2     | cp 28    | bm f3d2            | 2kr4/2pnbpp1/1p1qp3/p7/3P2P1/1B2PN1r/PPQ2P2/2KR2R1 w - - | Crafty Test Pos.152
// 151  | Failed     | c2e2     | cp 111   | bm e4d2            | 1n6/1kq1bp2/1pp1p3/p5p1/B2PN1P1/4P2r/PPQ2P2/1KR5 w - - | Crafty Test Pos.153
// 152  | Failed     | e4d2     | cp 70    | bm c1h1            | 6r1/1k1qb3/1pp1pp2/p2n2p1/3PN1P1/PB2P1Q1/1P3P2/1KR5 w - - | Crafty Test Pos.154
// 153  | Failed     | g6h5     | cp 39    | bm h1h7            | 8/1k4r1/1pp1qbQ1/p2n1pp1/3P4/PB2P3/1P1N1P2/1K5R w - - | Crafty Test Pos.155
// 154  | Failed     | h8h7     | cp 100   | bm e5c6            | 7Q/1k2b3/1pp1q3/p2nNp2/3P2p1/PB2P3/1P1K1P2/8 w - - | Crafty Test Pos.156
// 155  | Failed     | f4f5     | cp 124   | bm g4f5            | 8/1kn5/1p6/p7/5PQ1/PB6/1P3q2/2K5 w - - | Crafty Test Pos.157
// 156  | Failed     | d3e3     | cp 350   | bm d3c3            | 8/8/1p6/3k4/p4P2/P2K4/1P6/8 w - - | Crafty Test Pos.158
// 157  | Failed     | g1f3     | cp 11    | bm b1d2            | rnbqk2r/ppp1p1b1/3p1n1p/5pp1/3P4/2P1P1B1/PP3PPP/RN1QKBNR w KQkq - | Crafty Test Pos.159
// 158  | Success    | c3d4     | cp -21   | bm c3d4            | r1bq1rk1/ppp3b1/2n4p/3p2pn/3pPp2/1QPB1P2/PP1N1BPP/R3K1NR w KQ - | Crafty Test Pos.160
// 159  | Success    | f2f1     | cp -107  | bm f2f1            | r4rk1/ppp5/1q5p/n4bpn/4Np2/2NB1P2/PPQ2KPP/R6R w - - | Crafty Test Pos.161
// 160  | Failed     | d1d3     | cp -213  | bm d1a4            | 6k1/ppp5/1q5p/n2nr1p1/4Np2/5P2/PP2B1PP/3Q1K1R w - - | Crafty Test Pos.162
// 161  | Success    | f2e2     | cp -711  | bm f2e2            | 8/ppQ1n1k1/2nN3p/6p1/5p2/5P2/P3rKPP/7q w - - | Crafty Test Pos.163
// 162  | Failed     | g1f3     | cp 79    | bm f2f4            | rnbq1rk1/1pp1ppbp/p2p1np1/8/2PPP3/2N1B3/PP2BPPP/R2QK1NR w KQ - | Crafty Test Pos.164
// 163  | Failed     | a2a4     | cp 91    | bm c3d1            | rnbq1rk1/2n2pbp/p3p1p1/1p1pP3/3P1P2/2N1BN2/PP2B1PP/2R1QRK1 w - - | Crafty TestPos.165
// 164  | Failed     | c6c1     | cp 82    | bm f1c1            | n1bq1rk1/r2n1pbp/2RQp1p1/p2pP3/1p1P1P2/5N2/PP1BBNPP/5RK1 w - - | Crafty Test Pos.166
// 165  | Failed     | d3b5     | cp 87    | bm c2c8            | n1r2bk1/rb3p1p/1n2p1p1/p2pP3/1p1P1P2/3BBN2/PPR2NPP/2R3K1 w - - | Crafty Test Pos.167
// 166  | Failed     | c1d2     | cp 97    | bm a3b4            | 6k1/1bn2p1p/1n2pPp1/p2pN3/1p1P1P2/P2B4/1P4PP/2B3K1 w - - | Crafty Test Pos.168
// 167  | Success    | f4f5     | cp 1028  | bm f4f5            | 7k/1bn4p/5PpN/3pp3/1B1P1P2/3B4/1n4PP/6K1 w - - | Crafty Test Pos.169
// 168  | Success    | e4e5     | cp 40    | bm e4e5            | rnbqkb1r/pp2p3/2pp1n1p/5pp1/3PPB2/2P5/PP1N1PPP/R2QKBNR w KQkq - | Crafty Test Pos.170
// 169  | Failed     | f1c4     | cp -9    | bm e5d6            | rn1q1br1/ppk1p3/2ppb2p/4P1n1/3P1Q2/2P1NN2/PP3PPP/R3KB1R w KQ - | Crafty Test Pos.171
// 170  | Success    | f4f3     | cp -294  | bm f4f3            | rk3br1/pp2q3/2npn2p/1N1p4/5Q1P/2P5/PP2BPP1/R3K2R w KQ - | Crafty Test Pos.172
// 171  | Success    | f3d4     | cp 60    | bm f3d4            | r1bqk1nr/pp1nppbp/3p2p1/8/2PpP3/2N2N2/PP2BPPP/R1BQK2R w KQkq - | Crafty Test Pos.173
// 172  | Success    | d4f3     | cp 60    | bm d4f3            | r2qr1k1/1p1bppbp/p1np1np1/8/2PNPP2/2N1B2P/PP1QB1P1/R4RK1 w - - | Crafty Test Pos.174
// 173  | Success    | b4b5     | cp 192   | bm b4b5            | 1r2r1k1/1q1bppbp/ppnp2p1/4P2n/1PP2P2/P1NBBN1P/5QP1/2R2RK1 w - - | Crafty Test Pos.175
// 174  | Failed     | d3e2     | cp 339   | bm f2h4            | 1r1nr1kb/3qpp2/1p1p2p1/1P2P1N1/5Pb1/P1NBB3/5Q2/2RR2K1 w - - | Crafty Test Pos.176
// 175  | Success    | g1f3     | cp 28    | bm g1f3            | rnbqk2r/pp2p1bp/2p2ppn/3pP3/3P1P2/2P5/PP1N2PP/R1BQKBNR w KQkq - | Crafty Test Pos.177
// 176  | Failed     | c1h6     | cp 39    | bm a2a4            | r1b2rk1/1p2p1bp/1qn3pn/p2pP3/3P4/1N3N2/PP2B1PP/R1BQ1R1K w - - | Crafty Test Pos.178
// 177  | Failed     | e2c4     | cp 31    | bm d2c1            | r4rk1/1p2p1bp/2n3pn/p3P3/P2Pp3/q4N2/3BB1PP/1R1Q1R1K w - - | Crafty Test Pos.179
// 178  | Success    | c4d5     | cp 61    | bm c4d5            | 5r2/rp2p1kp/2n3p1/p3P3/P1BP4/1R3P2/7P/3R3K w - - | Crafty Test Pos.180
// 179  | Failed     | c5a5     | cp 75    | bm g2f2            | 8/3rp1kp/1rp3p1/p1R1P3/P2P4/5P2/6KP/3R4 w - - | Crafty Test Pos.181
// 180  | Failed     | e3e2     | cp -50   | bm e3f4            | 6k1/4R2p/4p1p1/p2pP3/3P4/r3KP2/7P/8 w - - | Crafty Test Pos.182
// 181  | Failed     | a6a7     | cp -41   | bm a6f6            | 8/5k2/R5pp/3pP3/p2r1PKP/8/8/8 w - - | Crafty Test Pos.183
// 182  | Success    | d7d5     | cp 21    | bm d7d5            | 5k2/3R4/4P2p/3p1K1p/p1r2P2/8/8/8 w - - | Crafty Test Pos.184
// 183  | Failed     | e6e7     | cp 314   | bm h5h3            | 6k1/8/4P3/4KP1R/2r5/p6p/8/8 w - - | Crafty Test Pos.185
// 184  | Success    | e1g1     | cp 12    | bm e1g1            | rnbqk2r/ppp3bp/3p1np1/4pp2/2P5/2NP1NP1/PP2PPBP/R1BQK2R w KQkq - | Crafty Test Pos.186
// 185  | Success    | b4b5     | cp 1     | bm b4b5            | r2q1rk1/1pp3bp/2npbnp1/4ppB1/1PP5/2NP1NP1/3QPPBP/1R3RK1 w - - | Crafty Test Pos.187
// 186  | Failed     | b5b6     | cp -33   | bm b5c6            | 5rk1/1pq1n1bp/2ppbnp1/1P2ppB1/N1P5/3P2P1/2NQPPBP/1R4K1 w - - | Crafty Test Pos.188
// 187  | Success    | g3f4     | cp -2    | bm g3f4            | 6k1/r2nn1b1/2ppb1pp/5p2/2P2p2/2NP2P1/2NBP1BP/1R4K1 w - - | Crafty Test Pos.189
// 188  | Failed     | e2e3     | cp -40   | bm d1b2            | 6k1/4n3/1n2bbpp/5p2/1N1p1P2/3P4/4PKBP/2BN4 w - - | Crafty Test Pos.190
// 189  | Failed     | f1g1     | cp -67   | bm d2e3            | 8/5bk1/5bpp/5p2/N2p1P2/3Pn3/3BP1BP/5K2 w - - | Crafty Test Pos.191
// 190  | Success    | f4g5     | cp -195  | bm f4g5            | 8/6k1/7p/2bB1pp1/5P2/4p3/4P1KP/8 w - - | Crafty Test Pos.192
// 191  | Failed     | e4e5     | cp 55    | bm d2d4            | r1bq1rk1/pp1pppbp/2n2np1/1Bp5/4P3/2P2N2/PP1P1PPP/RNBQR1K1 w - - | Crafty Test Pos.193
// 192  | Failed     | g2g4     | cp 57    | bm c1a3            | r2q1rk1/p3ppbp/2p3p1/3pPb2/3P4/2P2N1P/P4PP1/R1BQR1K1 w - - | Crafty Test Pos.194
// 193  | Success    | f3f4     | cp 58    | bm f3f4            | 1r4k1/prq1ppbp/2p3p1/2BpP3/b2P3N/2P1QP1P/P5P1/R1R3K1 w - - | Crafty Test Pos.195
// 194  | Success    | h6h7     | cp 199   | bm h6h7            | 1r6/pr2pk1p/2p2b1Q/1bBp1p2/3P4/2P2NqP/P5P1/R3R1K1 w - - | Crafty Test Pos.196
// 195  | Failed     | d4e5     | cp 338   | bm f5e5            | 1r1k4/pr2p3/2R5/1b1pqQ2/3P4/2P4P/P5P1/R5K1 w - - | Crafty Test Pos.197
// 196  | Success    | f7e7     | cp 403   | bm f7e7            | 8/p1k1pQ2/1rb5/3p4/3P2P1/2P4P/P6K/8 w - - | Crafty Test Pos.198
// 197  | Success    | b4a3     | cp 571   | bm b4a3            | 2k5/p2b4/8/8/1QPP2P1/6KP/P3r3/8 w - - | Crafty Test Pos.199
// 198  | Success    | e2e3     | cp 29    | bm e2e3            | rnbqk2r/ppp2pp1/4pb1p/3p4/2PP4/2N2N2/PP2PPPP/R2QKB1R w KQkq - | Crafty Test Pos.200
// 199  | Failed     | c4d3     | cp 69    | bm d2c3            | r1bq1rk1/1pp2pp1/p4b1p/3Ppn2/2B1N3/4PN2/PP1Q1PPP/R4RK1 w - - | Crafty Test Pos.201
// 200  | Failed     | f1e1     | cp 94    | bm g1h1            | 2r1r1k1/1p3pp1/pQ1n2qp/3P4/4p1b1/1B2P3/PP1N1PPP/1R3RK1 w - - | Crafty Test Pos.202
// 201  | Failed     | g2g3     | cp 115   | bm e1e2            | 2r3k1/1p3pp1/pQ1n2qp/3P4/4pPr1/1B2P3/PP4PP/1R2R1K1 w - - | Crafty Test Pos.203
// 202  | Failed     | e4c2     | cp 144   | bm e4d3            | 6k1/1p4p1/p2n2r1/3P1p1p/3QBPq1/4P1P1/PP5P/1R4K1 w - - | Crafty Test Pos.204
// 203  | Success    | e7b7     | cp 280   | bm e7b7            | 8/1p2Q1pk/6rq/p2P1p2/5P1p/3BP1nP/PP3K2/1R6 w - - | Crafty Test Pos.205
// 204  | Failed     | b1c1     | cp 211   | bm d5d6            | 8/6pk/5r2/2QP1p2/5P1p/1P1BP1nP/P5q1/1R2K3 w - - | Crafty Test Pos.206
// 205  | Failed     | c1c4     | cp 191   | bm c5f2            | 8/3P2pk/8/2Q2p2/4qP1p/1P1r3P/PK6/2R5 w - - | Crafty Test Pos.207
// 206  | Failed     | b3b4     | cp 811   | bm c3c4            | 3Q4/2R3pk/8/q4p2/5P1p/1PK5/3Q4/7r w - - | Crafty Test Pos.208
// 207  | Success    | f3e5     | cp 129   | bm f3e5            | r1b1kbnr/p2pqppp/1pn5/2pQ4/2B5/4PN2/PPP2PPP/RNB1K2R w KQkq - | Crafty Test Pos.209
// 208  | Success    | c1g5     | cp 229   | bm c1g5            | rkb3nr/p2p2p1/1p4n1/2p1q2p/2B1P2P/2N2Q2/PPP2PP1/R1B1K2R w KQ - | Crafty Test Pos.210
// 209  | Failed     | d1d2     | cp 229   | bm e1d2            | rkb2r2/3p2p1/pp2n3/2p1q1BB/3nP2P/2N1Q1P1/PPP2P2/3RK2R w K - | Crafty Test Pos.211
// 210  | Failed     | a2a3     | cp 190   | bm d2c1            | r4r2/kb3qpR/p2p4/1pp3P1/3nPPB1/2N1Q1P1/PPPK4/7R w - - | Crafty Test Pos.212
// 211  | Failed     | c2d1     | cp 289   | bm c2b1            | 4rr2/k5pR/p2p4/2p3P1/5PB1/1nPQ2P1/P1K3q1/7R w - - | Crafty Test Pos.213
// 212  | Success    | g4d7     | cp 700   | bm g4d7            | 5r2/6R1/p2R4/n1p3P1/k4PB1/2P3P1/2K5/4r3 w - - | Crafty Test Pos.214
// 213  | Failed     | c1d2     | cp -63   | bm d4d5            | rnbq1rk1/ppp1bppp/3p1n2/4p3/2PP4/P1N1P3/1P2NPPP/R1BQKB1R w KQ - | Crafty Test Pos.215
// 214  | Success    | e3e4     | cp 21    | bm e3e4            | 2rqbrk1/pp1nbppp/3p1n2/3Pp3/8/P1N1PPN1/1P2B1PP/R1BQ1R1K w - - | Crafty Test Pos.216
// 215  | Failed     | h1g1     | cp 84    | bm d1d2            | 2rqbr1k/1p1nbpp1/p2p1n1p/P2PpN2/1P2P3/2N1BP2/4B1PP/2RQ1R1K w - - | Crafty TestPos.217
// 216  | Success    | f5d6     | cp 208   | bm f5d6            | q3bbrk/1pB2pp1/p2p3p/P2PpN1n/1P2P3/5P2/3QB1PP/2R4K w - - | Crafty Test Pos.218
// 217  | Success    | e2b5     | cp 259   | bm e2b5            | q1r4k/1pB2bp1/1Q3p1p/Pp1Pp2n/4P3/5PP1/4B2P/2R4K w - - | Crafty Test Pos.219
//
// 218  | Failed     | d1e2     | cp 18    | bm b5d7            | rn1qkb1r/pp1b1ppp/5n2/1Bpp4/3P4/5N2/PPPN1PPP/R1BQK2R w KQkq - | Crafty Test Pos.220
// 219  | Failed     | g5f6     | cp -8    | bm g5h4            | r2q1rk1/pp3pp1/5n1p/2b3B1/3p4/3Q1N2/PPP2PPP/R4RK1 w - - | Crafty Test Pos.221
// 220  | Success    | d1d3     | cp -15   | bm d1d3            | 3rr1k1/pp3p2/7p/2b2p2/3p4/5N2/PPP2PPP/2RR1K2 w - - | Crafty Test Pos.222
// 221  | Failed     | b2b4     | cp 7     | bm e5d3            | 3r4/2r2pk1/1p5p/pRb1Np2/3p4/P7/1PP2PPP/2R2K2 w - - | Crafty Test Pos.223
// 222  | Success    | c3d4     | cp 72    | bm c3d4            | 2r5/5p2/1p3k1p/1RbR1p2/p2pr3/P1PN2P1/1P3P1P/5K2 w - - | Crafty Test Pos.224
// 223  | Success    | f1e2     | cp 174   | bm f1e2            | 8/5p2/1p5p/r2k1p2/PR1b4/1P1N2P1/5P1P/5K2 w - - | Crafty Test Pos.225
// 224  | Success    | b2c4     | cp 270   | bm b2c4            | 8/8/5p1p/1P1k1p1P/1b3P2/1P4P1/1N2K3/8 w - - | Crafty Test Pos.226
// 225  | Failed     | f3e5     | cp 13    | bm d2d3            | r1bq1rk1/ppp2ppp/2p2n2/4p3/1b2P3/2N2N2/PPPP1PPP/R1BQ1RK1 w - - | Crafty Test Pos.227
// 226  | Failed     | c1d2     | cp -30   | bm d1e2            | r4rk1/ppp1qppp/2p1b3/2b1p3/4P1P1/3P1N1P/PPP3PK/R1BQ1R2 w - - | Crafty Test Pos.228
// 227  | Failed     | g3h2     | cp 10    | bm h3h4            | r4rk1/p1p3p1/1p2bp1p/q1p1p3/P3P1P1/1P1PQNKP/2P2RP1/R7 w - - | Crafty Test Pos.229
// 228  | Success    | e2b2     | cp -92   | bm e2b2            | 5rk1/p1p3p1/1p6/2p1p1p1/P3P3/1Pq1rNK1/4QRP1/3R4 w - - | Crafty Test Pos.230
// 229  | Failed     | f3h4     | cp 288   | bm f3g5            | 2r3k1/p1pR2p1/1p6/4R3/r1p1P1K1/5N2/6P1/8 w - - | Crafty Test Pos.231
// 230  | Success    | g5h7     | cp 392   | bm g5h7            | 3r3k/p1R4r/1p6/4P1N1/2p3K1/6P1/8/8 w - - | Crafty Test Pos.232
// 231  | Failed     | d4f3     | cp 94    | bm f1e2            | rnbqkb1r/1p2pp1p/p2p1np1/8/P2NP3/2N5/1PP2PPP/R1BQKB1R w KQkq - | Crafty Test Pos.233
// 232  | Failed     | h1g1     | cp 26    | bm f4f5            | 2rqr1k1/1p2ppbp/p1npbnp1/8/P3PP2/RNN1B3/1PP1B1PP/3Q1R1K w - - | Crafty Test Pos.234
// 233  | Success    | b3b4     | cp 44    | bm b3b4            | 2r1r3/1p1qppk1/p2p1bp1/P2Pn3/8/1RP1B3/1P2B1PP/3Q1R1K w - - | Crafty Test Pos.235
// 234  | Failed     | b3b6     | cp 48    | bm f1f7            | 8/1prqppk1/p5p1/P2Pp3/7b/1QP5/1P2B1PP/5RK1 w - - | Crafty Test Pos.236
// 235  | Failed     | f3d5     | cp -16   | bm b4e4            | 8/2q1p1k1/pb4p1/4p3/1Q6/2P2B2/1P4PP/7K w - - | Crafty Test Pos.237
// 236  | Success    | h7f5     | cp -25   | bm h7f5            | 7k/4p2B/1b6/p3p1p1/8/2P5/1P3qPP/1Q5K w - - | Crafty Test Pos.238
// 237  | Success    | f5e4     | cp -29   | bm f5e4            | 7k/4p3/8/p1b1pBp1/2P5/1P6/5qPP/1Q5K w - - | Crafty Test Pos.239
// 238  | Success    | e1g1     | cp 111   | bm e1g1            | r1b1kbnr/1pqp1ppp/p1n1p3/8/3NP3/2N5/PPP1BPPP/R1BQK2R w KQkq - | Crafty Test Pos.240
// 239  | Success    | a6d3     | cp 100   | bm a6d3            | 1r3rk1/2qp1ppp/B1p1pn2/8/1b2P3/4B3/PPP2PPP/R2Q1RK1 w - - | Crafty Test Pos.241
// 240  | Success    | f3e4     | cp 0     | bm f3e4            | 1r3rk1/5ppp/2p5/3nq3/5p2/1P1B1Q2/P1P3PP/R4RK1 w - - | Crafty Test Pos.242
// 241  | Success    | f4d4     | cp 110   | bm f4d4            | 5rk1/5p1p/6q1/3Q2p1/P4R2/1P6/2r3PP/4R1K1 w - - | Crafty Test Pos.243
// 242  | Failed     | b4a5     | cp -100  | bm d2c2            | 8/7p/2q2pk1/p4rp1/1P1Q4/8/2rR2PP/1R4K1 w - - | Crafty Test Pos.244
// 243  | Success    | a6a7     | cp 298   | bm a6a7            | 8/7p/P4p1k/6p1/1P6/8/1q3rPP/R5QK w - - | Crafty Test Pos.245
// 244  | Failed     | f1f6     | cp 429   | bm e6f7            | 8/7k/4Qp1p/6p1/1q6/8/7P/5R1K w - - | Crafty Test Pos.246
// 245  | Success    | e1g1     | cp 111   | bm e1g1            | r1b1kbnr/1pqp1ppp/p1n1p3/8/3NP3/2N5/PPP1BPPP/R1BQK2R w KQkq - | Crafty Test Pos.247
// 246  | Failed     | e4e5     | cp 57    | bm a1d1            | r1b1k1nr/2B2pbp/p1p1p3/3p4/4P3/2N5/PPP1BPP1/R4R1K w kq - | Crafty Test Pos.248
// 247  | Failed     | g3d6     | cp 23    | bm e4d5            | 2b2rk1/1r3pbp/p1p1p1n1/3p4/N3P3/2P3B1/PP2BPP1/3R1R1K w - - | Crafty Test Pos.249
// 248  | Success    | h1g1     | cp -23   | bm h1g1            | 2b3k1/5rbp/pBr1p1n1/3p1p2/N7/2P5/PP2BPP1/3RR2K w - - | Crafty Test Pos.250
// 249  | Success    | b6d7     | cp 57    | bm b6d7            | 1r3b2/1b2rk1p/pN2p1n1/B4p1B/2p5/1P6/P4PP1/3RRK2 w - - | Crafty Test Pos.251
// 250  | Success    | a4a5     | cp 69    | bm a4a5            | 1r6/4k2p/2b1p1n1/p4p1B/R7/1P6/5PP1/4RK2 w - - | Crafty Test Pos.252
// 251  | Failed     | g2g1     | cp 159   | bm e1e5            | 8/4k3/4p2R/1b3pnB/8/1r4P1/5PK1/4R3 w - - | Crafty Test Pos.253
// 252  | Failed     | e1f1     | cp 139   | bm c1a1            | 8/4k3/4p2R/3b1p2/4n3/6P1/1r2BP2/2R1K3 w - - | Crafty Test Pos.254
// 253  | Success    | d3b1     | cp 265   | bm d3b1            | 8/4k3/2n1p2R/8/8/1b1B2P1/5P2/1r3K2 w - - | Crafty Test Pos.255
// 254  | Success    | g4g5     | cp 365   | bm g4g5            | 8/5n2/2b2k2/4p3/5PPR/4K3/8/1B6 w - - | Crafty Test Pos.256
// 255  | Failed     | b2b3     | cp -10   | bm f3h4            |r1bqk2r/pp1n1ppp/2pb1n2/3pp3/8/3P1NP1/PPPNPPBP/R1BQ1RK1 w kq - | Crafty Test Pos.257
// 256  | Failed     | h4f5     | cp -107  | bm d1e2            |r2r2k1/pp3ppp/1qpbbn2/2n1p3/4P2N/5PP1/PPP3BP/R1BQRN1K w - - | Crafty Test Pos.258
// 257  | Success    | g5e3     | cp -156  | bm g5e3            |r2r2k1/pp3pp1/2p1bn1p/q3p1B1/4P1PN/2b2P2/P1N1Q1BP/R3R2K w - - | Crafty Test Pos.259
// 258  | Failed     | e3c4     | cp -20   | bm e1g1            |4r1k1/p4pp1/1ppr1n2/4pPB1/q7/4NP2/P3Q1BP/4R2K w - - | Crafty Test Pos.260
// 259  | Failed     | e1c3     | cp -61   | bm e3g4            |4rk2/p4pp1/1pp5/4pP2/3r1q2/4NP2/P5BP/4Q1RK w - - | Crafty Test Pos.261
// 260  | Success    | e3d5     | cp -225  | bm e3d5            |4r3/p1k2pp1/1p4q1/2p5/3r1p2/4N3/P4QBP/6RK w - - | Crafty Test Pos.262
// 261  | Failed     | f7f6     | cp -399  | bm f7f8            |8/p1kq1Q2/1p6/2p5/1r3p2/4r3/P5BP/6RK w - - | Crafty Test Pos.263
// 262  | Failed     | c7h7     | cp -355  | bm c7f7            |8/p1Q5/1p6/2p3k1/3q4/4r3/r5BP/5R1K w - - | Crafty Test Pos.264
// 263  | Success    | g1f3     | cp 25    | bm g1f3            |rnbqk2r/ppp2ppp/5n2/8/1bBP4/2N5/PP3PPP/R1BQK1NR w KQkq - | Crafty Test Pos.265
// 264  | Failed     | d1b3     | cp 41    | bm f2f4            |r2q1rk1/pppn1ppp/6b1/4P3/1bB3P1/2N1B2P/PP3P2/R2Q1RK1 w - - | Crafty Test Pos.266
// 265  | Failed     | b3f7     | cp -196  | bm c3e4            | r4r2/pp3pkp/1np5/4q3/1b4P1/1BN1BQ1P/PP6/5RK1 w - - | Crafty Test Pos.267
// 266  | Success    | h6f5     | cp -151  | bm h6f5            | 5r1k/Bp5p/5p1N/3pq3/1br3P1/7P/PP3Q2/5RK1 w - - | Crafty Test Pos.268
// 267  | Success    | a7b8     | cp -69   | bm a7b8            | 7k/Bp2r2p/7N/3pp3/6P1/P6P/1P4K1/4b3 w - - | Crafty Test Pos.269
// 268  | Failed     | f1f2     | cp -203  | bm f1e2            | 8/1p4kp/8/4N1P1/3pp3/bP5P/8/5K2 w - - | Crafty Test Pos.270
// 269  | Failed     | f1b5     | cp 135   | bm d1d2            | r1bqkb1r/pp3ppp/2nppn2/6B1/3NP3/2N5/PPP2PPP/R2QKB1R w KQkq - | Crafty Test Pos.271
// 270  | Failed     | b2b4     | cp 25    | bm f4d2            | r2qk2r/1p3pp1/p1b1pn1p/b2p4/4PB2/P1N2P2/1PP3PP/2KRQB1R w kq - | Crafty Test Pos.272
// 271  | Success    | d2d4     | cp -15   | bm d2d4            | r2qr1k1/5pp1/2b2n1p/pp6/2Bp3Q/P4P2/NPPR2PP/2K4R w - - | Crafty Test Pos.273
// 272  | Success    | b2c1     | cp -62   | bm b2c1            | 1r4k1/5pp1/2bR1n1p/p7/8/4rP2/NKP3PP/5B1R w - - | Crafty Test Pos.274
// 273  | Failed     | c6d6     | cp -75   | bm c6b6            | r5k1/5pp1/2R2n1p/8/p7/5PP1/2P1K2P/7B w - - | Crafty Test Pos.275
// 274  | Success    | g2h3     | cp -139  | bm g2h3            | 6k1/5pp1/7p/8/8/2K2PP1/2P3BP/6r1 w - - | Crafty Test Pos.276
// 275  | Success    | d3e3     | cp -171  | bm d3e3            | 6k1/5p2/8/8/4B2p/3K1P2/2P2r2/8 w - - | Crafty Test Pos.277
// 276  | Failed     | g2d5     | cp -333  | bm f2f3            | 6k1/6r1/8/5p2/5P2/8/2P2KBp/8 w - - | Crafty Test Pos.278
// 277  | Failed     | c1g5     | cp 39    | bm g2g3            | r2qkb1r/pp1npppp/2p2n2/5b2/2QP4/2N2N2/PP2PPPP/R1B1KB1R w KQkq - | Crafty Test Pos.279
// 278  | Success    | c3e4     | cp -28   | bm c3e4            | 3r1rk1/pp2bppp/1qp1pn2/5b2/3Pn3/2N1P1P1/PP2QPBP/R1BRN1K1 w - - | Crafty Test Pos.280
// 279  | Failed     | c1e3     | cp 1     | bm c1f4            | 2nr1rk1/pp2bppp/1qp1p1b1/8/P2PP3/1P1R1PP1/2N1Q1BP/R1B3K1 w - - | Crafty Test Pos.281
// 280  | Success    | d3e3     | cp -25   | bm d3e3            | 5rk1/pp2nppp/1q4b1/4p3/PP2P3/3QbPP1/2N3BP/3R2K1 w - - | Crafty Test Pos.282
// 281  | Success    | c3b3     | cp 12    | bm c3b3            | 8/p3nkpp/1p3p2/4p3/PP1qP3/2Q1NPP1/7P/6K1 w - - | Crafty Test Pos.283
// 282  | Success    | f3g4     | cp -117  | bm f3g4            | 5k2/p3n1p1/1p2qp2/4p3/PP2P1p1/5P1K/3Q2NP/8 w - - | Crafty Test Pos.284
// 283  | Failed     | e3f5     | cp -118  | bm e3c4            | 5k2/p3n1p1/1p3p2/8/PP1p2P1/4N3/6KP/8 w - - | Crafty Test Pos.285
// 284  | Failed     | d6c8     | cp -141  | bm d6f5            | 8/p3k1p1/1p1N1p2/1P6/P7/3K1n2/8/8 w - - | Crafty Test Pos.286
// 285  | Success    | e2e3     | cp 1     | bm e2e3            | rn1qk2r/pbp1bppp/1p2pn2/3p2B1/2PP4/P1N2N2/1P2PPPP/R2QKB1R w KQkq - | Crafty Test Pos.287
// 286  | Failed     | c1d1     | cp 31    | bm d4c5            |r2q1rk1/1b2bppp/pp3n2/2pp4/3PnB2/P1NBPN2/1PQ2PPP/2R2RK1 w - - | Crafty Test Pos.288
// 287  | Failed     | b2a2     | cp 12    | bm b2b3            |2r1qrk1/1b2bppp/p4n2/1p1p4/1PnN1B2/P1NBP3/1Q3PPP/2RR2K1 w - - | Crafty Test Pos.289
// 288  | Success    | c3e2     | cp 120   | bm c3e2            |3rqrk1/1b3pp1/p4b1p/1p3B2/1PnP1N2/PQN5/5PPP/2RR2K1 w - - | Crafty Test Pos.290
// 289  | Success    | d4e2     | cp -2    | bm d4e2            |3r2k1/5pp1/p2r3p/1p6/1PnN1q2/P1Q5/5PPP/3RR1K1 w - - | Crafty Test Pos.291
// 290  | Failed     | c1a1     | cp -113  | bm h2h3            |6k1/6p1/4p2p/1p6/1qn5/6N1/5PPP/2Q3K1 w - - | Crafty Test Pos.292
// 291  | Success    | b7e7     | cp -130  | bm b7e7            |8/1Q4pk/4p2p/1p1q4/2n5/7P/5PP1/5NK1 w - - | Crafty Test Pos.293
// 292  | Success    | d3b1     | cp -27   | bm d3b1            |8/6p1/4p1kp/1p3n2/3q4/3QN2P/5PP1/6K1 w - - | Crafty Test Pos.294
// 293  | Failed     | e2e3     | cp 77    | bm h2h3            |rnbqk2r/1p2ppbp/p1p2np1/2Pp4/3P1B2/2N2N2/PP2PPPP/R2QKB1R w KQkq - | Crafty TestPos.295
// 294  | Success    | f4e3     | cp -3    | bm f4e3            |r1bqnrk1/1p4bp/p1p3p1/2nPp3/5B2/2NB1N1P/PP3PP1/2RQK2R w K - | Crafty Test Pos.296
// 295  | Success    | g2f3     | cp -52   | bm g2f3            |r2q1rk1/1p4bp/p1p3p1/8/6Q1/4Bb1P/PP3PP1/2R2RK1 w - - | Crafty Test Pos.297
// 296  | Success    | d3d8     | cp 201   | bm d3d8            | 6k1/1p3rbp/p1p3p1/8/2Q5/3RBP1P/1q3P2/6K1 w - - | Crafty Test Pos.298
// 297  | Failed     | d8f8     | cp 235   | bm h4h5            | 3R1bk1/1p3r1p/p1p1Q1pB/8/5P1P/q7/5PK1/8 w - - | Crafty Test Pos.299
// 298  | Failed     | e2e4     | cp 42    | bm c1g5            | r1bqk2r/pp2ppbp/2np1np1/2p5/8/2NP1NP1/PPP1PPBP/R1BQ1RK1 w kq - | Crafty Test Pos.300
// 299  | Success    | g5e4     | cp 55    | bm g5e4            | 1r1q1r2/3bppkp/2np2p1/ppp3N1/4n3/P2P2P1/1PPQPPBP/R4RK1 w - - | Crafty Test Pos.301
// 300  | Failed     | a6b6     | cp 24    | bm c1a1            | 5r2/2q1ppk1/Rrnp2pp/2p5/1p2N1b1/3PP1P1/1PP2PBP/2R1Q1K1 w - - | Crafty Test Pos.302
// 301  | Success    | e4d2     | cp -41   | bm e4d2            | 2r3k1/1b2pp2/1qn3pp/2pp4/1p2N3/1P1PP1P1/R1P2PBP/Q5K1 w - - | Crafty Test Pos.303
// 302  | Success    | a2a8     | cp -17   | bm a2a8            | r2q2k1/1b3p2/2n3p1/2ppp2p/1p6/1P1PPNPP/R1P2PB1/Q5K1 w - - | Crafty Test Pos.304
// 303  | Success    | g5f6     | cp -37   | bm g5f6            | b7/3n1k2/5pp1/2ppp1P1/1p4N1/1P1PP3/2P2PB1/6K1 w - - | Crafty Test Pos.305
// 304  | Failed     | h3c8     | cp -60   | bm f4g5            | b7/8/5k2/2pp2p1/1p3P2/1P1P2KB/2P5/8 w - - | Crafty Test Pos.306
// 305  | Failed     | f1e2     | cp -45   | bm d3d4            | r1bqk2r/pp2bppp/2np1n2/2p1p3/Q3P3/2PP1N2/PP1N1PPP/R1B1KB1R w KQkq - | Crafty Test Pos.307
// 306  | Success    | c2b3     | cp -21   | bm c2b3            | 2rq1rk1/1p1bbppp/p2p1n2/3Pp3/4P3/5N2/PPQN1PPP/R1B2RK1 w - - | Crafty Test Pos.308
// 307  | Failed     | d2b3     | cp -61   | bm g2g4            | 2rq2k1/1p4pp/3p1b2/pbnPpr2/8/Q1R2N1P/PP1N1PP1/R1B3K1 w - - | Crafty Test Pos.309
// 308  | Failed     | g2f3     | cp -427  | bm d2f3            | 2q3k1/1p4pp/3p1b2/pb1Pp3/8/1P1nQr1P/P2N1PK1/R1B5 w - - | Crafty Test Pos.310
// 309  | Success    | h2h1     | cp -1337 | bm h2h1            | 6k1/1p1b2pp/3p4/p2Pbq2/P7/1P3p1P/5P1K/4R1N1 w - - | Crafty Test Pos.311
// 310  | Success    | e1g1     | cp 37    | bm e1g1            | r1b1kb1r/pp3ppp/2n1pn2/2pq4/3P4/2P2N2/PP2BPPP/RNBQK2R w KQkq - | Crafty Test Pos.312
// 311  | Failed     | a1c1     | cp 55    | bm d1b3            | r1bq1rk1/1p2bppp/p1n1p3/3n4/3P4/2N2NB1/PP2BPPP/R2Q1RK1 w - - | Crafty Test Pos.313
// 312  | Failed     | e2d3     | cp 30    | bm d4d5            | 2rqr1k1/1b2bppp/p1n1pn2/1p6/1P1P1B2/P1N2N2/Q3BPPP/2RR2K1 w - - | Crafty Test Pos.314
// 313  | Failed     | d7c7     | cp -62   | bm f3d2            | b3r1k1/3R1ppp/p1n5/1p6/1P6/4BN2/4rPPP/2R3K1 w - - | Crafty Test Pos.315
// 314  | Failed     | d7d6     | cp -34   | bm d7a7            | 4r1k1/3R1pp1/p3b2p/1p6/1n6/4P3/3N1KPP/1R6 w - - | Crafty Test Pos.316
// 315  | Success    | d8b8     | cp -21   | bm d8b8            | 3R4/1b3ppk/p6p/1pn5/8/4P1P1/3NK2P/8 w - - | Crafty Test Pos.317
// 316  | Success    | d2e4     | cp 105   | bm d2e4            | 8/6pk/3R3p/1p3p2/4n3/4K1Pb/3N3P/8 w - - | Crafty Test Pos.318
// 317  | Failed     | b6b7     | cp 85    | bm b6e6            | 8/6k1/1R5p/6p1/4p3/4KbP1/7P/8 w - - | Crafty Test Pos.319
// 318  | Failed     | c3d5     | cp 45    | bm c1g5            | r1bqkb1r/pp3ppp/2np1n2/1N2p3/4P3/2N5/PPP2PPP/R1BQKB1R w KQkq - | Crafty Test Pos.320
// 319  | Failed     | h5h7     | cp 148   | bm g1h1            |r2qkbr1/5p1p/p1npb3/1p1Np2Q/4Pp2/N2B4/PPP2PPP/R4RK1 w q - | Crafty Test Pos.321
// 320  | Success    | c4b5     | cp -289  | bm c4b5            |2r1kb2/5p1p/p1npb2q/1p1Np3/2P1P3/N2B1pP1/PP3P1P/R5RK w - - | Crafty Test Pos.322
// 321  | Success    | f5e4     | cp -337  | bm f5e4            |4kb2/5p1p/n2p4/3rpB2/8/5NP1/PP3q1P/R5RK w - - | Crafty Test Pos.323
// 322  | Failed     | b2b4     | cp -421  | bm e1f1            |3k1b2/7p/2Bp3q/2n1pp2/8/P4NP1/1P5P/4R2K w - - | Crafty Test Pos.324
// 323  | Success    | g1f3     | cp 31    | bm g1f3            |r1bqk2r/pp1nbppp/2p1pn2/3p4/2PP4/2N1P3/PPQB1PPP/R3KBNR w KQkq - | Crafty Test Pos.325
// 324  | Failed     | e3e4     | cp 18    | bm a1c1            |r2q1rk1/pb1n1pp1/1p2pn1p/2b5/8/2N1PN2/PPQBBPPP/R2R2K1 w - - | Crafty Test Pos.326
// 325  | Failed     | f2f4     | cp 21    | bm f2f3            |rqr3k1/1b3pp1/p3pn1p/1p2b3/7Q/P1N1P3/1P1BBPPP/2RR2K1 w - - | Crafty Test Pos.327
// 326  | Failed     | h2h4     | cp 20    | bm e2d3            |r1r5/5kp1/p3p2p/1p1b1p2/3BnP2/P3P3/1P2B1PP/2RR2K1 w - - | Crafty Test Pos.328
// 327  | Failed     | d3b5     | cp 23    | bm d3b1            |8/4nkp1/p3p2p/1b2Bp2/5P2/1P1BP3/6PP/6K1 w - - | Crafty Test Pos.329
// 328  | Failed     | e4e5     | cp 28    | bm g2g3            |8/2B2k2/2n1p1p1/5p1p/4PP2/1b5P/5KP1/1B6 w - - | Crafty Test Pos.330
// 329  | Success    | f1e2     | cp -18   | bm f1e2            |r1bqk2r/ppp2p1p/2np1np1/4p3/1bP5/2NPPN2/PP1B1PPP/R2QKB1R w KQkq - | Crafty TestPos.331
// 330  | Failed     | f3d5     | cp 65    | bm c4d5            |r2q1rk1/ppp2p2/2np1npp/3b4/2P5/P1B1PB2/1P3PPP/R2Q1RK1 w - - | Crafty Test Pos.332
// 331  | Failed     | d1c1     | cp 48    | bm f3f4            | 2r2rk1/p3qp2/2pp1npp/4n3/4P3/P1B2P2/1P2B1PP/R2Q1RK1 w - - | Crafty Test Pos.333
// 332  | Success    | c1g5     | cp -48   | bm c1g5            | rnbq1rk1/p1pp1ppp/1p2pn2/8/2PP4/P1Q5/1P2PPPP/R1B1KBNR w KQ - | Crafty Test Pos.334
// 333  | Success    | c3d4     | cp -45   | bm c3d4            | 2rq1rk1/p2n1pp1/bp1ppn1p/8/2Pp3B/P1QBPP2/1P2N1PP/3RK2R w K - | Crafty Test Pos.335
// 334  | Failed     | d1d8     | cp -2    | bm b3c4            | 2rr4/p4pk1/bp2p2p/8/2p1P3/PP2P3/4N1PP/3RK2R w K - | Crafty Test Pos.336
// 335  | Failed     | e2d3     | cp -114  | bm b3c3            | 8/5pk1/p3p2p/1p6/4r3/PR2P3/4K1PP/8 w - - | Crafty Test Pos.337
// 336  | Success    | d3d6     | cp -57   | bm d3d6            | 8/8/p3p3/1p3p1p/r1k5/P2RPKPP/8/8 w - - | Crafty Test Pos.338
// 337  | Success    | b5f5     | cp -239  | bm b5f5            | 8/8/8/pR3p1p/7P/1p1KP1P1/1k6/r7 w - - | Crafty Test Pos.339
// 338  | Success    | e1g1     | cp -30   | bm e1g1            | r1bqkb1r/pp3ppp/2n2n2/2pp4/3P4/5NP1/PP2PPBP/RNBQK2R w KQkq - | Crafty Test Pos.340
// 339  | Failed     | g2d5     | cp 1     | bm d1d5            | r4rk1/pp3ppp/1qn2b2/3b4/8/1N4P1/PP2PPBP/R2Q1RK1 w - - | Crafty Test Pos.341
// 340  | Failed     | d1d8     | cp 16    | bm f5e4            | 2rr2k1/ppq1npp1/7p/5Q1P/8/bN2P1P1/P4PB1/1R1R2K1 w - - | Crafty Test Pos.342
// 341  | Failed     | a2a3     | cp 22    | bm d4f5            | 6k1/2r1npp1/pb1r3p/1p5P/3NB3/3RP1P1/P4P2/3R2K1 w - - | Crafty Test Pos.343
// 342  | Failed     | d7d6     | cp 94    | bm d7c7            | 5k2/2rR1pp1/1b5p/p2B3P/1p6/4P1P1/P4PK1/8 w - - | Crafty Test Pos.344
// 343  | Failed     | g4g5     | cp 61    | bm g4h5            | 3b1k2/8/6p1/p6p/1p2PPP1/1B3K2/P7/8 w - - | Crafty Test Pos.345
// 344  | Failed     | c6d7     | cp 99    | bm e4d5            | 8/8/2B2P2/p1b1P1kp/1p2K3/8/P7/8 w - - | Crafty Test Pos.346
// 345  | Failed     | e8c6     | cp 20    | bm e6e7            | 4B3/5K2/4PP2/p5k1/1p1b4/7p/P7/8 w - - | Crafty Test Pos.347
// ====================================================================================================================================
// Successful: 162 (46 %)
// Failed:     183 (53 %)
// Skipped:    0   (0 %)
// Not tested: 0   (0 %)
// Test time: 1.728.274 ms
func _TestCraftyTests(t *testing.T) {
	ts, _ := NewTestSuite("testsets/crafty_test.epd", 5 * time.Second, 0)
	ts.RunTests()
}

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

package movegen

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/position"
)

// ///////////////////////////////////////////////////////////////
// Perft tests from https://www.chessprogramming.org/Perft_Results
// ///////////////////////////////////////////////////////////////

//noinspection GoImportUsedAsName
func TestStandardPerft(t *testing.T) {

	maxDepth := 5
	var perft Perft
	assert := assert.New(t)

	var results = [10][6]uint64{
		// @formatter:off
		// N             Nodes         Captures           EP          Checks           Mates
		{ 0,                 1,               0,           0,              0,              0 },
		{ 1,                20,               0,           0,              0,              0 },
		{ 2,               400,               0,           0,              0,              0 },
		{ 3,             8_902,              34,           0,             12,              0 },
		{ 4,           197_281,           1_576,           0,            469,              8 },
		{ 5,         4_865_609,          82_719,         258,         27_351,            347 },
		{ 6,       119_060_324,       2_812_008,       5_248,        809_099,         10_828 },
		{ 7,     3_195_901_860,     108_329_926,     319_617,     33_103_848,        435_767 },
		{ 8,    84_998_978_956,   3_523_740_106,   7_187_977,    968_981_593,      9_852_036 },
		{ 9, 2_439_530_234_167, 125_208_536_153, 319_496_827, 36_095_901_903,    400_191_963 }}
	// @formatter:on

	for i := 1; i <= maxDepth; i++ {
		perft.StartPerft(position.StartFen, i, false)
		assert.Equal(results[i][1], perft.Nodes)
		assert.Equal(results[i][2], perft.CaptureCounter)
		assert.Equal(results[i][3], perft.EnpassantCounter)
		assert.Equal(results[i][4], perft.CheckCounter)
		assert.Equal(results[i][5], perft.CheckMateCounter)
	}
}

// Performing PERFT Test for Depth 6
// -----------------------------------------
// Time         : 28.532 ms
// NPS          :   - go test -coverprofile=coverage.txt -covermode=atomic172.724 nps
// Results:
//   Nodes     : 119.060.324
//   Captures  : 2.812.008
//   EnPassant : 5.248
//   Checks    : 809.099
//   CheckMates: 10.828
//   Castles   : 0
//   Promotions: 0
// -----------------------------------------
// Finished PERFT Test for Depth 6
//noinspection GoImportUsedAsName
func TestStandardPerftOd(t *testing.T) {

	maxDepth := 5
	var perft Perft
	assert := assert.New(t)

	var results = [10][6]uint64{
		// @formatter:off
		// N             Nodes         Captures           EP          Checks           Mates
		{ 0,                 1,               0,           0,              0,              0 },
		{ 1,                20,               0,           0,              0,              0 },
		{ 2,               400,               0,           0,              0,              0 },
		{ 3,             8_902,              34,           0,             12,              0 },
		{ 4,           197_281,           1_576,           0,            469,              8 },
		{ 5,         4_865_609,          82_719,         258,         27_351,            347 },
		{ 6,       119_060_324,       2_812_008,       5_248,        809_099,         10_828 },
		{ 7,     3_195_901_860,     108_329_926,     319_617,     33_103_848,        435_767 },
		{ 8,    84_998_978_956,   3_523_740_106,   7_187_977,    968_981_593,      9_852_036 },
		{ 9, 2_439_530_234_167, 125_208_536_153, 319_496_827, 36_095_901_903,    400_191_963 }}
	// @formatter:on

	for i := 1; i <= maxDepth; i++ {
		perft.StartPerft(position.StartFen, i, true)
		assert.Equal(results[i][1], perft.Nodes)
		assert.Equal(results[i][2], perft.CaptureCounter)
		assert.Equal(results[i][3], perft.EnpassantCounter)
		assert.Equal(results[i][4], perft.CheckCounter)
		assert.Equal(results[i][5], perft.CheckMateCounter)
	}
}

//noinspection GoImportUsedAsName
func TestKiwipetePerft(t *testing.T) {

	maxDepth := 4
	var perft Perft
	assert := assert.New(t)

	var kiwipete = [10][8]uint64{
		// @formatter:off
		// N             Nodes         Captures           EP          Checks           Mates		Castles		Promotions
		{ 0,                 1,               0,           0,              0,              0, 			  0,             0 },
		{ 1,                48,               8,           0,              0,              0, 			  2,             0 },
		{ 2,             2_039,             351,           1,              3,              0,		     91,             0 },
		{ 3,            97_862,          17_102,          45,            993,              1, 	      3_162,             0 },
		{ 4,         4_085_603,         757_163,       1_929,         25_523,             43, 		128_013,        15_172 },
		{ 5,       193_690_690,      35_043_416,      73_365,      3_309_887,         30_171, 	  4_993_637,         8_392 },
		{ 6,     8_031_647_685,   1_558_445_089,   3_577_504,     92_238_050,        360_003, 	184_513_607,    56_627_920 }}

	for depth := 1; depth <= maxDepth; depth++ {
		perft.StartPerft("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ", depth, true)
		assert.Equal(kiwipete[depth][1], perft.Nodes)
		assert.Equal(kiwipete[depth][2], perft.CaptureCounter)
		assert.Equal(kiwipete[depth][3], perft.EnpassantCounter)
		assert.Equal(kiwipete[depth][4], perft.CheckCounter)
		assert.Equal(kiwipete[depth][5], perft.CheckMateCounter)
		assert.Equal(kiwipete[depth][6], perft.CastleCounter)
		assert.Equal(kiwipete[depth][7], perft.PromotionCounter)
	}
}

//noinspection GoImportUsedAsName
func TestMirrorPerft(t *testing.T) {

	maxDepth := 5
	var perft Perft
	assert := assert.New(t)

	var mirrorPerft = [10][8]uint64{
		// @formatter:off
		// N             Nodes         Captures           EP          Checks           Mates		Castles		Promotions
		{ 0,                 1,               0,           0,              0,              0, 			  0,             0 },
		{ 1,     		     6,               0,           0,              0,              0, 	          0,             0 },
		{ 2,     		   264,              87,           0,             10,              0,	          6,            48 },
		{ 3,              9467,            1021,           4,             38,             22, 			  0,           120 },
		{ 4,            422333,          131393,           0,          15492,              5, 		   7795,         60032 },
		{ 5,          15833292,         2046173,        6512,         200568,          50562, 	          0,        329464 },
		{ 6,         706045033,       210369132,   		 212,       26973664,          81076, 	   10882006,      81102984 }}


	// white
	for depth := 1; depth <= maxDepth; depth++ {
		perft.StartPerft("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq -", depth, false)
		assert.Equal(mirrorPerft[depth][1], perft.Nodes)
		assert.Equal(mirrorPerft[depth][2], perft.CaptureCounter)
		assert.Equal(mirrorPerft[depth][3], perft.EnpassantCounter)
		assert.Equal(mirrorPerft[depth][4], perft.CheckCounter)
		assert.Equal(mirrorPerft[depth][5], perft.CheckMateCounter)
		assert.Equal(mirrorPerft[depth][6], perft.CastleCounter)
		assert.Equal(mirrorPerft[depth][7], perft.PromotionCounter)
	}

	// mirrored
	for depth := 1; depth <= maxDepth; depth++ {
		perft.StartPerft("r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ -", depth, false)
		assert.Equal(mirrorPerft[depth][1], perft.Nodes)
		assert.Equal(mirrorPerft[depth][2], perft.CaptureCounter)
		assert.Equal(mirrorPerft[depth][3], perft.EnpassantCounter)
		assert.Equal(mirrorPerft[depth][4], perft.CheckCounter)
		assert.Equal(mirrorPerft[depth][5], perft.CheckMateCounter)
		assert.Equal(mirrorPerft[depth][6], perft.CastleCounter)
		assert.Equal(mirrorPerft[depth][7], perft.PromotionCounter)
	}
}


//noinspection GoImportUsedAsName
func TestPos5Perft(t *testing.T) {

	maxDepth := 4
	var perft Perft
	assert := assert.New(t)

	var kiwipete = [10][2]uint64{
		// @formatter:off
		// N             Nodes
		{ 0,                 1 },
		{ 1,                44 },
		{ 2,             1_486 },
		{ 3,            62_379 },
		{ 4,         2_103_487 },
		{ 5,        89_941_194 }}

	for depth := 1; depth <= maxDepth; depth++ {
		perft.StartPerft("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ -", depth, false)
		assert.Equal(kiwipete[depth][1], perft.Nodes)
	}
}



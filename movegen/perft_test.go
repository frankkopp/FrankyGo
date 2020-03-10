/*
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

	. "github.com/frankkopp/FrankyGo/types"
)

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

//noinspection GoImportUsedAsName
func Test_StandardPerft(t *testing.T) {
	Init()
	maxDepth := 6
	var perft Perft
	assert := assert.New(t)

	for i := 1; i <= maxDepth; i++ {
		perft.StartPerft(StartFen, i)
		assert.Equal(results[i][1], perft.Nodes)
		assert.Equal(results[i][2], perft.CaptureCounter)
		assert.Equal(results[i][3], perft.EnpassantCounter)
		assert.Equal(results[i][4], perft.CheckCounter)
		assert.Equal(results[i][5], perft.CheckMateCounter)
	}
}


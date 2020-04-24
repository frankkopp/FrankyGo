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

package types

import "testing"

func TestRankType(t *testing.T) {
	tests := []struct {
		value    Rank
		expected int
	}{
		{Rank1, 0},
		{Rank8, 7},
		{RankNone, 8},
		{Rank(100), 100},
	}
	var got int
	for _, test := range tests {
		got = int(test.value)
		if test.expected != got {
			t.Errorf("rank %s == %d expected. Got %d", test.value.String(), test.expected, got)
		} else {
			t.Logf("rank %s == %d", test.value.String(), got)
		}
	}
}

func TestValidRank(t *testing.T) {
	tests := []struct {
		value    Rank
		expected bool
	}{
		{Rank1, true},
		{Rank8, true},
		{RankNone, false},
		{Rank(100), false},
	}
	var got bool
	for _, test := range tests {
		got = test.value.IsValid()
		if test.expected != got {
			t.Errorf("rank.valid(%s) %t expected. Got %t", test.value.String(), test.expected, got)
		} else {
			t.Logf("rank.valid(%s) == %t", test.value.String(), got)
		}
	}
}

func TestRankStr(t *testing.T) {
	tests := []struct {
		value    Rank
		expected string
	}{
		{Rank1, "1"},
		{Rank8, "8"},
		{RankNone, "-"},
		{Rank(100), "-"},
	}
	var got string
	for _, test := range tests {
		got = test.value.String()
		if test.expected != got {
			t.Errorf("rank label %s is expected. Got %s", test.expected, got)
		} else {
			t.Logf("rank label %s is %s", test.value.String(), got)
		}
	}
}

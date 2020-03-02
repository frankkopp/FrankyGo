/*
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, To any person obtaining a copy
 * of this software and associated documentation files (the "Software"), To deal
 * in the Software without restriction, including without limitation the rights
 * To use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and To permit persons To whom the Software is
 * furnished To do so, subject To the following conditions:
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
			t.Errorf("rank %s == %d expected. Got %d", test.value.Str(), test.expected, got)
		} else {
			t.Logf("rank %s == %d", test.value.Str(), got)
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
			t.Errorf("rank.valid(%s) %t expected. Got %t", test.value.Str(), test.expected, got)
		} else {
			t.Logf("rank.valid(%s) == %t", test.value.Str(), got)
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
		got = test.value.Str()
		if test.expected != got {
			t.Errorf("rank label %s is expected. Got %s", test.expected, got)
		} else {
			t.Logf("rank label %s is %s", test.value.Str(), got)
		}
	}
}

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

package types

import "testing"

func TestSquareType(t *testing.T) {
	tests := []struct {
		value    Square
		expected int
	}{
		{SqA1, 0},
		{SqH8, 63},
		{SqNone, 64},
		{Square(100), 100},
	}
	var got int
	for _, test := range tests {
		got = int(test.value)
		if test.expected != got {
			t.Errorf("square %s == %d expected. Got %d", test.value.str(), test.expected, got)
		} else {
			t.Logf("square %s == %d", test.value.str(), got)
		}
	}
}

func TestValidSquare(t *testing.T) {
	tests := []struct {
		value    Square
		expected bool
	}{
		{SqA1, true},
		{SqH8, true},
		{SqNone, false},
		{Square(100), false},
	}
	var got bool
	for _, test := range tests {
		got = test.value.isValid()
		if test.expected != got {
			t.Errorf("square.valid(%s) %t expected. Got %t", test.value.str(), test.expected, got)
		} else {
			t.Logf("square.valid(%s) == %t", test.value.str(), got)
		}
	}
}

func TestSquareStr(t *testing.T) {
	tests := []struct {
		value    Square
		expected string
	}{
		{SqA1, "a1"},
		{SqH8, "h8"},
		{SqNone, "--"},
		{Square(100), "--"},
	}
	var got string
	for _, test := range tests {
		got = test.value.str()
		if test.expected != got {
			t.Errorf("square label %s is expected. Got %s", test.expected, got)
		} else {
			t.Logf("square label %s is %s", test.value.str(), got)
		}
	}
}

func TestSquareFromFileRank(t *testing.T) {
	tests := []struct {
		file   File
		rank   Rank
		square Square
	}{
		{FileA, Rank1, SqA1},
		{FileH, Rank8, SqH8},
		{FileNone, RankNone, SqNone},
		{FileA, Rank(50), SqNone},
	}
	var got Square
	for _, test := range tests {
		got = squareOf(test.file, test.rank)
		if test.square != got {
			t.Errorf("squareOf(%s, %s) == %s is expected. Got %s", test.file.str(), test.rank.str(), test.square.str(), got.str())
		} else {
			t.Logf("squareOf(%s, %s) == %s", test.file.str(), test.rank.str(), got.str())
		}
	}
}

func TestSquareDir(t *testing.T) {
	tests := []struct {
		square   Square
		dir      Direction
		expected Square
	}{
		{SqA1, North, SqA2},
		{SqA1, East, SqB1},
		{SqA1, South, SqNone},
		{SqA1, West, SqNone},
		{SqH8, North, SqNone},
		{SqH8, East, SqNone},
		{SqH8, South, SqH7},
		{SqH8, West, SqG8},
	}
	var got Square
	for _, test := range tests {
		got = test.square.to(test.dir)
		if test.expected != got {
			t.Errorf("Square to %s of %s should be %s. Is %s",
				test.dir.str(), test.square.str(), test.expected.str(), got.str())
		} else {
			t.Logf("Square to %s of %s is %s.", test.dir.str(), test.square.str(), got.str())
		}
	}
}

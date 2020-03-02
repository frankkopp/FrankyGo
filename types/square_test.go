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
			t.Errorf("square %s == %d expected. Got %d", test.value.Str(), test.expected, got)
		} else {
			t.Logf("square %s == %d", test.value.Str(), got)
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
		got = test.value.IsValid()
		if test.expected != got {
			t.Errorf("square.valid(%s) %t expected. Got %t", test.value.Str(), test.expected, got)
		} else {
			t.Logf("square.valid(%s) == %t", test.value.Str(), got)
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
		got = test.value.Str()
		if test.expected != got {
			t.Errorf("square label %s is expected. Got %s", test.expected, got)
		} else {
			t.Logf("square label %s is %s", test.value.Str(), got)
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
		got = SquareOf(test.file, test.rank)
		if test.square != got {
			t.Errorf("SquareOf(%s, %s) == %s is expected. Got %s", test.file.Str(), test.rank.Str(), test.square.Str(), got.Str())
		} else {
			t.Logf("SquareOf(%s, %s) == %s", test.file.Str(), test.rank.Str(), got.Str())
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
		got = test.square.To(test.dir)
		if test.expected != got {
			t.Errorf("Square To %s of %s should be %s. Is %s",
				test.dir.Str(), test.square.Str(), test.expected.Str(), got.Str())
		} else {
			t.Logf("Square To %s of %s is %s.", test.dir.Str(), test.square.Str(), got.Str())
		}
	}
}

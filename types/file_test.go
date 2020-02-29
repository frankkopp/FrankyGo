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

func TestFileType(t *testing.T) {
	tests := []struct {
		value    File
		expected int
	}{
		{FileA, 0},
		{FileH, 7},
		{FileNone, 8},
		{File(100), 100},
	}
	var got int
	for _, test := range tests {
		got = int(test.value)
		if test.expected != got {
			t.Errorf("file %s == %d expected. Got %d", test.value.str(), test.expected, got)
		} else {
			t.Logf("file %s == %d", test.value.str(), got)
		}
	}
}

func TestValidFile(t *testing.T) {
	tests := []struct {
		value    File
		expected bool
	}{
		{FileA, true},
		{FileH, true},
		{FileNone, false},
		{File(100), false},
	}
	var got bool
	for _, test := range tests {
		got = test.value.isValid()
		if test.expected != got {
			t.Errorf("file.valid(%s) %t expected. Got %t", test.value.str(), test.expected, got)
		} else {
			t.Logf("file.valid(%s) == %t", test.value.str(), got)
		}
	}
}

func TestFileStr(t *testing.T) {
	tests := []struct {
		value    File
		expected string
	}{
		{FileA, "a"},
		{FileH, "h"},
		{FileNone, "-"},
		{File(100), "-"},
	}
	var got string
	for _, test := range tests {
		got = test.value.str()
		if test.expected != got {
			t.Errorf("file label %s is expected. Got %s", test.expected, got)
		} else {
			t.Logf("file label %s is %s", test.value.str(), got)
		}
	}
}

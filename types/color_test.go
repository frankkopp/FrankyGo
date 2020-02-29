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

func TestColorType(t *testing.T) {
	tests := []struct {
		value    Color
		expected int
	}{
		{White, 0},
		{Black, 1},
	}
	var got int
	for _, test := range tests {
		got = int(test.value)
		if test.expected != got {
			t.Errorf("Color %s == %d expected. Got %d", test.value.str(), test.expected, got)
		} else {
			t.Logf("Color %s == %d", test.value.str(), got)
		}
	}
}

func TestColorFlip(t *testing.T) {
	tests := []struct {
		value    Color
		expected Color
	}{
		{White, Black},
		{Black, White},
	}
	for _, test := range tests {
		got := test.value.flip()
		if test.expected != got {
			t.Errorf("Color %s.flip() == %s expected. Got %s", test.value.str(), test.expected.str(), got.str())
		} else {
			t.Logf("Color %s.flip() == %s", test.value.str(), got.str())
		}
	}
}


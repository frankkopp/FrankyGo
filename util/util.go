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

package util

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Abs - non branching Abs function for determine the absolute value of an int64
func Abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

// Abs64 - non branching Abs function for determine the absolute value of an int64
func Abs64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

// Min returns the smaller of the given integers
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Min64 returns the smaller of the given 64-bit integers
func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Max returns the bigger of the given integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Max64 returns the bigger of the given 64-bit integers
func Max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// TimeTrack is convenient way to measure timings of function.
// Usage: defer util.TimeTrack(time.Now(), "some text")
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	_, _ = out.Printf("%s took %d ns\n", name, elapsed.Nanoseconds())
}
var out = message.NewPrinter(language.German)

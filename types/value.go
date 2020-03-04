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

// Value represents the positional value of a chess position
type Value int16

// Constants for values
const (
	ValueZero               Value = 0
	ValueDraw               Value = 0
	ValueOne                Value = 1
	ValueInf                Value = 15_000
	ValueNA                 Value = -ValueInf - 1
	ValueMax                Value = 10_000
	ValueMin                Value = - ValueMax
	ValueCheckMate          Value = ValueMax
	ValueCheckMateThreshold Value = ValueCheckMate - MaxDepth - 1
)

// IsValid checks if value is within valid range (between Min and Max)
func (v Value) IsValid() bool {
	return v >= ValueMin && v <= ValueMax
}

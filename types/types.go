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

// Package types contains various user defined data types and their corresponding
// functionality we need for the chess engine.
package types

// MilliSec used for time value in milli sec
// Could be large therefore 64-bit
type MilliSec uint64

// Key is used for zobrist keys in chess positions.
// Zobrist keys need all 64 bits for distribution
type Key uint64

//Move 16 bits for move encoding - 16 bits for sort value
type Move uint32

// SqLength number of squares on a board
const SqLength int = 64

var initialized = false

// Init initializes pre computed data structures e.g. bitboards, etc.
// Keeps an initialized flag To avoid multiple executions.
func Init() {
	if initialized {
		return
	}

	// bitboards
	initBb()

	initialized = true
}

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

// Types are various user defined data types and their corresponding
// functionality we need for the chess engine.
package types

type MilliSec uint64 // could be large
type Key uint64      // zobrist key need all 64 bits for distribution
type Move uint32     // 16 bits for move encoding - 16 bits for sort value

// Global constants
const SqLength int = 64

var initialized = false

// Initialize pre computed data structure e.g. bitboards, etc.
// Keeps an initialized flag to avoid multiple executions.
func Init() {
	if initialized {
		return
	}

	// bitboards
	initBb()

	initialized = true
}

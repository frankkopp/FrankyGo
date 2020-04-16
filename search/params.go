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

package search

import (
	"math"
)

// This file contain data structures and functions to support the search with
// static or pre-computed parameters. Mostly for params too complex to be
// part of the search configuration

// lmrR is a lookup table for late move reductions in the dimensions
// depth and moves searched
var lmrR [32][64]int

// LmrReduction returns the search depth reduction for LMR
// depended on depth and moves searched
func LmrReduction(depth int, movesSearched int) int {
	if depth >= 32 || movesSearched >= 64 {
		return lmrR[31][63]
	}
	return lmrR[depth][movesSearched]
}

// prepare the pre-computed values
func init() {
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			switch {
			case i <= 3:
				lmrR[i][j] = 1
			case j <= 3:
				lmrR[i][j] = 1
			default:
				lmrR[i][j] = int(math.Round(((float64(i) * 0.7) * (float64(j) * 0.005)) + float64(1.0)))
			}
		}
	}
	// printLmr()
}

func printLmr() {
	for i := 3; i < 32; i++ {
		for j := 3; j < 64; j++ {
			out.Printf("%2d ", lmrR[i][j])
		}
		out.Println()
	}
}

var lmp [16]int

func init() {
	for i := 0; i < 16; i++ {
		// from Crafty
		lmp[i] = 3 + int(math.Pow(float64(i) + 0.5, 1.9))
		// out.Printf("%2d ", lmp[i])
	}
}

// LmpMovesSearched returns a depth dependent value for moves searched
// for late Move Prunings
func LmpMovesSearched(depth int) int {
	if depth >= 16 {
		return lmp[15]
	}
	return lmp[depth]
}

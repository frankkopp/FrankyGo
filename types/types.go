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

// Package types contains various user defined data types and their corresponding
// functionality we need for the chess engine.
// Many of these would be perfect enum candidates but GO does not provide enums
package types

import (
	"github.com/frankkopp/FrankyGo/logging"
)

var log = logging.GetLog("types")

var initialized = false

// Init initializes pre computed data structures e.g. bitboards, etc.
// Keeps an initialized flag to avoid multiple executions.
func init() {
	if initialized {
		return
	}
	log.Debug("Initializing data types")
	// bitboards
	initBb()
	// pos values
	initPosValues()
	initialized = true
}

const (
	// SqLength number of squares on a board
	SqLength int = 64

	// MaxDepth max search depth
	MaxDepth = 128

	// MaxMoves max number of moves for a game
	MaxMoves = 512

	// KB = 1.024 bytes
	KB uint64 = 1024

	// MB = KB * KB
	MB uint64 = KB * KB

	// GB = KB * MB
	GB uint64 = KB * MB

	// GamePhaseMax maximum game phase value. Game phase is used to
	// determine if we are in the beginning or end phase of a chess game
	// Game phase is calculated be the number of officers on the board
	// with this maximum
	GamePhaseMax = 24
)

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

package search

import (
	"time"

	"github.com/frankkopp/FrankyGo/moveslice"
)

// SearchLimits is data structure to hold all information about how
// a search of the chess games shall be controlled.
// Search needs to read these an determine the necessary limits.
// E.g. time controlled game or not
type SearchLimits struct {
	//  time control
	WhiteTime   time.Duration
	BlackTime   time.Duration
	WhiteInc    time.Duration
	BlackInc    time.Duration
	MoveTime    time.Duration
	MovesToGo   int
	TimeControl bool

	// extra limits
	Depth int
	Nodes int64
	Moves moveslice.MoveSlice

	// no time control
	Mate     int
	Ponder   bool
	Infinite bool
}

// NewSearchLimits creates a new empty SearchLimits
// instance and returns a pointer to it
func NewSearchLimits() *SearchLimits {
	return &SearchLimits{}
}

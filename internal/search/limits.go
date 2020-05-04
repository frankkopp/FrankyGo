//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package search

import (
	"time"

	"github.com/frankkopp/FrankyGo/internal/moveslice"
)

// Limits is data structure to hold all information about how
// a search of the chess games shall be controlled.
// Search needs to read these an determine the necessary limits.
// E.g. time controlled game or not
type Limits struct {
	// no time control
	Infinite bool
	Ponder   bool
	Mate     int

	// extra limits
	Depth int
	Nodes uint64
	Moves moveslice.MoveSlice

	//  time control
	TimeControl bool
	WhiteTime   time.Duration
	BlackTime   time.Duration
	WhiteInc    time.Duration
	BlackInc    time.Duration
	MoveTime    time.Duration

	// parameter
	MovesToGo   int
}

// NewSearchLimits creates a new empty Limits
// instance and returns a pointer to it
func NewSearchLimits() *Limits {
	return &Limits{}
}

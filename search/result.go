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

	. "github.com/frankkopp/FrankyGo/types"
)

// //////////////////////////////////////////////////////
// Result
// //////////////////////////////////////////////////////

// Result stores the result of a search. If BestMove is not MoveNone
// it can be assumed that all values are valid.
type Result struct {
	BestMove    Move
	PonderMove  Move
	searchTime  time.Duration
	searchDepth int
	extraDepth  int
}

func (searchResult *Result) String() string {
	return out.Sprintf("best move = %s, ponder move = %s, search time = %d ms, search dept = %d/%d",
		searchResult.BestMove.StringUci(), searchResult.PonderMove.StringUci(), searchResult.searchTime.Milliseconds(),
		searchResult.searchDepth, searchResult.extraDepth)
}

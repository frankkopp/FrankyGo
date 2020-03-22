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

// Package uciInterface defines the functions necessary to be implemented
// in a class to be used as a uci handler for the search class.
// This is necessary as GO does not allow circular imports.
// uci is importing Search to hold an instance of Search and Search needs
// a call back reference to a uci handler to be able to send UCI
// information to the UCI ui.
package uciInterface

import (
	"time"

	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/types"
)

// UciDriver the defines an interface for the search to be able to send
// uci protocol messages through a uciHandler which implements this interface
type UciDriver interface{
	SendReadyOk()
	SendInfoString(info string)
	SendIterationEndInfo(depth int, seldepth int, value types.Value, nodes uint64, nps uint64, time time.Duration, pv moveslice.MoveSlice)
	SendAspirationResearchInfo(depth int, seldepth int, value types.Value, valueType types.ValueType, nodes uint64, nps uint64, time time.Duration, pv moveslice.MoveSlice)
	SendCurrentRootMove(currMove types.Move, moveNumber int)
	SendSearchUpdate(depth int, seldepth int, nodes uint64, nps uint64, time time.Duration, hashfull int)
	SendCurrentLine(moveList moveslice.MoveSlice)
	SendResult(bestMove types.Move, ponderMove types.Move)
}

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
)

var logTest = logging.GetLog("test")

func TestSearch_IsReady(t *testing.T) {
	search := NewSearch()
	search.IsReady()
}

func TestSetupTimeControl(t *testing.T) {
	s := NewSearch()
	p := position.NewPosition()
	sl := &Limits{
		Infinite:    false,
		Ponder:      false,
		Mate:        0,
		Depth:       0,
		Nodes:       0,
		Moves:       nil,
		TimeControl: true,
		WhiteTime:   60 * time.Second,
		BlackTime:   60 * time.Second,
		WhiteInc:    2 * time.Second,
		BlackInc:    2 * time.Second,
		MoveTime:    0,
		MovesToGo:   20,
	}
	timeLimit := s.setupTimeControl(p,sl)
	assert.EqualValues(t, 4500, timeLimit.Milliseconds())

	p = position.NewPosition()
	sl = &Limits{
		Infinite:    false,
		Ponder:      false,
		Mate:        0,
		Depth:       0,
		Nodes:       0,
		Moves:       nil,
		TimeControl: true,
		WhiteTime:   60 * time.Second,
		BlackTime:   60 * time.Second,
		WhiteInc:    2 * time.Second,
		BlackInc:    2 * time.Second,
		MoveTime:    0,
		MovesToGo:   0,
	}
	timeLimit = s.setupTimeControl(p,sl)
	assert.EqualValues(t, 3150, timeLimit.Milliseconds())

	// game phase 0
	p = position.NewPositionFen("8/2P1P1P1/3PkP2/8/4K3/8/8/8 w - - 0 1")
	sl = &Limits{
		Infinite:    false,
		Ponder:      false,
		Mate:        0,
		Depth:       0,
		Nodes:       0,
		Moves:       nil,
		TimeControl: true,
		WhiteTime:   60 * time.Second,
		BlackTime:   60 * time.Second,
		WhiteInc:    0 * time.Second,
		BlackInc:    0 * time.Second,
		MoveTime:    0,
		MovesToGo:   0,
	}
	timeLimit = s.setupTimeControl(p,sl)
	assert.EqualValues(t, 5400, timeLimit.Milliseconds())
}

func TestWaitWhileSearching(t *testing.T) {
	search := NewSearch()
	p := position.NewPosition()
	sl := NewSearchLimits()
	start := time.Now()
	// FIXME: Prototype
	search.StartSearch(*p, *sl)
	logTest.Debug("Search started...waiting to finish")
	search.WaitWhileSearching()
	logTest.Debug("Search finished")
	elapsed := time.Since(start)
	out.Printf("Time %d ms\n", elapsed.Milliseconds())
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(2_000))
}

func TestIsSearching(t *testing.T) {
	search := NewSearch()
	p := position.NewPosition()
	sl := NewSearchLimits()
	start := time.Now()
	// FIXME: Prototype
	search.StartSearch(*p, *sl)
	logTest.Debug("Check searching in 1 sec")
	time.Sleep(time.Second)
	assert.True(t, search.IsSearching())
	logTest.Debugf("Is searching = %v", search.IsSearching())
	search.WaitWhileSearching()
	elapsed := time.Since(start)
	out.Printf("Time %d ms\n", elapsed.Milliseconds())
	assert.False(t, search.IsSearching())
	logTest.Debugf("Is searching = %v", search.IsSearching())
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(2_000))
}

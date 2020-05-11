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

package evaluator

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

func TestSize(t *testing.T) {
	sizeof := unsafe.Sizeof(cacheEntry{})
	out.Println(sizeof)
	assert.EqualValues(t, 16, sizeof)
}

func TestNewPawnCache(t *testing.T) {
	pc := newPawnCache()
	assert.EqualValues(t, 0, pc.len())
	assert.EqualValues(t, 0, pc.hits)
	assert.EqualValues(t, 0, pc.misses)
	assert.EqualValues(t, 0, pc.replace)
}

func TestPutGet(t *testing.T) {
	pc := newPawnCache()
	assert.EqualValues(t, 0, pc.len())
	assert.EqualValues(t, 0, pc.hits)
	assert.EqualValues(t, 0, pc.misses)
	assert.EqualValues(t, 0, pc.replace)

	p := position.NewPosition()

	pc.put(p.PawnKey(), &Score{
		MidGameValue: 1,
		EndGameValue: 11,
	})
	assert.EqualValues(t, 1, pc.len())
	assert.EqualValues(t, 0, pc.hits)
	assert.EqualValues(t, 0, pc.misses)
	assert.EqualValues(t, 0, pc.replace)

	p.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	pc.put(p.PawnKey(), &Score{
		MidGameValue: 2,
		EndGameValue: 22,
	})
	assert.EqualValues(t, 2, pc.len())
	assert.EqualValues(t, 0, pc.hits)
	assert.EqualValues(t, 0, pc.misses)
	assert.EqualValues(t, 0, pc.replace)

	// hit
	e := pc.getEntry(p.PawnKey())
	assert.NotNil(t, e)
	assert.EqualValues(t, 2, e.score.MidGameValue)
	assert.EqualValues(t, 22, e.score.EndGameValue)
	assert.EqualValues(t, 1, pc.hits)
	assert.EqualValues(t, 0, pc.misses)

	p.UndoMove()

	// hit
	e = pc.getEntry(p.PawnKey())
	assert.NotNil(t, e)
	assert.EqualValues(t, 1, e.score.MidGameValue)
	assert.EqualValues(t, 11, e.score.EndGameValue)
	assert.EqualValues(t, 2, pc.hits)
	assert.EqualValues(t, 0, pc.misses)

	// miss
	p.DoMove(CreateMove(SqD2, SqD4, Normal, PtNone))
	e = pc.getEntry(p.PawnKey())
	assert.Nil(t, e)
	assert.EqualValues(t, 2, pc.hits)
	assert.EqualValues(t, 1, pc.misses)

	pc.clear()
	assert.EqualValues(t, 0, pc.len())
	assert.EqualValues(t, 0, pc.hits)
	assert.EqualValues(t, 0, pc.misses)
	assert.EqualValues(t, 0, pc.replace)
}

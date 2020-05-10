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

	"github.com/stretchr/testify/assert"

	. "github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

func TestEvalPiecePawnsCache(t *testing.T) {
	Settings.Eval.UsePawnEval = true
	Settings.Eval.UsePawnCache = true

	e := NewEvaluator()
	Settings.Eval.Tempo = 0
	p := position.NewPosition()
	var score *Score
	e.InitEval(p)

	assert.EqualValues(t, 0, e.pawnCache.len())
	assert.EqualValues(t, 0, e.pawnCache.hits)
	assert.EqualValues(t, 0, e.pawnCache.misses)

	score = e.evaluatePawns()
	assert.EqualValues(t, 1, e.pawnCache.len())
	assert.EqualValues(t, 0, e.pawnCache.hits)
	assert.EqualValues(t, 1, e.pawnCache.misses)

	score2 := e.evaluatePawns()
	assert.EqualValues(t, 1, e.pawnCache.len())
	assert.EqualValues(t, 1, e.pawnCache.hits)
	assert.EqualValues(t, 1, e.pawnCache.misses)

	assert.EqualValues(t, score, score2)
}

func TestEvalPiecePawns(t *testing.T) {
	Settings.Eval.UsePawnEval = true
	Settings.Eval.UsePawnCache = false

	e := NewEvaluator()
	Settings.Eval.Tempo = 0
	p := position.NewPosition()
	var score *Score
	e.InitEval(p)

	score = e.evaluatePawns()
	out.Printf("Pawns: %s\n", score)

}

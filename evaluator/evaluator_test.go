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

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

func TestMaterial(t *testing.T) {
	e:=NewEvaluator()
	config.Settings.Eval.Tempo = 0

	// Startpos
	p := position.NewPosition()
	gpf := float32(p.GamePhase() / GamePhaseMax)
	assert.EqualValues(t,0, e.material(p, gpf))
	assert.EqualValues(t,0, e.positional(p, gpf))

	p.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	gpf = float32(p.GamePhase() / GamePhaseMax)
	assert.EqualValues(t,0, e.material(p, gpf))
	assert.EqualValues(t,55, e.positional(p, gpf))

	p.DoMove(CreateMove(SqD7, SqD5, Normal, PtNone))
	gpf = float32(p.GamePhase() / GamePhaseMax)
	assert.EqualValues(t,0, e.material(p, gpf))
	assert.EqualValues(t,0, e.positional(p, gpf))

	p.DoMove(CreateMove(SqE4, SqD5, Normal, PtNone))
	gpf = float32(p.GamePhase() / GamePhaseMax)
	assert.EqualValues(t,100, e.material(p, gpf))
	assert.EqualValues(t,30, e.positional(p, gpf))

	// TODO - this will change with additional evaluations
	assert.EqualValues(t,-130, e.Evaluate(p))
}

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

func TestPawnBitboards(t *testing.T) {
	e := NewEvaluator()
	p := position.NewPosition("r3k2r/1ppn3p/2q1q1nb/4P2N/2q1Pp2/B5RP/pbp2PP1/1R4K1 w kq - 0 1")
	e.InitEval(p)
	e.evaluatePawns()

	p = position.NewPosition("1r4k1/PBP2pp1/b5rp/2Q1pP2/4p2n/2Q1Q1NB/1PPN3P/R3K2R w - - 0 1")
	e.InitEval(p)
	e.evaluatePawns()
}

// this test does not really show reality as the cache entry will be cached in the L1 of the CPU
// usually this would not be the case during a search and we need to test this again in real
// world scenario.
func Benchmark(b *testing.B) {
	benchmarks := []struct {
		name  string
		cache bool
	}{
		{"w/o cache", false},
		{"with cache", true},
	}

	p := position.NewPosition("r3k2r/1ppn3p/2q1q1nb/4P2N/2q1Pp2/B5RP/pbp2PP1/1R4K1 w kq - 0 1")
	e := NewEvaluator()
	e.InitEval(p)

	for _, bm := range benchmarks {
		Settings.Eval.UsePawnCache = bm.cache
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				e.evaluatePawns()
			}
		})
	}
}


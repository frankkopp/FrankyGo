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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/pkg/movegen"
	"github.com/frankkopp/FrankyGo/pkg/position"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

func TestLeastValuablePiece(t *testing.T) {
	p := position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3")
	attacksTo := p.AttacksTo(SqE5, Black)

	logTest.Debug("All attackers\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2339760743907840, attacksTo)

	lva := getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqG6, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqD7, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqB2, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqE6, lva)
	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqNone, lva)
}

func TestSee(t *testing.T) {
	p := position.NewPosition("1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - -")
	move, _ := movegen.NewMoveGen().GetMoveFromUci(p, "d3e5")
	seeScore := see(p, move)
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, -220, seeScore)

	p = position.NewPosition("1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - -")
	move, _ = movegen.NewMoveGen().GetMoveFromUci(p, "e1e5")
	seeScore = see(p, move)
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, 100, seeScore)

	p = position.NewPosition("5q1k/8/8/8/RRQ2nrr/8/8/K7 w - -")
	move, _ = movegen.NewMoveGen().GetMoveFromUci(p, "c4f4")
	seeScore = see(p, move)
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, -580, seeScore)

	p = position.NewPosition("k6q/3n1n2/3b4/4p3/3P1P2/3N1N2/8/K7 w - -")
	move, _ = movegen.NewMoveGen().GetMoveFromUci(p, "d3e5")
	seeScore = see(p, move)
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, 100, seeScore)

	p = position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R2R1K1 b kq e3")
	move, _ = movegen.NewMoveGen().GetMoveFromUci(p, "a2b1Q")
	seeScore = see(p, move)
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, 500, seeScore)
}

func TestTimingSee(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("./bin")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof

	config.Settings.Search.UseBook = false
	config.Settings.Search.UseSEE = true

	p := position.NewPosition("k6q/3n1n2/3b4/4p3/3P1P2/3N1N2/8/K7 w - -")
	move, _ := movegen.NewMoveGen().GetMoveFromUci(p, "d3e5")

	const rounds = 5
	const iterations uint64 = 10_000_000

	seeScore := ValueNA
	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			seeScore = see(p, move)
		}
		elapsed := time.Since(start)
		out.Printf("Test took %s for %d iterations\n", elapsed, iterations)
		out.Printf("Test took %d ns per iteration\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Iterations per sec %d\n", int64(iterations*1e9)/elapsed.Nanoseconds())
	}
	logTest.Debug("See score:", seeScore)
	assert.EqualValues(t, 100, seeScore)

}

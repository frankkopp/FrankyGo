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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

func TestAttacks(t *testing.T) {
	p := position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
	a := NewAttacks()
	a.Compute(p)
	assert.Equal(t, p.ZobristKey(), a.Zobrist)
	assert.EqualValues(t, SqF1.Bb()|SqG1.Bb(), a.From[White][SqH1]&^p.OccupiedBb(White))
	assert.EqualValues(t, SqD8.Bb()|SqE7.Bb()|SqF8.Bb(), a.From[Black][SqE8]&^p.OccupiedBb(Black))
	assert.EqualValues(t, SqC6.Bb()|SqH5.Bb(), a.To[Black][SqE5]&p.OccupiedBb(Black))
}

func TestCompareWithPseudo(t *testing.T) {
	p := position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
	a := NewAttacks()
	a.nonPawnAttacks(p)
	for sq := SqA1; sq <= SqH8; sq++ {
		if p.GetPiece(sq) == PieceNone || p.GetPiece(sq).TypeOf() == Pawn {
			continue
		}
		c := p.GetPiece(sq).ColorOf()
		pt := p.GetPiece(sq).TypeOf()

		// compare the Attacks build with Attack and Magic bitboards
		// to the attacks calculated with the lopp in the local function
		magicAttacks := a.From[c][sq]
		nonMagicAttacks := buildAttacks(p, pt, sq)

		out.Println("Non Magic Attacks:\n", magicAttacks.StringBoard())
		out.Println("Build Attacks:\n", nonMagicAttacks.StringBoard())

		assert.EqualValues(t, magicAttacks, nonMagicAttacks)

		out.Println("==================================================")
	}
}

// to compare magic bitboard attacks with loop generated attacks
func buildAttacks(p *position.Position, pt PieceType, sq Square) Bitboard{
	occupiedAll := p.OccupiedAll()
	attacks := BbZero
	pseudoTo := GetPseudoAttacks(pt, sq) // & ^myPieces
	// iterate over all target squares of the piece
	if pt < Bishop { // king, knight
		attacks = pseudoTo
	} else {
		for tmp := pseudoTo; tmp != BbZero; {
			to := tmp.PopLsb()
			if Intermediate(sq, to)&occupiedAll == 0 {
				attacks.PushSquare(to)
			}
		}
	}
	return attacks
}

func Test_TimingAttacks(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof

	p := position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
	a := NewAttacks()

	const rounds = 5
	const iterations uint64 = 10_000_000

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			a.Clear()
			a.Compute(p)
		}
		elapsed := time.Since(start)
		out.Printf("Test took %s for %d iterations\n", elapsed, iterations)
		out.Printf("Test took %d ns per iteration\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Iterations per sec %d\n", int64(iterations*1e9)/elapsed.Nanoseconds())
	}
	_ = a
}

func Benchmark_NonPawnAttacks(b *testing.B) {
	p := position.NewPosition("6k1/p1qb1p1p/1p3np1/2b2p2/2B5/2P3N1/PP2QPPP/4N1K1 b - -")
	a := NewAttacks()

	f1 := func() {
		a.Clear()
		a.Compute(p)
	}

	benchmarks := []struct {
		name string
		f    func()
	}{
		{"New Instance", f1},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.f()
			}
		})
	}
	_ = a
}

func BenchmarkAttacks_ClearNewVsClear(b *testing.B) {
	a := NewAttacks()
	f1 := func() { a = NewAttacks() }
	f2 := func() { a.Clear() }
	benchmarks := []struct {
		name string
		f    func()
	}{
		{"New Instance", f1},
		{"Clear", f2},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.f()
			}
		})
	}
	_ = a
}

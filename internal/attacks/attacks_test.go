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

package attacks

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	logging2 "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

var logTest *logging2.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	config.Setup()
	logTest = logging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

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

		// out.Println("Non Magic Attacks:\n", magicStringBoard())
		// out.Println("Build Attacks:\n", nonMagicStringBoard())

		assert.EqualValues(t, magicAttacks, nonMagicAttacks)

		// out.Println("==================================================")
	}
}

func TestAttacksTo(t *testing.T) {
	var p *position.Position
	var attacksTo Bitboard

	p = position.NewPosition("2brr1k1/1pq1b1p1/p1np1p1p/P1p1p2n/1PNPPP2/2P1BNP1/4Q1BP/R2R2K1 w - -")
	attacksTo = AttacksTo(p, SqE5, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 740294656, attacksTo)

	attacksTo = AttacksTo(p, SqF1, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 20552, attacksTo)

	attacksTo = AttacksTo(p, SqD4, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 3407880, attacksTo)

	attacksTo = AttacksTo(p, SqD4, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 4483945857024, attacksTo)

	attacksTo = AttacksTo(p, SqD6, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 582090251837636608, attacksTo)

	attacksTo = AttacksTo(p, SqF8, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 5769111122661605376, attacksTo)

	p = position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3")
	attacksTo = AttacksTo(p, SqE5, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2339760743907840, attacksTo)

	attacksTo = AttacksTo(p, SqB1, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 1280, attacksTo)

	attacksTo = AttacksTo(p, SqG3, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 40960, attacksTo)

	attacksTo = AttacksTo(p, SqE4, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 4398113619968, attacksTo)
}

func TestRevealedAttacks(t *testing.T) {
	p := position.NewPosition("1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - -")
	occ := p.OccupiedAll()

	sq := SqE5

	attacksTo := AttacksTo(p, sq, White) | AttacksTo(p, sq, Black)
	logTest.Debug("Direct\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2286984186302464, attacksTo)

	// take away bishop on f6
	attacksTo.PopSquare(SqF6)
	occ.PopSquare(SqF6)

	attacksTo |= RevealedAttacks(p, sq, occ, White) | RevealedAttacks(p, sq, occ, Black)
	logTest.Debug("Revealed\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, Bitboard(9225623836668989440), attacksTo)

	// take away rook on e2
	attacksTo.PopSquare(SqE2)
	occ.PopSquare(SqE2)

	attacksTo |= RevealedAttacks(p, sq, occ, White) | RevealedAttacks(p, sq, occ, Black)
	logTest.Debug("Revealed\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, Bitboard(9225623836668985360), attacksTo)
}

// to compare magic bitboard attacks with loop generated attacks
func buildAttacks(p *position.Position, pt PieceType, sq Square) Bitboard {
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

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

package movegen

import (
	"log"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/movelist"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

func TestConstruction(t *testing.T) {
	Init()
	mg := New()
	log.Printf("%s type of %s", mg.String(), reflect.TypeOf(mg))
}

func Test_movegen_generatePawnMoves(t *testing.T) {
	Init()
	mg := New()
	pos := position.NewFen("1kr3nr/pp1pP1P1/2p1p3/3P1p2/1n1bP3/2P5/PP3PPP/RNBQKBNR w KQ -")
	moves := movelist.MoveList{}
	moves.SetMinCapacity(6) // 2^6 = 64

	mg.generatePawnMoves(&pos, GenCap, &moves)
	assert.Equal(t, 9, moves.Len())

	moves.Clear()
	mg.generatePawnMoves(&pos, GenNonCap, &moves)
	assert.Equal(t, 16, moves.Len())

	moves.Clear()
	mg.generatePawnMoves(&pos, GenAll, &moves)
	assert.Equal(t, 25, moves.Len())

	sort.Stable(&moves)
	log.Printf("Moves: %d\n", moves.Len())
	l := moves.Len()
	for i := 0; i < l; i++ {
		log.Printf("Move: %s\n", moves.At(i))
	}
}

func Test_movegen_generateCastling(t *testing.T) {
	Init()
	mg := New()
	pos := position.NewFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R w KQkq -")
	moves := movelist.MoveList{}
	moves.SetMinCapacity(6) // 2^6 = 64

	mg.generateCastling(&pos, GenAll, &moves)
	assert.Equal(t, 2, moves.Len())
	assert.Equal(t, "e1g1 e1c1", moves.StringUci())
	moves.Clear()

	pos = position.NewFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R b KQkq -")
	mg.generateCastling(&pos, GenAll, &moves)
	assert.Equal(t, 2, moves.Len())
	assert.Equal(t, "e8g8 e8c8", moves.StringUci())

}

func Test_movegen_generateKingMoves(t *testing.T) {
	Init()
	mg := New()
	moves := movelist.MoveList{}
	moves.SetMinCapacity(6) // 2^6 = 64

	pos := position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateKingMoves(&pos, GenAll, &moves)
	assert.Equal(t, 3, moves.Len())
	assert.Equal(t, "e1d2 e1d1 e1f1", moves.StringUci())
	moves.Clear()

	pos = position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateKingMoves(&pos, GenAll, &moves)
	assert.Equal(t, 3, moves.Len())
	assert.Equal(t, "e8d7 e8d8 e8f8", moves.StringUci())
}

func Test_movegen_generateMoves(t *testing.T) {
	Init()
	mg := New()
	moves := movelist.MoveList{}
	moves.SetMinCapacity(6) // 2^6 = 64

	pos := position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateMoves(&pos, GenCap, &moves)
	assert.Equal(t, 7, moves.Len())
	assert.Equal(t, "f3d2 f3e5 d7e5 d7b6 d7f6 b5c6 e2d2", moves.StringUci())
	moves.Clear()

	pos = position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(&pos, GenNonCap, &moves)
	assert.Equal(t, 28, moves.Len())
	assert.Equal(t, "d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6g4 f6d5 f6h5 f6g8 b4a3 b4a5 b4c5 b4d6 b7a6 "+
		"b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
	moves.Clear()

	pos = position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(&pos, GenAll, &moves)
	assert.Equal(t, 34, moves.Len())
	assert.Equal(t, "d2f3 d2e4 d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6e4 f6d7 f6g4 f6d5 f6h5 f6g8 b4c3 "+
		"b4a3 b4a5 b4c5 b4d6 b7a6 b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7d7 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
}

func Test_movegen_GeneratePseudoLegalMoves(t *testing.T) {
	Init()
	mg := New()

	pos := position.New()
	moves := mg.GeneratePseudoLegalMoves(&pos, GenAll)
	assert.Equal(t, 20, moves.Len())
	assert.Equal(t, "e2e4 d2d4 g1f3 b1c3 h2h3 a2a3 h2h4 g2g4 f2f4 e2e3 d2d3 c2c4 b2b4 a2a4 g2g3 b2b3 f2f3 "+
		"c2c3 g1h3 b1a3", moves.StringUci())
	// l := mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	log.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	pos = position.NewFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	moves = mg.GeneratePseudoLegalMoves(&pos, GenAll)
	assert.Equal(t, 40, moves.Len())
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 d7e5 f3e5 d7b6 e2d2 e1d2 e1g1 e1c1 d3d4 f3d4 d7c5 h1f1 a1d1 a1c1 b5c4 f3g5 h2h3 e2e3 a2a3 c3c4 h2h4 g2g4 a2a4 e1f1 g2g3 e2d1 b2b3 b5a6 b5a4 e2f1 h1g1 a1b1 e1d1 d7f8 f3h4 d7b8 f3g1", moves.StringUci())
	// l = mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	log.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	pos = position.NewFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves = mg.GeneratePseudoLegalMoves(&pos, GenAll)
	assert.Equal(t, 86, moves.Len())
	moves.Clear()

	pos = position.NewFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	moves = mg.GeneratePseudoLegalMoves(&pos, GenAll)
	assert.Equal(t, 218, moves.Len())
	moves.Clear()
}

func Test(t *testing.T) {
	out := message.NewPrinter(language.German)
	Init()
	const rounds = 5
	const iterations uint64 = 500_000

	mg := New()
	pos := position.NewFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	var moves *movelist.MoveList

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			moves = mg.GeneratePseudoLegalMoves(&pos, GenAll)
			moves.Clear()
		}
		elapsed := time.Since(start)
		out.Printf("GeneratePseudoLegalMoves took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		out.Printf("GeneratePseudoLegalMoves took %d ns\n", elapsed.Nanoseconds()/int64(iterations))
		generated := uint64(86) * iterations
		out.Printf("GeneratePseudoLegalMoves %d mps\n", (generated*uint64(1_000_000_000))/uint64(elapsed.Nanoseconds()))
	}
}

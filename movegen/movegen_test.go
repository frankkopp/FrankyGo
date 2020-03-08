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

	"github.com/stretchr/testify/assert"

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

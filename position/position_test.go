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

package position

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/frankkopp/FrankyGo/types"
)

func TestPositionCreation(t *testing.T) {
	Init()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p := NewFen(fen)
	fmt.Print(p.String())
	assert.Equal(t, SqA1.Bb()|SqH1.Bb()|SqA8.Bb()|SqH8.Bb(), p.piecesBb[White][Rook] | p.piecesBb[Black][Rook])
	assert.Equal(t, SqB1.Bb()|SqG1.Bb()|SqB8.Bb()|SqG8.Bb(), p.piecesBb[White][Knight] | p.piecesBb[Black][Knight])
	assert.Equal(t, SqC1.Bb()|SqF1.Bb()|SqC8.Bb()|SqF8.Bb(), p.piecesBb[White][Bishop] | p.piecesBb[Black][Bishop])
	assert.Equal(t, SqD1.Bb()|SqD8.Bb(), p.piecesBb[White][Queen] | p.piecesBb[Black][Queen])
	assert.Equal(t, SqE1.Bb()|SqE8.Bb(), p.piecesBb[White][King] | p.piecesBb[Black][King])
	assert.Equal(t, Rank2_Bb|Rank7_Bb, p.piecesBb[White][Pawn] | p.piecesBb[Black][Pawn])
	assert.Equal(t, White, p.nextPlayer)
	assert.Equal(t, CastlingAny, p.castlingRights)
	assert.Equal(t, SqNone, p.enPassantSquare)
	assert.Equal(t, 0, p.halfMoveClock)
	assert.Equal(t, 1, p.nextHalfMoveNumber)
	assert.Equal(t, 0, p.material[White] - p.material[Black])
	assert.Equal(t, 0, p.materialNonPawn[White] - p.materialNonPawn[Black])
	assert.Equal(t, 0, p.psqMidValue[White] - p.psqMidValue[Black])
	assert.Equal(t, 0, p.psqEndValue[White] - p.psqEndValue[Black])
	assert.Equal(t, fen, p.StringFen())

	fmt.Println()

	fen = "r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3 0 14"
	p = NewFen(fen)
	fmt.Print(p.String())
	assert.Equal(t, SqB1.Bb()|SqG3.Bb()|SqA8.Bb()|SqH8.Bb(), p.piecesBb[White][Rook] | p.piecesBb[Black][Rook])
	assert.Equal(t, SqD7.Bb()|SqG6.Bb(), p.piecesBb[White][Knight] | p.piecesBb[Black][Knight])
	assert.Equal(t, SqB2.Bb(), p.piecesBb[White][Bishop] | p.piecesBb[Black][Bishop])
	assert.Equal(t, SqC4.Bb()|SqC6.Bb()|SqE6.Bb(), p.piecesBb[White][Queen] | p.piecesBb[Black][Queen])
	assert.Equal(t, SqG1.Bb()|SqE8.Bb(), p.piecesBb[White][King] | p.piecesBb[Black][King])
	assert.Equal(t,
		SqA2.Bb()|SqB7.Bb()|SqC2.Bb()|SqC7.Bb()|SqE4.Bb()|SqE5.Bb()|SqF2.Bb()|SqF4.Bb()|SqG2.Bb()|SqH2.Bb()|SqH7.Bb(),
		p.piecesBb[White][Pawn] | p.piecesBb[Black][Pawn])
	assert.Equal(t, Black, p.nextPlayer)
	assert.Equal(t, CastlingBlack, p.castlingRights)
	assert.Equal(t, SqE3, p.enPassantSquare)
	assert.Equal(t, 0, p.halfMoveClock)
	assert.Equal(t, 28, p.nextHalfMoveNumber)
	assert.Equal(t, -3770, p.material[White] - p.material[Black])
	assert.Equal(t, -3670, p.materialNonPawn[White] - p.materialNonPawn[Black])
	assert.Equal(t, 113, p.psqMidValue[White] - p.psqMidValue[Black])
	assert.Equal(t, -165, p.psqEndValue[White] - p.psqEndValue[Black])
	assert.Equal(t, fen, p.StringFen())
}

func Benchmark(b *testing.B) {
	Init()
	for i := 0; i < b.N; i++ {
		NewFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	}
}


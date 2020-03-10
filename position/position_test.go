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
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	. "github.com/frankkopp/FrankyGo/types"

	"github.com/stretchr/testify/assert"
)

func TestPositionCreation(t *testing.T) {
	Init()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p := NewFen(fen)
	// fmt.Print(p.String())
	assert.Equal(t, SqA1.Bb()|SqH1.Bb()|SqA8.Bb()|SqH8.Bb(), p.piecesBb[White][Rook]|p.piecesBb[Black][Rook])
	assert.Equal(t, SqB1.Bb()|SqG1.Bb()|SqB8.Bb()|SqG8.Bb(), p.piecesBb[White][Knight]|p.piecesBb[Black][Knight])
	assert.Equal(t, SqC1.Bb()|SqF1.Bb()|SqC8.Bb()|SqF8.Bb(), p.piecesBb[White][Bishop]|p.piecesBb[Black][Bishop])
	assert.Equal(t, SqD1.Bb()|SqD8.Bb(), p.piecesBb[White][Queen]|p.piecesBb[Black][Queen])
	assert.Equal(t, SqE1.Bb()|SqE8.Bb(), p.piecesBb[White][King]|p.piecesBb[Black][King])
	assert.Equal(t, Rank2_Bb|Rank7_Bb, p.piecesBb[White][Pawn]|p.piecesBb[Black][Pawn])
	assert.Equal(t, White, p.nextPlayer)
	assert.Equal(t, CastlingAny, p.castlingRights)
	assert.Equal(t, SqNone, p.enPassantSquare)
	assert.Equal(t, 0, p.halfMoveClock)
	assert.Equal(t, 1, p.nextHalfMoveNumber)
	assert.Equal(t, Value(0), p.material[White]-p.material[Black])
	assert.Equal(t, Value(0), p.materialNonPawn[White]-p.materialNonPawn[Black])
	assert.Equal(t, Value(0), p.psqMidValue[White]-p.psqMidValue[Black])
	assert.Equal(t, Value(0), p.psqEndValue[White]-p.psqEndValue[Black])
	assert.Equal(t, fen, p.StringFen())

	fmt.Println()

	fen = "r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3 0 14"
	p = NewFen(fen)
	// fmt.Print(p.String())
	assert.Equal(t, SqB1.Bb()|SqG3.Bb()|SqA8.Bb()|SqH8.Bb(), p.piecesBb[White][Rook]|p.piecesBb[Black][Rook])
	assert.Equal(t, SqD7.Bb()|SqG6.Bb(), p.piecesBb[White][Knight]|p.piecesBb[Black][Knight])
	assert.Equal(t, SqB2.Bb(), p.piecesBb[White][Bishop]|p.piecesBb[Black][Bishop])
	assert.Equal(t, SqC4.Bb()|SqC6.Bb()|SqE6.Bb(), p.piecesBb[White][Queen]|p.piecesBb[Black][Queen])
	assert.Equal(t, SqG1.Bb()|SqE8.Bb(), p.piecesBb[White][King]|p.piecesBb[Black][King])
	assert.Equal(t,
		SqA2.Bb()|SqB7.Bb()|SqC2.Bb()|SqC7.Bb()|SqE4.Bb()|SqE5.Bb()|SqF2.Bb()|SqF4.Bb()|SqG2.Bb()|SqH2.Bb()|SqH7.Bb(),
		p.piecesBb[White][Pawn]|p.piecesBb[Black][Pawn])
	assert.Equal(t, Black, p.nextPlayer)
	assert.Equal(t, CastlingBlack, p.castlingRights)
	assert.Equal(t, SqE3, p.enPassantSquare)
	assert.Equal(t, 0, p.halfMoveClock)
	assert.Equal(t, 28, p.nextHalfMoveNumber)
	assert.Equal(t, Value(-3770), p.material[White]-p.material[Black])
	assert.Equal(t, Value(-3670), p.materialNonPawn[White]-p.materialNonPawn[Black])
	assert.Equal(t, Value(113), p.psqMidValue[White]-p.psqMidValue[Black])
	assert.Equal(t, Value(-165), p.psqEndValue[White]-p.psqEndValue[Black])
	assert.Equal(t, fen, p.StringFen())
}

func TestPositionEquality(t *testing.T) {
	Init()

	// equal
	p1 := New()
	p2 := NewFen(StartFen)
	assert.Equal(t, p1, p2)

	// not equal
	p3 := NewFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3 0 14")
	assert.NotEqual(t, p1, p3)

	// copy
	p3 = p2
	assert.Equal(t, p1, p3)
	p3.castlingRights.Remove(CastlingWhiteOO) // change to p3
	assert.NotEqual(t, p1, p3)
	assert.Equal(t, p1, p2)                // p2 from which p3 is copied is unchanged
	p3.castlingRights.Add(CastlingWhiteOO) // undo change
	assert.Equal(t, p1, p3)
}

func TestPosition_DoUndoMove(t *testing.T) {
	Init()
	p := New()
	p.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	p.DoMove(CreateMove(SqD7, SqD5, Normal, PtNone))
	p.DoMove(CreateMove(SqE4, SqD5, Normal, PtNone))
	p.DoMove(CreateMove(SqD8, SqD5, Normal, PtNone))
	p.DoMove(CreateMove(SqB1, SqC3, Normal, PtNone))
	p.UndoMove()
	p.UndoMove()
	p.UndoMove()
	p.UndoMove()
	p.UndoMove()
	assert.Equal(t, StartFen, p.StringFen())
}

func TestPosition_DoMoveNormal(t *testing.T) {
	Init()
	var fen string
	var position Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqC4, SqD4, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/3qPp2/B5R1/p1p2PPP/1R4K1 w kq - 1 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqC4, SqE4, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/4qp2/B5R1/p1p2PPP/1R4K1 w kq - 0 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w kq -"
	position = NewFen(fen)
	move = CreateMove(SqG3, SqG6, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1R1/8/2q1Pp2/B7/p1p2PPP/1R4K1 b kq - 0 1", position.StringFen())
}

func TestPosition_DoMoveCastling(t *testing.T) {
	Init()
	var fen string
	var position Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqE8, SqG8, Castling, PtNone)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r4rk1/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w - - 1 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqE8, SqC8, Castling, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "2kr3r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w - - 1 2", position.StringFen())
}

func TestPosition_DoMoveEnPassant(t *testing.T) {
	Init()
	var fen string
	var position Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqF4, SqE3, EnPassant, PtNone)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q5/B3p1R1/p1p2PPP/1R4K1 w kq - 0 2", position.StringFen())
}

func TestPosition_DoMovePromotion(t *testing.T) {
	Init()
	var fen string
	var position Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqA2, SqA1, Promotion, Queen)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/2p2PPP/qR4K1 w kq - 0 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	move = CreateMove(SqA2, SqB1, Promotion, Rook)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/2p2PPP/1r4K1 w kq - 0 2", position.StringFen())
}

func TestPosition_IsAttacked(t *testing.T) {
	Init()
	var fen string
	var position Position

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/6R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)

	// pawns
	assert.True(t, position.IsAttacked(SqG3, White))
	assert.True(t, position.IsAttacked(SqE3, White))
	assert.True(t, position.IsAttacked(SqB1, Black))
	assert.True(t, position.IsAttacked(SqE4, Black))
	assert.True(t, position.IsAttacked(SqE3, Black))

	// knight
	assert.True(t, position.IsAttacked(SqE5, Black))
	assert.True(t, position.IsAttacked(SqF4, Black))
	assert.False(t, position.IsAttacked(SqG1, Black))

	// sliding
	assert.True(t, position.IsAttacked(SqG6, White))
	assert.True(t, position.IsAttacked(SqA5, Black))

	fen = "rnbqkbnr/1ppppppp/8/p7/Q1P5/8/PP1PPPPP/RNB1KBNR b KQkq - 1 2"
	position = NewFen(fen)

	// king
	assert.True(t, position.IsAttacked(SqD1, White))
	assert.False(t, position.IsAttacked(SqE1, Black))

	// rook
	assert.True(t, position.IsAttacked(SqA5, Black))
	assert.False(t, position.IsAttacked(SqA4, Black))

	// queen
	assert.False(t, position.IsAttacked(SqE8, White))
	assert.True(t, position.IsAttacked(SqD7, White))
	assert.False(t, position.IsAttacked(SqE8, White))

	// en passant
	fen = "rnbqkbnr/1pp1pppp/p7/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6"
	position = NewFen(fen)
	assert.True(t, position.IsAttacked(SqD5, White))

	fen = "rnbqkbnr/1pp1pppp/p7/2Pp4/8/8/PP1PPPPP/RNBQKBNR w KQkq d6"
	position = NewFen(fen)
	assert.True(t, position.IsAttacked(SqD5, White))

	fen = "rnbqkbnr/pppp1ppp/8/8/3Pp3/7P/PPP1PPP1/RNBQKBNR b - d3"
	position = NewFen(fen)
	assert.True(t, position.IsAttacked(SqD4, Black))

	fen = "rnbqkbnr/pppp1ppp/8/8/2pP4/7P/PPP1PPP1/RNBQKBNR b - d3"
	position = NewFen(fen)
	assert.True(t, position.IsAttacked(SqD4, Black))

	// bug tests
	fen = "r1bqk1nr/pppp1ppp/2nb4/1B2B3/3pP3/8/PPP2PPP/RN1QK1NR b KQkq -"
	position = NewFen(fen)
	assert.False(t, position.IsAttacked(SqE8, White))
	assert.False(t, position.IsAttacked(SqE1, Black))

	fen = "rnbqkbnr/ppp1pppp/8/1B6/3Pp3/8/PPP2PPP/RNBQK1NR b KQkq -"
	position = NewFen(fen)
	assert.True(t, position.IsAttacked(SqE8, White))
	assert.False(t, position.IsAttacked(SqE1, Black))

	fen = "8/1pk2p2/2p5/5p2/8/1pp2Q2/5K2/8 w - -"
	position = NewFen(fen)
	assert.False(t, position.IsAttacked(SqF7, White))
	assert.False(t, position.IsAttacked(SqB7, White))
	assert.False(t, position.IsAttacked(SqB3, White))
}

func TestPosition_IsLegalMoves(t *testing.T) {
	Init()
	var fen string
	var position Position

	// no o-o castling / o-o-o is allowed
	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	assert.True(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))

	// in check - no castling at all
	fen = "r3k2r/1ppn3p/2q1qNn1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))

}

func TestPosition_WasLegalMove(t *testing.T) {
	Init()
	var fen string
	var position Position

	// no o-o castling
	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position = NewFen(fen)
	position.DoMove(CreateMove(SqE8, SqG8, Castling, PtNone)) // illegal as king crosses attacked square
	assert.False(t, position.WasLegalMove())
	position.UndoMove()
	position.DoMove(CreateMove(SqE8, SqC8, Castling, PtNone)) // legal
	assert.True(t, position.WasLegalMove())

	// in check - no castling at all
	fen = "r3k2r/1ppn3p/2q1qNn1/8/2q1Pp2/6R1/p1p2PPP/1R4K1 b kq -"
	position = NewFen(fen)
	position.DoMove(CreateMove(SqE8, SqG8, Castling, PtNone)) // illegal as king crosses attacked square
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	position.UndoMove()
	position.DoMove(CreateMove(SqE8, SqC8, Castling, PtNone))
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))
}

//noinspection GoUnhandledErrorResult
func Test_TimingDoUndo(t *testing.T) {
	out := message.NewPrinter(language.German)
	Init()
	const rounds = 5
	const iterations uint64 = 10_000_000

	// prepare moves
	e2e4 := CreateMove(SqE2, SqE4, Normal, PtNone)
	d7d5 := CreateMove(SqD7, SqD5, Normal, PtNone)
	e4d5 := CreateMove(SqE4, SqD5, Normal, PtNone)
	d8d5 := CreateMove(SqD8, SqD5, Normal, PtNone)
	b1c3 := CreateMove(SqB1, SqC3, Normal, PtNone)

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		p := New()
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			p.DoMove(e2e4)
			p.DoMove(d7d5)
			p.DoMove(e4d5)
			p.DoMove(d8d5)
			p.DoMove(b1c3)
			p.UndoMove()
			p.UndoMove()
			p.UndoMove()
			p.UndoMove()
			p.UndoMove()
		}
		elapsed := time.Since(start)
		out.Printf("DoMove/UndoMove took %d ns for %d iterations with 5 do/undo pairs\n", elapsed.Nanoseconds(), iterations)
		out.Printf("DoMove/UndoMove took %d ns per do/undo pair\n", elapsed.Nanoseconds()/int64(iterations*5))
		out.Printf("Positions per sec %d pps\n",int64(iterations*5*1e9)/elapsed.Nanoseconds())
	}
}

// ///////////////////////////////////
// BENCHMARKS

func Benchmark_New(b *testing.B) {
	Init()
	for i := 0; i < b.N; i++ {
		NewFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	}
}

func Benchmark_DoUndo(b *testing.B) {
	Init()
	p := New()
	for i := 0; i < b.N; i++ {
		p.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
		p.DoMove(CreateMove(SqD7, SqD5, Normal, PtNone))
		p.DoMove(CreateMove(SqE4, SqD5, Normal, PtNone))
		p.DoMove(CreateMove(SqD8, SqD5, Normal, PtNone))
		p.DoMove(CreateMove(SqB1, SqC3, Normal, PtNone))
		p.UndoMove()
		p.UndoMove()
		p.UndoMove()
		p.UndoMove()
		p.UndoMove()
	}
}

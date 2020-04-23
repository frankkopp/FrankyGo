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

package position

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/op/go-logging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	. "github.com/frankkopp/FrankyGo/internal/types"

	"github.com/stretchr/testify/assert"
)

var out = message.NewPrinter(language.German)
var logTest *logging.Logger

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
	logTest = myLogging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func TestPositionCreation(t *testing.T) {

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, _ := NewPositionFen(fen)
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
	p, _ = NewPositionFen(fen)
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
	assert.Equal(t, Value(118), p.psqMidValue[White]-p.psqMidValue[Black])
	assert.Equal(t, Value(-165), p.psqEndValue[White]-p.psqEndValue[Black])
	assert.Equal(t, fen, p.StringFen())
}

func TestPositionEquality(t *testing.T) {

	// equal
	p1 := NewPosition()
	p2, _ := NewPositionFen(StartFen)
	assert.Equal(t, p1, p2)

	// not equal
	p3, _ := NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3 0 14")
	assert.NotEqual(t, p1, p3)

	// copy
	*p3 = *p2
	assert.Equal(t, *p1, *p3)
	p3.castlingRights.Remove(CastlingWhiteOO) // change to p3
	assert.NotEqual(t, *p1, *p3)
	assert.Equal(t, *p1, *p2)              // p2 from which p3 is copied is unchanged
	p3.castlingRights.Add(CastlingWhiteOO) // undo change
	assert.Equal(t, *p1, *p3)
}

func TestPosition_DoUndoMove(t *testing.T) {

	p := NewPosition()
	startZobrist := p.ZobristKey()
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
	assert.Equal(t, startZobrist, p.ZobristKey())
}

func TestPosition_DoMoveNormal(t *testing.T) {

	var fen string
	var position *Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqC4, SqD4, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/3qPp2/B5R1/p1p2PPP/1R4K1 w kq - 1 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqC4, SqE4, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/4qp2/B5R1/p1p2PPP/1R4K1 w kq - 0 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w kq -"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqG3, SqG6, Normal, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1R1/8/2q1Pp2/B7/p1p2PPP/1R4K1 b kq - 0 1", position.StringFen())
}

func TestPosition_DoMoveCastling(t *testing.T) {

	var fen string
	var position *Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqE8, SqG8, Castling, PtNone)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r4rk1/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w - - 1 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqE8, SqC8, Castling, PtNone)
	position.DoMove(move)
	// log.Println(position.String())
	assert.Equal(t, "2kr3r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 w - - 1 2", position.StringFen())
}

func TestPosition_DoMoveEnPassant(t *testing.T) {

	var fen string
	var position *Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqF4, SqE3, EnPassant, PtNone)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q5/B3p1R1/p1p2PPP/1R4K1 w kq - 0 2", position.StringFen())
}

func TestPosition_DoMovePromotion(t *testing.T) {

	var fen string
	var position *Position
	var move Move

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqA2, SqA1, Promotion, Queen)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/2p2PPP/qR4K1 w kq - 0 2", position.StringFen())

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	move = CreateMove(SqA2, SqB1, Promotion, Rook)
	position.DoMove(move) // would be illegal as King crosses attacked square
	// log.Println(position.String())
	assert.Equal(t, "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/2p2PPP/1r4K1 w kq - 0 2", position.StringFen())
}

func TestPosition_IsAttacked(t *testing.T) {

	var fen string
	var position *Position

	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/6R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)

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
	position, _ = NewPositionFen(fen)

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
	position, _ = NewPositionFen(fen)
	assert.True(t, position.IsAttacked(SqD5, White))

	fen = "rnbqkbnr/1pp1pppp/p7/2Pp4/8/8/PP1PPPPP/RNBQKBNR w KQkq d6"
	position, _ = NewPositionFen(fen)
	assert.True(t, position.IsAttacked(SqD5, White))

	fen = "rnbqkbnr/pppp1ppp/8/8/3Pp3/7P/PPP1PPP1/RNBQKBNR b - d3"
	position, _ = NewPositionFen(fen)
	assert.True(t, position.IsAttacked(SqD4, Black))

	fen = "rnbqkbnr/pppp1ppp/8/8/2pP4/7P/PPP1PPP1/RNBQKBNR b - d3"
	position, _ = NewPositionFen(fen)
	assert.True(t, position.IsAttacked(SqD4, Black))

	// bug tests
	fen = "r1bqk1nr/pppp1ppp/2nb4/1B2B3/3pP3/8/PPP2PPP/RN1QK1NR b KQkq -"
	position, _ = NewPositionFen(fen)
	assert.False(t, position.IsAttacked(SqE8, White))
	assert.False(t, position.IsAttacked(SqE1, Black))

	fen = "rnbqkbnr/ppp1pppp/8/1B6/3Pp3/8/PPP2PPP/RNBQK1NR b KQkq -"
	position, _ = NewPositionFen(fen)
	assert.True(t, position.IsAttacked(SqE8, White))
	assert.False(t, position.IsAttacked(SqE1, Black))

	fen = "8/1pk2p2/2p5/5p2/8/1pp2Q2/5K2/8 w - -"
	position, _ = NewPositionFen(fen)
	assert.False(t, position.IsAttacked(SqF7, White))
	assert.False(t, position.IsAttacked(SqB7, White))
	assert.False(t, position.IsAttacked(SqB3, White))
}

func TestPosition_IsLegalMoves(t *testing.T) {

	var fen string
	var position *Position

	// no o-o castling / o-o-o is allowed
	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	assert.True(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))

	// in check - no castling at all
	fen = "r3k2r/1ppn3p/2q1qNn1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))

}

func TestPosition_WasLegalMove(t *testing.T) {

	var fen string
	var position *Position

	// no o-o castling
	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	position.DoMove(CreateMove(SqE8, SqG8, Castling, PtNone)) // illegal as king crosses attacked square
	assert.False(t, position.WasLegalMove())
	position.UndoMove()
	position.DoMove(CreateMove(SqE8, SqC8, Castling, PtNone)) // legal
	assert.True(t, position.WasLegalMove())

	// in check - no castling at all
	fen = "r3k2r/1ppn3p/2q1qNn1/8/2q1Pp2/6R1/p1p2PPP/1R4K1 b kq -"
	position, _ = NewPositionFen(fen)
	position.DoMove(CreateMove(SqE8, SqG8, Castling, PtNone)) // illegal as king crosses attacked square
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqG8, Castling, PtNone)))
	position.UndoMove()
	position.DoMove(CreateMove(SqE8, SqC8, Castling, PtNone))
	assert.False(t, position.IsLegalMove(CreateMove(SqE8, SqC8, Castling, PtNone)))
}

func TestPositionGivesCheck(t *testing.T) {

	// DIRECT CHECKS

	// Pawns
	p := NewPosition("4r3/1pn3k1/4p1b1/p1Pp1P1r/3P2NR/1P3B2/3K2P1/4R3 w - -")
	move := CreateMove(SqF5, SqF6, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("5k2/4pp2/1N2n1p1/r3P2p/P5PP/2rR1K2/P7/3R4 b - -")
	move = CreateMove(SqH5, SqG4, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// promotion
	p = NewPosition("1k3r2/1p1bP3/2p2p1Q/Ppb5/4Rp1P/2q2N1P/5PB1/6K1 w - -")
	move = CreateMove(SqE7, SqF8, Promotion, Queen)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("1r3r2/1p1bP2k/2p2n2/p1Pp4/P2N1PpP/1R2p3/1P2P1BP/3R2K1 w - -")
	move = CreateMove(SqE7, SqF8, Promotion, Knight)
	assert.True(t, p.GivesCheck(move))

	// Knights
	p = NewPosition("5k2/4pp2/1N2n1p1/r3P2p/P5PP/2rR1K2/P7/3R4 w - -")
	move = CreateMove(SqB6, SqD7, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("5k2/4pp2/1N2n1p1/r3P2p/P5PP/2rR1K2/P7/3R4 b - -")
	move = CreateMove(SqE6, SqD4, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// Rooks
	p = NewPosition("5k2/4pp2/1N2n1pp/r3P3/P5PP/2rR4/P3K3/3R4 w - -")
	move = CreateMove(SqD3, SqD8, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("5k2/4pp2/1N2n1pp/r3P3/P5PP/2rR4/P3K3/3R4 b - -")
	move = CreateMove(SqC3, SqC2, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// blocked opponent piece - no check
	p = NewPosition("5k2/4pp2/1N2n1pp/r3P3/P5PP/2rR4/P2RK3/8 b - -")
	move = CreateMove(SqC3, SqC2, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))

	// blocked own piece - no check
	p = NewPosition("5k2/4pp2/1N2n1pp/r3P3/P5PP/2rR4/P2nK3/3R4 b - -")
	move = CreateMove(SqC3, SqC2, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))

	// Bishop
	p = NewPosition("6k1/3q2b1/p1rrnpp1/P3p3/2B1P3/1p1R3Q/1P4PP/1B1R3K w - -")
	move = CreateMove(SqC4, SqE6, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// Queen
	p = NewPosition("5k2/4pp2/1N2n1pp/r3P3/P5PP/2qR4/P3K3/3R4 b - -")
	move = CreateMove(SqC3, SqC2, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("6k1/3q2b1/p1rrnpp1/P3p3/2B1P3/1p1R3Q/1P4PP/1B1R3K w - -")
	move = CreateMove(SqH3, SqE6, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("6k1/p3q2p/1n1Q2pB/8/5P2/6P1/PP5P/3R2K1 b - -")
	move = CreateMove(SqE7, SqE3, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// no check
	p = NewPosition("6k1/p3q2p/1n1Q2pB/8/5P2/6P1/PP5P/3R2K1 b - -")
	move = CreateMove(SqE7, SqE4, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))

	// Castling checks
	p = NewPosition("r4k1r/8/8/8/8/8/8/R3K2R w KQ -")
	move = CreateMove(SqE1, SqG1, Castling, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("r2k3r/8/8/8/8/8/8/R3K2R w KQ -")
	move = CreateMove(SqE1, SqC1, Castling, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("r3k2r/8/8/8/8/8/8/R4K1R b kq -")
	move = CreateMove(SqE8, SqG8, Castling, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("r3k2r/8/8/8/8/8/8/R2K3R b kq -")
	move = CreateMove(SqE8, SqC8, Castling, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("r6r/8/8/8/8/8/8/2k1K2R w K -")
	move = CreateMove(SqE1, SqG1, Castling, PtNone)
	assert.True(t, p.GivesCheck(move))

	// en passant checks
	p = NewPosition("8/3r1pk1/p1R2p2/1p5p/r2Pp3/PRP3P1/4KP1P/8 b - d3")
	move = CreateMove(SqE4, SqD3, EnPassant, PtNone)
	assert.True(t, p.GivesCheck(move))

	// REVEALED CHECKS
	p = NewPosition("6k1/8/3P1bp1/2BNp3/8/1Q3P1q/7r/1K2R3 w - -")
	move = CreateMove(SqD5, SqE7, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("6k1/8/3P1bp1/2BNp3/8/1Q3P1q/7r/1K2R3 w - -")
	move = CreateMove(SqD5, SqC7, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("1Q1N2k1/8/3P1bp1/2B1p3/8/5P1q/7r/1K2R3 w - -")
	move = CreateMove(SqD8, SqE6, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("1R1N2k1/8/3P1bp1/2B1p3/8/5P1q/7r/1K2R3 w - -")
	move = CreateMove(SqD8, SqE6, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	// revealed by en passant capture
	p = NewPosition("8/b2r1pk1/p1R2p2/1p5p/r2Pp3/PRP3P1/5K1P/8 b - d3")
	move = CreateMove(SqE4, SqD3, EnPassant, PtNone)
	assert.True(t, p.GivesCheck(move))

	// Misc
	p = NewPosition("2r1r3/pb1n1kpn/1p1qp3/6p1/2PP4/8/P2Q1PPP/3R1RK1 w - -")
	move = CreateMove(SqF2, SqF4, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))

	p = NewPosition("2r1r1k1/pb3pp1/1p1qpn2/4n1p1/2PP4/6KP/P2Q1PP1/3RR3 b - -")
	move = CreateMove(SqE5, SqD3, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q1NNQQ2/1p6/qk3KB1 b - -")
	move = CreateMove(SqB1, SqC2, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("8/8/8/8/8/5K2/R7/7k w - -")
	move = CreateMove(SqA2, SqH2, Normal, PtNone)
	assert.True(t, p.GivesCheck(move))

	p = NewPosition("r1bqkb1r/ppp1pppp/2n2n2/1B1P4/8/8/PPPP1PPP/RNBQK1NR w KQkq -")
	move = CreateMove(SqD5, SqC6, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))

	p = NewPosition("rnbq1bnr/pppkpppp/8/3p4/3P4/3Q4/PPP1PPPP/RNB1KBNR w KQ -")
	move = CreateMove(SqD3, SqH7, Normal, PtNone)
	assert.False(t, p.GivesCheck(move))
}

func TestPosition_CheckRepetitions(t *testing.T) {
	// test 1
	position := NewPosition()
	position.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	position.DoMove(CreateMove(SqE7, SqE5, Normal, PtNone))
	// takes 3 loops to get to repetition
	for i := 0; i <= 2; i++ {
		position.DoMove(CreateMove(SqG1, SqF3, Normal, PtNone))
		position.DoMove(CreateMove(SqB8, SqC6, Normal, PtNone))
		position.DoMove(CreateMove(SqF3, SqG1, Normal, PtNone))
		position.DoMove(CreateMove(SqC6, SqB8, Normal, PtNone))
	}
	assert.True(t, position.CheckRepetitions(2))

	// test 2
	position, _ = NewPositionFen("6k1/p3q2p/1n1Q2pB/8/5P2/6P1/PP5P/3R2K1 b - -")
	position.DoMove(CreateMove(SqE7, SqE3, Normal, PtNone))
	position.DoMove(CreateMove(SqG1, SqG2, Normal, PtNone))
	// takes 3 loops to get to repetition
	for i := 0; i <= 2; i++ {
		position.DoMove(CreateMove(SqE3, SqE2, Normal, PtNone))
		position.DoMove(CreateMove(SqG2, SqG1, Normal, PtNone))
		position.DoMove(CreateMove(SqE2, SqE3, Normal, PtNone))
		position.DoMove(CreateMove(SqG1, SqG2, Normal, PtNone))
	}
	assert.True(t, position.CheckRepetitions(2))
}

func TestPosition_DoNullMove(t *testing.T) {
	var fen string
	var position *Position

	// no o-o castling / o-o-o is allowed
	fen = "r3k2r/1ppn3p/2q1q1n1/8/2q1Pp2/B5R1/p1p2PPP/1R4K1 b kq e3"
	position, _ = NewPositionFen(fen)
	p1 := *position
	position.DoNullMove()
	position.UndoNullMove()
	assert.Equal(t, p1.StringFen(), position.StringFen())
	assert.Equal(t, p1.ZobristKey(), position.ZobristKey())
}

func TestPosition_CheckInsufficientMaterial(t *testing.T) {
	// 	both sides have a bare king
	position, _ := NewPositionFen("8/3k4/8/8/8/8/4K3/8 w - -")
	assert.True(t, position.HasInsufficientMaterial())

	// 	one side has a king and a minor piece against a bare king
	// 	both sides have a king and a minor piece each
	position, _ = NewPositionFen("8/3k4/8/8/8/2B5/4K3/8 w - -")
	assert.True(t, position.HasInsufficientMaterial())
	position, _ = NewPositionFen("8/8/4K3/8/8/2b5/4k3/8 b - -")
	assert.True(t, position.HasInsufficientMaterial())

	// 	both sides have a king and a bishop, the bishops being the same color
	position, _ = NewPositionFen("8/8/3BK3/8/8/2b5/4k3/8 b - -")
	assert.True(t, position.HasInsufficientMaterial())
	position, _ = NewPositionFen("8/8/2B1K3/8/8/8/2b1k3/8 b - -")
	assert.True(t, position.HasInsufficientMaterial())
	position, _ = NewPositionFen("8/8/4K3/2B5/8/8/2b1k3/8 b - -")
	assert.True(t, position.HasInsufficientMaterial())

	// one side has two bishops a mate can be forced
	position, _ = NewPositionFen("8/8/2B1K3/2B5/8/8/2n1k3/8 b - -")
	assert.False(t, position.HasInsufficientMaterial())

	// 	two knights against the bare king
	position, _ = NewPositionFen("8/8/2NNK3/8/8/8/4k3/8 w - -")
	assert.True(t, position.HasInsufficientMaterial())
	position, _ = NewPositionFen("8/8/2nnk3/8/8/8/4K3/8 w - -")
	assert.True(t, position.HasInsufficientMaterial())

	// 	the weaker side has a minor piece against two knights
	position, _ = NewPositionFen("8/8/2n1kn2/8/8/8/4K3/4B3 w - -")
	assert.True(t, position.HasInsufficientMaterial())

	// 	two bishops draw against a bishop
	position, _ = NewPositionFen("8/8/3bk1b1/8/8/8/4K3/4B3 w - -")
	assert.True(t, position.HasInsufficientMaterial())

	// 	two minor pieces against one draw, except when the stronger side has a bishop pair
	position, _ = NewPositionFen("8/8/3bk1b1/8/8/8/4K3/4N3 w - -")
	assert.False(t, position.HasInsufficientMaterial())
	position, _ = NewPositionFen("8/8/3bk1n1/8/8/8/4K3/4N3 w - -")
	assert.True(t, position.HasInsufficientMaterial())

}

// DoMove/UndoMove took 3.065.041.200 ns for 10.000.000 iterations with 5 do/undo pairs
// DoMove/UndoMove took 61 ns per do/undo pair
// Positions per sec 16.312.994 pps
//
//noinspection GoUnhandledErrorResult
func TestTimingDoUndo(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

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
		p := NewPosition()
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
		out.Printf("Positions per sec %d pps\n", int64(iterations*5*1e9)/elapsed.Nanoseconds())
	}
}

var res bool

func TestTimingMatvsPop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const rounds = 5
	const iterations uint64 = 1_000_000_000

	p, _ := NewPositionFen("6k1/p3q2p/1n1Q2pB/8/5P2/6P1/PP5P/3R2K1 b - -")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			test := (p.materialNonPawn[White] < 2*Bishop.ValueOf() && p.materialNonPawn[Black] <= Bishop.ValueOf()) ||
				(p.materialNonPawn[White] <= Bishop.ValueOf() && p.materialNonPawn[Black] < 2*Bishop.ValueOf())
			res = test
		}
		elapsed := time.Since(start)
		out.Printf("Test took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		out.Printf("Test took %d ns per test\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Test per sec %d tps\n", iterations*1e9/uint64(elapsed.Nanoseconds()))
	}
}

func TestTimingMatvsPop2(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const rounds = 5
	const iterations uint64 = 1_000_000_000

	p, _ := NewPositionFen("6k1/p3q2p/1n1Q2pB/8/5P2/6P1/PP5P/3R2K1 b - -")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			test := (p.piecesBb[White][Bishop].PopCount()+p.piecesBb[White][Knight].PopCount() == 2 &&
				p.piecesBb[Black][Bishop].PopCount()+p.piecesBb[Black][Knight].PopCount() == 1) ||
				(p.piecesBb[Black][Bishop].PopCount()+p.piecesBb[Black][Knight].PopCount() == 2 &&
					p.piecesBb[White][Bishop].PopCount()+p.piecesBb[White][Knight].PopCount() == 1)
			res = test
		}
		elapsed := time.Since(start)
		out.Printf("Test took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		out.Printf("Test took %d ns per test\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Test per sec %d tps\n", (iterations*1e9)/uint64(elapsed.Nanoseconds()))
	}
}

func TestTimingIsAttacked(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	out := message.NewPrinter(language.German)

	const rounds = 5
	const iterations uint64 = 10_000_000

	p, _ := NewPositionFen("r5k1/p1qb1p1p/1p3np1/2b2p2/2B5/2P3N1/PP2QPPP/R3N1K1 b - -")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		test := false
		for i := uint64(0); i < iterations; i++ {
			for sq := SqA1; sq <= SqH8; sq++ {
				test = p.IsAttacked(sq, White)
				test = p.IsAttacked(sq, Black)
			}
			res = test
		}
		elapsed := time.Since(start)
		out.Printf("Test took %s for %d iterations\n", elapsed, iterations)
		out.Printf("Test took %d ns per test\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Tests per sec %d tps\n", iterations*1e9/uint64(elapsed.Nanoseconds()))
	}
}

//
// func BenchmarkIsAttackedVariations(b *testing.B) {
//
// 	p, _ := NewPositionFen("r5k1/p1qb1p1p/1p3np1/2b2p2/2B5/2P3N1/PP2QPPP/R3N1K1 b - -")
// 	// p = NewPosition()
//
// 	f1 := func() {
// 		for sq := SqA1; sq <= SqH8; sq++ {
// 			res = GetAttacksBb(Bishop, sq, p.OccupiedAll()) & p.piecesBb[White][Bishop] > 0 ||
// 				GetAttacksBb(Rook, sq, p.OccupiedAll()) & p.piecesBb[White][Rook] > 0 ||
// 				GetAttacksBb(Queen, sq, p.OccupiedAll()) & p.piecesBb[White][Queen] > 0
// 		}
// 	}
//
// 	f2 := func() {
// 		// for sq := SqA1; sq <= SqH8; sq++ {
// 		// 	res = (GetPseudoAttacks(Rook, sq)&p.piecesBb[White][Rook] != 0 || (GetPseudoAttacks(Rook, sq)&p.piecesBb[White][Queen] != 0)) &&
// 		// 		(((GetMovesOnRank(sq, p.OccupiedAll()) |
// 		// 			GetMovesOnFileRotated(sq, p.occupiedBbL90[White]|p.occupiedBbL90[Black])) &
// 		// 			(p.piecesBb[White][Rook] | p.piecesBb[White][Queen])) != 0) &&
// 		// 		(GetPseudoAttacks(Bishop, sq)&p.piecesBb[White][Bishop] != 0 || (GetPseudoAttacks(Bishop, sq)&p.piecesBb[White][Queen] != 0)) &&
// 		// 		(((GetMovesDiagUpRotated(sq, p.occupiedBbR45[White]|p.occupiedBbR45[Black]) |
// 		// 			GetMovesDiagDownRotated(sq, p.occupiedBbL45[White]|p.occupiedBbL45[Black])) &
// 		// 			(p.piecesBb[White][Bishop] | p.piecesBb[White][Queen])) != 0)
// 		//
// 		// 	res = (GetPseudoAttacks(Rook, sq)&p.piecesBb[Black][Rook] != 0 || (GetPseudoAttacks(Rook, sq)&p.piecesBb[Black][Queen] != 0)) &&
// 		// 		(((GetMovesOnRank(sq, p.OccupiedAll()) |
// 		// 			GetMovesOnFileRotated(sq, p.occupiedBbL90[Black]|p.occupiedBbL90[Black])) &
// 		// 			(p.piecesBb[Black][Rook] | p.piecesBb[Black][Queen])) != 0) &&
// 		// 		(GetPseudoAttacks(Bishop, sq)&p.piecesBb[Black][Bishop] != 0 || (GetPseudoAttacks(Bishop, sq)&p.piecesBb[Black][Queen] != 0)) &&
// 		// 		(((GetMovesDiagUpRotated(sq, p.occupiedBbR45[Black]|p.occupiedBbR45[Black]) |
// 		// 			GetMovesDiagDownRotated(sq, p.occupiedBbL45[Black]|p.occupiedBbL45[Black])) &
// 		// 			(p.piecesBb[Black][Bishop] | p.piecesBb[Black][Queen])) != 0)
// 		// }
// 	}
//
// 	benchmarks := []struct {
// 		name string
// 		f    func()
// 	}{
// 		{"Magic", f1},
// 		{"NonMagic", f2},
// 	}
// 	for _, bm := range benchmarks {
// 		b.Run(bm.name, func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				bm.f()
// 			}
// 		})
// 	}
// }

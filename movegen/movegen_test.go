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

package movegen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var logTest = logging.GetLog("test")

func TestMovegen_String(t *testing.T) {

	mg := NewMoveGen()
	out.Println(mg.String())
}

func Test_movegen_generatePawnMoves(t *testing.T) {
	mg := NewMoveGen()
	pos := position.NewPositionFen("1kr3nr/pp1pP1P1/2p1p3/3P1p2/1n1bP3/2P5/PP3PPP/RNBQKBNR w KQ -")
	moves := moveslice.MoveSlice{}

	mg.generatePawnMoves(pos, GenCap, &moves)
	assert.Equal(t, 9, len(moves))

	moves.Clear()
	mg.generatePawnMoves(pos, GenNonCap, &moves)
	assert.Equal(t, 16, len(moves))

	moves.Clear()
	mg.generatePawnMoves(pos, GenAll, &moves)
	assert.Equal(t, 25, len(moves))

	// moves.Sort()
	// fmt.Printf("Moves: %d\n", len(moves))
	// l := len(moves)
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("Move: %s\n", moves.At(i))
	// }
}

func Test_movegen_generateCastling(t *testing.T) {
	mg := NewMoveGen()
	pos := position.NewPositionFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R w KQkq -")
	moves := moveslice.MoveSlice{}

	mg.generateCastling(pos, GenAll, &moves)
	assert.Equal(t, 2, len(moves))
	assert.Equal(t, "e1g1 e1c1", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R b KQkq -")
	mg.generateCastling(pos, GenAll, &moves)
	assert.Equal(t, 2, len(moves))
	assert.Equal(t, "e8g8 e8c8", moves.StringUci())

}

func Test_movegen_generateKingMoves(t *testing.T) {
	mg := NewMoveGen()
	moves := moveslice.MoveSlice{}

	pos := position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateKingMoves(pos, GenAll, &moves)
	assert.Equal(t, 3, len(moves))
	assert.Equal(t, "e1d2 e1d1 e1f1", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateKingMoves(pos, GenAll, &moves)
	assert.Equal(t, 3, len(moves))
	assert.Equal(t, "e8d7 e8d8 e8f8", moves.StringUci())
}

func Test_movegen_generateMoves(t *testing.T) {
	mg := NewMoveGen()
	moves := moveslice.MoveSlice{}

	pos := position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateMoves(pos, GenCap, &moves)
	assert.Equal(t, 7, len(moves))
	assert.Equal(t, "f3d2 f3e5 d7e5 d7b6 d7f6 b5c6 e2d2", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(pos, GenNonCap, &moves)
	assert.Equal(t, 28, len(moves))
	assert.Equal(t, "d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6g4 f6d5 f6h5 f6g8 b4a3 b4a5 b4c5 b4d6 b7a6 "+
		"b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(pos, GenAll, &moves)
	assert.Equal(t, 34, len(moves))
	assert.Equal(t, "d2f3 d2e4 d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6e4 f6d7 f6g4 f6d5 f6h5 f6g8 b4c3 "+
		"b4a3 b4a5 b4c5 b4d6 b7a6 b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7d7 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
}

func TestOnDemand(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()

	var moves = moveslice.NewMoveSlice(100)
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 20, len(moves))
	assert.Equal(t, "e2e4 d2d4 h2h3 a2a3 h2h4 g2g4 f2f4 e2e3 d2d3 c2c4 b2b4 a2a4 g2g3 "+
		"b2b3 f2f3 c2c3 g1f3 b1c3 g1h3 b1a3", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 40, len(moves))
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 d7e5 f3e5 d7b6 e2d2 e1d2 d3d4 h2h3 a2a3 c3c4 h2h4 g2g4 "+
		"a2a4 g2g3 b2b3 e1g1 e1c1 f3d4 d7c5 h1f1 a1d1 a1c1 b5c4 f3g5 e2e3 e2d1 b5a6 b5a4 e2f1 h1g1 a1b1 "+
		"d7f8 f3h4 d7b8 f3g1 e1f1 e1d1", moves.StringUci())
	moves.Clear()

	// 86
	pos = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 86, len(moves))
	moves.Clear()

	// 218
	pos = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 218, len(moves))
	moves.Clear()

}

func Test_movegen_GeneratePseudoLegalMoves(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()
	moves := mg.GeneratePseudoLegalMoves(pos, GenAll)
	assert.Equal(t, 20, len(*moves))
	assert.Equal(t, "e2e4 d2d4 g1f3 b1c3 h2h3 a2a3 h2h4 g2g4 f2f4 e2e3 d2d3 c2c4 b2b4 a2a4 g2g3 b2b3 f2f3 "+
		"c2c3 g1h3 b1a3", moves.StringUci())
	// l := mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll)
	assert.Equal(t, 40, len(*moves))
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 d7e5 f3e5 d7b6 e2d2 e1d2 e1g1 e1c1 d3d4 f3d4 d7c5 h1f1 a1d1 "+
		"a1c1 b5c4 f3g5 h2h3 e2e3 a2a3 c3c4 h2h4 g2g4 a2a4 e1f1 g2g3 e2d1 b2b3 b5a6 b5a4 e2f1 h1g1 a1b1 e1d1 "+
		"d7f8 f3h4 d7b8 f3g1", moves.StringUci())
	// l = mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	// 86 moves
	pos = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll)
	assert.Equal(t, 86, len(*moves))
	moves.Clear()

	// 218 moves
	pos = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll)
	assert.Equal(t, 218, len(*moves))
	moves.Clear()
}

func Test_movegen_GenerateLegalMoves(t *testing.T) {
	mg := NewMoveGen()

	pos := position.NewPosition()
	moves := mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 20, len(*moves))
	assert.Equal(t, "e2e4 d2d4 g1f3 b1c3 h2h3 a2a3 h2h4 g2g4 f2f4 e2e3 d2d3 c2c4 b2b4 a2a4 g2g3 b2b3 f2f3 "+
		"c2c3 g1h3 b1a3", moves.StringUci())
	moves.Clear()

	pos = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 38, len(*moves))
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 d7e5 f3e5 d7b6 e2d2 e1d2 e1c1 d3d4 f3d4 d7c5 h1f1 a1d1 a1c1 b5c4 "+
		"f3g5 h2h3 e2e3 a2a3 c3c4 h2h4 g2g4 a2a4 g2g3 e2d1 b2b3 b5a6 b5a4 e2f1 h1g1 a1b1 e1d1 d7f8 f3h4 d7b8 f3g1",
		moves.StringUci())
	// l = mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	// 86 moves
	pos = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 83, len(*moves))
	assert.Equal(t, "c2b1Q a2b1Q c2b1N a2b1N f4g3 b2a3 f4e3 a8a3 d7e5 g6e5 b2e5 e6e5 c6e4 c4e4 c2b1R a2b1R "+
		"c2b1B a2b1B e8c8 c2c1Q a2a1Q c2c1N a2a1N h8f8 a8d8 a8c8 e8d8 d7b6 g6e7 e6f7 e6e7 a8a7 a8a6 a8a5 a8a4 h7h6 "+
		"d7f8 d7f6 d7c5 g6f8 e6g8 e6f6 e6d6 e6f5 e6d5 e6g4 e6h3 c6d6 c6b6 c6a6 c6d5 c6c5 c6b5 c6a4 c4a6 c4d5 c4c5 "+
		"c4b5 c4b4 c4a4 c4b3 c4e2 c4f1 b2d4 b2c3 b2c1 b2a1 d7b8 g6h4 c4d4 c4d3 c4c3 h7h5 b7b5 h8g8 a8b8 b7b6 e8f7 "+
		"f4f3 c2c1R a2a1R c2c1B a2a1B", moves.StringUci())
	moves.Clear()

	// 218 moves
	pos = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 218, len(*moves))
	moves.Clear()
}

func TestHasLegalMoves(t *testing.T) {

	mg := NewMoveGen()

	// check mate position
	pos := position.NewPositionFen("rn2kbnr/pbpp1ppp/8/1p2p1q1/4K3/3P4/PPP1PPPP/RNBQ1BNR w kq -")
	assert.False(t, mg.HasLegalMove(pos))
	assert.True(t, pos.HasCheck())

	// stale mate position
	pos = position.NewPositionFen("7k/5K2/6Q1/8/8/8/8/8 b - -")
	assert.False(t, mg.HasLegalMove(pos))
	assert.False(t, pos.HasCheck())

	// only en passant
	pos = position.NewPositionFen("8/8/8/8/5Pp1/6P1/7k/K3BQ2 b - f3")
	assert.True(t, mg.HasLegalMove(pos))
	assert.False(t, pos.HasCheck())
}

func TestMovegen_GetMoveFromUci(t *testing.T) {

	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	mg := NewMoveGen()

	// invalid pattern
	move := mg.GetMoveFromUci(pos, "8888")
	assert.Equal(t, MoveNone, move)

	// valid move
	move = mg.GetMoveFromUci(pos, "b7b5")
	assert.Equal(t, CreateMove(SqB7, SqB5, Normal, PtNone), move)

	// invalid move
	move = mg.GetMoveFromUci(pos, "a7a5")
	assert.Equal(t, MoveNone, move)

	// valid promotion
	move = mg.GetMoveFromUci(pos, "a2a1Q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid promotion (we allow lower case promotions)
	move = mg.GetMoveFromUci(pos, "a2a1q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid castling
	move = mg.GetMoveFromUci(pos, "e8c8")
	assert.Equal(t, CreateMove(SqE8, SqC8, Castling, PtNone), move)

	// invalid castling
	move = mg.GetMoveFromUci(pos, "e8g8")
	assert.Equal(t, MoveNone, move)
}

func TestMovegen_GetMoveFromSan(t *testing.T) {

	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	mg := NewMoveGen()

	// invalid pattern
	move := mg.GetMoveFromSan(pos, "33")
	assert.Equal(t, MoveNone, move)

	// valid move
	move = mg.GetMoveFromSan(pos, "b5")
	assert.Equal(t, CreateMove(SqB7, SqB5, Normal, PtNone), move)

	// invalid move
	move = mg.GetMoveFromSan(pos, "a5")
	assert.Equal(t, MoveNone, move)

	// valid promotion
	move = mg.GetMoveFromSan(pos, "a1Q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid promotion (we allow lower case promotions)
	move = mg.GetMoveFromSan(pos, "a1q")
	assert.Equal(t, MoveNone, move)

	// valid castling
	move = mg.GetMoveFromSan(pos, "O-O-O")
	assert.Equal(t, CreateMove(SqE8, SqC8, Castling, PtNone), move)

	// invalid castling
	move = mg.GetMoveFromSan(pos, "O-O")
	assert.Equal(t, MoveNone, move)

	// ambiguous
	move = mg.GetMoveFromSan(pos, "Ne5")
	assert.Equal(t, MoveNone, move)
	move = mg.GetMoveFromSan(pos, "Nde5")
	assert.Equal(t, CreateMove(SqD7, SqE5, Normal, PtNone), move)
	move = mg.GetMoveFromSan(pos, "Nge5")
	assert.Equal(t, CreateMove(SqG6, SqE5, Normal, PtNone), move)
	move = mg.GetMoveFromSan(pos, "N7e5")
	assert.Equal(t, CreateMove(SqD7, SqE5, Normal, PtNone), move)
	move = mg.GetMoveFromSan(pos, "N6e5")
	assert.Equal(t, CreateMove(SqG6, SqE5, Normal, PtNone), move)
	move = mg.GetMoveFromSan(pos, "ab1Q")
	assert.Equal(t, CreateMove(SqA2, SqB1, Promotion, Queen), move)
	move = mg.GetMoveFromSan(pos, "cb1Q")
	assert.Equal(t, CreateMove(SqC2, SqB1, Promotion, Queen), move)
}

func TestOnDemandKillerPv(t *testing.T) {

	mg := NewMoveGen()
	var moves = moveslice.NewMoveSlice(100)

	// 86
	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	mg.StoreKiller(mg.GetMoveFromUci(pos, "g6h4"))
	mg.StoreKiller(mg.GetMoveFromUci(pos, "b7b6"))
	mg.SetPvMove(mg.GetMoveFromUci(pos, "a2b1Q")) // changes c2b1Q a2b1Q to a2b1Q c2b1Q
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 86, len(moves))
	assert.Equal(t, "a2b1Q c2b1Q c2b1N a2b1N f4g3 f4e3 c2b1R a2b1R c2b1B a2b1B b2a3 a8a3 d7e5 "+
		"g6e5 b2e5 e6e5 c6e4 c4e4 b7b6 c2c1Q a2a1Q c2c1N a2a1N h7h6 h7h5 b7b5 f4f3 c2c1R a2a1R c2c1B "+
		"a2a1B e8g8 e8c8 g6h4 h8f8 a8d8 a8c8 d7b6 g6e7 e6f7 e6e7 a8a7 a8a6 a8a5 a8a4 d7f8 d7f6 d7c5 "+
		"g6f8 e6g8 e6f6 e6d6 e6f5 e6d5 e6g4 e6h3 c6d6 c6b6 c6a6 c6d5 c6c5 c6b5 c6a4 c4a6 c4d5 c4c5 "+
		"c4b5 c4b4 c4a4 c4b3 c4e2 c4f1 b2d4 b2c3 b2c1 b2a1 d7b8 c4d4 c4d3 c4c3 h8g8 a8b8 e8f8 e8d8 "+
		"e8f7 e8e7", moves.StringUci())
	// c2b1Q >a2b1Q< c2b1N a2b1N f4g3 f4e3 c2b1R a2b1R c2b1B a2b1B b2a3 a8a3 d7e5 g6e5 b2e5 e6e5 c6e4
	// c4e4 c2c1Q a2a1Q c2c1N a2a1N h7h6 h7h5 b7b5 >b7b6< f4f3 c2c1R a2a1R c2c1B a2a1B e8g8 e8c8 h8f8
	// a8d8 a8c8 d7b6 g6e7 e6f7 e6e7 a8a7 a8a6 a8a5 a8a4 d7f8 d7f6 d7c5 g6f8 e6g8 e6f6 e6d6 e6f5
	// e6d5 e6g4 e6h3 c6d6 c6b6 c6a6 c6d5 c6c5 c6b5 c6a4 c4a6 c4d5 c4c5 c4b5 c4b4 c4a4 c4b3 c4e2
	// c4f1 b2d4 b2c3 b2c1 b2a1 d7b8 >g6h4< c4d4 c4d3 c4c3 h8g8 a8b8 e8f8 e8d8 e8f7 e8e7
	moves.Clear()

	// 48 kiwipete
	pos = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")
	mg.StoreKiller(mg.GetMoveFromUci(pos, "d2g5"))
	mg.StoreKiller(mg.GetMoveFromUci(pos, "b2b3"))
	mg.SetPvMove(mg.GetMoveFromUci(pos, "e2a6"))
	for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
		moves.PushBack(move)
	}
	assert.Equal(t, 48, len(moves))
	assert.Equal(t, "e2a6 d5e6 g2h3 e5f7 e5d7 e5g6 f3f6 f3h3 b2b3 d5d6 a2a3 g2g4 a2a4 g2g3 e1g1 "+
		"e1c1 d2g5 e5c4 e5d3 h1f1 a1d1 a1c1 e5c6 e2c4 e2d3 d2f4 d2e3 f3f4 f3e3 f3d3 c3b5 e2b5 f3f5 e5g4 "+
		"f3g4 f3g3 f3h5 e2d1 d2h6 h1g1 a1b1 c3a4 c3d1 c3b1 e2f1 d2c1 e1f1 e1d1", moves.StringUci())
	// d5e6 g2h3 >e2a6< e5f7 e5d7 e5g6 f3f6 f3h3 d5d6 a2a3 g2g4 a2a4 g2g3 b2b3 >e1g1< e1c1 e5c4 e5d3 h1f1
	// a1d1 a1c1 e5c6 e2c4 e2d3 d2f4 d2e3 f3f4 f3e3 f3d3 c3b5 e2b5 >d2g5< f3f5 e5g4 f3g4 f3g3 f3h5 e2d1
	// d2h6 h1g1 a1b1 c3a4 c3d1 c3b1 e2f1 d2c1 e1f1 e1d1
	moves.Clear()

}

// MoveList
// GeneratePseudoLegalMoves took 6.948.781.000 ns for 1.000.000 iterations
// GeneratePseudoLegalMoves took 6.948 ns
// GeneratePseudoLegalMoves 86.000.000 generated 12.376.271 mps
// Move Array
// GeneratePseudoLegalMoves took 3.050.877.900 ns for 1.000.000 iterations
// GeneratePseudoLegalMoves took 3.050 ns
// GeneratePseudoLegalMoves 86.000.000 generated 28.188.607 mps
// MoveArray with ForEach cleanup of moves
// GeneratePseudoLegalMoves took 1.965.578.500 ns for 1.000.000 iterations
// GeneratePseudoLegalMoves took 1.965 ns
// GeneratePseudoLegalMoves 86.000.000 generated 43.753.022 mps
//noinspection GoUnhandledErrorResult
func TestTimingPseudoMoveGen(t *testing.T) {
	out := message.NewPrinter(language.German)

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves := mg.GeneratePseudoLegalMoves(pos, GenAll)

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			moves.Clear()
			moves = mg.GeneratePseudoLegalMoves(pos, GenAll)
		}
		elapsed := time.Since(start)
		out.Printf("GeneratePseudoLegalMoves took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		generated := uint64(len(*moves)) * iterations
		out.Printf("%d moves generated in %d ns: %d mps\n",
			generated,
			elapsed.Nanoseconds()/int64(iterations),
			(generated*uint64(1_000_000_000))/uint64(elapsed.Nanoseconds()))
	}
}

func TestTimingOnDemandMoveGen(t *testing.T) {
	out := message.NewPrinter(language.German)

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		generated := uint64(0)
		for i := uint64(0); i < iterations; i++ {
			generated = 0
			mg.ResetOnDemand()
			for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
				generated++
			}
		}
		elapsed := time.Since(start)
		out.Printf("GeneratePseudoLegalMoves took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		generated = generated * iterations
		out.Printf("%d moves generated in %d ns: %d mps\n",
			generated,
			elapsed.Nanoseconds()/int64(iterations),
			(generated*uint64(1_000_000_000))/uint64(elapsed.Nanoseconds()))
	}
}

func TestTimingOnDemandRealMoveGen(t *testing.T) {
	out := message.NewPrinter(language.German)

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	pos := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	k1 := mg.GetMoveFromUci(pos, "g6h4")
	k2 := mg.GetMoveFromUci(pos, "b7b6")
	pv := mg.GetMoveFromUci(pos, "a2b1Q")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		generated := uint64(0)
		for i := uint64(0); i < iterations; i++ {
			generated = 0
			mg.ResetOnDemand()
			mg.StoreKiller(k1)
			mg.StoreKiller(k2)
			mg.SetPvMove(pv)
			for move := mg.GetNextMove(pos, GenAll); move != MoveNone; move = mg.GetNextMove(pos, GenAll) {
				generated++
			}
		}
		elapsed := time.Since(start)
		out.Printf("GeneratePseudoLegalMoves took %d ns for %d iterations\n", elapsed.Nanoseconds(), iterations)
		generated = generated * iterations
		out.Printf("%d moves generated in %d ns: %d mps\n",
			generated,
			elapsed.Nanoseconds()/int64(iterations),
			(generated*uint64(1_000_000_000))/uint64(elapsed.Nanoseconds()))
	}
}

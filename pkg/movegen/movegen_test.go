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

package movegen

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/pkg/position"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

var logTest *logging.Logger

// make tests run in the projects root directory.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests.
func TestMain(m *testing.M) {
	config.Setup()
	logTest = myLogging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func TestMovegenString(t *testing.T) {
	mg := NewMoveGen()
	out.Println(mg.String())
	assert.EqualValues(t, "MoveGen: { OnDemand Stage: { 0 }, PV Move: Move: { MoveNone } Killer Move 1: "+
		"Move: { MoveNone } Killer Move 2: Move: { MoveNone } }", mg.String())
}

func TestMovegenGeneratePawnMoves(t *testing.T) {

	mg := NewMoveGen()
	pos, _ := position.NewPositionFen("1kr3nr/pp1pP1P1/2p1p3/3P1p2/1n1bP3/2P5/PP3PPP/RNBQKBNR w KQ -")
	moves := moveslice.MoveSlice{}

	mg.generatePawnMoves(pos, GenNonQuiet, false, BbZero, &moves)
	assert.Equal(t, 11, moves.Len())

	moves.Clear()
	mg.generatePawnMoves(pos, GenQuiet, false, BbZero, &moves)
	assert.Equal(t, 14, moves.Len())

	moves.Clear()
	mg.generatePawnMoves(pos, GenAll, false, BbZero, &moves)
	assert.Equal(t, 25, moves.Len())

	// moves.Sort()
	// fmt.Printf("Moves: %d\n", moves.Len())
	// l := moves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("Move: %s\n", moves.At(i))
	// }
}

func TestMovegenGenerateCastling(t *testing.T) {
	mg := NewMoveGen()
	pos, _ := position.NewPositionFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R w KQkq -")
	moves := moveslice.MoveSlice{}

	mg.generateCastling(pos, GenAll, &moves)
	assert.Equal(t, 2, moves.Len())
	assert.Equal(t, "e1g1 e1c1", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbppqppp/1pn2n2/1B2p3/1b2P3/N1PP1N2/PP1BQPPP/R3K2R b KQkq -")
	mg.generateCastling(pos, GenAll, &moves)
	assert.Equal(t, 2, moves.Len())
	assert.Equal(t, "e8g8 e8c8", moves.StringUci())

}

func TestMovegenGenerateKingMoves(t *testing.T) {
	mg := NewMoveGen()
	moves := moveslice.MoveSlice{}

	pos, _ := position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateKingMoves(pos, GenAll, false, BbZero, &moves)
	assert.Equal(t, 3, moves.Len())
	assert.Equal(t, "e1d2 e1d1 e1f1", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateKingMoves(pos, GenAll, false, BbZero, &moves)
	assert.Equal(t, 3, moves.Len())
	assert.Equal(t, "e8d7 e8d8 e8f8", moves.StringUci())
}

func TestMovegenGenerateMoves(t *testing.T) {
	mg := NewMoveGen()
	moves := moveslice.MoveSlice{}

	pos, _ := position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	mg.generateMoves(pos, GenNonQuiet, false, BbZero, &moves)
	assert.Equal(t, 7, moves.Len())
	assert.Equal(t, "f3d2 f3e5 d7e5 d7b6 d7f6 b5c6 e2d2", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(pos, GenQuiet, false, BbZero, &moves)
	assert.Equal(t, 28, moves.Len())
	assert.Equal(t, "d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6g4 f6d5 f6h5 f6g8 b4a3 b4a5 b4c5 b4d6 b7a6 "+
		"b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R b KQkq -")
	mg.generateMoves(pos, GenAll, false, BbZero, &moves)
	assert.Equal(t, 34, moves.Len())
	assert.Equal(t, "d2f3 d2e4 d2b1 d2f1 d2b3 d2c4 c6d4 c6a5 c6b8 c6d8 f6e4 f6d7 f6g4 f6d5 f6h5 f6g8 b4c3 "+
		"b4a3 b4a5 b4c5 b4d6 b7a6 b7c8 a8b8 a8c8 a8d8 h8f8 h8g8 e7d7 e7c5 e7d6 e7e6 e7d8 e7f8", moves.StringUci())
}

func TestOnDemand(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()

	var moves = moveslice.NewMoveSlice(100)
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 20, moves.Len())
	assert.Equal(t, "d2d4 e2e4 a2a3 h2h3 a2a4 b2b4 c2c4 f2f4 g2g4 h2h4 d2d3 e2e3 b2b3 g2g3 c2c3 f2f3 b1c3 "+
		"g1f3 b1a3 g1h3", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 40, moves.Len())
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 f3e5 d7e5 d7b6 e2d2 e1d2 d3d4 a2a3 h2h3 a2a4 g2g4 h2h4 c3c4 b2b3 "+
		"g2g3 e1g1 e1c1 f3d4 d7c5 a1c1 a1d1 h1f1 b5c4 f3g5 e2e3 e2d1 b5a4 b5a6 a1b1 h1g1 e2f1 f3g1 f3h4 d7f8 d7b8 "+
		"e1f1 e1d1", moves.StringUci())
	moves.Clear()

	// 86
	pos, _ = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 86, moves.Len())
	moves.Clear()

	// 218
	pos, _ = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 218, moves.Len())
	moves.Clear()

}

func TestOnDemandPromNonQuiet(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()

	var moves = moveslice.NewMoveSlice(100)
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 20, moves.Len())
	assert.Equal(t, "d2d4 e2e4 a2a3 h2h3 a2a4 b2b4 c2c4 f2f4 g2g4 h2h4 d2d3 e2e3 b2b3 g2g3 c2c3 f2f3 b1c3 "+
		"g1f3 b1a3 g1h3", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 40, moves.Len())
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 f3e5 d7e5 d7b6 e2d2 e1d2 d3d4 a2a3 h2h3 a2a4 g2g4 h2h4 c3c4 b2b3 "+
		"g2g3 e1g1 e1c1 f3d4 d7c5 a1c1 a1d1 h1f1 b5c4 f3g5 e2e3 e2d1 b5a4 b5a6 a1b1 h1g1 e2f1 f3g1 f3h4 d7f8 d7b8 "+
		"e1f1 e1d1", moves.StringUci())
	moves.Clear()

	// 218
	pos, _ = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 218, moves.Len())
	moves.Clear()

	// 86
	pos, _ = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
		// fmt.Printf("%s\n", move.String())
	}
	assert.Equal(t, 86, moves.Len())
	assert.Equal(t, "c2b1Q a2b1Q a2a1Q c2c1Q c2b1N a2b1N f4g3 a2a1N c2c1N f4e3 c2b1R a2b1R c2b1B a2b1B b2a3 a8a3 g6e5 d7e5 b2e5 e6e5 c4e4 c6e4 f4f3 h7h6 b7b5 h7h5 b7b6 a2a1R c2c1R a2a1B c2c1B e8g8 e8c8 d7c5 a8c8 a8d8 h8f8 d7f6 b2d4 g6e7 d7b6 b2c3 c4c5 c4d5 c6c5 c6d5 c6d6 e6d5 e6f5 e6d6 e6f6 e6e7 e6f7 c4d4 a8a4 a8a5 a8a6 a8a7 c4e2 c4b3 c4c3 c4d3 c4b4 c4b5 c6b5 c6b6 e6g4 c4a4 c6a4 b2c1 a8b8 h8g8 c4f1 c4a6 c6a6 e6h3 e6g8 g6f8 d7f8 b2a1 d7b8 g6h4 e8f8 e8e7 e8f7 e8d8", moves.StringUci())
	moves.Clear()

	// 86
	mg.ResetOnDemand()
	pos, _ = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
	}
	assert.Equal(t, 86, moves.Len())
	assert.Equal(t, "c2b1Q a2b1Q a2a1Q c2c1Q c2b1N a2b1N f4g3 a2a1N c2c1N f4e3 c2b1R a2b1R c2b1B a2b1B b2a3 a8a3 g6e5 d7e5 b2e5 e6e5 c4e4 c6e4 f4f3 h7h6 b7b5 h7h5 b7b6 a2a1R c2c1R a2a1B c2c1B e8g8 e8c8 d7c5 a8c8 a8d8 h8f8 d7f6 b2d4 g6e7 d7b6 b2c3 c4c5 c4d5 c6c5 c6d5 c6d6 e6d5 e6f5 e6d6 e6f6 e6e7 e6f7 c4d4 a8a4 a8a5 a8a6 a8a7 c4e2 c4b3 c4c3 c4d3 c4b4 c4b5 c6b5 c6b6 e6g4 c4a4 c6a4 b2c1 a8b8 h8g8 c4f1 c4a6 c6a6 e6h3 e6g8 g6f8 d7f8 b2a1 d7b8 g6h4 e8f8 e8e7 e8f7 e8d8", moves.StringUci())
	moves.Clear()

}

func TestMovegenGeneratePseudoLegalMoves(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()
	moves := mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 20, len(*moves))
	assert.Equal(t, "d2d4 e2e4 b1c3 g1f3 a2a3 h2h3 a2a4 b2b4 c2c4 f2f4 g2g4 h2h4 d2d3 e2e3 b2b3 g2g3 c2c3 "+
		"f2f3 b1a3 g1h3", moves.StringUci())
	// l := mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 40, len(*moves))
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 f3e5 d7e5 d7b6 e2d2 e1d2 e1g1 e1c1 d3d4 f3d4 d7c5 a1c1 a1d1 h1f1 "+
		"b5c4 a2a3 h2h3 f3g5 e2e3 a2a4 g2g4 h2h4 c3c4 e1f1 b2b3 g2g3 e2d1 b5a4 b5a6 a1b1 h1g1 e2f1 e1d1 f3g1 f3h4 "+
		"d7f8 d7b8", moves.StringUci())
	// l = mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	// 86 moves
	pos, _ = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 86, len(*moves))
	moves.Clear()

	// 218 moves
	pos, _ = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 218, len(*moves))
	moves.Clear()
}

func TestMovegenGenerateLegalMoves(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition()
	moves := mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 20, len(*moves))
	assert.Equal(t, "d2d4 e2e4 b1c3 g1f3 a2a3 h2h3 a2a4 b2b4 c2c4 f2f4 g2g4 h2h4 d2d3 e2e3 b2b3 g2g3 c2c3 "+
		"f2f3 b1a3 g1h3", moves.StringUci())
	moves.Clear()

	pos, _ = position.NewPositionFen("r3k2r/pbpNqppp/1pn2n2/1B2p3/1b2P3/2PP1N2/PP1nQPPP/R3K2R w KQkq -")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 38, len(*moves))
	assert.Equal(t, "c3b4 d7f6 f3d2 b5c6 f3e5 d7e5 d7b6 e2d2 e1d2 e1c1 d3d4 f3d4 d7c5 a1c1 a1d1 h1f1 b5c4 "+
		"a2a3 h2h3 f3g5 e2e3 a2a4 g2g4 h2h4 c3c4 b2b3 g2g3 e2d1 b5a4 b5a6 a1b1 h1g1 e2f1 e1d1 f3g1 f3h4 d7f8 d7b8",
		moves.StringUci())
	// l = mg.pseudoLegalMoves.Len()
	// for i := 0; i < l; i++ {
	// 	fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	// }
	moves.Clear()

	// 86 moves
	pos, _ = position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 83, len(*moves))
	assert.Equal(t, "c2b1Q a2b1Q a2a1Q c2c1Q c2b1N a2b1N f4g3 a2a1N c2c1N f4e3 b2a3 a8a3 g6e5 d7e5 b2e5 e6e5 c4e4 c6e4 c2b1R a2b1R c2b1B a2b1B e8c8 d7c5 a8c8 a8d8 h8f8 d7f6 b2d4 f4f3 h7h6 g6e7 d7b6 b2c3 c4c5 c4d5 c6c5 c6d5 c6d6 e6d5 e6f5 e6d6 e6f6 e6e7 e6f7 c4d4 b7b5 h7h5 a8a4 a8a5 a8a6 a8a7 c4e2 c4b3 c4c3 c4d3 c4b4 c4b5 c6b5 c6b6 e6g4 b7b6 c4a4 c6a4 b2c1 a8b8 h8g8 c4f1 c4a6 c6a6 e6h3 e6g8 e8f7 e8d8 g6f8 d7f8 b2a1 d7b8 g6h4 a2a1R c2c1R a2a1B c2c1B", moves.StringUci())
	moves.Clear()

	// 218 moves
	pos, _ = position.NewPositionFen("R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - -")
	moves = mg.GenerateLegalMoves(pos, GenAll)
	assert.Equal(t, 218, len(*moves))
	moves.Clear()
}

func TestHasLegalMoves(t *testing.T) {

	mg := NewMoveGen()

	// check mate position
	pos, _ := position.NewPositionFen("rn2kbnr/pbpp1ppp/8/1p2p1q1/4K3/3P4/PPP1PPPP/RNBQ1BNR w kq -")
	assert.False(t, mg.HasLegalMove(pos))
	assert.True(t, pos.HasCheck())

	// stale mate position
	pos, _ = position.NewPositionFen("7k/5K2/6Q1/8/8/8/8/8 b - -")
	assert.False(t, mg.HasLegalMove(pos))
	assert.False(t, pos.HasCheck())

	// only en passant
	pos, _ = position.NewPositionFen("8/8/8/8/5Pp1/6P1/7k/K3BQ2 b - f3")
	assert.True(t, mg.HasLegalMove(pos))
	assert.False(t, pos.HasCheck())
}

func TestMovegenGetMoveFromUci(t *testing.T) {

	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	mg := NewMoveGen()

	// invalid pattern
	move, _ := mg.GetMoveFromUci(pos, "8888")
	assert.Equal(t, MoveNone, move)

	// valid move
	move, _ = mg.GetMoveFromUci(pos, "b7b5")
	assert.Equal(t, CreateMove(SqB7, SqB5, Normal, PtNone), move)

	// invalid move
	move, _ = mg.GetMoveFromUci(pos, "a7a5")
	assert.Equal(t, MoveNone, move)

	// valid promotion
	move, _ = mg.GetMoveFromUci(pos, "a2a1Q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid promotion (we allow lower case promotions)
	move, _ = mg.GetMoveFromUci(pos, "a2a1q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid castling
	move, _ = mg.GetMoveFromUci(pos, "e8c8")
	assert.Equal(t, CreateMove(SqE8, SqC8, Castling, PtNone), move)

	// invalid castling
	move, _ = mg.GetMoveFromUci(pos, "e8g8")
	assert.Equal(t, MoveNone, move)
}

func TestMovegenGetMoveFromSan(t *testing.T) {

	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	mg := NewMoveGen()

	// invalid pattern
	move, _ := mg.GetMoveFromSan(pos, "33")
	assert.Equal(t, MoveNone, move)

	// valid move
	move, _ = mg.GetMoveFromSan(pos, "b5")
	assert.Equal(t, CreateMove(SqB7, SqB5, Normal, PtNone), move)

	// invalid move
	move, _ = mg.GetMoveFromSan(pos, "a5")
	assert.Equal(t, MoveNone, move)

	// valid promotion
	move, _ = mg.GetMoveFromSan(pos, "a1Q")
	assert.Equal(t, CreateMove(SqA2, SqA1, Promotion, Queen), move)

	// valid promotion (we allow lower case promotions)
	move, _ = mg.GetMoveFromSan(pos, "a1q")
	assert.Equal(t, MoveNone, move)

	// valid castling
	move, _ = mg.GetMoveFromSan(pos, "O-O-O")
	assert.Equal(t, CreateMove(SqE8, SqC8, Castling, PtNone), move)

	// invalid castling
	move, _ = mg.GetMoveFromSan(pos, "O-O")
	assert.Equal(t, MoveNone, move)

	// ambiguous
	move, _ = mg.GetMoveFromSan(pos, "Ne5")
	assert.Equal(t, MoveNone, move)
	move, _ = mg.GetMoveFromSan(pos, "Nde5")
	assert.Equal(t, CreateMove(SqD7, SqE5, Normal, PtNone), move)
	move, _ = mg.GetMoveFromSan(pos, "Nge5")
	assert.Equal(t, CreateMove(SqG6, SqE5, Normal, PtNone), move)
	move, _ = mg.GetMoveFromSan(pos, "N7e5")
	assert.Equal(t, CreateMove(SqD7, SqE5, Normal, PtNone), move)
	move, _ = mg.GetMoveFromSan(pos, "N6e5")
	assert.Equal(t, CreateMove(SqG6, SqE5, Normal, PtNone), move)
	move, _ = mg.GetMoveFromSan(pos, "ab1Q")
	assert.Equal(t, CreateMove(SqA2, SqB1, Promotion, Queen), move)
	move, _ = mg.GetMoveFromSan(pos, "cb1Q")
	assert.Equal(t, CreateMove(SqC2, SqB1, Promotion, Queen), move)
}

func TestOnDemandKillerPv(t *testing.T) {
	mg := NewMoveGen()
	var moves = moveslice.NewMoveSlice(100)

	// 86
	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moveFromUci, _ := mg.GetMoveFromUci(pos, "g6h4")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "b7b6")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "a2b1Q")
	mg.SetPvMove(moveFromUci) // changes c2b1Q a2b1Q to a2b1Q c2b1Q
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
		// fmt.Println(move.String())
	}
	assert.Equal(t, 86, moves.Len())
	assert.Equal(t, "a2b1Q c2b1Q a2a1Q c2c1Q c2b1N a2b1N a2a1N c2c1N f4g3 f4e3 c2b1R a2b1R c2b1B a2b1B b2a3 a8a3 g6e5 d7e5 b2e5 e6e5 c4e4 c6e4 b7b6 f4f3 h7h6 b7b5 h7h5 a2a1R c2c1R a2a1B c2c1B e8g8 e8c8 g6h4 d7c5 a8c8 a8d8 h8f8 d7f6 b2d4 b2c3 c4c5 c4d5 c6c5 c6d5 c6d6 e6d5 e6f5 e6d6 e6f6 g6e7 d7b6 e6e7 e6f7 c4d4 a8a4 a8a5 a8a6 a8a7 c4e2 c4b3 c4c3 c4d3 c4b4 c4b5 c6b5 c6b6 e6g4 c4a4 c6a4 a8b8 h8g8 b2c1 c4f1 c4a6 c6a6 e6h3 e6g8 g6f8 d7f8 b2a1 d7b8 e8f8 e8e7 e8f7 e8d8", moves.StringUci())
	moves.Clear()

	// 48 kiwipete
	pos, _ = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")
	moveFromUci, _ = mg.GetMoveFromUci(pos, "d2g5")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "b2b3")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "e2a6")
	mg.SetPvMove(moveFromUci)
	for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
		moves.PushBack(move)
		// fmt.Println(move.String())
	}
	assert.Equal(t, 48, moves.Len())
	assert.Equal(t, "e2a6 g2h3 d5e6 e5g6 e5d7 e5f7 f3f6 f3h3 b2b3 a2a3 d5d6 a2a4 g2g4 g2g3 e1g1 e1c1 d2g5 "+
		"e5d3 e5c4 a1c1 a1d1 h1f1 e5c6 d2e3 d2f4 e2d3 e2c4 c3b5 e2b5 f3d3 f3e3 f3f4 f3f5 e5g4 f3g3 f3g4 f3h5 d2h6 "+
		"e2d1 a1b1 h1g1 c3b1 c3d1 c3a4 d2c1 e2f1 e1f1 e1d1", moves.StringUci())
	moves.Clear()

}

func TestPseudoLegalPVKiller(t *testing.T) {

	mg := NewMoveGen()
	var moves = moveslice.NewMoveSlice(100)

	// 86
	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	moveFromUci, _ := mg.GetMoveFromUci(pos, "a2b1Q")
	mg.SetPvMove(moveFromUci) // changes c2b1Q a2b1Q to a2b1Q c2b1Q
	moveFromUci, _ = mg.GetMoveFromUci(pos, "g6h4")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "b7b6")
	mg.StoreKiller(moveFromUci)
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 86, moves.Len())
	assert.Equal(t, "a2b1Q c2b1Q a2a1Q c2c1Q c2b1N a2b1N f4g3 a2a1N c2c1N f4e3 b2a3 a8a3 g6e5 d7e5 b2e5 e6e5 c4e4 c6e4 c2b1R a2b1R c2b1B a2b1B b7b6 g6h4 e8g8 e8c8 d7c5 f4f3 a8c8 a8d8 h8f8 d7f6 b2d4 h7h6 b2c3 c4c5 c4d5 c6c5 c6d5 c6d6 e6d5 e6f5 e6d6 e6f6 g6e7 d7b6 e6e7 e6f7 c4d4 b7b5 h7h5 a8a4 a8a5 a8a6 a8a7 c4e2 c4b3 c4c3 c4d3 c4b4 c4b5 c6b5 c6b6 e6g4 e8f8 c4a4 c6a4 a8b8 h8g8 b2c1 c4f1 c4a6 c6a6 e6h3 e6g8 e8e7 e8f7 g6f8 d7f8 b2a1 e8d8 d7b8 a2a1R c2c1R a2a1B c2c1B", moves.StringUci())
	moves.Clear()

	// 48 kiwipete
	pos, _ = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")
	moveFromUci, _ = mg.GetMoveFromUci(pos, "e2a6")
	mg.SetPvMove(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "d2g5")
	mg.StoreKiller(moveFromUci)
	moveFromUci, _ = mg.GetMoveFromUci(pos, "b2b3")
	mg.StoreKiller(moveFromUci)
	moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	assert.Equal(t, 48, moves.Len())
	assert.Equal(t, "e2a6 g2h3 d5e6 e5g6 e5d7 e5f7 f3f6 f3h3 b2b3 d2g5 e1g1 e1c1 e5d3 e5c4 a1c1 a1d1 h1f1 e5c6 d2e3 d2f4 e2d3 e2c4 a2a3 d5d6 c3b5 e2b5 f3d3 f3e3 f3f4 f3f5 a2a4 g2g4 e1f1 e5g4 f3g3 f3g4 g2g3 f3h5 d2h6 e2d1 a1b1 h1g1 e1d1 c3b1 c3d1 c3a4 d2c1 e2f1", moves.StringUci())
	moves.Clear()

}

func TestEvasion(t *testing.T) {
	mg := NewMoveGen()
	var p *position.Position
	var pseudoLegalMoves, evasionMoves, legalMoves *moveslice.MoveSlice

	// TODO: real tests

	p = position.NewPosition("r3k2r/1pp4p/2q1qNn1/3nP3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq -")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()

	p = position.NewPosition("5k2/8/8/8/8/8/6p1/3K1R2 b - -")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()

	p = position.NewPosition("5k2/8/8/8/8/6p1/5R2/3K4 b - -")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()

	p = position.NewPosition("8/8/8/3k4/4Pp2/8/8/3K4 b - e3")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()

	p = position.NewPosition("8/8/8/3k2n1/8/8/6B1/3K4 b - -")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()

	p = position.NewPosition("5k2/3N4/8/8/8/8/6p1/3K1R2 b - - 1 1 ")
	pseudoLegalMoves = mg.GeneratePseudoLegalMoves(p, GenAll, false).Clone()
	evasionMoves = mg.GeneratePseudoLegalMoves(p, GenAll, true).Clone()
	legalMoves = mg.GenerateLegalMoves(p, GenAll).Clone()
	out.Printf("PseudoLegal: %3d %s\n", pseudoLegalMoves.Len(), pseudoLegalMoves.StringUci())
	out.Printf("Evasion    : %3d %s\n", evasionMoves.Len(), evasionMoves.StringUci())
	out.Printf("Legal      : %3d %s\n", legalMoves.Len(), legalMoves.StringUci())
	out.Println()
}

func TestTimingPseudoMoveGen(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	// kiwi pete
	pos := position.NewPosition("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - ")
	moves := mg.GeneratePseudoLegalMoves(pos, GenAll, false)

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			moves.Clear()
			moves = mg.GeneratePseudoLegalMoves(pos, GenAll, false)
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

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		generated := uint64(0)
		for i := uint64(0); i < iterations; i++ {
			generated = 0
			mg.ResetOnDemand()
			for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
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

// GeneratePseudoLegalMoves took 1.655.000.800 ns for 1.000.000 iterations
// 86.000.000 moves generated in 1.655 ns: 51.963.721 mps.
func TestTimingOnDemandRealMoveGen(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("./bin")).Stop()

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const rounds = 5
	const iterations uint64 = 1_000_000

	mg := NewMoveGen()
	pos, _ := position.NewPositionFen("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	k1, _ := mg.GetMoveFromUci(pos, "g6h4")
	k2, _ := mg.GetMoveFromUci(pos, "b7b6")
	pv, _ := mg.GetMoveFromUci(pos, "a2b1Q")

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
			for move := mg.GetNextMove(pos, GenAll, false); move != MoveNone; move = mg.GetNextMove(pos, GenAll, false) {
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

func TestTimingGenerateMovesNew(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http :8080 ./main ./prof.null/cpu.pprof

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	mg := NewMoveGen()
	p := position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
	result := Value(0)

	const rounds = 5
	const iterations uint64 = 10_000_000

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			mg.pseudoLegalMoves.Clear()
			mg.generateMoves(p, GenAll, false, BbZero, mg.pseudoLegalMoves)
		}
		elapsed := time.Since(start)
		out.Printf("Test took %s for %d iterations\n", elapsed, iterations)
		out.Printf("Test took %d ns per iteration\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Iterations per sec %d\n", int64(iterations*1e9)/elapsed.Nanoseconds())
	}
	_ = result
}

func TestDebug(t *testing.T) {

	mg := NewMoveGen()

	pos := position.NewPosition("1k1r4/pp1b1R2/3q2pp/4p3/2B5/4Q3/PPP2B2/2K5 b - -")
	moves := mg.GeneratePseudoLegalMoves(pos, GenAll, false)
	l := mg.pseudoLegalMoves.Len()
	for i := 0; i < l; i++ {
		fmt.Printf("%d. %s\n", i+1, moves.At(i).String())
	}
	moves.Clear()
}

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
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	logging2 "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	. "github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var logTest *logging2.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	out.Println("Test Main Setup Tests ====================")
	Setup()
	logTest = logging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

//noinspection GoStructInitializationWithoutFieldNames
func TestEvaluator_valueFromScore(t *testing.T) {
	e := NewEvaluator()
	v:= Value(0)

	e.gamePhaseFactor = 1.0
	e.score = Score{ 10, 0 }
	v = e.value()
	assert.EqualValues(t, 10, v)
	e.gamePhaseFactor = 0.0
	v = e.value()
	assert.EqualValues(t, 0, v)
	e.gamePhaseFactor = 0.5
	v = e.value()
	assert.EqualValues(t, 5, v)

	e.gamePhaseFactor = 1.0
	e.score = Score{ 50, 50 }
	v = e.value()
	assert.EqualValues(t, 50, v)
	e.gamePhaseFactor = 0.0
	v = e.value()
	assert.EqualValues(t, 50, v)
	e.gamePhaseFactor = 0.5
	v = e.value()
	assert.EqualValues(t, 50, v)
}

func TestLazyEval(t *testing.T) {
	e := NewEvaluator()
	Settings.Eval.Tempo = 0
	Settings.Eval.UseLazyEval = true
	Settings.Eval.UseAttacksInEval = false
	Settings.Eval.UseAdvancedPieceEval = false
	Settings.Eval.UseKingEval = false
	p := position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/B5R1/pbp2PPP/1R4K1 b kq e3")
	value := e.Evaluate(p)
	out.Println(value)
	assert.EqualValues(t, 3332, value)
	p = position.NewPosition("5r1k/1q6/8/8/8/8/6P1/7K b - - 0 1 ")
	value = e.Evaluate(p)
	out.Println(value)
	assert.EqualValues(t, 1293, value)
}

// func TestEvaluator_evalPieceKnights(t *testing.T) {
// 	e := NewEvaluator()
// 	Settings.Eval.Tempo = 0
// 	p := position.NewPosition("rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w KQkq -")
// 	var score *Score
//
// 	// Knights
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Knight)
// 	out.Printf("White Knight: %s\n", score)
// 	assert.EqualValues(t, Score{}, *score)
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 6, e.mobility[White])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[White][Knight].StringBoard())
// 	assert.EqualValues(t, SqA3.Bb()|SqC3.Bb()|SqF3.Bb()|SqH3.Bb()|SqD2.Bb()|SqE2.Bb(), e.attacks[White][Knight])
//
// 	score = e.evalPiece(Black, Knight)
// 	out.Printf("Black Knight: %s\n", score)
// 	assert.EqualValues(t, Score{}, *score)
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 6, e.mobility[Black])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[Black][Knight].StringBoard())
// 	assert.EqualValues(t, SqA6.Bb()|SqC6.Bb()|SqF6.Bb()|SqH6.Bb()|SqD7.Bb()|SqE7.Bb(), e.attacks[Black][Knight])
// }
//
// func TestEvaluator_evalPieceBishop(t *testing.T) {
// 	e := NewEvaluator()
// 	Settings.Eval.Tempo = 0
// 	var score *Score
//
// 	// Bishop
// 	p := position.NewPosition("rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w KQkq -")
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Bishop)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[White][Bishop].StringBoard(), e.attacks[White][Bishop])
// 	out.Printf("White Bishop: %s\n", score)
// 	assert.EqualValues(t, Score{20,20}, *score)
// 	assert.EqualValues(t, 20, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 14, e.mobility[White])
// 	assert.EqualValues(t, 142121081854464, e.attacks[White][Bishop])
//
// 	score = e.evalPiece(Black, Bishop)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[White][Bishop].StringBoard(), e.attacks[White][Bishop])
// 	out.Printf("Black Bishop: %s\n", score)
// 	assert.EqualValues(t, Score{20,20}, *score)
// 	assert.EqualValues(t, 20, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 14, e.mobility[Black])
// 	assert.EqualValues(t, 25501128917581824, e.attacks[Black][Bishop])
//
// 	p = position.NewPosition("4k3/1p3bb1/p1p3p1/8/3P4/1P6/2PPP1P1/2B1KB2 w - -")
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Bishop)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[White][Bishop].StringBoard(), e.attacks[White][Bishop])
// 	out.Printf("White Bishop: %s\n", score)
// 	assert.EqualValues(t, Score{-5, -50}, *score)
// 	assert.EqualValues(t, -41, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 2, e.mobility[White])
// 	assert.EqualValues(t, 66048, e.attacks[White][Bishop])
//
// 	score = e.evalPiece(Black, Bishop)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[Black][Bishop].StringBoard(), e.attacks[Black][Bishop])
// 	out.Printf("Black Bishop: %s\n", score)
// 	assert.EqualValues(t, Score{95,0}, *score)
// 	assert.EqualValues(t, 15, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 11, e.mobility[Black])
// 	assert.EqualValues(t, Bitboard(16141094681823019008), e.attacks[Black][Bishop])
// }
//
// func TestEvaluator_evalPieceRook(t *testing.T) {
// 	e := NewEvaluator()
// 	Settings.Eval.Tempo = 0
// 	var score *Score
//
// 	// Rook
// 	p := position.NewPosition("rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w KQkq -")
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Rook)
// 	out.Printf("White Rook: %s\n", score)
// 	assert.EqualValues(t, Score{50,0}, *score)
// 	assert.EqualValues(t, 50, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 14, e.mobility[White])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[White][Rook].StringBoard())
// 	assert.EqualValues(t, uint64(9331882296111890688), e.attacks[White][Rook])
//
// 	score = e.evalPiece(Black, Rook)
// 	out.Printf("Black Rook: %s\n", score)
// 	assert.EqualValues(t, Score{50,0}, *score)
// 	assert.EqualValues(t, 50, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 14, e.mobility[Black])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[Black][Rook].StringBoard())
// 	assert.EqualValues(t, uint64(36452665219187073), e.attacks[Black][Rook])
//
// 	p = position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Rook)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[White][Rook].StringBoard(), e.attacks[White][Rook])
// 	out.Printf("White Rook: %s\n", score)
// 	assert.EqualValues(t, Score{-15,0}, *score)
// 	assert.EqualValues(t, -15, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 9, e.mobility[White])
// 	assert.EqualValues(t, uint64(282578800148834), e.attacks[White][Rook])
//
// 	score = e.evalPiece(Black, Rook)
// 	fmt.Printf("Attacks     : \n%s\n%d\n", e.attacks[Black][Rook].StringBoard(), e.attacks[Black][Rook])
// 	out.Printf("Black Rook: %s\n", score)
// 	assert.EqualValues(t, Score{-34,6}, *score)
// 	assert.EqualValues(t, -34, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 3, e.mobility[Black])
// 	assert.EqualValues(t, uint64(7061644215716937728), e.attacks[Black][Rook])
// }
//
// func TestEvaluator_evalPieceQueen(t *testing.T) {
// 	e := NewEvaluator()
// 	Settings.Eval.Tempo = 0
// 	p := position.NewPosition("rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w KQkq -")
// 	var score *Score
//
// 	// Queen
// 	e.InitEval(p)
// 	score = e.evalPiece(White, Queen)
// 	out.Printf("White Queen: %s\n", score)
// 	assert.EqualValues(t, Score{10,10}, *score)
// 	assert.EqualValues(t, 10, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 14, e.mobility[White])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[White][Queen].StringBoard())
// 	assert.EqualValues(t, uint64(578721933553179648), e.attacks[White][Queen])
//
// 	score = e.evalPiece(Black, Queen)
// 	out.Printf("Black Queen: %s\n", score)
// 	assert.EqualValues(t, Score{10,10}, *score)
// 	assert.EqualValues(t, 10, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 14, e.mobility[Black])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[Black][Queen].StringBoard())
// 	assert.EqualValues(t, uint64(7927794651105288), e.attacks[Black][Queen])
// }
//
// func TestEvaluator_evalPieceKing(t *testing.T) {
// 	e := NewEvaluator()
// 	Settings.Eval.Tempo = 0
// 	p := position.NewPosition("6k1/p1qb1p1p/1p3np1/2b2p2/2B5/2P3N1/PP2QPPP/4N1K1 b - -")
// 	var score *Score
//
//
// 	e.InitEval(p)
// 	// prepare attacks
// 	// TODO pawns
// 	e.evalPiece(White, Knight)
// 	e.evalPiece(Black, Knight)
// 	e.evalPiece(White, Bishop)
// 	e.evalPiece(Black, Bishop)
// 	e.evalPiece(White, Rook)
// 	e.evalPiece(Black, Rook)
// 	e.evalPiece(White, Queen)
// 	e.evalPiece(Black, Queen)
//
// 	// King
// 	score = e.evalKing(White)
// 	out.Printf("White King: %s\n", score)
// 	assert.EqualValues(t, Score{10,10}, *score)
// 	assert.EqualValues(t, 10, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[White])
// 	assert.EqualValues(t, 14, e.mobility[White])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[White][King].StringBoard())
// 	assert.EqualValues(t, uint64(578721933553179648), e.attacks[White][King])
//
// 	score = e.evalKing(Black)
// 	out.Printf("Black King: %s\n", score)
// 	assert.EqualValues(t, Score{10,10}, *score)
// 	assert.EqualValues(t, 10, score.ValueFromScore(p.GamePhaseFactor()))
// 	out.Printf("Mobility    : %d\n", e.mobility[Black])
// 	assert.EqualValues(t, 14, e.mobility[Black])
// 	fmt.Printf("Attacks     : \n%s\n", e.attacks[Black][King].StringBoard())
// 	assert.EqualValues(t, uint64(7927794651105288), e.attacks[Black][King])
// }

func TestStartPosZeroEval(t *testing.T) {
	Settings.Eval.Tempo = 0
	p := position.NewPosition()
	e := NewEvaluator()
	v := e.Evaluate(p)
	out.Println("Value =", v)
	assert.EqualValues(t, 0, e.Evaluate(p))
}

func TestMirroredZeroEval(t *testing.T) {
	Settings.Eval.Tempo = 0
	p := position.NewPosition("r1bq1rk1/pppp1pp1/2n2n1p/1B2p3/1b2P3/2N2N1P/PPPP1PP1/R1BQ1RK1 w - -")
	e := NewEvaluator()
	v := e.Evaluate(p)
	out.Println("Value =", v)
	assert.EqualValues(t, 0, e.Evaluate(p))
}

// func TestTestFensCheck(t *testing.T) {
// 	e := NewEvaluator()
// 	var p *position.Position
// 	for i, fen := range testutil.Fens {
// 		p, _ = position.NewPositionFen(fen)
// 		e.InitEval(p)
// 		out.Printf("%d: %s\n", i+1, e.Report())
// 	}
// }
//
// func TestManualFenCheck(t *testing.T) {
// 	var fen string
// 	e := NewEvaluator()
// 	var p *position.Position
//
// 	fen = "r2qk2r/pppn1ppp/3bpn2/1P1p3b/8/4PN1P/PBPPBPP1/RN1QK2R w KQkq - 1 8 "
// 	p, _ = position.NewPositionFen(fen)
// 	e.InitEval(p)
// 	out.Printf("%s\n", e.Report())
//
// 	fen = "r2q1rk1/pppn1ppp/3bpn2/1P1p3b/8/2N1PN1P/PBPPBPP1/R2Q1RK1 b - - 4 9 "
// 	p, _ = position.NewPositionFen(fen)
// 	e.InitEval(p)
// 	out.Printf("%s\n", e.Report())
// }

func Test_TimingEval(t *testing.T) {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath("../bin")).Stop()
	// go tool pprof -http :8080 ./main ./prof.null/cpu.pprof

	out := message.NewPrinter(language.German)
	p := position.NewPosition("r1b1k2r/pppp1ppp/2n2n2/1Bb1p2q/4P3/2NP1N2/1PP2PPP/R1BQK2R w KQkq -")
	e := NewEvaluator()
	result := Value(0)

	const rounds = 5
	const iterations uint64 = 10_000_000

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			result = e.Evaluate(p)
		}
		elapsed := time.Since(start)
		out.Printf("Test took %s for %d iterations\n", elapsed, iterations)
		out.Printf("Test took %d ns per iteration\n", elapsed.Nanoseconds()/int64(iterations))
		out.Printf("Iterations per sec %d\n", int64(iterations*1e9)/elapsed.Nanoseconds())
	}
	_ = result
}


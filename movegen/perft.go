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
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var out = message.NewPrinter(language.German)

const (
	ns uint64 = 1_000_000_000
)

type Perft struct {
	Nodes            uint64
	CheckCounter     uint64
	CheckMateCounter uint64
	CaptureCounter   uint64
	EnpassantCounter uint64
	CastleCounter    uint64
	PromotionCounter uint64
}

// StartPerft is using the "normal" move generation and doesn't divide the
// the perft depths.
//noinspection GoUnhandledErrorResult
func (p *Perft) StartPerft(fen string, depth int) {

	p.resetCounter()
	pos := position.NewFen(fen)
	mgList := make([]movegen, depth+1)
	for i := 0; i <= depth; i++ {
		mgList[i] = New()
	}

	out.Printf("Performing PERFT Test for Depth %d\n", depth)
	out.Printf("-----------------------------------------\n")

	start := time.Now()

	result := p.miniMax(depth, &pos, &mgList)
	p.Nodes = result

	elapsed := time.Since(start)

	out.Printf("Time         : %d ms\n", elapsed.Milliseconds())
	out.Printf("NPS          : %d nps\n", (p.Nodes*ns)/uint64(elapsed.Nanoseconds()+1))
	out.Printf("Results:\n")
	out.Printf("   Nodes     : %d\n", p.Nodes)
	out.Printf("   Captures  : %d\n", p.CaptureCounter)
	out.Printf("   EnPassant : %d\n", p.EnpassantCounter)
	out.Printf("   Checks    : %d\n", p.CheckCounter)
	out.Printf("   CheckMates: %d\n", p.CheckMateCounter)
	out.Printf("   Castles   : %d\n", p.CastleCounter)
	out.Printf("   Promotions: %d\n", p.PromotionCounter)
	out.Printf("-----------------------------------------\n")

}

func (p *Perft) StartPerftOD(fen string, maxdepth int) {

}

func (p *Perft) StartPerftDevide(fen string, maxdepth int) {

}

func (p *Perft) resetCounter() {
	p.Nodes = 0
	p.CheckCounter = 0
	p.CheckMateCounter = 0
	p.CaptureCounter = 0
	p.EnpassantCounter = 0
	p.CastleCounter = 0
	p.PromotionCounter = 0
}

func (p *Perft) miniMax(depth int, positionPtr *position.Position, mgListPtr *[]movegen) uint64 {
	totalNodes := uint64(0)
	movegens := *mgListPtr
	// moves to search recursively
	movesPtr := movegens[depth].GeneratePseudoLegalMoves(positionPtr, GenAll)
	for _, move := range *movesPtr {
		if depth > 1 {
			positionPtr.DoMove(move)
			if positionPtr.WasLegalMove() {
				totalNodes += p.miniMax(depth-1, positionPtr, mgListPtr)
			}
			positionPtr.UndoMove()
		} else {
			capture := positionPtr.GetPiece(move.To()) != PieceNone
			enpassant := move.MoveType() == EnPassant
			castling := move.MoveType() == Castling
			promotion := move.MoveType() == Promotion
			positionPtr.DoMove(move)
			if positionPtr.WasLegalMove() {
				totalNodes++
				if enpassant {
					p.EnpassantCounter++
					p.CaptureCounter++
				}
				if capture {
					p.CaptureCounter++
				}
				if castling {
					p.CastleCounter++
				}
				if promotion {
					p.PromotionCounter++
				}
				if positionPtr.HasCheck() {
					p.CheckCounter++
				}
				if !movegens[0].HasLegalMove(positionPtr) {
					p.CheckMateCounter++
				}
			}
			positionPtr.UndoMove()
		}
	}
	return totalNodes
}

func (p *Perft) miniMaxOD(d int, position *position.Position, moveGenList *[]movegen) uint64 {

	return 0
}

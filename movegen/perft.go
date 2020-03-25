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
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var out = message.NewPrinter(language.German)

// Perft is class to test move generation of the chess engine.
type Perft struct {
	Nodes            uint64
	CheckCounter     uint64
	CheckMateCounter uint64
	CaptureCounter   uint64
	EnpassantCounter uint64
	CastleCounter    uint64
	PromotionCounter uint64
	stopFlag         bool
}

// NewPerft creates a new empty Perft instance
func NewPerft() *Perft {
	return &Perft{}
}

// Stop can be used when perft has been started
// in a goroutine to stop the currently running
// perft test
func (p *Perft) Stop() {
	p.stopFlag = true
}

// StartPerftMulti is using normal or on demand move generation and doesn't
// divide the the perft depths. It iterates through the given start to end depths.
// If this has been started in a go routine it can be stopped via Stop()
func (p *Perft) StartPerftMulti(fen string, startDepth int, endDepth int, onDemandFlag bool) {
	p.stopFlag = false
	for i := startDepth; i <= endDepth; i++ {
		if p.stopFlag {
			out.Print("Perft multi depth stopped\n")
			return
		}
		p.StartPerft(fen, i, onDemandFlag)
	}
}

// StartPerft is using normal or on demand move generation and doesn't
// divide the the perft depths.
// If this has been started in a go routine it can be stopped via Stop()
func (p *Perft) StartPerft(fen string, depth int, onDemandFlag bool) {
	p.stopFlag = false

	// set 1 as minimum
	if depth <= 0 {
		depth = 1
	}

	// prepare
	p.resetCounter()
	posPtr := position.NewPositionFen(fen)
	mgList := make([]*Movegen, depth+1)
	for i := 0; i <= depth; i++ {
		mgList[i] = NewMoveGen()
	}

	out.Printf("Performing PERFT Test for Depth %d\n", depth)
	out.Printf("-----------------------------------------\n")

	result := uint64(0)

	// the actual perft call
	start := time.Now()
	if onDemandFlag {
		result = p.miniMaxOD(depth, posPtr, &mgList)
	} else {
		result = p.miniMax(depth, posPtr, &mgList)
	}
	elapsed := time.Since(start)

	if result == 0 {
		out.Print("Perft stopped\n")
		return
	}

	p.Nodes = result

	out.Printf("Time         : %d ms\n", elapsed.Milliseconds())
	out.Printf("NPS          : %d nps\n", (p.Nodes*uint64(time.Second.Nanoseconds()))/uint64(elapsed.Nanoseconds()+1))
	out.Printf("Results:\n")
	out.Printf("   Nodes     : %d\n", p.Nodes)
	out.Printf("   Captures  : %d\n", p.CaptureCounter)
	out.Printf("   EnPassant : %d\n", p.EnpassantCounter)
	out.Printf("   Checks    : %d\n", p.CheckCounter)
	out.Printf("   CheckMates: %d\n", p.CheckMateCounter)
	out.Printf("   Castles   : %d\n", p.CastleCounter)
	out.Printf("   Promotions: %d\n", p.PromotionCounter)
	out.Printf("-----------------------------------------\n")
	out.Printf("Finished PERFT Test for Depth %d\n\n", depth)
}

func (p *Perft) miniMax(depth int, positionPtr *position.Position, mgListPtr *[]*Movegen) uint64 {
	totalNodes := uint64(0)
	movegens := *mgListPtr
	// moves to search recursively
	movesPtr := movegens[depth].GeneratePseudoLegalMoves(positionPtr, GenAll)
	for _, move := range *movesPtr {
		if p.stopFlag {
			return 0
		}
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

func (p *Perft) miniMaxOD(depth int, positionPtr *position.Position,  mgListPtr *[]*Movegen) uint64 {
	totalNodes := uint64(0)
	movegens := *mgListPtr
	// moves to search recursively
	mg := movegens[depth]
	for move := mg.GetNextMove(positionPtr, GenAll); move != MoveNone; move = mg.GetNextMove(positionPtr, GenAll) {
		if p.stopFlag {
			return 0
		}
		if depth > 1 {
			positionPtr.DoMove(move)
			if positionPtr.WasLegalMove() {
				totalNodes += p.miniMaxOD(depth-1, positionPtr, mgListPtr)
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

func (p *Perft) resetCounter() {
	p.Nodes = 0
	p.CheckCounter = 0
	p.CheckMateCounter = 0
	p.CaptureCounter = 0
	p.EnpassantCounter = 0
	p.CastleCounter = 0
	p.PromotionCounter = 0
}

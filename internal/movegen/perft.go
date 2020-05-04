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
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
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
func (perft *Perft) Stop() {
	perft.stopFlag = true
}

// StartPerftMulti is using normal or on demand move generation and doesn't
// divide the the perft depths. It iterates through the given start to end depths.
// If this has been started in a go routine it can be stopped via Stop()
func (perft *Perft) StartPerftMulti(fen string, startDepth int, endDepth int, onDemandFlag bool) {
	perft.stopFlag = false
	for i := startDepth; i <= endDepth; i++ {
		if perft.stopFlag {
			out.Print("Perft multi depth stopped\n")
			return
		}
		perft.StartPerft(fen, i, onDemandFlag)
	}
}

// StartPerft is using normal or on demand move generation and doesn't
// divide the the perft depths.
// If this has been started in a go routine it can be stopped via Stop()
func (perft *Perft) StartPerft(fen string, depth int, onDemandFlag bool) {
	perft.stopFlag = false

	// set 1 as minimum
	if depth <= 0 {
		depth = 1
	}

	// prepare
	perft.resetCounter()
	posPtr, _ := position.NewPositionFen(fen)
	mgList := make([]*Movegen, depth+1)
	for i := 0; i <= depth; i++ {
		mgList[i] = NewMoveGen()
	}

	out.Printf("Performing PERFT Test for Depth %d\n", depth)
	out.Printf("FEN: %s\n", fen)
	out.Printf("-----------------------------------------\n")

	result := uint64(0)

	// the actual perft call
	start := time.Now()
	if onDemandFlag {
		result = perft.miniMaxOD(depth, posPtr, &mgList)
	} else {
		result = perft.miniMax(depth, posPtr, &mgList)
	}
	elapsed := time.Since(start)

	if result == 0 {
		out.Print("Perft stopped\n")
		return
	}

	perft.Nodes = result

	out.Printf("Time         : %s\n", elapsed)
	out.Printf("NPS          : %d nps\n", (perft.Nodes*uint64(time.Second.Nanoseconds()))/uint64(elapsed.Nanoseconds()+1))
	out.Printf("Results:\n")
	out.Printf("   Nodes     : %d\n", perft.Nodes)
	out.Printf("   Captures  : %d\n", perft.CaptureCounter)
	out.Printf("   EnPassant : %d\n", perft.EnpassantCounter)
	out.Printf("   Checks    : %d\n", perft.CheckCounter)
	out.Printf("   CheckMates: %d\n", perft.CheckMateCounter)
	out.Printf("   Castles   : %d\n", perft.CastleCounter)
	out.Printf("   Promotions: %d\n", perft.PromotionCounter)
	out.Printf("-----------------------------------------\n")
	out.Printf("Finished PERFT Test for Depth %d\n\n", depth)
}

func (perft *Perft) miniMax(depth int, p *position.Position, mgListPtr *[]*Movegen) uint64 {
	totalNodes := uint64(0)
	movegens := *mgListPtr
	// moves to search recursively
	movesPtr := movegens[depth].GeneratePseudoLegalMoves(p, GenAll, p.HasCheck())
	for _, move := range *movesPtr {
		if perft.stopFlag {
			return 0
		}
		if depth > 1 {
			p.DoMove(move)
			if p.WasLegalMove() {
				totalNodes += perft.miniMax(depth-1, p, mgListPtr)
			}
			p.UndoMove()
		} else {
			capture := p.GetPiece(move.To()) != PieceNone
			enpassant := move.MoveType() == EnPassant
			castling := move.MoveType() == Castling
			promotion := move.MoveType() == Promotion
			p.DoMove(move)
			if p.WasLegalMove() {
				totalNodes++
				if enpassant {
					perft.EnpassantCounter++
					perft.CaptureCounter++
				}
				if capture {
					perft.CaptureCounter++
				}
				if castling {
					perft.CastleCounter++
				}
				if promotion {
					perft.PromotionCounter++
				}
				if p.HasCheck() {
					perft.CheckCounter++
				}
				if !movegens[0].HasLegalMove(p) {
					perft.CheckMateCounter++
				}
			}
			p.UndoMove()
		}
	}
	return totalNodes
}

func (perft *Perft) miniMaxOD(depth int, p *position.Position, mgListPtr *[]*Movegen) uint64 {
	totalNodes := uint64(0)
	movegens := *mgListPtr
	// moves to search recursively
	mg := movegens[depth]
	hasCheck := p.HasCheck()
	for move := mg.GetNextMove(p, GenAll, hasCheck); move != MoveNone; move = mg.GetNextMove(p, GenAll, hasCheck) {
		if perft.stopFlag {
			return 0
		}
		if depth > 1 {
			p.DoMove(move)
			if p.WasLegalMove() {
				totalNodes += perft.miniMaxOD(depth-1, p, mgListPtr)
			}
			p.UndoMove()
		} else {
			capture := p.GetPiece(move.To()) != PieceNone
			enpassant := move.MoveType() == EnPassant
			castling := move.MoveType() == Castling
			promotion := move.MoveType() == Promotion
			p.DoMove(move)
			if p.WasLegalMove() {
				totalNodes++
				if enpassant {
					perft.EnpassantCounter++
					perft.CaptureCounter++
				}
				if capture {
					perft.CaptureCounter++
				}
				if castling {
					perft.CastleCounter++
				}
				if promotion {
					perft.PromotionCounter++
				}
				if p.HasCheck() {
					perft.CheckCounter++
				}
				if !movegens[0].HasLegalMove(p) {
					perft.CheckMateCounter++
				}
			}
			p.UndoMove()
		}
	}
	return totalNodes
}

func (perft *Perft) resetCounter() {
	perft.Nodes = 0
	perft.CheckCounter = 0
	perft.CheckMateCounter = 0
	perft.CaptureCounter = 0
	perft.EnpassantCounter = 0
	perft.CastleCounter = 0
	perft.PromotionCounter = 0
}

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

package search

import (
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

// //////////////////////////////////////////////////////
// Statistics
// //////////////////////////////////////////////////////

// Statistics are extra data and stats not essential for a functioning search
type Statistics struct {
	TTHit  uint64
	TTMiss uint64

	Evaluations uint64
	EvalFromTT  uint64

	BestMoveChange       uint64
	AspirationResearches uint64

	BetaCuts     uint64
	BetaCuts1st  uint64

	Razorings   uint64
	RfpPrunings uint64
	QFpPrunings uint64
	FpPrunings  uint64

	NMPMateAlpha uint64
	NMPMateBeta  uint64

	ThreatExtension uint64
	CheckExtension  uint64

	LmpCuts       uint64
	LmrResearches uint64
	LmrReductions uint64

	TTMoveUsed uint64
	NoTTMove   uint64
	TTCuts     uint64
	TTNoCuts   uint64

	IIDmoves    uint64
	IIDsearches uint64

	LeafPositionsEvaluated uint64
	Checkmates             uint64
	Stalemates             uint64
	RootPvsResearches      uint64
	PvsResearches          uint64
	NullMoveCuts           uint64
	StandpatCuts           uint64
	Mdp                    uint64
	CheckInQS              uint64

	CurrentIterationDepth    int
	CurrentSearchDepth       int
	CurrentExtraSearchDepth  int
	CurrentVariation         moveslice.MoveSlice
	CurrentRootMoveIndex     int
	CurrentRootMove          Move
	CurrentBestRootMove      Move
	CurrentBestRootMoveValue Value
}

func (s *Statistics) String() string {
	return out.Sprintf("%+v", *s)
}

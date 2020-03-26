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

package search

import (
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var trace = false
var slog = getSearchLog()

func (s *Search) rootSearch(position *position.Position, depth int, alpha Value, beta Value) {
	if trace {
		slog.Debugf("Ply %2.d Depth %2.d start: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
		defer slog.Debugf("Ply %2.d Depth %2.d end: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
	}

	// in root search search all moves and store value in root
	// moves for sorting for next iteration
	// best move is stored in pv[0][0]
	// best value is stored in pv[0][0].value
	// the next iteration begins with the best move of the last
	// iteration so we can be sure pv[0][0] sill be set with the
	// last best move from the previous iteration independent of
	// the value. Any better move found is really better and will
	// replace pv[0][0] and also will be sorted first in the
	// next iteration

	// prepare root node search
	bestNodeValue := ValueNA

	// prepare move loop
	var value Value

	// ///////////////////////////////////////////////////////
	// MOVE LOOP
	for i, m := range *s.rootMoves {

		position.DoMove(m)
		s.nodesVisited++
		s.statistics.CurrentVariation.PushBack(m)
		s.statistics.CurrentRootMoveIndex = i
		s.statistics.CurrentRootMove = m

		// check repetition and 50 moves
		if s.checkDrawRepAnd50(position, 2) {
			value = ValueDraw
		} else {
			value = -s.search(position, depth-1, 1, -beta, -alpha)
			// iterationDepth, PLY_ROOT, alpha, beta, Do_Null_Move, nodeType
		}

		s.statistics.CurrentVariation.PopBack()
		position.UndoMove()

		// we want to do at least have one complete search with depth 1
		// later we can stop any time - any new best moves will have been
		// stored in pv[0]
		if s.stopConditions() && depth > 1 {
			return
		}

		// set the value into he root move to later be able to sort
		// root moves according to value
		s.rootMoves.Set(i, m.SetValue(value))

		// Did we find a better move for this node (not ply)?
		// For the first move this is always the case.
		if value > bestNodeValue {
			// new best value
			bestNodeValue = value
			// we have a new pv[0][0] - store pv+1 tp pv
			savePV(m, s.pv[1], s.pv[0])
		}
	}
	// MOVE LOOP
	// ///////////////////////////////////////////////////////

}

func (s *Search) search(position *position.Position, depth int, ply int, alpha Value, beta Value) Value {
	if trace {
		slog.Debugf("%0*s Ply %2.d Depth %2.d start:  %s", ply, "", ply, depth, s.statistics.CurrentVariation.StringUci())
		defer slog.Debugf("%0*s Ply %2.d Depth %2.d end:  %s",ply, "", ply, depth, s.statistics.CurrentVariation.StringUci())
	}

	// Check if search should be stopped
	if s.stopConditions() {
		return ValueNA
	}

	// Leaf node when depth == 0 or max ply has been reached
	if depth == 0 || ply >= MaxDepth {
		return s.qsearch(position, depth, ply, alpha, beta)
	}

	// prepare node search
	bestNodeValue := ValueNA
	bestNodeMove := MoveNone // used to store in the TT
	myMg := s.mg[ply]
	myMg.ResetOnDemand()
	s.pv[ply].Clear()

	// prepare move loop
	var value Value
	movesSearched := 0

	// ///////////////////////////////////////////////////////
	// MOVE LOOP
	for move := myMg.GetNextMove(position, movegen.GenAll); move != MoveNone; move = myMg.GetNextMove(position, movegen.GenAll) {

		position.DoMove(move)

		// check if legal move or skip (root moves are always legal)
		if !position.WasLegalMove() {
			position.UndoMove()
			continue
		}

		// we only count legal moves
		s.nodesVisited++
		s.statistics.CurrentVariation.PushBack(move)
		s.sendSearchUpdateToUci()

		// check repetition and 50 moves
		if s.checkDrawRepAnd50(position, 2) {
			value = ValueDraw
		} else {
			value = -s.search(position, depth-1, ply+1, -beta, -alpha)
			// iterationDepth, PLY_ROOT, alpha, beta, Do_Null_Move, nodeType
		}

		movesSearched++
		s.statistics.CurrentVariation.PopBack()
		position.UndoMove()

		// check if we should stop the search
		if s.stopConditions() {
			return ValueNA
		}

		// Did we find a better move for this node (not ply)?
		// For the first move this is always the case.
		if value > bestNodeValue {
			// these are only valid for this node
			// not for all of the ply (not yet clear if >alpha)
			bestNodeValue = value
			bestNodeMove = move
			// Did we find a better move than in previous nodes in ply
			// then this is our new PV and best move for this ply.
			// If we never find a better alpha this means all moves in
			// this node are worse then other moves in other nodes which
			// raised alpha - meaning we have a better move from another
			// node we would play. We will return alpha and store a alpha
			// node in TT with no best move for TT.
			if value > alpha {
				// we have a new best move for the ply
				savePV(move, s.pv[ply+1], s.pv[ply])
				// If we found a move that is better or equal than beta
				// this means that the opponent can/will avoid this
				// position altogether so we can stop search this node.
				// We will not know if our best move is really the
				// best move or how good it really is (value is an lower bound)
				// as we cut off the rest of the search of the node here.
				// We will safe the move as a killer to be able to search it
				// earlier in another node of the ply.
				if value >= beta {
					break
				}
				// We found a move between alpha and beta which means we
				// really have found the best move so far in the ply which
				// can be forced (opponent can't avoid it).
				// We raise alpha so the successive searches in this ply
				// need to find even better moves or dismiss the moves.
				alpha = value
			}
		}
	}
	// MOVE LOOP
	// ///////////////////////////////////////////////////////

	// if we did not have at least one legal move
	// then we might have a mate or in quiescence
	// only quite moves
	if movesSearched == 0 && !s.stopConditions() {
		if position.HasCheck() {
			bestNodeValue = -ValueCheckMate
		} else {
			bestNodeValue = ValueDraw
		}
	}

	// store TT
	// TODO
	_ = bestNodeMove

	return bestNodeValue
}

func (s *Search) qsearch(position *position.Position, depth int, ply int, alpha Value, beta Value) Value {
	// TODO
	return s.evaluate(position)
}

func (s *Search) evaluate(position *position.Position) Value {
	// TODO
	return position.Material(position.NextPlayer()) - position.Material(position.NextPlayer().Flip())
}

// savePV adds the given move as first move to a cleared dest and the appends
// all src moves to dest
func savePV(move Move, src *moveslice.MoveSlice, dest *moveslice.MoveSlice) {
	dest.Clear()
	dest.PushBack(move)
	*dest = append(*dest, *src...)
}

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
	golog "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var trace = false

var slog = getSearchTraceLog()

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
		defer slog.Debugf("%0*s Ply %2.d Depth %2.d end  :  %s", ply, "", ply, depth, s.statistics.CurrentVariation.StringUci())
	}

	// Check if search should be stopped
	if s.stopConditions() {
		return ValueNA
	}

	// Leaf node when depth == 0 or max ply has been reached
	if depth == 0 || ply >= MaxDepth {
		return s.qsearch(position, ply, alpha, beta)
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

		// ///////////////////////////////////////////////////////
		// DO MOVE
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
		}

		movesSearched++
		s.statistics.CurrentVariation.PopBack()
		position.UndoMove()
		// UNDO MOVE
		// ///////////////////////////////////////////////////////

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
	// then we might have a mate or stalemate
	if movesSearched == 0 && !s.stopConditions() {
		if position.HasCheck() {
			bestNodeValue = -ValueCheckMate + Value(ply)
		} else {
			bestNodeValue = ValueDraw
		}
	}

	// store TT
	// TODO
	_ = bestNodeMove

	return bestNodeValue
}

func (s *Search) qsearch(position *position.Position, ply int, alpha Value, beta Value) Value {

	if !config.Settings.Search.UseQuiescence {
		return s.evaluate(position)
	}

	if s.statistics.CurrentExtraSearchDepth < ply {
		s.statistics.CurrentExtraSearchDepth = ply
	}

	// prepare node search
	bestNodeValue := ValueNA
	bestNodeMove := MoveNone // used to store in the TT
	myMg := s.mg[ply]
	myMg.ResetOnDemand()
	s.pv[ply].Clear()
	hasCheck := position.HasCheck()

	if !hasCheck {
		// get an evaluation for the position
		staticEval := s.evaluate(position)
		// Standpat Cut
		if staticEval >= beta {
			return staticEval
		}
		if staticEval > alpha {
			alpha = staticEval
		}
		bestNodeValue = staticEval
	}

	// prepare move loop
	var value Value
	movesSearched := 0

	// if in check we search all moves
	// this is in fact a search extension for checks
	var mode movegen.GenMode
	if hasCheck {
		mode = movegen.GenAll
	} else {
		mode = movegen.GenCap
	}

	// ///////////////////////////////////////////////////////
	// MOVE LOOP
	for move := myMg.GetNextMove(position, mode); move != MoveNone; move = myMg.GetNextMove(position, mode) {

		// reduce number of moves searched in quiescence
		// by looking at good captures only
		if !hasCheck && !s.goodCapture(position, move) {
			continue
		}

		// ///////////////////////////////////////////////////////
		// DO MOVE
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

		// check repetition and 50 moves when in check
		// otherwise only capturing moves are generated
		// which break repetition and 50-moves rule anyway
		if hasCheck && s.checkDrawRepAnd50(position, 2) {
			value = ValueDraw
		} else {
			value = -s.qsearch(position, ply+1, -beta, -alpha)
		}

		movesSearched++
		s.statistics.CurrentVariation.PopBack()
		position.UndoMove()
		// UNDO MOVE
		// ///////////////////////////////////////////////////////

		// check if we should stop the search
		if s.stopConditions() {
			return ValueNA
		}

		// see search for documentation
		if value > bestNodeValue {
			bestNodeValue = value
			bestNodeMove = move
			if value > alpha {
				savePV(move, s.pv[ply+1], s.pv[ply])
				if value >= beta {
					break
				}
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
		// if we have a mate we had a check before and therefore
		// generated all move. We can be sure this is a mate.
		if position.HasCheck() {
			bestNodeValue = -ValueCheckMate + Value(ply)
		}
		// if we do not have mate we had no check and
		// therefore might have quiet moves which we
		// did not generate.
		// We return the standpat value in this case
		// with bestNodeValue which we have set to the
		// static eval earlier
	}

	// TODO TT store
	_ = bestNodeMove

	return bestNodeValue
}

func (s *Search) evaluate(position *position.Position) Value {
	return s.eval.Evaluate(position)
}

func (s *Search) goodCapture(p *position.Position, move Move) bool {
	// Lower value piece captures higher value piece
	// With a margin to also look at Bishop x Knight
	return p.GetPiece(move.From()).ValueOf() + 50 < p.GetPiece(move.To()).ValueOf() ||
	// all recaptures should be looked at
	(p.LastMove() != MoveNone && p.LastMove().To() == move.To() && p.LastCapturedPiece() != PieceNone) ||
	// undefended pieces captures are good
	// If the defender is "behind" the attacker this will not be recognized
	// here This is not too bad as it only adds a move to qsearch which we
	// could otherwise ignore
	!p.IsAttacked(move.To(), p.NextPlayer().Flip())
	// Check SEE score of higher value pieces to low value pieces
	// || (SearchConfig::USE_QS_SEE && (Attacks::see(position, move) > 0));
}

// savePV adds the given move as first move to a cleared dest and the appends
// all src moves to dest
func savePV(move Move, src *moveslice.MoveSlice, dest *moveslice.MoveSlice) {
	dest.Clear()
	dest.PushBack(move)
	*dest = append(*dest, *src...)
}

// getSearchTraceLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
// for usage in the search itself
func getSearchTraceLog() *logging.Logger {
	searchLog := logging.MustGetLogger("search")

	searchLogFormat := logging.MustStringFormatter(`%{time:15:04:05.000} %{level:-7.7s}:  %{message}`)

	backend1 := logging.NewLogBackend(os.Stdout, "", golog.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, searchLogFormat)
	searchBackEnd := logging.AddModuleLevel(backend1Formatter)
	searchBackEnd.SetLevel(logging.Level(config.SearchLogLevel), "")

	// File backend
	programName, _ := os.Executable()
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")
	var logPath string
	if filepath.IsAbs(config.Settings.Log.LogPath) {
		logPath = config.Settings.Log.LogPath
	} else {
		dir, _ := os.Getwd()
		logPath = dir + "/" + config.Settings.Log.LogPath
	}
	searchLogFilePath := logPath + "/" + exeName + "_searchlog.log"
	searchLogFilePath = filepath.Clean(searchLogFilePath)

	// create file backend
	var err error
	searchLogFile, err := os.OpenFile(searchLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// buf := bufio.NewWriter(searchLogFile)

	// we use either Stdout or file - if file is valid we use only file
	if err != nil {
		golog.Println("Logfile could not be created:", err)
		searchLog.SetBackend(searchBackEnd)
	} else {
		backend2 := logging.NewLogBackend(searchLogFile, "", golog.Lmsgprefix)
		backend2Formatter := logging.NewBackendFormatter(backend2, searchLogFormat)
		searchBackEnd2 := logging.AddModuleLevel(backend2Formatter)
		searchBackEnd2.SetLevel(logging.DEBUG, "")
		// multi := logging2.SetBackend(uciBackEnd1, searchBackEnd2)
		searchLog.SetBackend(searchBackEnd2)
		searchLog.Infof("Log %s started at %s:", searchLogFile.Name(), time.Now().String())
	}
	return searchLog
}

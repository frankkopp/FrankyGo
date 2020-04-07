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

	. "github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/transpositiontable"
	. "github.com/frankkopp/FrankyGo/types"
)

var trace = false

var slog = getSearchTraceLog()

// rootSearch starts the actual recursive alpha beta search with the root moves for the first ply.
// As root moves are treated a little different this separate function supports readability
// as mixing it with the normal search would require quite some "if ply==0" statements.
func (s *Search) rootSearch(position *position.Position, depth int, alpha Value, beta Value) {
	if trace {
		slog.Debugf("Ply %-2.d Depth %-2.d start: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
		defer slog.Debugf("Ply %-2.d Depth %-2.d end: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
	}

	// In root search we search all moves and store the value
	// into the  root moves themselves for sorting in the
	// next iteration
	// best move is stored in pv[0][0]
	// best value is stored in pv[0][0].value
	// The next iteration begins with the best move of the last
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
			// ///////////////////////////////////////////////////////////////////
			// PVS
			// Initial PVS move are search without PVS uses full search window.
			if !Settings.Search.UsePVS || i == 0 {
				value = -s.search(position, depth-1, 1, -beta, -alpha, true, true)
			} else {
				// Null window search after the initial PV search.
				value = -s.search(position, depth-1, 1, -alpha-1, -alpha, false, true)
				// If this move improved alpha without exceeding beta we do a proper full window
				// search to get an accurate score.
				if value > alpha && value < beta && !s.stopConditions() {
					s.statistics.RootPvsResearches++
					value = -s.search(position, depth-1, 1, -beta, -alpha, true, true)
				}
			}
			// ///////////////////////////////////////////////////////////////////
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

func (s *Search) search(position *position.Position, depth int, ply int, alpha Value, beta Value, isPV bool, doNull bool) Value {
	if trace {
		slog.Debugf("%0*s Ply %-2.d Depth %-2.d a:%-6.d b:%-6.d pv:%-6.v start:  %s", ply, "", ply, depth, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
		defer slog.Debugf("%0*s Ply %-2.d Depth %-2.d a:%-6.d b:%-6.d pv:%-6.v end  :  %s", ply, "", ply, depth, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
	}

	// Check if search should be stopped
	if s.stopConditions() {
		return ValueNA
	}

	// Leaf node when depth == 0 or max ply has been reached
	if depth == 0 || ply >= MaxDepth {
		return s.qsearch(position, ply, alpha, beta, isPV)
	}

	// Mate Distance Pruning
	// Did we already find a shorter mate then ignore
	// this one.
	if Settings.Search.UseMDP {
		if alpha < -ValueCheckMate+Value(ply) {
			alpha = -ValueCheckMate + Value(ply)
		}
		if beta > ValueCheckMate-Value(ply) {
			beta = ValueCheckMate - Value(ply)
		}
		if alpha >= beta {
			s.statistics.Mdp++
			return alpha
		}
	}

	// prepare node search
	bestNodeValue := ValueNA
	bestNodeMove := MoveNone // used to store in the TT
	ttMove := MoveNone
	ttType := ALPHA

	// TT Lookup
	// Results of searches are stored in the TT to be used to
	// avoid searching positions several times. If a position
	// is stored in the TT we retrieve a pointer to the entry.
	// We use the stored move as a best move from previous searches
	// and search it first (through setting PV move in move gen).
	// If we have a value from a similar or deeper search we check
	// if the value is usable. Exact values mean that the previously
	// stored result already was a precise result and we do not
	// need to search the position again. We can stop  searching
	// this branch and return the value.
	// Alpha or Beta entries will only be used if they improve
	// the current values.
	// TODO : Some engine treat the cut for alpha and beta nodes
	//  differently for PV and non PV nodes - needs more testing
	//  if this is relevant
	var ttEntry *transpositiontable.TtEntry
	if Settings.Search.UseTT {
		ttEntry = s.tt.Probe(position.ZobristKey())
		if ttEntry != nil { // tt hit
			s.statistics.TTHit++
			ttMove = ttEntry.Move.MoveOf()
			if int(ttEntry.Depth) >= depth {
				ttValue := valueFromTT(ttEntry.Move.ValueOf(), ply)
				cut := false
				switch {
				case !ttValue.IsValid():
					cut = false
				case ttEntry.Type == EXACT:
					cut = true
				case ttEntry.Type == ALPHA && ttValue <= alpha:
					cut = true
				case ttEntry.Type == BETA && ttValue >= beta:
					cut = true
				}
				if cut && Settings.Search.UseTTValue {
					s.getPVLine(position, s.pv[ply], depth)
					s.statistics.TTCuts++
					return ttValue
				} else {
					s.statistics.TTNoCuts++
				}
			}
		} else {
			s.statistics.TTMiss++
		}
	}

	// NULL MOVE PRUNING
	// https://www.chessprogramming.org/Null_Move_Pruning
	// Under the assumption the in most chess position it would be better
	// do make a move than to not make a move we can assume that if
	// our positional value after a null move is already above beta (>beta)
	// it would be above beta when doing a move in any case.
	// Certain situations need to be considered though:
	// - Zugzwang - it would be better not to move
	// - in check - this would lead to an illegal situation where the king is captured
	// - recursive null moves should be avoided
	if Settings.Search.UseNullMove {
		if isPV &&
			doNull &&
			depth >= Settings.Search.NmpDepth &&
			position.MaterialNonPawn(position.NextPlayer()) > 0 &&
			!position.HasCheck() {

			// determine depth reduction
			// ICCA Journal, Vol. 22, No. 3
			// Ernst A. Heinz, Adaptive Null-Move Pruning, postscript
			// http://people.csail.mit.edu/heinz/ps/adpt_null.ps.gz
			r := Settings.Search.NmpReduction
			if depth > 8 || (depth > 6 && position.GamePhase() >= 3) {
				r += 1
			}
			newDepth := depth - r - 1
			// double check that depth does not get negative
			if newDepth < 0 {
				newDepth = 0
			}

			// do null move search
			position.DoNullMove()
			s.nodesVisited++
			nValue := -s.search(position, newDepth, ply+1, -beta, -beta+1, isPV, false)
			position.UndoNullMove()

			// check if we should stop the search
			if s.stopConditions() {
				return ValueNA
			}

			// if the value is higher than beta even after making two
			// moves it is not worth searching and it will be cut
			if nValue >= beta {
				s.statistics.NullMoveCuts++
				// Store TT
				if Settings.Search.UseTT {
					s.storeTT(position, depth, ply, ttMove, nValue, BETA)
				}
				return nValue
			}
		}
	}

	// Internal Iterative Deepening (IID)
	// https://www.chessprogramming.org/Internal_Iterative_Deepening
	// Used when no best move from the tt is available from a previous
	// searches.  IID is used to find a good move to search first by
	// searching the current position to a reduced depth, and using
	// the best move of that search as the first move at the real depth.
	// TODO Does not make a big difference in search tree size - needs to be tested
	if Settings.Search.UseIID {
		if depth >= Settings.Search.IIDDepth &&
			ttMove != MoveNone && // no move from TT
			doNull && // avoid in null move search
			isPV { // TODO test if this is necessary

			// get the new depth and make sure it is >0
			newDepth := depth - Settings.Search.IIDReduction
			if newDepth < 0 {
				newDepth = 0
			}

			// do the actual reduced search
			s.search(position, newDepth, ply, alpha, beta, isPV, true)
			s.statistics.IIDsearches++

			// check if we should stop the search
			if s.stopConditions() {
				return ValueNA
			}

			// get the best move from the reduced search if available
			if s.pv[ply].Len() > 0 {
				s.statistics.IIDmoves++
				ttMove = (*s.pv[ply])[0]
			}
		}
	}

	// reset search
	// !important to do this after IID!
	// or IID need to do this itself
	myMg := s.mg[ply]
	myMg.ResetOnDemand()
	s.pv[ply].Clear()

	// PV Move Sort
	// When we received a best move for the position from the
	// TT we set it as PV move in the movegen so it will be
	// searched first.
	if Settings.Search.UseTTMove {
		if ttMove != MoveNone {
			s.statistics.TTMoveUsed++
			myMg.SetPvMove(ttMove)
		} else {
			s.statistics.NoTTMove++
		}
	}

	// prepare move loop
	var value Value
	movesSearched := 0

	// ///////////////////////////////////////////////////////
	// MOVE LOOP
	for move := myMg.GetNextMove(position, movegen.GenAll); move != MoveNone; move = myMg.GetNextMove(position, movegen.GenAll) {

		// ///////////////////////////////////////////////////////
		// DO MOVE
		position.DoMove(move)

		// check if legal move or skip
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
			// ///////////////////////////////////////////////////////////////////
			// PVS
			// Initial PVS move are search without PVS uses full search window.
			// https://www.chessprogramming.org/Principal_Variation_Search
			if !Settings.Search.UsePVS || movesSearched == 0 {
				value = -s.search(position, depth-1, ply+1, -beta, -alpha, true, true)
			} else {
				// Null window search after the initial PV search.
				value = -s.search(position, depth-1, ply+1, -alpha-1, -alpha, false, true)
				// If this move improved alpha without exceeding beta we do a proper full window
				// search to get an accurate score.
				if value > alpha && value < beta && !s.stopConditions() {
					s.statistics.PvsResearches++
					value = -s.search(position, depth-1, ply+1, -beta, -alpha, true, true)
				}
			}
			// ///////////////////////////////////////////////////////////////////
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
			// These "best" values are only valid for this node
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
				// best move or how good it really is (value is a lower bound)
				// as we cut off the rest of the search of the node here.
				// We will safe the move as a killer to be able to search it
				// earlier in another node of the ply.
				if value >= beta {
					s.statistics.BetaCuts++
					// s.log.Debugf("Beta Cuts on %d th move\n", movesSearched)
					if movesSearched == 1 {
						s.statistics.BetaCuts1st++
					}
					if Settings.Search.UseKiller {
						myMg.StoreKiller(move)
					}
					ttType = BETA
					break
				}
				// We found a move between alpha and beta which means we
				// really have found the best move so far in the ply which
				// can be forced (opponent can't avoid it).
				// We raise alpha so the successive searches in this ply
				// need to find even better moves or dismiss the moves.
				alpha = value
				ttType = EXACT
			}
		}
	}
	// MOVE LOOP
	// ///////////////////////////////////////////////////////

	// If we did not have at least one legal move
	// then we might have a mate or stalemate
	if movesSearched == 0 && !s.stopConditions() {
		if position.HasCheck() { // mate
			s.statistics.Checkmates++
			bestNodeValue = -ValueCheckMate + Value(ply)
		} else { // stalemate
			s.statistics.Stalemates++
			bestNodeValue = ValueDraw
		}
		// this is in any case an exact value
		ttType = EXACT
	}

	// Store TT
	// Store search result for this node into the transposition table
	if Settings.Search.UseTT {
		s.storeTT(position, depth, ply, bestNodeMove, bestNodeValue, ttType)
	}

	return bestNodeValue
}

func (s *Search) qsearch(position *position.Position, ply int, alpha Value, beta Value, isPV bool) Value {
	if trace {
		slog.Debugf("%0*s Ply %-2.d QSearch     a:%-6.d b:%-6.d pv:%-6.v start:  %s", ply, "", ply, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
		defer slog.Debugf("%0*s Ply %-2.d QSearch     a:%-6.d b:%-6.d pv:%-6.v end  :  %s", ply, "", ply, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
	}

	if !Settings.Search.UseQuiescence {
		return s.evaluate(position)
	}

	if s.statistics.CurrentExtraSearchDepth < ply {
		s.statistics.CurrentExtraSearchDepth = ply
	}

	// Mate Distance Pruning
	// Did we already find a shorter mate then ignore
	// this one.
	if Settings.Search.UseMDP {
		if alpha < -ValueCheckMate+Value(ply) {
			alpha = -ValueCheckMate + Value(ply)
		}
		if beta > ValueCheckMate-Value(ply) {
			beta = ValueCheckMate - Value(ply)
		}
		if alpha >= beta {
			s.statistics.Mdp++
			return alpha
		}
	}

	// prepare node search
	bestNodeValue := ValueNA
	ttType := ALPHA
	ttMove := MoveNone
	hasCheck := position.HasCheck()

	// if in check we simply do a normal search (all moves) in qsearch
	if !hasCheck {
		// get an evaluation for the position
		staticEval := s.evaluate(position)
		// Quiescence StandPat
		// Use evaluation as a standing pat (lower bound)
		// https://www.chessprogramming.org/Quiescence_Search#Standing_Pat
		// Assumption is that there is at least on move which would improve the
		// current position. So if we are already >beta we don't need to look at it.
		if Settings.Search.UseQSStandpat && staticEval > alpha {
			if staticEval >= beta {
				s.statistics.StandpatCuts++
				return staticEval
			}
			alpha = staticEval
		}
		bestNodeValue = staticEval
	}

	// TT Lookup
	var ttEntry *transpositiontable.TtEntry
	if Settings.Search.UseQSTT {
		ttEntry = s.tt.Probe(position.ZobristKey())
		if ttEntry != nil { // tt hit
			s.statistics.TTHit++
			ttMove = ttEntry.Move.MoveOf()
			ttValue := valueFromTT(ttEntry.Move.ValueOf(), ply)
			cut := false
			switch {
			case !ttValue.IsValid():
				cut = false
			case ttEntry.Type == EXACT:
				cut = true
			case ttEntry.Type == ALPHA && ttValue <= alpha:
				cut = true
			case ttEntry.Type == BETA && ttValue >= beta:
				cut = true
			}
			if cut && Settings.Search.UseTTValue {
				s.statistics.TTCuts++
				return ttValue
			} else {
				s.statistics.TTNoCuts++
			}
		} else {
			s.statistics.TTMiss++
		}
	}

	// prepare node search
	bestNodeMove := MoveNone // used to store in the TT
	myMg := s.mg[ply]
	myMg.ResetOnDemand()
	s.pv[ply].Clear()

	// PV Move Sort
	// When we received a best move for the position from the
	// TT we set it as PV move in the movegen so it will be
	// searched first.
	if Settings.Search.UseQSTT {
		if ttMove != MoveNone {
			s.statistics.TTMoveUsed++
			myMg.SetPvMove(ttMove)
		} else {
			s.statistics.NoTTMove++
		}
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

		// check if legal move or skip
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
			value = -s.qsearch(position, ply+1, -beta, -alpha, isPV)
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

		// see search function above for documentation
		if value > bestNodeValue {
			bestNodeValue = value
			bestNodeMove = move
			if value > alpha {
				savePV(move, s.pv[ply+1], s.pv[ply])
				if value >= beta {
					ttType = BETA
					break
				}
				alpha = value
				ttType = EXACT
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
			s.statistics.Checkmates++
			bestNodeValue = -ValueCheckMate + Value(ply)
			ttType = EXACT
		}
		// if we do not have mate we had no check and
		// therefore might have only quiet moves which
		// we did not generate.
		// We return the standpat value in this case
		// which we have set to bestNodeValue in the
		// static eval earlier
	}

	// Store TT
	if Settings.Search.UseQSTT {
		s.storeTT(position, 1, ply, bestNodeMove, bestNodeValue, ttType)
	}

	return bestNodeValue
}

// call evaluation on the position
func (s *Search) evaluate(position *position.Position) Value {
	s.statistics.LeafPositionsEvaluated++
	return s.eval.Evaluate(position)
}

// reduce the number of moves searched in quiescence search by trying
// to only look at good captures. Might be improved with SEE in the
// future
func (s *Search) goodCapture(p *position.Position, move Move) bool {
	// Lower value piece captures higher value piece
	// With a margin to also look at Bishop x Knight
	return p.GetPiece(move.From()).ValueOf()+50 < p.GetPiece(move.To()).ValueOf() ||
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

// storeTT stores a position into the TT
func (s *Search) storeTT(p *position.Position, depth int, ply int, move Move, value Value, valueType ValueType) {
	s.tt.Put(p.ZobristKey(), move, int8(depth), valueToTT(value, ply), valueType, false)
}

// getPVLine fills the given pv move list with the pv move starting from the given
// depth as long as these position are in the TT
func (s *Search) getPVLine(p *position.Position, pv *moveslice.MoveSlice, depth int) {
	// Recursion-less reading of the chain of pv moves
	pv.Clear()
	counter := 0
	ttMatch := s.tt.GetEntry(p.ZobristKey())
	for ttMatch != nil && ttMatch.Move != MoveNone && counter < depth {
		pv.PushBack(ttMatch.Move)
		p.DoMove(ttMatch.Move)
		counter++
		ttMatch = s.tt.GetEntry(p.ZobristKey())
	}
	for i := 0; i < counter; i++ {
		p.UndoMove()
	}
}

// correct the value for mate distance when storing to TT
func valueToTT(value Value, ply int) Value {
	if value.IsCheckMateValue() {
		if value > 0 {
			value = value + Value(ply)
		} else {
			value = value - Value(ply)
		}
	}
	return value
}

// correct the value for mate distance when reading from TT
func valueFromTT(value Value, ply int) Value {
	if value.IsCheckMateValue() {
		if value > 0 {
			value = value - Value(ply)
		} else {
			value = value + Value(ply)
		}
	}
	return value
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
	searchBackEnd.SetLevel(logging.Level(SearchLogLevel), "")

	// File backend
	programName, _ := os.Executable()
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")
	var logPath string
	if filepath.IsAbs(Settings.Log.LogPath) {
		logPath = Settings.Log.LogPath
	} else {
		dir, _ := os.Getwd()
		logPath = dir + "/" + Settings.Log.LogPath
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

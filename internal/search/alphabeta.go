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
	"fmt"
	golog "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"

	. "github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/movegen"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/position"
	"github.com/frankkopp/FrankyGo/internal/transpositiontable"
	. "github.com/frankkopp/FrankyGo/internal/types"
	"github.com/frankkopp/FrankyGo/internal/util"
)

var trace = false

// rootSearch starts the actual recursive alpha beta search with the root moves for the first ply.
// As root moves are treated a little different this separate function supports readability
// as mixing it with the normal search would require quite some "if ply==0" statements.
func (s *Search) rootSearch(position *position.Position, depth int, alpha Value, beta Value) {
	if trace {
		s.slog.Debugf("Ply %-2.d Depth %-2.d start: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
		defer s.slog.Debugf("Ply %-2.d Depth %-2.d end: %s", 0, depth, s.statistics.CurrentVariation.StringUci())
	}

	// In root search we search all moves and store the value
	// into the root moves themselves for sorting in the
	// next iteration
	// best move is stored in pv[0][0]
	// best value is stored in pv[0][0].value
	// The next iteration begins with the best move of the last
	// iteration so we can be sure pv[0][0] will be set with the
	// last best move from the previous iteration independent of
	// the value. Any better move found is really better and will
	// replace pv[0][0] and also will be sorted first in the
	// next iteration

	// prepare root node search
	bestNodeValue := ValueNA
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
			// First move in a node is an assumed PV and searched with full search window
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

		// we want to do at least one complete search with depth 1
		// After that we can stop any time - any new best moves will
		// have been stored in pv[0]
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

// search is the normal alpha beta search after the root move ply (ply > 0)
// it will be called recursively until the remaining depth == 0 and we would
// enter quiescence search. Search consumes about 60% of the search time and
// all major prunings are done here. Quiescence search uses about 40% of the
// search time and has less options for pruning as not all moves are searched.
func (s *Search) search(p *position.Position, depth int, ply int, alpha Value, beta Value, isPV bool, doNull bool) Value {
	if trace {
		s.slog.Debugf("%0*s Ply %-2.d Depth %-2.d a:%-6.d b:%-6.d pv:%-6.v start:  %s", ply, "", ply, depth, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
		defer s.slog.Debugf("%0*s Ply %-2.d Depth %-2.d a:%-6.d b:%-6.d pv:%-6.v end  :  %s", ply, "", ply, depth, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
	}

	// Check if search should be stopped
	if s.stopConditions() {
		return ValueNA
	}

	// Enter quiescence search when depth == 0 or max ply has been reached
	if depth == 0 || ply >= MaxDepth {
		return s.qsearch(p, ply, alpha, beta, isPV)
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
	us := p.NextPlayer()
	bestNodeValue := ValueNA
	bestNodeMove := MoveNone // used to store in the TT
	ttMove := MoveNone
	ttType := ALPHA
	hasCheck := p.HasCheck()
	matethreat := false

	// TT Lookup
	// Results of searches are stored in the TT to be used to
	// avoid searching positions several times. If a position
	// is stored in the TT we retrieve a pointer to the entry.
	// We use the stored move as a best move from previous searches
	// and search it first (through setting PV move in move gen).
	// If we have a value from a similar or deeper search we check
	// if the value is usable. Exact values mean that the previously
	// stored result already was a precise result and we do not
	// need to search the position again. We can stop searching
	// this branch and return the value.
	// Alpha or Beta entries will only be used if they improve
	// the current values.
	// TODO : Some engines treat the cut for alpha and beta nodes
	//  differently for PV and non PV nodes - needs more testing
	//  if this is relevant
	var ttEntry *transpositiontable.TtEntry
	if Settings.Search.UseTT {
		ttEntry = s.tt.Probe(p.ZobristKey())
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
					s.getPVLine(p, s.pv[ply], depth)
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

	// Reverse Futility Pruning, (RFP, Static Null Move Pruning)
	// https://www.chessprogramming.org/Reverse_Futility_Pruning
	// Anticipate likely alpha low in the next ply by a beta cut
	// off before making and evaluating the move
	if Settings.Search.UseRFP &&
		doNull &&
		depth <= 3 &&
		!isPV &&
		!hasCheck {
		// get an evaluation for the position
		staticEval := s.evaluate(p, ply)
		margin := rfp[depth]
		if staticEval-margin >= beta {
			s.statistics.RfpPrunings++
			return staticEval - margin // fail-hard: beta / fail-soft: staticEval - evalMargin;
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
		if doNull &&
			!isPV &&
			depth >= Settings.Search.NmpDepth &&
			p.MaterialNonPawn(us) > 0 &&
			!hasCheck {
			// possible other criteria: eval > beta

			// determine depth reduction
			// ICCA Journal, Vol. 22, No. 3
			// Ernst A. Heinz, Adaptive Null-Move Pruning, postscript
			// http://people.csail.mit.edu/heinz/ps/adpt_null.ps.gz
			r := Settings.Search.NmpReduction
			if depth > 8 || (depth > 6 && p.GamePhase() >= 3) {
				r += 1
			}
			newDepth := depth - r - 1
			// double check that depth does not get negative
			if newDepth < 0 {
				newDepth = 0
			}

			// do null move search
			p.DoNullMove()
			s.nodesVisited++
			nValue := -s.search(p, newDepth, ply+1, -beta, -beta+1, false, false)
			p.UndoNullMove()

			// check if we should stop the search
			if s.stopConditions() {
				return ValueNA
			}

			// flag for mate threats
			if nValue > ValueCheckMateThreshold {
				// although this player did not make a move the value still is
				// a mate - very good! Just adjust the value to not return an
				// unproven mate
				s.statistics.NMPMateBeta++
				nValue = ValueCheckMateThreshold
			} else if nValue < ValueCheckMateThreshold {
				// the player did not move a got mated ==> mate threat
				s.statistics.NMPMateAlpha++
				matethreat = true
			}

			// if the value is higher than beta even after not making
			// a move it is not worth searching as it will very likely
			// be above beta if we make a move
			if nValue >= beta {
				s.statistics.NullMoveCuts++
				// Store TT
				if Settings.Search.UseTT {
					s.storeTT(p, depth, ply, ttMove, nValue, BETA)
				}
				return nValue
			}
		}
	}

	// Internal Iterative Deepening (IID)
	// https://www.chessprogramming.org/Internal_Iterative_Deepening
	// Used when no best move from the tt is available from a previous
	// searches. IID is used to find a good move to search first by
	// searching the current position to a reduced depth, and using
	// the best move of that search as the first move at the real depth.
	// Does not make a big difference in search tree size when move
	// order already is good.
	if Settings.Search.UseIID {
		if depth >= Settings.Search.IIDDepth &&
			ttMove != MoveNone && // no move from TT
			doNull && // avoid in null move search
			isPV {

			// get the new depth and make sure it is >0
			newDepth := depth - Settings.Search.IIDReduction
			if newDepth < 0 {
				newDepth = 0
			}

			// do the actual reduced search
			s.search(p, newDepth, ply, alpha, beta, isPV, true)
			s.statistics.IIDsearches++

			// check if we should stop the search
			if s.stopConditions() {
				return ValueNA
			}

			// get the best move from the reduced search if available
			if s.pv[ply].Len() > 0 {
				s.statistics.IIDmoves++
				ttMove = (*s.pv[ply])[0].MoveOf()
			}
		}
	}

	// reset search
	// !important to do this after IID!
	myMg := s.mg[ply]
	myMg.ResetOnDemand()
	s.pv[ply].Clear()

	// PV Move Sort
	// When we received a best move for the position from the
	// TT or IID we set it as PV move in the movegen so it will
	// be searched first.
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
	for move := myMg.GetNextMove(p, movegen.GenAll, hasCheck);
		move != MoveNone; move = myMg.GetNextMove(p, movegen.GenAll, hasCheck) {

		from := move.From()
		to := move.To()

		if false { // DEBUG
			err := false
			msg := ""
			switch {
			case !move.IsValid():
				msg = fmt.Sprintf("Position DoMove: Invalid move %s", move.String())
				err = true
			case p.GetPiece(from) == PieceNone:
				msg = fmt.Sprintf("Position DoMove: No piece on %s for move %s", p.GetPiece(from).String(), move.StringUci())
				err = true
			case p.GetPiece(from).ColorOf() != us:
				msg = fmt.Sprintf("Position DoMove: Piece to move does not belong to next player %s", p.GetPiece(from).String())
				err = true
			case p.GetPiece(to).TypeOf() == King:
				msg = fmt.Sprintf("Position DoMove: King cannot be captured!")
				err = true
			}
			if err {
				s.log.Criticalf("Search              : Depth %d Ply %d alpha %d beta %d isPv %t doNull %t\n", depth, ply, alpha, beta, isPV, doNull)
				s.log.Criticalf("Position            : %s\n", p.StringFen())
				s.log.Criticalf("Move                : %s\n", move.String())
				s.log.Criticalf("Moves Searched      : %d\n", movesSearched)
				s.log.Criticalf("ttMove              : %s\n", ttMove.String())
				s.log.Criticalf("bestMove            : %s\n", bestNodeMove.String())
				s.log.Criticalf("MoveGen PV          : %s\n", myMg.PvMove())
				s.log.Criticalf("MoveGen K1          : %s\n", myMg.KillerMoves()[0])
				s.log.Criticalf("MoveGen K2          : %s\n", myMg.KillerMoves()[1])
				s.log.Criticalf("MoveGen Moves       : %s\n", myMg.GeneratePseudoLegalMoves(p, movegen.GenAll, false).StringUci())
				s.log.Criticalf(msg)
				panic(msg)
			}
		} // DEBUG

		// prepare newDepth
		newDepth := depth - 1
		lmrDepth := newDepth
		extension := 0

		givesCheck := p.GivesCheck(move)

		// Here we try some search extensions. This has to be done
		// very carefully as it usually is more effective to prune
		// than to extend.
		if Settings.Search.UseExt {
			// The check extensions is a bit redundant as our QS search
			// searches all moves anyway when in check. But with this
			// extension we hope to profit from using the prunings
			// of the normal search which are not available in
			// qsearch.
			if Settings.Search.UseCheckExt && givesCheck {
				s.statistics.CheckExtension++
				extension = 1
			}

			// If we have found a mate threat during Null Move Search
			// we extend normal search by one ply to try to find
			// a way out.
			// Deactivated in config as this grows the search tree
			// too much.
			if Settings.Search.UseThreatExt && matethreat {
				s.statistics.ThreatExtension++
				extension = 1
			}

			// With this turned off we still can use extension to
			// at least avoid reductions for these moves.
			if Settings.Search.UseExtAddDepth {
				newDepth += extension
			}
		}

		// ///////////////////////////////////////////////////////
		// Forward Pruning
		// FP will only be done when the move is not
		// interesting - no check, no capture, etc.
		if !isPV &&
			extension == 0 &&
			move != ttMove &&
			move != (*myMg.KillerMoves())[0] &&
			move != (*myMg.KillerMoves())[1] &&
			move.MoveType() != Promotion &&
			!p.IsCapturingMove(move) &&
			!hasCheck && // pre move
			!givesCheck && // post move
			!matethreat { // from pre move null move check

			// to check in futility pruning what material delta we have
			materialEval := p.Material(us) - p.Material(us.Flip())
			moveGain := p.GetPiece(to).ValueOf()

			// Futility Pruning
			// Using an array of margin values for each depth
			// we try to prune moves if they seem not worth
			// searching any further. They are so far below
			// alpha that we can assume a beta cutoff in the
			// next iteration anyway.
			// This is a typical forward pruning technique
			// which might lead to errors.
			// Limited Razoring / Extended FP are covered by this.
			// TODO: needs testing and tuning
			// TODO: Crafty excepts moves were passed pawns are far ahead.
			if Settings.Search.UseFP && depth < 7 {
				futilityMargin := fp[depth]
				if materialEval+moveGain+futilityMargin <= alpha {
					if materialEval+moveGain > bestNodeValue {
						bestNodeValue = materialEval + moveGain
					}
					s.statistics.FpPrunings++
					continue
				}
			}

			// LMP - Late Move Pruning
			// aka Move Count Based Pruning
			// TODO: dangerous needs testing and tuning
			if Settings.Search.UseLmp {
				if movesSearched >= LmpMovesSearched(depth) {
					s.statistics.LmpCuts++
					continue
				}
			}

			// LMR
			// Late Move Reduction assumes that later moves a rarely
			// exceeding alpha and therefore the search is reduced in
			// depth. This is in effect a soft transition into
			// quiescence search as we usually try the pv move and
			// capturing moves first. In quiescence only capturing
			// moves are searched anyway.
			// newDepth is the "standard" new depth (depth - 1)
			// lmrDepth is set to newDepth and only reduced
			// if conditions apply.
			// TODO: needs testing and tuning
			if Settings.Search.UseLmr {
				// compute reduction from depth and move searched
				if depth >= Settings.Search.LmrDepth &&
					movesSearched >= Settings.Search.LmrMovesSearched {
					lmrDepth -= LmrReduction(depth, movesSearched)
					s.statistics.LmrReductions++
				}
				// make sure not to become negative
				if lmrDepth < 0 {
					lmrDepth = 0
				}
			}
		}
		// ///////////////////////////////////////////////////////

		// ///////////////////////////////////////////////////////
		// DO MOVE
		p.DoMove(move)

		// check if legal move or skip
		if !p.WasLegalMove() {
			p.UndoMove()
			continue
		}

		// we only count legal moves
		s.nodesVisited++
		s.statistics.CurrentVariation.PushBack(move)
		s.sendSearchUpdateToUci()

		// check repetition and 50 moves
		if s.checkDrawRepAnd50(p, 2) {
			value = ValueDraw

		} else {

			// ///////////////////////////////////////////////////////
			// PVS
			// First move in Node will be search with the full window. Due to move
			// ordering we assume this is the PV. Every other move is searched with
			// a null window as we only try to prove that the move is bad (<alpha)
			// or that the move is too good (>beta). If this prove fails we need
			// to research the move again with a full window.
			// https://www.chessprogramming.org/Principal_Variation_Search
			if !Settings.Search.UsePVS || movesSearched == 0 {
				value = -s.search(p, newDepth, ply+1, -beta, -alpha, true, true)
			} else {
				// Null window search after the initial PV search.
				// As depth we use a potentially reduced depth if Late Move Reduction
				// conditions have been met above.
				value = -s.search(p, lmrDepth, ply+1, -alpha-1, -alpha, false, true)
				// If this move improved alpha without exceeding beta we do a proper full window
				// search to get an accurate score.
				// Without LMR we check for value > alpha && value < beta
				// With LMR we re-search when value > alpha
				if value > alpha && !s.stopConditions() {
					// did we actually have a LMR reduction?
					if lmrDepth < newDepth {
						s.statistics.LmrResearches++
						value = -s.search(p, newDepth, ply+1, -beta, -alpha, true, true)
					} else if value < beta {
						s.statistics.PvsResearches++
						value = -s.search(p, newDepth, ply+1, -beta, -alpha, true, true)
					}
				}
			}
			// ///////////////////////////////////////////////////////
		}

		movesSearched++
		s.statistics.CurrentVariation.PopBack()
		p.UndoMove()
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
					// Count beta cuts
					s.statistics.BetaCuts++
					// Count beta cuts on first move
					if movesSearched == 1 {
						s.statistics.BetaCuts1st++
					}
					// store move which caused a beta cut off in this ply
					if Settings.Search.UseKiller && !p.IsCapturingMove(move) {
						myMg.StoreKiller(move)
					}
					// counter for moves which caused a beta cut off
					// we use 1 << depth as an increment to favor deeper searches
					// a more repetitions
					if Settings.Search.UseHistoryCounter {
						s.history.HistoryCount[us][from][to] += 1 << depth
					}
					// store a successful counter move to the previous opponent move
					if Settings.Search.UseCounterMoves {
						lastMove := p.LastMove()
						if lastMove != MoveNone {
							s.history.CounterMoves[lastMove.From()][lastMove.To()] = move
						}
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
		// no beta cutoff - decrease historyCounter for the move
		// we decrease it by only half the increase amount
		if Settings.Search.UseHistoryCounter {
			s.history.HistoryCount[us][from][to] -= 1 << depth
			if s.history.HistoryCount[us][from][to] < 0 {
				s.history.HistoryCount[us][from][to] = 0
			}
		}
	}
	// MOVE LOOP
	// ///////////////////////////////////////////////////////

	// If we did not have at least one legal move
	// then we might have a mate or stalemate
	if movesSearched == 0 && !s.stopConditions() {
		if p.HasCheck() { // mate
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
		s.storeTT(p, depth, ply, bestNodeMove, bestNodeValue, ttType)
	}

	return bestNodeValue
}

// qsearch is a simplified search to counter the horizon effect in depth based
// searches. It continues the search into deeper branches as long as there are
// so called non quiet moves (usually capture, checks, promotions). Only if the
// position is relatively quiet we will compute an evaluation of the position
// to return to the previous depth.
// Look for non quiet moves is supported be the move generator which only
// generates captures or promotions in qsearch (when not in check) and also
// by SEE (Static Exchange Evaluation) to determine winning captured sequences.
func (s *Search) qsearch(p *position.Position, ply int, alpha Value, beta Value, isPV bool) Value {
	if trace {
		s.slog.Debugf("%0*s Ply %-2.d QSearch     a:%-6.d b:%-6.d pv:%-6.v start:  %s", ply, "", ply, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
		defer s.slog.Debugf("%0*s Ply %-2.d QSearch     a:%-6.d b:%-6.d pv:%-6.v end  :  %s", ply, "", ply, alpha, beta, isPV, s.statistics.CurrentVariation.StringUci())
	}

	if s.statistics.CurrentExtraSearchDepth < ply {
		s.statistics.CurrentExtraSearchDepth = ply
	}

	// if we have deactivated qsearch or we have reached our maximum depth
	// we evaluate the position and return the value
	if !Settings.Search.UseQuiescence || ply >= MaxDepth {
		return s.evaluate(p, ply)
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
	hasCheck := p.HasCheck()

	// if in check we simply do a normal search (all moves) in qsearch
	if !hasCheck {
		// get an evaluation for the position
		staticEval := s.evaluate(p, ply)
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
		ttEntry = s.tt.Probe(p.ZobristKey())
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
		s.statistics.CheckInQS++
		mode = movegen.GenAll
	} else {
		mode = movegen.GenNonQuiet
	}

	// ///////////////////////////////////////////////////////
	// MOVE LOOP
	for move := myMg.GetNextMove(p, mode, hasCheck);
		move != MoveNone; move = myMg.GetNextMove(p, mode, hasCheck) {

		// reduce number of moves searched in quiescence
		// by looking at good captures only
		if !hasCheck && !s.goodCapture(p, move) {
			continue
		}

		// ///////////////////////////////////////////////////////
		// DO MOVE
		p.DoMove(move)

		// check if legal move or skip
		if !p.WasLegalMove() {
			p.UndoMove()
			continue
		}

		// we only count legal moves
		s.nodesVisited++
		s.statistics.CurrentVariation.PushBack(move)
		s.sendSearchUpdateToUci()

		// check repetition and 50 moves when in check
		// otherwise only capturing moves are generated
		// which break repetition and 50-moves rule anyway
		if hasCheck && s.checkDrawRepAnd50(p, 2) {
			value = ValueDraw
		} else {
			value = -s.qsearch(p, ply+1, -beta, -alpha, isPV)
		}

		movesSearched++
		s.statistics.CurrentVariation.PopBack()
		p.UndoMove()
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
					// Count beta cuts
					s.statistics.BetaCuts++
					// Count beta cuts on first move
					if movesSearched == 1 {
						s.statistics.BetaCuts1st++
					}
					// counter for moves which caused a beta cut off
					// we use 1 << depth as an increment to favor deeper searches
					// a more repetitions
					if Settings.Search.UseHistoryCounter {
						s.history.HistoryCount[p.NextPlayer()][move.From()][move.To()] += 1 << 1
					}
					// store a successful counter move to the previous opponent move
					if Settings.Search.UseCounterMoves {
						lastMove := p.LastMove()
						if lastMove != MoveNone {
							s.history.CounterMoves[lastMove.From()][lastMove.To()] = move
						}
					}
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
		if p.HasCheck() {
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
		s.storeTT(p, 1, ply, bestNodeMove, bestNodeValue, ttType)
	}

	return bestNodeValue
}

// call evaluation on the position
func (s *Search) evaluate(position *position.Position, ply int) Value {
	s.statistics.LeafPositionsEvaluated++

	var value = ValueNA

	if Settings.Search.UseTT && Settings.Search.UseEvalTT {
		ttEntry := s.tt.Probe(position.ZobristKey())
		if ttEntry != nil { // tt hit
			s.statistics.TTHit++
			s.statistics.EvaluationsFromTT++
			value = valueFromTT(ttEntry.Move.ValueOf(), ply)
		}
	}

	if value == ValueNA {
		s.statistics.Evaluations++
		value = s.eval.Evaluate(position)
	}

	if Settings.Search.UseTT && Settings.Search.UseEvalTT {
		s.storeTT(position, 0, ply, MoveNone, value, EXACT)
	}

	return value
}

// reduce the number of moves searched in quiescence search by trying
// to only look at good captures. Might be improved with SEE in the
// future
func (s *Search) goodCapture(p *position.Position, move Move) bool {
	if Settings.Search.UseSEE {
		// Check SEE score of higher value pieces to low value pieces
		return see(p, move) > 0
	} else {
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
	}
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
		pv.PushBack(ttMatch.Move.MoveOf())
		p.DoMove(ttMatch.Move.MoveOf())
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
	searchLog.SetBackend(searchBackEnd)

	// File backend
	programName, _ := os.Executable()
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")

	// find log path
	logPath, err := util.ResolveFolder(Settings.Log.LogPath)
	if err != nil {
		golog.Println("Log folder could not be found:", err)
		return searchLog
	}
	searchLogFilePath := filepath.Join(logPath, exeName+"_search.log")

	// create file backend
	searchLogFile, err := os.OpenFile(searchLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		golog.Println("Logfile could not be created:", err)
		return searchLog
	}
	backend2 := logging.NewLogBackend(searchLogFile, "", golog.Lmsgprefix)
	backend2Formatter := logging.NewBackendFormatter(backend2, searchLogFormat)
	searchBackEnd2 := logging.AddModuleLevel(backend2Formatter)
	searchBackEnd2.SetLevel(logging.DEBUG, "")
	// multi := logging2.SetBackend(uciBackEnd1, searchBackEnd2)
	searchLog.SetBackend(searchBackEnd2)
	searchLog.Infof("Log %s started at %s:", searchLogFile.Name(), time.Now().String())
	return searchLog
}

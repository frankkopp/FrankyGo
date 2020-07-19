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

package evaluator

import (
	. "github.com/frankkopp/FrankyGo/internal/config"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

func (e *Evaluator) evaluatePawns() *Score {
	tmpScore.MidGameValue = 0
	tmpScore.EndGameValue = 0

	// look on cache table
	if Settings.Eval.UsePawnCache {
		entry := e.pawnCache.getEntry(e.position.PawnKey())
		if entry != nil {
			tmpScore.MidGameValue += entry.score.MidGameValue
			tmpScore.EndGameValue += entry.score.EndGameValue
			return &tmpScore
		}
	}

	// no cache hit - calculate
	for c := White; c <= Black; c++ {
		// out.Println(c.String())

		ourPawns := e.position.PiecesBb(c, Pawn)
		theirPawns := e.position.PiecesBb(c.Flip(), Pawn)

		isolated := BbZero
		doubled := BbZero
		passed := BbZero
		blocked := BbZero
		phalanx := BbZero // both pawns are counted
		supported := BbZero

		pawns := ourPawns
		for pawns != 0 {
			square := pawns.PopLsb()
			neighbours := ourPawns & square.NeighbourFilesMask()

			// isolated pawns
			if neighbours == BbZero {
				isolated |= square.Bb()
			}

			// get a mask for the forward squares
			var rayForward Bitboard
			if c == White {
				rayForward = square.Ray(N)
			} else {
				rayForward = square.Ray(S)
			}

			// doubled pawn
			doubled |= ^square.Bb() & ourPawns & rayForward

			// passed pawns - no opp pawns in front or on neighbouring
			// files and also no own pawns in front
			if (theirPawns&square.PassedPawnMask(c))|(ourPawns&rayForward) == BbZero {
				passed |= square.Bb()
			}

			// blocked pawns
			if ((ourPawns | theirPawns) & rayForward) != BbZero {
				blocked |= square.Bb()
			}

			// pawns as neighbours in a row = phalanx
			phalanx |= ourPawns & neighbours & square.RankBb()

			// pawn as neighbours in the row forward = supported pawns
			supported |= ourPawns & neighbours & square.To(c.MoveDirection()).RankBb()
		}

		// out.Printf("Isolated : %d %s\n", isolated.PopCount(), isolated.StringGrouped())
		// out.Printf("Doubled  : %d %s\n", doubled.PopCount(), doubled.StringGrouped())
		// out.Printf("Passed   : %d %s\n", passed.PopCount(), passed.StringGrouped())
		// out.Printf("Blocked  : %d %s\n", blocked.PopCount(), blocked.StringGrouped())
		// out.Printf("Phalanx  : %d %s\n", phalanx.PopCount(), phalanx.StringGrouped())
		// out.Printf("Supported: %d %s\n", supported.PopCount(), supported.StringGrouped())

		tmpMid := int16(isolated.PopCount()) * Settings.Eval.PawnIsolatedMidMalus
		tmpEnd := int16(isolated.PopCount()) * Settings.Eval.PawnIsolatedEndMalus
		tmpMid += int16(doubled.PopCount()) * Settings.Eval.PawnDoubledMidMalus
		tmpEnd += int16(doubled.PopCount()) * Settings.Eval.PawnDoubledEndMalus
		tmpMid += int16(passed.PopCount()) * Settings.Eval.PawnPassedMidBonus
		tmpEnd += int16(passed.PopCount()) * Settings.Eval.PawnPassedEndBonus
		tmpMid += int16(blocked.PopCount()) * Settings.Eval.PawnBlockedMidMalus
		tmpEnd += int16(blocked.PopCount()) * Settings.Eval.PawnBlockedEndMalus
		tmpMid += int16(phalanx.PopCount()) * Settings.Eval.PawnPhalanxMidBonus
		tmpEnd += int16(phalanx.PopCount()) * Settings.Eval.PawnPhalanxEndBonus
		tmpMid += int16(supported.PopCount()) * Settings.Eval.PawnSupportedMidBonus
		tmpEnd += int16(supported.PopCount()) * Settings.Eval.PawnSupportedEndBonus

		// add it to total score
		if c == White {
			tmpScore.MidGameValue += tmpMid
			tmpScore.EndGameValue += tmpEnd
		} else {
			tmpScore.MidGameValue -= tmpMid
			tmpScore.EndGameValue -= tmpEnd
		}

		// e.log.Debugf("Raw pawn eval for %s: mid:%d end:%d", c.String(), tmpMid, tmpEnd)
	}

	// store in cache
	if Settings.Eval.UsePawnCache {
		e.pawnCache.put(e.position.PawnKey(), &tmpScore)
	}

	// e.log.Debugf("Pawn Eval: %d/%d", tmpScore.MidGameValue, tmpScore.EndGameValue)

	return &tmpScore
}

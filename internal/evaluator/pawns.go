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
	. "github.com/frankkopp/FrankyGo/internal/config"
	. "github.com/frankkopp/FrankyGo/internal/types"
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

		out.Println(c.String())

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
			if ourPawns&rayForward != BbZero {
				doubled |= square.Bb()
			}

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
			if ourPawns&neighbours&square.To(-c.MoveDirection()).RankBb() != BbZero {
				supported |= square.Bb()
			}
		}

		// out.Printf("Isolated :\n%s", isolated.StringBoard())
		// out.Printf("Doubled  :\n%s", doubled.StringBoard())
		// out.Printf("Passed   :\n%s", passed.StringBoard())
		// out.Printf("Blocked  :\n%s", blocked.StringBoard())
		// out.Printf("Phalanx  :\n%s", phalanx.StringBoard())
		// out.Printf("Supported:\n%s", supported.StringBoard())

	}

	// store in cache
	if Settings.Eval.UsePawnCache {
		e.pawnCache.put(e.position.PawnKey(), &tmpScore)
	}

	return &tmpScore
}

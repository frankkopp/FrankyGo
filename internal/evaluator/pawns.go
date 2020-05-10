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
	// DEBUG
	//  Prototype
	tmpScore.MidGameValue = int16(e.position.PiecesBb(White, Pawn).PopCount() - e.position.PiecesBb(Black, Pawn).PopCount())
	tmpScore.EndGameValue = tmpScore.MidGameValue

	// store in cache
	if Settings.Eval.UsePawnCache {
		e.pawnCache.put(e.position.PawnKey(), &tmpScore)
	}

	return &tmpScore
}

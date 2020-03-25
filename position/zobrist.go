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

package position

import (
	. "github.com/frankkopp/FrankyGo/types"
)

// helper data structure for Zobrist IDs for chess positions
type zobrist struct {
	pieces         [PieceLength][SqLength]Key
	castlingRights [CastlingRightsLength]Key
	enPassantFile  [8]Key
	nextPlayer     Key
}

var zobristBase = zobrist{}

func initZobrist() {
	// Zobrist Key initialization
	r := NewRandom(1070372)
	for pc := PieceNone; pc < PieceLength; pc++ {
		for sq := SqA1; sq <= SqH8; sq++ {
			zobristBase.pieces[pc][sq] = Key(r.Rand64())
		}
	}
	for cr := CastlingNone; cr <= CastlingAny; cr++ {
		zobristBase.castlingRights[cr] = Key(r.Rand64())
	}
	for f := FileA; f <= FileH; f++ {
		zobristBase.enPassantFile[f] = Key(r.Rand64())
	}
	zobristBase.nextPlayer = Key(r.Rand64())
}

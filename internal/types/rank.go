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

package types

// Rank represents a chess board rank 1-8
type Rank uint8

// Rank represents a chess board rank 1-8
//noinspection GoUnusedConst
const (
	Rank1      Rank = iota
	Rank2      Rank = iota
	Rank3      Rank = iota
	Rank4      Rank = iota
	Rank5      Rank = iota
	Rank6      Rank = iota
	Rank7      Rank = iota
	Rank8      Rank = iota
	RankNone   Rank = iota
	RankLength      = RankNone
)

// IsValid checks if f represents a valid file
func (r Rank) IsValid() bool {
	return r < RankNone
}

// Bb returns a Bitboard of the given rank
func (r Rank) Bb() Bitboard {
	return rankBb[r]
}

const rankLabels string = "12345678"

// String returns a string letter for the file (e.g. a - h)
// if r is not a valid rank returns "-"
func (r Rank) String() string {
	if r > Rank8 {
		return "-"
	}
	return string(rankLabels[r])
}

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

package types

// Rank represents a chess board rank 1-8
//  Rank1      Rank = 0
//  Rank2      Rank = 1
//  Rank3      Rank = 2
//  Rank4      Rank = 3
//  Rank5      Rank = 4
//  Rank6      Rank = 5
//  Rank7      Rank = 6
//  Rank8      Rank = 7
//  RankNone   Rank = 8
//  RankLength      = RankNone
type Rank uint8

// Rank represents a chess board rank 1-8
//noinspection GoUnusedConst
const (
	Rank1      Rank = 0
	Rank2      Rank = 1
	Rank3      Rank = 2
	Rank4      Rank = 3
	Rank5      Rank = 4
	Rank6      Rank = 5
	Rank7      Rank = 6
	Rank8      Rank = 7
	RankNone   Rank = 8
	RankLength      = RankNone
)

// IsValid checks if f represents a valid file
func (r Rank) IsValid() bool {
	return r < RankNone
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

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

import (
	"fmt"
)

// Square represent exactly on square on a chess board.
//  SqA1   Square = iota // 0
//	SqB1   Square = iota // 1
//	SqC1   Square = iota
//	SqD1   Square = iota
//  ...
//  SqG8   Square = iota
//	SqH8   Square = iota // 63
//	SqNone Square = iota // 64
type Square uint8

//noinspection GoUnusedConst
const (
	SqA1 Square = iota // 0
	SqB1               // 1
	SqC1               // 2
	SqD1
	SqE1
	SqF1
	SqG1
	SqH1
	SqA2
	SqB2
	SqC2
	SqD2
	SqE2
	SqF2
	SqG2
	SqH2
	SqA3
	SqB3
	SqC3
	SqD3
	SqE3
	SqF3
	SqG3
	SqH3
	SqA4
	SqB4
	SqC4
	SqD4
	SqE4
	SqF4
	SqG4
	SqH4
	SqA5
	SqB5
	SqC5
	SqD5
	SqE5
	SqF5
	SqG5
	SqH5
	SqA6
	SqB6
	SqC6
	SqD6
	SqE6
	SqF6
	SqG6
	SqH6
	SqA7
	SqB7
	SqC7
	SqD7
	SqE7
	SqF7
	SqG7
	SqH7
	SqA8
	SqB8
	SqC8
	SqD8
	SqE8
	SqF8
	SqG8
	SqH8   // 63
	SqNone // 64
)

// IsValid checks a value of type square if it represents a valid
// square on a chess board (e.q. sq < 64).
func (sq Square) IsValid() bool {
	return sq < SqNone
}

// FileOf returns the file of the square
func (sq Square) FileOf() File {
	return File(sq & 7)
}

// RankOf returns the rank of the square
func (sq Square) RankOf() Rank {
	return Rank(sq >> 3)
}

// MakeSquare returns a square based on the string given or SqNone if
// no valid square could be read from the string
func MakeSquare(s string) Square {
	file := File(s[0] - 'a')
	rank := Rank(s[1] - '1')
	if !file.IsValid() || !rank.IsValid() {
		return SqNone
	}
	return SquareOf(file, rank)
}

// SquareOf returns a square from file and rank
// Returns SqNone for invalid files or ranks
func SquareOf(f File, r Rank) Square {
	if !f.IsValid() || !r.IsValid() {
		return SqNone
	}
	return Square((int(r) << 3) + int(f))
}

// To returns the square on the chess board in the given direction
func (sq Square) To(d Direction) Square {
	// Precomputed
	// order:  North, East, South, West, Northeast, Southeast, Southwest, Northwest
	switch d {
	case North:
		return sqTo[sq][0]
	case East:
		return sqTo[sq][1]
	case South:
		return sqTo[sq][2]
	case West:
		return sqTo[sq][3]
	case Northeast:
		return sqTo[sq][4]
	case Southeast:
		return sqTo[sq][5]
	case Southwest:
		return sqTo[sq][6]
	case Northwest:
		return sqTo[sq][7]
	default:
		panic(fmt.Sprintf("Invalid direction %d", d))
	}
}

// String returns a string of the file letter and rank number (e.g. e5)
// if the sq is not a valid square returns "-"
func (sq Square) String() string {
	if !sq.IsValid() {
		return "-"
	}
	return sq.FileOf().String() + sq.RankOf().String()
}

// ///////////////////////////////////////
// Initialization
// ///////////////////////////////////////

var sqTo [SqLength][8]Square

func init() {
	for sq := SqA1; sq < SqNone; sq++ {
		for i, dir := range Directions {
			sqTo[sq][i] = sq.toPreCompute(dir)
		}
	}
}

func (sq Square) toPreCompute(d Direction) Square {
	// overflow To south or north are easily detected <0 ot >63
	// east and west need check
	switch d {
	case North:
		sq += Square(d)
	case East:
		if sq.FileOf() < FileH {
			sq += Square(d)
		} else {
			return SqNone
		}
	case South:
		sq += Square(d)
	case West:
		if sq.FileOf() > FileA {
			sq += Square(d)
		} else {
			return SqNone
		}
	case Northeast:
		if sq.FileOf() < FileH {
			sq += Square(d)
		} else {
			return SqNone
		}
	case Southeast:
		if sq.FileOf() < FileH {
			sq += Square(d)
		} else {
			return SqNone
		}
	case Southwest:
		if sq.FileOf() > FileA {
			sq += Square(d)
		} else {
			return SqNone
		}
	case Northwest:
		if sq.FileOf() > FileA {
			sq += Square(d)
		} else {
			return SqNone
		}
	default:
		panic(fmt.Sprintf("Invalid direction %d", d))
	}
	if sq.IsValid() {
		return sq
	} else {
		return SqNone
	}
}

/*
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

import "fmt"

// A square represent exactly on square on a chess board.
type Square uint8

//noinspection GoUnusedConst
const (
	SqA1   Square = iota // 0
	SqB1   Square = iota // 1
	SqC1   Square = iota
	SqD1   Square = iota
	SqE1   Square = iota
	SqF1   Square = iota
	SqG1   Square = iota
	SqH1   Square = iota
	SqA2   Square = iota
	SqB2   Square = iota
	SqC2   Square = iota
	SqD2   Square = iota
	SqE2   Square = iota
	SqF2   Square = iota
	SqG2   Square = iota
	SqH2   Square = iota
	SqA3   Square = iota
	SqB3   Square = iota
	SqC3   Square = iota
	SqD3   Square = iota
	SqE3   Square = iota
	SqF3   Square = iota
	SqG3   Square = iota
	SqH3   Square = iota
	SqA4   Square = iota
	SqB4   Square = iota
	SqC4   Square = iota
	SqD4   Square = iota
	SqE4   Square = iota
	SqF4   Square = iota
	SqG4   Square = iota
	SqH4   Square = iota
	SqA5   Square = iota
	SqB5   Square = iota
	SqC5   Square = iota
	SqD5   Square = iota
	SqE5   Square = iota
	SqF5   Square = iota
	SqG5   Square = iota
	SqH5   Square = iota
	SqA6   Square = iota
	SqB6   Square = iota
	SqC6   Square = iota
	SqD6   Square = iota
	SqE6   Square = iota
	SqF6   Square = iota
	SqG6   Square = iota
	SqH6   Square = iota
	SqA7   Square = iota
	SqB7   Square = iota
	SqC7   Square = iota
	SqD7   Square = iota
	SqE7   Square = iota
	SqF7   Square = iota
	SqG7   Square = iota
	SqH7   Square = iota
	SqA8   Square = iota
	SqB8   Square = iota
	SqC8   Square = iota
	SqD8   Square = iota
	SqE8   Square = iota
	SqF8   Square = iota
	SqG8   Square = iota
	SqH8   Square = iota // 63
	SqNone Square = iota // 64
)

// isValid checks a value of type square if it represents a valid
// square on a chess board (e.q. sq < 63).
func (sq Square) isValid() bool {
	return sq < SqNone
}

// Returns the file of the square
func (sq Square) FileOf() File {
	return File(sq & 7)
}

// Returns the rank of the square
func (sq Square) RankOf() Rank {
	return Rank(sq >> 3)
}

// returns a string of the file letter and rank number (e.g. e5)
// if the sq is not a valid square returns "--"
func (sq Square) str() string {
	if !sq.isValid() {
		return "--"
	}
	return sq.FileOf().str() + sq.RankOf().str()
}

// Return square from file and rank
// Returns SqNone for invalid files or ranks
func squareOf(f File, r Rank) Square {
	if !f.isValid() || !r.isValid() {
		return SqNone
	}
	return Square((int(r) << 3) + int(f))
}

// Returns the square on the chess board in the given direction
// TODO: should maybe pre-calculated
func (sq Square) to(d Direction) Square {
	// overflow to south or north are easily detected <0 ot >63
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
	if sq.isValid() {
		return sq
	} else {
		return SqNone
	}
}

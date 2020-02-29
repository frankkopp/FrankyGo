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

import (
	"FrankyGo/config"
	"fmt"
	"log"
)

// 64 bit for each square on the board
type Bitboard uint64

// various constant bitboards for convenience
//noinspection ALL
const (
	BbZero Bitboard = Bitboard(0)
	BbAll  Bitboard = ^BbZero
	BbOne  Bitboard = Bitboard(1)

	FileA_Bb Bitboard = 0x0101010101010101
	FileB_Bb Bitboard = FileA_Bb << 1
	FileC_Bb Bitboard = FileA_Bb << 2
	FileD_Bb Bitboard = FileA_Bb << 3
	FileE_Bb Bitboard = FileA_Bb << 4
	FileF_Bb Bitboard = FileA_Bb << 5
	FileG_Bb Bitboard = FileA_Bb << 6
	FileH_Bb Bitboard = FileA_Bb << 7

	Rank1_Bb Bitboard = 0xFF
	Rank2_Bb Bitboard = Rank1_Bb << (8 * 1)
	Rank3_Bb Bitboard = Rank1_Bb << (8 * 2)
	Rank4_Bb Bitboard = Rank1_Bb << (8 * 3)
	Rank5_Bb Bitboard = Rank1_Bb << (8 * 4)
	Rank6_Bb Bitboard = Rank1_Bb << (8 * 5)
	Rank7_Bb Bitboard = Rank1_Bb << (8 * 6)
	Rank8_Bb Bitboard = Rank1_Bb << (8 * 7)
)

// Internal square to bitboard array. Needs to be initialized
var sqBb [64]Bitboard

// Pre computes various bitboards to avoid runtime calculation
// Will only run once (checks an initialized flag)
func initBb() {
	for i := SqA1; i < SqNone; i++ {
		sqBb[i] = i.bitboard_()
	}
}

// Returns a Bitboard of the square by shifting the
// square onto an empty bitboards.
// Usually one would use Bitboard() after initializing with InitBb
func (sq Square) bitboard_() Bitboard {
	return Bitboard(uint64(1) << sq)
}

// Returns a Bitboard of the square by accessing the pre calculated
// square to bitboard array.
// Initialize with InitBb() before use. Throws panic otherwise.
func (sq Square) Bitboard() Bitboard {
	// assertion
	if config.DEBUG && !initialized {
		log.Printf("Warning: Bitboards not initialized. Using runtime calculation.\n")
		return sq.bitboard_()
	}
	return sqBb[sq]
}

// sets the corresponding bit of the bitboard_ for the square
func (b Bitboard) put(s Square) Bitboard {
	return b | s.Bitboard()
}

// sets the corresponding bit of the bitboard_ for the square
func (b Bitboard) remove(s Square) Bitboard {
	return b & s.Bitboard()
}

// returns a string representation of the 64 bits
func (b Bitboard) str() string {
	return fmt.Sprintf("%-0.64b", b)
}

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

package transpositiontable

import (
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

// TtEntry struct is the data structure for each entry in the transposition
// table. Each entry has 16-bytes (128-bits).
type TtEntry struct {
	// struct is partially bit encoded to make it more compact
	// and stay <= 16 byte
	key   Key    // 64-bit Zobrist Key
	move  uint16 // 16-bit move part of a Move - convert with Move(e.Move)
	eval  int16  // 16-bit evaluation value by static evaluator
	value int16  // 16-bit value during search
	vmeta uint16 // 16-bit depth 7-bit, vtype 2-bit, age 3-bit
	// depth 7-bit 0-127
	// vtype 3-bit 0-7   0=used 1=generated, not used, >1 older generation
	// age 2-bit None, Exact, Alpha (upper), Beta (lower)
}

const (
	// TtEntrySize is the size in bytes for each TtEntry
	TtEntrySize = 16 // 16 bytes

	ageMask    = uint16(0b0000_0000_0000_0111)
	vtypeMask  = uint16(0b0000_0000_0001_1000)
	vtypeShift = uint16(3)
	depthMask  = uint16(0b0000_1111_1110_0000)
	depthShift = uint16(5)
)

func (e *TtEntry) decreaseAge() {
	// age is stored in the last 2 bits --> we can just decrease
	if e.Age() > 0 {
		e.vmeta--
	}
}

func (e *TtEntry) increaseAge() {
	// age is stored in the last 2 bits --> we can just increase
	if e.Age() <= 7 {
		e.vmeta++
	}
}

func (e *TtEntry) Key() Key {
	return e.key
}

func (e *TtEntry) Move() Move {
	return Move(e.move)
}

func (e *TtEntry) Value() Value {
	return Value(e.value)
}

func (e *TtEntry) Eval() Value {
	return Value(e.eval)
}

func (e *TtEntry) Depth() int8 {
	return int8((e.vmeta & depthMask) >> depthShift)
}

func (e *TtEntry) Age() int8 {
	// last 3 bits
	return int8(e.vmeta & ageMask)
}

func (e *TtEntry) Vtype() ValueType {
	return ValueType((e.vmeta & vtypeMask) >> vtypeShift)
}

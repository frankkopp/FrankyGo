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
	"strings"

	"github.com/frankkopp/FrankyGo/internal/assert"
)

// Move is a 32bit unsigned int type for encoding chess moves as a primitive data type
// 16 bits for move encoding - 16 bits for sort value
//  MoveNone Move = 0
//  BITMAP 32-bit
//  |-value ------------------------|-Move -------------------------|
//  3 3 2 2 2 2 2 2 2 2 2 2 1 1 1 1 | 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 | 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
//  --------------------------------|--------------------------------
//                                  |                     1 1 1 1 1 1  to
//                                  |         1 1 1 1 1 1              from
//                                  |     1 1                          promotion piece type (pt-2 > 0-3)
//                                  | 1 1                              move type
//  1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 |                                  move sort value
type Move uint32

const (
	// MoveNone empty non valid move
	MoveNone Move = 0
)

// CreateMove returns an encoded Move instance
func CreateMove(from Square, to Square, t MoveType, promType PieceType) Move {
	if promType < Knight {
		promType = Knight
	}
	// promType will be reduced to 2 bits (4 values) Knight, Bishop, Rook, Queen
	// therefore we subtract the Knight value from the promType to get
	// value between 0 and 3 (0b00 - 0b11)
	return Move(to) |
		Move(from)<<fromShift |
		Move(promType-Knight)<<promTypeShift |
		Move(t)<<typeShift
}

// CreateMoveValue returns an encoded Move instance including a sort value
func CreateMoveValue(from Square, to Square, t MoveType, promType PieceType, value Value) Move {
	if promType < Knight {
		promType = Knight
	}
	// promType will be reduced to 2 bits (4 values) Knight, Bishop, Rook, Queen
	// therefore we subtract the Knight value from the promType to get
	// value between 0 and 3 (0b00 - 0b11)
	return Move(value-ValueNA)<<valueShift |
		Move(to) |
		Move(from)<<fromShift |
		Move(promType-Knight)<<promTypeShift |
		Move(t)<<typeShift
}

// MoveType returns the type of the move as defined in MoveType
// Normal, Promotion, EnPassant, Castling
func (m Move) MoveType() MoveType {
	return MoveType((m & moveTypeMask) >> typeShift)
}

// PromotionType returns the PieceType considered for promotion when
// move type is also MoveType.Promotion.
// Must be ignored when move type is not MoveType.Promotion.
func (m Move) PromotionType() PieceType {
	return PieceType((m&promTypeMask)>>promTypeShift) + Knight
}

// To returns the to-Square of the move
func (m Move) To() Square {
	return Square(m & toMask)
}

// From returns the from-Square of the move
func (m Move) From() Square {
	return Square((m & fromMask) >> fromShift)
}

// MoveOf returns the move without any value (least 16-bits)
func (m Move) MoveOf() Move {
	return m & moveMask
}

// ValueOf returns the sort value for the move used in the move generator
func (m Move) ValueOf() Value {
	return Value((m&valueMask)>>valueShift) + ValueNA
}

// SetValue encodes the given value into the high 16-bit of the move
func (m *Move) SetValue(v Value) Move {
	if assert.DEBUG {
		assert.Assert(v == ValueNA || v.IsValid(), "Invalid value value: %d", v)
	}
	// can't store a value on MoveNone
	if *m == MoveNone {
		return *m
	}
	// when saving a value to a move we shift value to a positive integer
	// (0-VALUE_NONE) and encode it into the move. For retrieving we then shift
	// the value back to a range from VALUE_NONE to VALUE_INF
	*m = *m&moveMask | Move(v-ValueNA)<<valueShift
	return *m
}

// IsValid check if the move has valid squares, promotion type and move type.
// MoveNone is not a valid move in this sense.
func (m Move) IsValid() bool {
	return m != MoveNone &&
		m.From().IsValid() &&
		m.To().IsValid() &&
		m.PromotionType().IsValid() &&
		m.MoveType().IsValid() &&
		(m.ValueOf() == ValueNA || m.ValueOf().IsValid())
}

// String string representation of a move which is UCI compatible
func (m Move) String() string {
	if m == MoveNone {
		return "Move: { MoveNone }"
	}
	return fmt.Sprintf("Move: { %-5s  type:%1s  prom:%1s  value:%-6d  (%d) }",
		m.StringUci(), m.MoveType().String(), m.PromotionType().Char(), m.ValueOf(), m)
}

// StringUci string representation of a move which is UCI compatible
func (m Move) StringUci() string {
	if m == MoveNone {
		return "NoMove"
	}
	var os strings.Builder
	os.WriteString(m.From().String())
	os.WriteString(m.To().String())
	if m.MoveType() == Promotion {
		os.WriteString(m.PromotionType().Char())
	}
	return os.String()
}

// StringBits returns a string with details of a Move
// E.g. Move { From[001100](e2) To[011100](e4) Prom[11](N) tType[00](n) value[0000000000000000](0) (796)}
func (m Move) StringBits() string {
	return fmt.Sprintf(
		"Move { From[%-0.6b](%s) To[%-0.6b](%s) Prom[%-0.2b](%s) tType[%-0.2b](%s) value[%-0.16b](%d) (%d)}",
		m.From(), m.From().String(),
		m.To(), m.To().String(),
		m.PromotionType(), (m.PromotionType()).Char(),
		m.MoveType(), m.MoveType().String(),
		m.ValueOf(), m.ValueOf(),
		m)
}

/* @formatter:off
   BITMAP 32-bit
   |-value ------------------------|-Move -------------------------|
   3 3 2 2 2 2 2 2 2 2 2 2 1 1 1 1 | 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0
   1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 | 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
   --------------------------------|--------------------------------
                                   |                     1 1 1 1 1 1  to
                                   |         1 1 1 1 1 1              from
                                   |     1 1                          promotion piece type (pt-2 > 0-3)
                                   | 1 1                              move type
   1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 |                                  move sort value
*/ // @formatter:on

const (
	fromShift     uint = 6
	promTypeShift uint = 12
	typeShift     uint = 14
	valueShift    uint = 16

	squareMask   Move = 0x3F
	toMask            = squareMask
	fromMask          = squareMask << fromShift
	promTypeMask Move = 3 << promTypeShift
	moveTypeMask Move = 3 << typeShift
	moveMask     Move = 0xFFFF               // first 16-bit
	valueMask    Move = 0xFFFF << valueShift // second 16-bit
)

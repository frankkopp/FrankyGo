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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePiece(t *testing.T) {
	type args struct {
		c  Color
		pt PieceType
	}
	tests := []struct {
		name string
		args args
		want Piece
	}{
		{"White King", args{White, King}, WhiteKing},
		{"White King", args{Black, King}, BlackKing},
		{"White King", args{White, Knight}, WhiteKnight},
		{"White King", args{Black, Knight}, BlackKnight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePiece(tt.args.c, tt.args.pt); got != tt.want {
				t.Errorf("MakePiece() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPiece_ValueOf(t *testing.T) {
	tests := []struct {
		name string
		p    Piece
		want Value
	}{
		{ "White King", WhiteKing, 2000},
		{ "White King", BlackKing, 2000},
		{ "White King", WhiteBishop, 330},
		{ "White King", BlackKnight, 320},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ValueOf(); got != tt.want {
				t.Errorf("ValueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPieceFromChar(t *testing.T) {
	assert.Equal(t, PieceNone, PieceFromChar(""))
	assert.Equal(t, PieceNone, PieceFromChar("nnn"))
	assert.Equal(t, PieceNone, PieceFromChar("-"))
	assert.Equal(t, WhiteKing, PieceFromChar("K"))
	assert.Equal(t, BlackKing, PieceFromChar("k"))
	assert.Equal(t, WhiteKnight, PieceFromChar("N"))
	assert.Equal(t, BlackKnight, PieceFromChar("n"))
}

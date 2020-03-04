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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateMove(t *testing.T) {
	type args struct {
		from     Square
		to       Square
		t        MoveType
		promType PieceType
	}
	tests := []struct {
		name string
		args args
		want Move
	}{
		{"e2e4", args{SqE2, SqE4, Normal, PtNone}, Move(796)},
		{"e1g1 castling", args{SqE1, SqG1, Castling, PtNone}, Move(49414)},
		{"a2a1Q", args{SqA2, SqA1, Promotion, Queen}, Move(29184)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMove(tt.args.from, tt.args.to, tt.args.t, tt.args.promType)
			fmt.Printf("%s\n", got.StrBits())
			if got != tt.want {
				t.Errorf("CreateMove() = \n%v, want \n%v", got.StrBits(), tt.want.StrBits())
			}
		})
	}
}

func TestMove_SetValue(t *testing.T) {
	m := CreateMove(SqE2, SqE4, Normal, PtNone)
	m.SetValue(999)
	assert.Equal(t, Value(999), m.ValueOf())

	m = CreateMove(SqE2, SqE4, Promotion, Queen)
	m.SetValue(ValueMax)
	assert.Equal(t, ValueMax, m.ValueOf())
}

func Test_Str(t *testing.T) {
	assert.Equal(t,"e2e4", CreateMove(SqE2, SqE4, Normal, PtNone).Str())
	assert.Equal(t,"e7e5", CreateMove(SqE7, SqE5, Normal, PtNone).Str())
	assert.Equal(t,"a2a1Q", CreateMove(SqA2, SqA1, Promotion, Queen).Str())
}

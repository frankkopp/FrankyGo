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

func TestCastlingRights_Has(t *testing.T) {
	assert := assert.New(t)
	var cr CastlingRights
	cr.Add(CastlingAny)
	assert.Equal(CastlingAny, cr)

	assert.True(cr.Has(CastlingWhiteOO))
	cr.Remove(CastlingWhiteOO)
	assert.Equal(0b1110, int(cr))
	assert.False(cr.Has(CastlingWhiteOO))

	assert.True(cr.Has(CastlingBlack))
	assert.True(cr.Has(CastlingBlackOO))
	assert.True(cr.Has(CastlingBlackOOO))
	cr.Remove(CastlingBlack)
	assert.False(cr.Has(CastlingBlack))
	assert.False(cr.Has(CastlingBlackOO))
	assert.False(cr.Has(CastlingBlackOOO))
	assert.True(cr.Has(CastlingWhiteOOO))
}

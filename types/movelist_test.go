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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Deque(t *testing.T) {
	var moveList = MoveList{}
	moveList.SetMinCapacity(8)
	moveList.PushBack(CreateMove(SqG1, SqF3, Normal, PtNone))
	moveList.PushBack(CreateMove(SqB8, SqC6, Normal, PtNone))
	moveList.PushFront(CreateMove(SqE7, SqE5, Normal, PtNone))
	moveList.PushFront(CreateMove(SqE2, SqE4, Normal, PtNone))
	assert.Equal(t, 4, moveList.Len())
}

func TestMoveList_String(t *testing.T) {
	var moveList = MoveList{}
	moveList.SetMinCapacity(8)
	moveList.PushBack(CreateMove(SqG1, SqF3, Normal, PtNone))
	moveList.PushBack(CreateMove(SqB8, SqC6, Normal, PtNone))
	moveList.PushFront(CreateMove(SqE7, SqE5, Normal, PtNone))
	moveList.PushFront(CreateMove(SqE2, SqE4, Normal, PtNone))
	assert.Equal(t, 4, moveList.Len())
	assert.Equal(t,"e2e4 e7e5 g1f3 b8c6", moveList.StringUci())
	assert.Equal(t,"MoveList: [4] { Move: { e2e4 type:n prom:N (796) }, Move: { e7e5 type:n prom:N (3364) }, Move: { g1f3 type:n prom:N (405) }, Move: { b8c6 type:n prom:N (3690) } }",
		moveList.String())
}
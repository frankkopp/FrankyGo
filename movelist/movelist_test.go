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

package movelist

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/types"
)

func Test_Deque(t *testing.T) {
	var moveList = MoveList{}
	moveList.SetMinCapacity(8)
	moveList.PushBack(types.CreateMove(types.SqG1, types.SqF3, types.Normal, types.PtNone))
	moveList.PushBack(types.CreateMove(types.SqB8, types.SqC6, types.Normal, types.PtNone))
	moveList.PushFront(types.CreateMove(types.SqE7, types.SqE5, types.Normal, types.PtNone))
	moveList.PushFront(types.CreateMove(types.SqE2, types.SqE4, types.Normal, types.PtNone))
	assert.Equal(t, 4, moveList.Len())
}

func TestMoveList_String(t *testing.T) {
	var moveList = MoveList{}
	moveList.SetMinCapacity(8)
	moveList.PushBack(types.CreateMove(types.SqG1, types.SqF3, types.Normal, types.PtNone))
	moveList.PushBack(types.CreateMove(types.SqB8, types.SqC6, types.Normal, types.PtNone))
	moveList.PushFront(types.CreateMove(types.SqE7, types.SqE5, types.Normal, types.PtNone))
	moveList.PushFront(types.CreateMove(types.SqE2, types.SqE4, types.Normal, types.PtNone))
	assert.Equal(t, 4, moveList.Len())
	assert.Equal(t,"e2e4 e7e5 g1f3 b8c6", moveList.StringUci())
	assert.Equal(t,"MoveList: [4] { Move: { e2e4 type:n prom:N value:-15001 (796) }, Move: { e7e5 type:n prom:N value:-15001 (3364) }, Move: { g1f3 type:n prom:N value:-15001 (405) }, Move: { b8c6 type:n prom:N value:-15001 (3690) } }",
		moveList.String())
}

func TestMoveList_Sort(t *testing.T) {
	var moveList = MoveList{}
	moveList.SetMinCapacity(8)
	moveList.PushBack(types.CreateMoveValue(types.SqG1, types.SqF3, types.Normal, types.PtNone, 111))
	moveList.PushBack(types.CreateMoveValue(types.SqB8, types.SqC6, types.Normal, types.PtNone, 333))
	moveList.PushBack(types.CreateMoveValue(types.SqE7, types.SqE5, types.Normal, types.PtNone, 222))
	moveList.PushBack(types.CreateMoveValue(types.SqA2, types.SqA3, types.Normal, types.PtNone, 222))
	moveList.PushBack(types.CreateMoveValue(types.SqE2, types.SqE4, types.Normal, types.PtNone, 444))
	assert.Equal(t, 5, moveList.Len())
	sort.Stable(&moveList)
	assert.Equal(t,"e2e4 b8c6 e7e5 a2a3 g1f3", moveList.StringUci())
}

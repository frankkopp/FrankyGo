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

package transpositiontable

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/franky_logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var logTest = franky_logging.GetLog("test")

func TestEntrySize(t *testing.T) {
	e := TtEntry{}
	assert.Equal(t, 16, unsafe.Sizeof(e))
	logTest.Debugf("Size of Entry %d bytes", unsafe.Sizeof(e))
}

func TestNew(t *testing.T) {

	tt := New(2)
	assert.Equal(t, uint64(131_072), tt.maxNumberOfEntries)
	assert.Equal(t, 131_072, cap(tt.data))

	tt = New(64)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = New(100)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = New(4_096)
	assert.Equal(t, uint64(268_435_456), tt.maxNumberOfEntries)
	assert.Equal(t, 268_435_456, cap(tt.data))

	tt = New(35_000)
	assert.Equal(t, uint64(2_147_483_648), tt.maxNumberOfEntries)
	assert.Equal(t, 2_147_483_648, cap(tt.data))
	assert.Equal(t, 2_147_483_648, len(tt.data))
	assert.Equal(t, 32_768*MB, tt.sizeInByte)
	for i, _ := range tt.data {
		tt.data[i].Key = position.Key(i)
	}
	assert.Equal(t, position.Key(0), tt.data[0].Key)
	assert.Equal(t, position.Key(2_147_483_647), tt.data[2_147_483_647].Key)
}

func TestGetAndProbe(t *testing.T) {
	// setup
	tt := New(64)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	pos := position.New()
	move := CreateMove(SqE2, SqE4, Normal, PtNone)
	tt.data[tt.hash(pos.ZobristKey())] = TtEntry{
		Key:        pos.ZobristKey(),
		Move:       move,
		Depth:      5,
		Age:        1,
		Type:       Vnone,
		MateThreat: false,
	}
	tt.numberOfEntries++

	// test to get unaltered entry
	e := tt.GetEntry(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.Key )
	assert.Equal(t, move, e.Move)
	assert.EqualValues(t, 5, e.Depth)
	assert.EqualValues(t, 1, e.Age)
	assert.Equal(t, Vnone, e.Type)

	// age must be reduced by 1
	e = tt.Probe(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.Key )
	assert.Equal(t, move, e.Move)
	assert.EqualValues(t, 5, e.Depth)
	assert.EqualValues(t, 0, e.Age)
	assert.Equal(t, Vnone, e.Type)

	// age does not go below 0
	e = tt.Probe(pos.ZobristKey())
	assert.EqualValues(t, 0, e.Age)

	// not in tt
	pos.DoMove(move)
	e = tt.Probe(pos.ZobristKey())
	assert.Nil(t, e)
}

func TestClear(t *testing.T) {
	// setup
	tt := New(1)

	pos := position.New()
	move := CreateMove(SqE2, SqE4, Normal, PtNone)
	tt.data[tt.hash(pos.ZobristKey())] = TtEntry{
		Key:        pos.ZobristKey(),
		Move:       move,
		Depth:      5,
		Age:        1,
		Type:       Vnone,
		MateThreat: false,
	}
	tt.numberOfEntries++

	e := tt.Probe(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.Key )
	assert.Equal(t, move, e.Move)
	assert.EqualValues(t, 5, e.Depth)
	assert.EqualValues(t, 0, e.Age)
	assert.Equal(t, Vnone, e.Type)
	assert.EqualValues(t, 1, tt.numberOfEntries)

	tt.Clear()

	// entry is gone
	e = tt.Probe(pos.ZobristKey())
	assert.Nil(t, e)
	assert.EqualValues(t, 0, tt.numberOfEntries)
}

func TestAge(t *testing.T) {
	// setup
	tt := New(5_000)

	logTest.Debug("Filling tt")
	startTime := time.Now()
	for i, _ := range tt.data {
		tt.data[i].Key = position.Key(i)
		tt.data[i].Age++
	}
	tt.data[0].Age = 0
	elapsed := time.Since(startTime)
	logTest.Debug(out.Sprintf("TT of %d elements filled in %d ms\n", len(tt.data), elapsed.Milliseconds()))

	// test
	assert.EqualValues(t, 0, tt.GetEntry(0).Age)
	assert.EqualValues(t, 1, tt.GetEntry(1).Age)
	assert.EqualValues(t, 1, tt.GetEntry(1_000).Age)
	assert.EqualValues(t, 1, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).Age)

	logTest.Debug("Ageing entries")
	tt.ageEntries()

	assert.EqualValues(t, 0, tt.GetEntry(0).Age)
	assert.EqualValues(t, 2, tt.GetEntry(1).Age)
	assert.EqualValues(t, 2, tt.GetEntry(1_000).Age)
	assert.EqualValues(t, 2, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).Age)
}


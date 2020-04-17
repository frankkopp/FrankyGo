/*
 * FrankyGo - UCI chess engine in GO for learning purposes
 *
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
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
	"unsafe"

	logging2 "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

var logTest *logging2.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	config.Setup()
	logTest = logging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func TestEntrySize(t *testing.T) {
	e := TtEntry{
		Key:        0,
		Move:       0,
		Depth:      0,
		Age:        0,
		Type:       0,
		MateThreat: false,
	}
	assert.EqualValues(t, 16, unsafe.Sizeof(e))
	logTest.Debugf("Size of Entry %d bytes", unsafe.Sizeof(e))
}

func TestNew(t *testing.T) {

	tt := NewTtTable(2)
	assert.Equal(t, uint64(131_072), tt.maxNumberOfEntries)
	assert.Equal(t, 131_072, cap(tt.data))
	logTest.Debug(tt.String())

	tt = NewTtTable(64)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = NewTtTable(100)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = NewTtTable(4_096)
	assert.Equal(t, uint64(268_435_456), tt.maxNumberOfEntries)
	assert.Equal(t, 268_435_456, cap(tt.data))

	// Too much for Travis
	// tt = NewTtTable(35_000)
	// assert.Equal(t, uint64(2_147_483_648), tt.maxNumberOfEntries)
	// assert.Equal(t, 2_147_483_648, cap(tt.data))
	// assert.Equal(t, 2_147_483_648, len(tt.data))
	// assert.Equal(t, 32_768*MB, tt.sizeInByte)
	// for i, _ := range tt.data {
	// 	tt.data[i].Key = position.Key(i)
	// }
	// assert.Equal(t, position.Key(0), tt.data[0].Key)
	// assert.Equal(t, position.Key(2_147_483_647), tt.data[2_147_483_647].Key)
}

func TestGetAndProbe(t *testing.T) {
	// setup

	tt := NewTtTable(64)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	pos := position.NewPosition()
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
	assert.Equal(t, pos.ZobristKey(), e.Key)
	assert.Equal(t, move, e.Move)
	assert.EqualValues(t, 5, e.Depth)
	assert.EqualValues(t, 1, e.Age)
	assert.Equal(t, Vnone, e.Type)

	// age must be reduced by 1
	e = tt.Probe(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.Key)
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
	tt := NewTtTable(1)

	pos := position.NewPosition()
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
	assert.Equal(t, pos.ZobristKey(), e.Key)
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
	tt := NewTtTable(5_000)

	logTest.Debug("Filling tt")
	startTime := time.Now()
	for i, _ := range tt.data {
		tt.numberOfEntries++
		tt.data[i].Key = position.Key(i)
		tt.data[i].Age++
	}
	tt.data[0].Age = 0
	tt.numberOfEntries--
	elapsed := time.Since(startTime)
	logTest.Debug(out.Sprintf("TT of %d elements filled in %d ms\n", len(tt.data), elapsed.Milliseconds()))
	logTest.Debug(tt.String())

	// test
	assert.EqualValues(t, 0, tt.GetEntry(0).Age)
	assert.EqualValues(t, 1, tt.GetEntry(1).Age)
	assert.EqualValues(t, 1, tt.GetEntry(1_000).Age)
	assert.EqualValues(t, 1, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).Age)

	logTest.Debug("Aging entries")
	tt.AgeEntries()

	assert.EqualValues(t, 0, tt.GetEntry(0).Age)
	assert.EqualValues(t, 2, tt.GetEntry(1).Age)
	assert.EqualValues(t, 2, tt.GetEntry(1_000).Age)
	assert.EqualValues(t, 2, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).Age)
}

func TestPut(t *testing.T) {
	// setup

	tt := NewTtTable(4)
	move := CreateMove(SqE2, SqE4, Normal, PtNone)

	// test of put and probe
	tt.Put(111, move, 4, Value(111), ALPHA, false)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 1, tt.Stats.numberOfPuts)
	e := tt.Probe(111)
	assert.EqualValues(t, 111, e.Key)
	assert.EqualValues(t, move, e.Move.MoveOf())
	assert.EqualValues(t, 111, e.Move.ValueOf())
	assert.EqualValues(t, 4, e.Depth)
	assert.EqualValues(t, ALPHA, e.Type)
	assert.EqualValues(t, 0, e.Age)
	assert.EqualValues(t, false, e.MateThreat)

	// test of put update and probe
	tt.Put(111, move, 5, Value(112), BETA, true)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 2, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 0, tt.Stats.numberOfCollisions)
	e = tt.Probe(111)
	assert.EqualValues(t, 111, e.Key)
	assert.EqualValues(t, move, e.Move.MoveOf())
	assert.EqualValues(t, 112, e.Move.ValueOf())
	assert.EqualValues(t, 5, e.Depth)
	assert.EqualValues(t, BETA, e.Type)
	assert.EqualValues(t, 0, e.Age)
	assert.EqualValues(t, true, e.MateThreat)

	// test of collision
	collisionKey := position.Key(111 + tt.maxNumberOfEntries)
	tt.Put(collisionKey, move, 6, Value(113), EXACT, false)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 3, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 1, tt.Stats.numberOfCollisions)
	assert.EqualValues(t, 1, tt.Stats.numberOfOverwrites)
	e = tt.Probe(collisionKey)
	assert.EqualValues(t, collisionKey, e.Key)
	assert.EqualValues(t, move, e.Move.MoveOf())
	assert.EqualValues(t, 113, e.Move.ValueOf())
	assert.EqualValues(t, 6, e.Depth)
	assert.EqualValues(t, EXACT, e.Type)
	assert.EqualValues(t, 0, e.Age)
	assert.EqualValues(t, false, e.MateThreat)

	// test of collision lower depth
	collisionKey2 := position.Key(111 + (tt.maxNumberOfEntries << 1))
	tt.Put(collisionKey2, move, 4, Value(114), BETA, true)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 4, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 2, tt.Stats.numberOfCollisions)
	assert.EqualValues(t, 1, tt.Stats.numberOfOverwrites)
	e = tt.Probe(collisionKey2)
	assert.Nil(t, e)
	e = tt.Probe(collisionKey)
	assert.EqualValues(t, collisionKey, e.Key)
	assert.EqualValues(t, move, e.Move.MoveOf())
	assert.EqualValues(t, 113, e.Move.ValueOf())
	assert.EqualValues(t, 6, e.Depth)
	assert.EqualValues(t, EXACT, e.Type)
	assert.EqualValues(t, 0, e.Age)
	assert.EqualValues(t, false, e.MateThreat)
}

func TestTimingTTe(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup
	tt := NewTtTable(1_024)
	move := CreateMove(SqE2, SqE4, Normal, PtNone)

	const rounds = 5
	const iterations uint64 = 50_000_000

	for r := 1; r <= rounds; r++ {
		out.Printf("Round %d\n", r)
		key := position.Key(rand.Uint64())
		depth := int8(rand.Int31n(128))
		value := Value(rand.Int31n(int32(ValueMax)))
		valueType := ValueType(rand.Int31n(4))
		start := time.Now()
		for i := uint64(0); i < iterations; i++ {
			tt.Put(key+position.Key(i), move, depth, value, valueType, false)
		}
		for i := uint64(0); i < iterations; i++ {
			key := position.Key(key + position.Key(2*i))
			_ = tt.Probe(key)
		}
		elapsed := time.Since(start)
		out.Println(tt.String())
		out.Printf("TimingTT took %d ns for %d iterations (1 put 1 probe)\n", elapsed.Nanoseconds(), iterations)
		out.Printf("1 put/probes in %d ns: %d tts\n",
			elapsed.Nanoseconds()/int64(iterations),
			(iterations*uint64(time.Second.Nanoseconds()))/uint64(elapsed.Nanoseconds()))

	}
}

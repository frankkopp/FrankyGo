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
		key:   0,
		move:  0,
		value: 0,
		eval:  0,
		depth: 0,
		age:   0,
		vtype: 0,
	}
	logTest.Debugf("Size of e.Key = %d bytes", unsafe.Sizeof(e.key))
	logTest.Debugf("Size of e.Move = %d bytes", unsafe.Sizeof(e.move))
	logTest.Debugf("Size of e.Value = %d bytes", unsafe.Sizeof(e.value))
	logTest.Debugf("Size of e.Eval = %d bytes", unsafe.Sizeof(e.eval))
	logTest.Debugf("Size of e.Depth = %d bytes", unsafe.Sizeof(e.depth))
	logTest.Debugf("Size of e.Age = %d bytes", unsafe.Sizeof(e.age))
	logTest.Debugf("Size of e.Type = %d bytes", unsafe.Sizeof(e.vtype))
	assert.EqualValues(t, 24, unsafe.Sizeof(e))
	logTest.Debugf("Size of Entry %d bytes", unsafe.Sizeof(e))
}

func TestNew(t *testing.T) {

	tt := NewTtTable(2)
	assert.Equal(t, uint64(0x10000), tt.maxNumberOfEntries)
	assert.Equal(t, 0x10000, cap(tt.data))
	logTest.Debug(tt.String())

	tt = NewTtTable(64)
	assert.Equal(t, uint64(0x200000), tt.maxNumberOfEntries)
	assert.Equal(t, 0x200000, cap(tt.data))

	tt = NewTtTable(100)
	assert.Equal(t, uint64(0x400000), tt.maxNumberOfEntries)
	assert.Equal(t, 0x400000, cap(tt.data))

	tt = NewTtTable(4_096)
	assert.Equal(t, uint64(0x8000000), tt.maxNumberOfEntries)
	assert.Equal(t, 0x8000000, cap(tt.data))

	// Too much for Travis
	// tt = NewTtTable(35_000)
	// assert.Equal(t, uint64(0x40000000), tt.maxNumberOfEntries)
	// assert.Equal(t, 0x40000000, cap(tt.data))
	// assert.Equal(t, 0x40000000, len(tt.data))
	// assert.Equal(t, uint64(0x600000000), tt.sizeInByte)
	// for i := range tt.data {
	// 	tt.data[i].Key = position.Key(i)
	// }
	// assert.Equal(t, position.Key(0), tt.data[0].Key)
	// assert.Equal(t, position.Key(1_073_741_823), tt.data[1_073_741_823].Key)
}

func TestGetAndProbe(t *testing.T) {
	// setup

	tt := NewTtTable(64)
	assert.Equal(t, uint64(0x200000), tt.maxNumberOfEntries)
	assert.Equal(t, 0x200000, cap(tt.data))

	pos := position.NewPosition()
	move := CreateMove(SqE2, SqE4, Normal, PtNone)
	tt.data[tt.hash(pos.ZobristKey())] = TtEntry{
		key:   pos.ZobristKey(),
		move:  uint16(move),
		value: int16(ValueNA),
		eval:  int16(ValueNA),
		depth: 5,
		age:   1,
		vtype: Vnone,
	}
	tt.numberOfEntries++

	// test to get unaltered entry
	e := tt.GetEntry(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.key)
	assert.Equal(t, move, e.Move())
	assert.EqualValues(t, 5, e.depth)
	assert.EqualValues(t, 1, e.age)
	assert.Equal(t, Vnone, e.vtype)

	// age must be reduced by 1
	e = tt.Probe(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.key)
	assert.Equal(t, move, e.Move())
	assert.EqualValues(t, 5, e.depth)
	assert.EqualValues(t, 0, e.age)
	assert.Equal(t, Vnone, e.vtype)

	// age does not go below 0
	e = tt.Probe(pos.ZobristKey())
	assert.EqualValues(t, 0, e.age)

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
		key:   pos.ZobristKey(),
		move:  uint16(move),
		value: int16(ValueNA),
		eval:  int16(ValueNA),
		depth: 5,
		age:   1,
		vtype: Vnone,
	}
	tt.numberOfEntries++

	e := tt.Probe(pos.ZobristKey())
	assert.Equal(t, pos.ZobristKey(), e.key)
	assert.Equal(t, move, e.Move())
	assert.EqualValues(t, 5, e.depth)
	assert.EqualValues(t, 0, e.age)
	assert.Equal(t, Vnone, e.vtype)
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
	for i := range tt.data {
		tt.numberOfEntries++
		tt.data[i].key = position.Key(i)
		tt.data[i].age++
	}
	tt.data[0].age = 0
	tt.numberOfEntries--
	elapsed := time.Since(startTime)
	logTest.Debug(out.Sprintf("TT of %d elements filled in %d ms\n", len(tt.data), elapsed.Milliseconds()))
	logTest.Debug(tt.String())

	// test
	assert.EqualValues(t, 0, tt.GetEntry(0).age)
	assert.EqualValues(t, 1, tt.GetEntry(1).age)
	assert.EqualValues(t, 1, tt.GetEntry(1_000).age)
	assert.EqualValues(t, 1, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).age)

	logTest.Debug("Aging entries")
	tt.AgeEntries()

	assert.EqualValues(t, 0, tt.GetEntry(0).age)
	assert.EqualValues(t, 2, tt.GetEntry(1).age)
	assert.EqualValues(t, 2, tt.GetEntry(1_000).age)
	assert.EqualValues(t, 2, tt.GetEntry(position.Key(tt.maxNumberOfEntries-1)).age)
}

func TestPut(t *testing.T) {
	// setup

	tt := NewTtTable(4)
	move := CreateMove(SqE2, SqE4, Normal, PtNone)

	// test of put and probe
	tt.Put(111, move, 4, Value(111), ALPHA, ValueNA)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 1, tt.Stats.numberOfPuts)
	e := tt.Probe(111)
	assert.EqualValues(t, 111, e.key)
	assert.EqualValues(t, move, e.Move())
	assert.EqualValues(t, 111, e.value)
	assert.EqualValues(t, 4, e.depth)
	assert.EqualValues(t, ALPHA, e.vtype)
	assert.EqualValues(t, 0, e.age)

	// test of put update and probe
	tt.Put(111, move, 5, Value(112), BETA, ValueNA)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 2, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 0, tt.Stats.numberOfCollisions)
	e = tt.Probe(111)
	assert.EqualValues(t, 111, e.key)
	assert.EqualValues(t, move, e.Move())
	assert.EqualValues(t, 112, e.value)
	assert.EqualValues(t, 5, e.depth)
	assert.EqualValues(t, BETA, e.vtype)
	assert.EqualValues(t, 0, e.age)

	// test of collision
	collisionKey := position.Key(111 + tt.maxNumberOfEntries)
	tt.Put(collisionKey, move, 6, Value(113), EXACT, ValueNA)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 3, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 1, tt.Stats.numberOfCollisions)
	assert.EqualValues(t, 1, tt.Stats.numberOfOverwrites)
	e = tt.Probe(collisionKey)
	assert.EqualValues(t, collisionKey, e.key)
	assert.EqualValues(t, move, e.Move())
	assert.EqualValues(t, 113, e.value)
	assert.EqualValues(t, 6, e.depth)
	assert.EqualValues(t, EXACT, e.vtype)
	assert.EqualValues(t, 0, e.age)

	// test of collision lower depth
	collisionKey2 := position.Key(111 + (tt.maxNumberOfEntries << 1))
	tt.Put(collisionKey2, move, 4, Value(114), BETA, ValueNA)
	assert.EqualValues(t, 1, tt.Len())
	assert.EqualValues(t, 4, tt.Stats.numberOfPuts)
	assert.EqualValues(t, 1, tt.Stats.numberOfUpdates)
	assert.EqualValues(t, 2, tt.Stats.numberOfCollisions)
	assert.EqualValues(t, 1, tt.Stats.numberOfOverwrites)
	e = tt.Probe(collisionKey2)
	assert.Nil(t, e)
	e = tt.Probe(collisionKey)
	assert.EqualValues(t, collisionKey, e.key)
	assert.EqualValues(t, move, e.Move())
	assert.EqualValues(t, 113, e.value)
	assert.EqualValues(t, 6, e.depth)
	assert.EqualValues(t, EXACT, e.vtype)
	assert.EqualValues(t, 0, e.age)
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
			tt.Put(key+position.Key(i), move, depth, value, valueType, ValueNA)
		}
		for i := uint64(0); i < iterations; i++ {
			key := position.Key(key + position.Key(2*i))
			_ = tt.Probe(key)
		}
		elapsed := time.Since(start)
		out.Println(tt.String())
		out.Printf("TimingTT took %s for %d iterations (1 put 1 probe)\n", elapsed, iterations)
		out.Printf("1 put/probes in %d ns: %d tts\n",
			elapsed.Nanoseconds()/int64(iterations),
			(iterations*uint64(time.Second.Nanoseconds()))/uint64(elapsed.Nanoseconds()))

	}
}

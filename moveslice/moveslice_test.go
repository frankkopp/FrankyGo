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

package moveslice

import (
	"math/rand"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/config"
	myLogging "github.com/frankkopp/FrankyGo/logging"
	. "github.com/frankkopp/FrankyGo/types"
)

var logTest *logging.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	config.Setup()
	logTest = myLogging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

var (
	e2e4 = CreateMoveValue(SqE2, SqE4, Normal, PtNone, 111)
	d7d5 = CreateMoveValue(SqD7, SqD5, Normal, PtNone, 222)
	e4d5 = CreateMoveValue(SqE4, SqD5, Normal, PtNone, 333)
	d8d5 = CreateMoveValue(SqD8, SqD5, Normal, PtNone, 444)
	b1c3 = CreateMoveValue(SqB1, SqC3, Normal, PtNone, 555)
)

func TestNew(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 0, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

}

func TestMoveArrayPushBack(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }

	for i := 0; i < 1_000_000; i++ {
		ma.PushBack(e2e4)
	}
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 1_000_005, len(*ma))
	assert.Equal(t, 1_163_264, cap(*ma))
}

func TestMoveArrayPopBack(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	assert.Panics(t, func(){ ma.PopBack() })

	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

	m1 := ma.PopBack()
	assert.Equal(t, b1c3, m1)
	m2 := ma.PopBack()
	assert.Equal(t, d8d5, m2)
	assert.Equal(t, 3, len(*ma))

	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }
}


func TestMoveArrayPushFront(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushFront(e2e4)
	ma.PushFront(d7d5)
	ma.PushFront(e4d5)
	ma.PushFront(d8d5)
	ma.PushFront(b1c3)

	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }
}

func TestMoveArrayPopFront(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	assert.Panics(t, func(){ ma.PopFront() })
	ma.PushFront(e2e4)
	ma.PushFront(d7d5)
	ma.PushFront(e4d5)
	ma.PushFront(d8d5)
	ma.PushFront(b1c3)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

	m1 := ma.PopFront()
	assert.Equal(t, b1c3, m1)
	m2 := ma.PopFront()
	assert.Equal(t, d8d5, m2)
	assert.Equal(t, 3, len(*ma))

	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }
}

func TestMoveArrayClear(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	ma.Clear()
	assert.Equal(t, 0, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
}

func TestMoveArrayAccess(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))

	assert.Equal(t, e2e4, ma.Front())
	assert.Equal(t, ma.At(0), ma.Front())
	assert.Equal(t, b1c3, ma.Back())
	assert.Equal(t, ma.At(len(*ma)-1), ma.Back())
	ma.Set(0, b1c3)
	assert.Equal(t, b1c3, ma.Front())
	assert.Equal(t, ma.At(0), ma.Front())
}

func TestMoveArrayString(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	// logTest.Debugf("String() = %s", ma.String())
	// logTest.Debugf("StringUci() = %s", ma.StringUci())
	assert.Equal(t, "e2e4 d7d5 e4d5 d8d5 b1c3", ma.StringUci())
}

func TestMoveArraySort(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	// logTest.Debugf("Len=%d", len(*ma))
	// logTest.Debugf("Cap=%d", cap(*ma))
	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	// logTest.Debugf("String() = %s", ma.String())
	// logTest.Debugf("StringUci() = %s", ma.StringUci())
	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }
	// fmt.Println("Sorted:")
	// ma.Sort()
	// logTest.Debugf("String() = %s", ma.String())
	// logTest.Debugf("StringUci() = %s", ma.StringUci())
	// for i, v := range *ma {
	// 	fmt.Println(i, v)
	// }
}


func TestMoveArraySortRandom(t *testing.T) {

	ma := NewMoveSlice(MaxMoves)
	items := 10_000

	// generate random moves
	for i := 0; i < items; i++ {
		ma.PushBack(Move(rand.Int31()))
	}

	// sort
	// start := time.Now()
	ma.Sort()
	// elapsed := time.Since(start)
	// out.Printf("%d ns\n", elapsed.Nanoseconds())

	// check
	tmp := ma.At(0)
	for i := 0; i < items; i++ {
		assert.True(t, tmp >= ma.At(i))
		tmp = ma.At(i)
	}

}


func TestMoveArrayFilter(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	assert.Equal(t, "e2e4 d7d5 e4d5 d8d5 b1c3", ma.StringUci())

	ma.Filter(func(i int) bool {
		return ma.At(i) != e4d5
	})

	assert.Equal(t, 4, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	assert.Equal(t, "e2e4 d7d5 d8d5 b1c3", ma.StringUci())
}

func TestMoveArray_FilterCopy(t *testing.T) {
	ma := NewMoveSlice(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	assert.Equal(t, "e2e4 d7d5 e4d5 d8d5 b1c3", ma.StringUci())

	ma2 := NewMoveSlice(cap(*ma))
	ma.FilterCopy(ma2, func(i int) bool {
		return ma.At(i) != e4d5
	})

	assert.Equal(t, 5, len(*ma))
	assert.Equal(t, MaxMoves, cap(*ma))
	assert.Equal(t, "e2e4 d7d5 e4d5 d8d5 b1c3", ma.StringUci())

	assert.Equal(t, 4, len(*ma2))
	assert.Equal(t, MaxMoves, cap(*ma2))
	assert.Equal(t, "e2e4 d7d5 d8d5 b1c3", ma2.StringUci())
}

func TestForEach(t *testing.T) {
	// fill array
	noOfItems := 1_000
	ma := NewMoveSlice(noOfItems)
	for i := 0; i < noOfItems; i++ {
		ma.PushBack(e2e4)
	}

	// counter and mutex
	var mux sync.Mutex
	var counter int

	// parallel execution
	ma.ForEachParallel(func(i int) {
		m := ma.At(i)
		f := m.From()
		t := m.To()
		mt := m.MoveType()
		pt := m.PromotionType()
		v := Value(999)
		ma.Set(i, CreateMoveValue(f, t, mt, pt, v))
		// simulate cpu intense calculation
		n := float64(100000)
		for n > 1 {
			n /= 1.000001
		}
		mux.Lock()
		counter++
		mux.Unlock()
	})

	assert.Equal(t, noOfItems, counter)
	assert.Equal(t, Value(999), ma.Front().ValueOf())
	assert.Equal(t, Value(999), ma.At(10).ValueOf())
	assert.Equal(t, Value(999), ma.At(100).ValueOf())
	assert.Equal(t, Value(999), ma.Back().ValueOf())

	// fmt.Printf("counter %d\n", counter)
	// ma.ForEach(func(i int) {
	// 	fmt.Printf("%d: %s\n", i, ma.At(i).String())
	// })
}


func TestGoLandWithVeryLongName(t *testing.T) {
	// fill array
	noOfItems := 1_000
	ma := NewMoveSlice(noOfItems)
	for i := 0; i < noOfItems; i++ {
		ma.PushBack(e2e4)
	}

	// fmt.Printf("counter %d\n", noOfItems)
	// ma.ForEach(func (i int) {
	// 	fmt.Printf("%d: %s\n", i, ma.At(i).String())
	// })
}


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

package movearray

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/frankkopp/FrankyGo/types"
)

var (
	e2e4 = CreateMoveValue(SqE2, SqE4, Normal, PtNone, 111)
	d7d5 = CreateMoveValue(SqD7, SqD5, Normal, PtNone, 222)
	e4d5 = CreateMoveValue(SqE4, SqD5, Normal, PtNone, 333)
	d8d5 = CreateMoveValue(SqD8, SqD5, Normal, PtNone, 444)
	b1c3 = CreateMoveValue(SqB1, SqC3, Normal, PtNone, 555)
)

func TestNew(t *testing.T) {
	ma := New(MaxMoves)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 0, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())
}

func TestMoveArray_PushBack(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())

	for i, v := range ma.data {
		fmt.Println(i, v)
	}

	for i := 0; i < 1_000_000; i++ {
		ma.PushBack(e2e4)
	}
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 1_000_005, ma.Len())
	assert.Equal(t, 1_163_264, ma.Cap())
}

func TestMoveArray_PopBack(t *testing.T) {
	ma := New(MaxMoves)
	assert.Panics(t, func(){ ma.PopBack() })

	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)

	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())

	m1 := ma.PopBack()
	assert.Equal(t, b1c3, m1)
	m2 := ma.PopBack()
	assert.Equal(t, d8d5, m2)
	assert.Equal(t, 3, ma.Len())

	for i, v := range ma.data {
		fmt.Println(i, v)
	}
}


func TestMoveArray_PushFront(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushFront(e2e4)
	ma.PushFront(d7d5)
	ma.PushFront(e4d5)
	ma.PushFront(d8d5)
	ma.PushFront(b1c3)

	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())

	for i, v := range ma.data {
		fmt.Println(i, v)
	}
}

func TestMoveArray_PopFront(t *testing.T) {
	ma := New(MaxMoves)
	assert.Panics(t, func(){ ma.PopFront() })
	ma.PushFront(e2e4)
	ma.PushFront(d7d5)
	ma.PushFront(e4d5)
	ma.PushFront(d8d5)
	ma.PushFront(b1c3)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())

	m1 := ma.PopFront()
	assert.Equal(t, b1c3, m1)
	m2 := ma.PopFront()
	assert.Equal(t, d8d5, m2)
	assert.Equal(t, 3, ma.Len())

	for i, v := range ma.data {
		fmt.Println(i, v)
	}
}

func TestMoveArray_Clear(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())
	ma.Clear()
	assert.Equal(t, 0, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())
}

func TestMoveArray_Access(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())

	assert.Equal(t, e2e4, ma.Front())
	assert.Equal(t, ma.At(0), ma.Front())
	assert.Equal(t, b1c3, ma.Back())
	assert.Equal(t, ma.At(ma.Len()-1), ma.Back())
	ma.Set(0, b1c3)
	assert.Equal(t, b1c3, ma.Front())
	assert.Equal(t, ma.At(0), ma.Front())
}

func TestMoveArray_String(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())
	log.Printf("String() = %s", ma.String())
	log.Printf("StringUci() = %s", ma.StringUci())
	assert.Equal(t, "e2e4 d7d5 e4d5 d8d5 b1c3", ma.StringUci())
}

func TestMoveArray_Sort(t *testing.T) {
	ma := New(MaxMoves)
	ma.PushBack(e2e4)
	ma.PushBack(d7d5)
	ma.PushBack(e4d5)
	ma.PushBack(d8d5)
	ma.PushBack(b1c3)
	log.Printf("Len=%d", ma.Len())
	log.Printf("Cap=%d", ma.Cap())
	assert.Equal(t, 5, ma.Len())
	assert.Equal(t, MaxMoves, ma.Cap())
	log.Printf("String() = %s", ma.String())
	log.Printf("StringUci() = %s", ma.StringUci())
	for i, v := range ma.data {
		fmt.Println(i, v)
	}
	fmt.Println("Sorted:")
	ma.Sort()
	log.Printf("String() = %s", ma.String())
	log.Printf("StringUci() = %s", ma.StringUci())
	for i, v := range ma.data {
		fmt.Println(i, v)
	}
}

func TestMoveArray_ForEach(t *testing.T) {
	// fill array
	noOfItems := 1_000
	ma := New(noOfItems)
	for i := 0; i < noOfItems; i++ {
		ma.PushBack(e2e4)
	}

	// counter and mutex
	var mux sync.Mutex
	var counter int

	// parallel execution
	ma.ForEachParallel(func(i int){
		m := ma.At(i)
		f := m.From()
		t := m.To()
		mt := m.MoveType()
		pt := m.PromotionType()
		v:= Value(999)
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

	fmt.Printf("Counter %d\n", counter)
	assert.Equal(t, noOfItems, counter)
	assert.Equal(t, Value(999), ma.Front().ValueOf())
	assert.Equal(t, Value(999), ma.At(10).ValueOf())
	assert.Equal(t, Value(999), ma.At(100).ValueOf())
	assert.Equal(t, Value(999), ma.Back().ValueOf())

	// ma.ForEach(func (i int) {
	// 	fmt.Printf("%d: %s\n", i, ma.At(i).String())
	// })
}

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

// Package moveslice provides an array (slice) facade to be used with
// chess moves.
package moveslice

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/frankkopp/FrankyGo/types"
)

// MoveSlice represents a data structure (go slice) for Move.
type MoveSlice []Move


// New creates a new move slice with the given capacity and 0 elements
// Is identical to MoveSlice(make([]Move, 0, cap))
func New(cap int) MoveSlice {
	return make([]Move, 0, cap)
}

// Len returns the number of moves currently stored in the slice
func (ms *MoveSlice) Len() int {
	return len(*ms)
}

// Cap returns the capacity of the slice
func (ms *MoveSlice) Cap() int {
	return cap(*ms)
}

// PushBack appends an element at the end of the slice
func (ms *MoveSlice) PushBack(m Move) {
	*ms = append(*ms, m)
}

// PopBack removes and returns the move from the back of the slice.
// If the slice is empty, the call panics.
func (ms *MoveSlice) PopBack() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: PopBack() called on empty slice")
	}
	backMove := (*ms)[len(*ms)-1]
	*ms = (*ms)[:len(*ms)-1]
	return backMove
}

// PushFront prepends an element at the beginning of the slice using
// the underlying array (does not create a new array)
func (ms *MoveSlice) PushFront(m Move) {
	*ms = append(*ms, MoveNone)
	copy((*ms)[1:], *ms)
	(*ms)[0] = m
}

// PopFront removes and returns the move from the front of the slice.
// If the slice is empty, the call panics.
// Shrinks the capacity of the slice and might lead to earlier
// re-allocations
func (ms *MoveSlice) PopFront() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: PopFront() called on empty slice")
	}
		frontMove := (*ms)[0]
	*ms = (*ms)[1:]
	return frontMove
}

// Front returns the move at the front of the slice.  This is the element
// that would be returned by PopFront(). This call panics if the slice is
// empty.
func (ms *MoveSlice) Front() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: Front() called when empty")
	}
	return (*ms)[0]
}

// Back returns the move at the back of the slice.  This is the element
// that would be returned by PopBack().  This call panics if the slice is
// empty.
func (ms *MoveSlice) Back() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: Back() called when empty")
	}
	return (*ms)[len(*ms)-1]
}

// At returns the move at index i in the slice without removing the move
// from the slice. At(0) refers to the first move and is the same as Front().
// At(Len()-1) refers to the last move and is the same as Back().
// Index will not be checked against bounds.
func (ms *MoveSlice) At(i int) Move {
	return (*ms)[i]
}

// Set puts a move at index i in the slice. Set shares the same purpose
// than At() but perform the opposite operation. The index i is the same
// index defined by At().
// Index will not be checked against bounds.
func (ms *MoveSlice) Set(i int, move Move) {
	(*ms)[i] = move
}

// Filter removes all elements from the MoveSlice for
// which the given call to func will return false.
// Rebuilds the data slice by looping over all elements
// and only re-adding elements for which the call to the
// given func is true. Reuses the underlying array
func (ms *MoveSlice) Filter(f func(index int) bool) {
	b := (*ms)[:0]
	for i, x := range *ms {
		if f(i) {
			b = append(b, x)
		}
	}
	*ms = b
}

// FilterCopy copies the MoveSlice into the given destination slice
// without the filtered elements. An element is filtered when
// the given call to func will return false for the element.
func (ms *MoveSlice) FilterCopy(dest *MoveSlice, f func(index int) bool) {
	for i, x := range *ms {
		if f(i) {
			*dest = append(*dest, x)
		}
	}
}

// ForEach simple range loop calling the given function on each element
// in stored order
func (ms *MoveSlice) ForEach(f func(index int)) {
	for index, _ := range *ms {
		f(index)
	}
}

// ForEachParallel simple loop over all elements calling a goroutine
// which calls the given func with the index of the current element
// as a parameter.
// Waits until all elements have been processed. There is not
// synchronization for the parallel execution. This needs to done
// in the provided function
func (ms *MoveSlice) ForEachParallel(f func(index int)) {
	sliceLength := len(*ms)
	var wg sync.WaitGroup
	wg.Add(sliceLength)
	for index, _ := range *ms {
		go func(i int) {
			defer wg.Done()
			f(i)
		}(index)
	}
	wg.Wait()
}

// Clear removes all moves from the slice, but retains the current capacity.
// This is useful when repeatedly reusing the slice at high frequency to avoid
// GC during reuse.
func (ms *MoveSlice) Clear() {
	*ms = (*ms)[:0]
}

// Sort sorts the moves from highest value to lowest value
// Uses InsertionSort as MoveSlices are mostly pre-sorted and small
func (ms *MoveSlice) Sort() {
	l := len(*ms)
	for i := 1; i < l; i++ {
		tmp := (*ms)[i]
		j := i
		for j > 0 && tmp > (*ms)[j-1] {
			(*ms)[j] = (*ms)[j-1]
			j--
		}
		(*ms)[j] = tmp
	}
}

// String returns a string representation of a slice of moves
func (ms *MoveSlice) String() string {
	var os strings.Builder
	size := len(*ms)
	os.WriteString(fmt.Sprintf("MoveList: [%d] { ", size))
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(", ")
		}
		m := ms.At(i)
		os.WriteString(m.String())
	}
	os.WriteString(" }")
	return os.String()
}

// StringUci returns a string with a space separated list
// of all moves i the list in UCI protocol format
func (ms *MoveSlice) StringUci() string {
	var os strings.Builder
	size := len(*ms)
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := ms.At(i)
		os.WriteString(m.StringUci())
	}
	return os.String()
}

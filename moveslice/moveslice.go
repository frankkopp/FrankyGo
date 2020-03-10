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

// Package moveslice provides a array (slice) facade to be used with
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


// New creates a new move array with the given capacity and 0 elements
// Is identical to MoveSlice(make([]Move, 0, cap))
func New(cap int) MoveSlice {
	return make([]Move, 0, cap)
}

// PushBack appends an element at the end of the array
func (ma *MoveSlice) PushBack(m Move) {
	*ma = append(*ma, m)
}

// PopBack removes and returns the move from the back of the queue.
// If the queue is empty, the call panics.
func (ma *MoveSlice) PopBack() Move {
	if len(*ma) <= 0 {
		panic("MoveSlice: PopBack() called on empty array")
	}
	backMove := (*ma)[len(*ma)-1]
	*ma = (*ma)[:len(*ma)-1]
	return backMove
}

// PushFront prepends an element at the beginning of the array using
// the underlying array (does not create a new one)
func (ma *MoveSlice) PushFront(m Move) {
	*ma = append(*ma, MoveNone)
	copy((*ma)[1:], *ma)
	(*ma)[0] = m
}

// PopFront removes and returns the move from the front of the array.
// If the array is empty, the call panics.
// Shrinks the capacity of the array and might lead to earlier
// re-allocations
func (ma *MoveSlice) PopFront() Move {
	if len(*ma) <= 0 {
		panic("MoveSlice: PopFront() called on empty array")
	}
		frontMove := (*ma)[0]
	*ma = (*ma)[1:]
	return frontMove
}

// Front returns the move at the front of the array.  This is the element
// that would be returned by PopFront(). This call panics if the array is
// empty.
func (ma *MoveSlice) Front() Move {
	if len(*ma) <= 0 {
		panic("MoveSlice: Front() called when empty")
	}
	return (*ma)[0]
}

// Back returns the move at the back of the array.  This is the element
// that would be returned by PopBack().  This call panics if the array is
// empty.
func (ma *MoveSlice) Back() Move {
	if len(*ma) <= 0 {
		panic("MoveSlice: Back() called when empty")
	}
	return (*ma)[len(*ma)-1]
}

// At returns the move at index i in the array without removing the move
// from the array. At(0) refers to the first move and is the same as Front().
// At(Len()-1) refers to the last move and is the same as Back().
// Index will not be checked against bounds.
func (ma *MoveSlice) At(i int) Move {
	return (*ma)[i]
}

// Set puts a move at index i in the queue. Set shares the same purpose
// than At() but perform the opposite operation. The index i is the same
// index defined by At().
// Index will not be checked against bounds.
func (ma *MoveSlice) Set(i int, move Move) {
	(*ma)[i] = move
}

// Filter removes all elements from the MoveSlice for
// which the given call to func will return false.
// Rebuilds the data slice by looping over all elements
// and only re-adding elements for which the call to the
// given func is true. Reuses the underlying array
func (ma *MoveSlice) Filter(f func(index int) bool) {
	b := (*ma)[:0]
	for i, x := range *ma {
		if f(i) {
			b = append(b, x)
		}
	}
	*ma = b
}

// FilterCopy copies the MoveSlice into the given destination array
// without the filtered elements. An element is filtered when
// the given call to func will return false for the element.
func (ma *MoveSlice) FilterCopy(dest *MoveSlice, f func(index int) bool) {
	for i, x := range *ma {
		if f(i) {
			*dest = append(*dest, x)
		}
	}
}

// ForEach simple range loop calling the given function on each element
// in stored order
func (ma *MoveSlice) ForEach(f func(index int)) {
	for index, _ := range *ma {
		f(index)
	}
}

// ForEachParallel simple loop over all elements calling a goroutine
// which calls the given func with the index of the current element
// as a parameter.
// Waits until all elements have been processed. There is not
// synchronization for the parallel execution. This needs to done
// in the provided function
func (ma *MoveSlice) ForEachParallel(f func(index int)) {
	sliceLength := len(*ma)
	var wg sync.WaitGroup
	wg.Add(sliceLength)
	for index, _ := range *ma {
		go func(i int) {
			defer wg.Done()
			f(i)
		}(index)
	}
	wg.Wait()
}

// Data allows access to the underlying slice which is good for range loops
// Use with care!
func (ma *MoveSlice) Data() []Move {
	return *ma
}

// Clear removes all moves from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.
func (ma *MoveSlice) Clear() {
	*ma = (*ma)[:0]
}

// Sort sorts the moves from highest value to lowest value
// Uses InsertionSort as MoveSlices are mostly pre-sorted and small
func (ma *MoveSlice) Sort() {
	l := len(*ma)
	for i := 1; i < l; i++ {
		tmp := (*ma)[i]
		j := i
		for j > 0 && tmp > (*ma)[j-1] {
			(*ma)[j] = (*ma)[j-1]
			j--
		}
		(*ma)[j] = tmp
	}
}

// String returns a string representation of a move list
func (ma *MoveSlice) String() string {
	var os strings.Builder
	size := len(*ma)
	os.WriteString(fmt.Sprintf("MoveList: [%d] { ", size))
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(", ")
		}
		m := ma.At(i)
		os.WriteString(m.String())
	}
	os.WriteString(" }")
	return os.String()
}

// StringUci returns a string with a space separated list
// of all moves i the list in UCI protocol format
func (ma *MoveSlice) StringUci() string {
	var os strings.Builder
	size := len(*ma)
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := ma.At(i)
		os.WriteString(m.StringUci())
	}
	return os.String()
}

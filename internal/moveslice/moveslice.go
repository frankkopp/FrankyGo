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

// Package moveslice provides helper functionality for slices
// of type Move (chess moves).
package moveslice

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/frankkopp/FrankyGo/internal/types"
)

// MoveSlice represents a data structure (go slice) for Move.
type MoveSlice []Move

// NewMoveSlice creates a new move slice with the given capacity
// and 0 elements.
// Is identical to MoveSlice(make([]Move, 0, cap))
func NewMoveSlice(cap int) *MoveSlice {
	moves := make([]Move, 0, cap)
	return (*MoveSlice)(&moves)
}

// Len returns the number of moves currently stored in the slice.
// Equivalent to len(ms)
func (ms *MoveSlice) Len() int {
	return len(*ms)
}

// Cap returns the capacity of the slice
// Equivalent to cap(ms)
func (ms *MoveSlice) Cap() int {
	return cap(*ms)
}

// PushBack appends an element at the end of the slice
// Equivalent to append(ms, m)
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
// Moves (copies) all elements by one index slot and adds the new move at
// the front.
func (ms *MoveSlice) PushFront(m Move) {
	*ms = append(*ms, MoveNone)
	copy((*ms)[1:], *ms)
	(*ms)[0] = m
}

// PopFront removes and returns the move from the front of the slice.
// If the slice is empty, the call panics.
// Shrinks the capacity of the slice as it only shifts the start of
// the slice within the underlying array. Might lead to earlier
// re-allocations
func (ms *MoveSlice) PopFront() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: PopFront() called on empty slice")
	}
	frontMove := (*ms)[0]
	*ms = (*ms)[1:]
	return frontMove
}

// Front returns the move at the front of the slice. This is the element
// that would be returned by ms[0].
// This call panics if the slice is empty.
func (ms *MoveSlice) Front() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: Front() called when empty")
	}
	return (*ms)[0]
}

// Back returns the move at the back of the slice. This is the element
// that would be returned by ms[len[ms)-1].
// This call panics if the slice is empty.
func (ms *MoveSlice) Back() Move {
	if len(*ms) <= 0 {
		panic("MoveSlice: Back() called when empty")
	}
	return (*ms)[len(*ms)-1]
}

// At returns the move at index i in the slice without removing the move
// from the slice. At(0) refers to the first move and is the same as Front().
// At(Len()-1) refers to the last move and is the same as Back().
// Index will be checked against bounds and panics if out of bounds
func (ms *MoveSlice) At(i int) Move {
	if len(*ms) == 0 || i < 0 || i >= len(*ms) {
		panic("MoveSlice: Index out of bounds")
	}
	return (*ms)[i]
}

// Set puts a move at index i in the slice. Set shares the same purpose
// than At() but performs the opposite operation. The index i is the same
// index defined by At().
// Index will be checked against bounds and panics if out of bounds
func (ms *MoveSlice) Set(i int, move Move) {
	if len(*ms) == 0 || i < 0 || i >= len(*ms) {
		panic("MoveSlice: Index out of bounds")
	}
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

// Clone copies the MoveSlice into a newly create MoveSlice
// doing a deep copy.
func (ms *MoveSlice) Clone() *MoveSlice {
	dest := make([]Move, ms.Len(), ms.Cap())
	copy(dest, *ms)
	return (*MoveSlice)(&dest)
}

// Equals returns true if all elements of the MoveSlice equals
// the elements of the other MoveSlice
func (ms *MoveSlice) Equals(other *MoveSlice) bool {
	if ms.Len() != other.Len() {
		return false
	}
	for i, m := range *ms {
		if m != (*other)[i] {
			return false
		}
	}
	return true
}

// ForEach simple range loop calling the given function on each element
// in stored order
func (ms *MoveSlice) ForEach(f func(index int)) {
	for index := range *ms {
		f(index)
	}
}

// ForEachParallel simple loop over all elements calling a goroutine
// which calls the given func with the index of the current element
// as a parameter.
// Waits until all elements have been processed. There is no
// synchronization for the parallel execution. This needs to done
// in the provided function if necessary
func (ms *MoveSlice) ForEachParallel(f func(index int)) {
	sliceLength := len(*ms)
	var wg sync.WaitGroup
	wg.Add(sliceLength)
	for index := range *ms {
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
	// *ms = nil
	*ms = (*ms)[:0]
}

// Sort will sort moves from highest Value to lowest Value .
// It uses a stable InsertionSort as MoveSlices are mostly pre-sorted and small.
// Sorts move only on the base of their Value (Move value == move&0xFFFF0000).
// Otherwise the order will not be changed.
func (ms *MoveSlice) Sort() {
	l := len(*ms)
	for i := 1; i < l; i++ {
		tmp := (*ms)[i]
		j := i
		for j > 0 && (tmp&0xFFFF0000) > ((*ms)[j-1]&0xFFFF0000) {
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
// of all moves in the list in UCI protocol format
func (ms *MoveSlice) StringUci() string {
	var os strings.Builder
	size := len(*ms)
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := (*ms)[i]
		os.WriteString(m.StringUci())
	}
	return os.String()
}

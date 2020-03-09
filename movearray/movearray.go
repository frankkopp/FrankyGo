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

// Package movearray provides a array (slice) facade to be used with
// chess moves.
package movearray

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/frankkopp/FrankyGo/types"
)

// MoveArray represents a data structure (go slice) for Move.
type MoveArray struct {
	data []Move
}

// New creates a new move array with the given capacity and 0 elements
func New(cap int) MoveArray {
	ma := MoveArray{}
	ma.data = make([]Move, 0, cap)
	return ma
}

// Len returns the number of moves currently stored in the array
func (ma *MoveArray) Len() int {
	return len(ma.data)
}

// Cap returns the capacity of the array
func (ma *MoveArray) Cap() int {
	return cap(ma.data)
}

// PushBack appends an element at the end of the array
func (ma *MoveArray) PushBack(m Move) {
	ma.data = append(ma.data, m)
}

// PopBack removes and returns the move from the back of the queue.
// If the queue is empty, the call panics.
func (ma *MoveArray) PopBack() Move {
	if len(ma.data) <= 0 {
		panic("MoveArray: PopBack() called on empty array")
	}
	backMove := ma.data[len(ma.data)-1]
	ma.data = ma.data[:len(ma.data)-1]
	return backMove
}

// PushFront prepends an element at the beginning of the array
func (ma *MoveArray) PushFront(m Move) {
	ma.data = append(ma.data, MoveNone)
	copy(ma.data[1:], ma.data)
	ma.data[0] = m
}

// PopFront removes and returns the move from the front of the array.
// If the array is empty, the call panics.
// Shrinks the capacity of the array and might lead to earlier
// re-allocations
func (ma *MoveArray) PopFront() Move {
	if len(ma.data) <= 0 {
		panic("MoveArray: PopFront() called on empty array")
	}
	frontMove := ma.data[0]
	ma.data = ma.data[1:]
	return frontMove
}

// Front returns the move at the front of the array.  This is the element
// that would be returned by PopFront(). This call panics if the array is
// empty.
func (ma *MoveArray) Front() Move {
	if len(ma.data) <= 0 {
		panic("MoveArray: Front() called when empty")
	}
	return ma.data[0]
}

// Back returns the move at the back of the array.  This is the element
// that would be returned by PopBack().  This call panics if the array is
// empty.
func (ma *MoveArray) Back() Move {
	if len(ma.data) <= 0 {
		panic("MoveArray: Back() called when empty")
	}
	return ma.data[len(ma.data)-1]
}

// At returns the move at index i in the array without removing the move
// from the array. At(0) refers to the first move and is the same as Front().
// At(Len()-1) refers to the last move and is the same as Back().
// Index will not be checked against bounds.
func (ma *MoveArray) At(i int) Move {
	return ma.data[i]
}

// Set puts a move at index i in the queue. Set shares the same purpose
// than At() but perform the opposite operation. The index i is the same
// index defined by At().
// Index will not be checked against bounds.
func (ma *MoveArray) Set(i int, move Move) {
	ma.data[i] = move
}

// Filter removes all elements from the MoveArray for
// which the given call to func will return false.
// Rebuilds the data slice by looping over all elements
// and only re-adding elements for which the call to the
// given func is true. Reuses the underlying array
func (ma *MoveArray) Filter(f func (index int) bool) {
	b := ma.data[:0]
	for i, x := range ma.data {
		if f(i) {
			b = append(b, x)
		}
	}
	ma.data = b
}

// FilterCopy copies the MoveArray into a new MoveArray
// without the filtered elements. AN element is filtered when
// the given call to func will return false for the element.
func (ma *MoveArray) FilterCopy(f func (index int) bool) *MoveArray {
	newArray := New(cap(ma.data))
	for i, x := range ma.data {
		if f(i) {
			newArray.data = append(newArray.data, x)
		}
	}
	return &newArray
}

// ForEach simple range loop calling the given function on each element
// in stored order
func (ma *MoveArray) ForEach(f func (index int)) {
	for index, _ := range ma.data  {
		f(index)
	}
}

// ForEachParallel simple loop over all elements calling a goroutine
// which calls the given func with the index of the current element
// as a parameter.
// Waits until all elements have been processed. There is not
// synchronization for the parallel execution. This needs to done
// in the provided function
func (ma *MoveArray) ForEachParallel(f func (index int)) {
	sliceLength := len(ma.data)
	var wg sync.WaitGroup
	wg.Add(sliceLength)
	for index, _ := range ma.data  {
		go func(i int) {
			defer wg.Done()
			f(i)
		}(index)
	}
	wg.Wait()
}

// Clear removes all moves from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.
func (ma *MoveArray) Clear() {
	ma.data = ma.data[:0]
}

// Sort sorts the moves from highest value to lowest value
// Uses InsertionSort as MoveArrays are mostly pre-sorted and small
func (ma *MoveArray) Sort() {
	l := len(ma.data)
	for i := 1; i < l; i++ {
		tmp := ma.data[i]
		j := i
		for  j > 0 && tmp > ma.data[j-1] {
			ma.data[j] = ma.data[j-1]
			j--
		}
		ma.data[j] = tmp
	}
}

// String returns a string representation of a move list
func (ma *MoveArray) String() string {
	var os strings.Builder
	size := ma.Len()
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
func (ma *MoveArray) StringUci() string {
	var os strings.Builder
	size := ma.Len()
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := ma.At(i)
		os.WriteString(m.StringUci())
	}
	return os.String()
}

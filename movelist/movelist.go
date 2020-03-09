/*
 * MIT License
 *
 * Copyright (c) 2018 Andrew J. Gillis (original deque - also MIT license)
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

// Package movelist adapts the original deque package from github.com/gammazero/deque
// to be used with Move only and adding sort interface and String(). Extending the
// original package was not possible as it did not allow to write to arbitrary elements.
// Adapting it also gives a chance to optimize for the intended usage for Move
package movelist

import (
	"fmt"
	"strings"

	. "github.com/frankkopp/FrankyGo/types"
)

// minCapacity is the smallest capacity that deque may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// MoveList represents a single instance of the MoveList data structure.
type MoveList struct {
	buf    []Move
	head   int
	tail   int
	count  int
	minCap int
}

// Len returns the number of elements currently stored in the queue.
func (ml *MoveList) Len() int {
	return ml.count
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (ml *MoveList) PushBack(elem Move) {
	ml.growIfFull()

	ml.buf[ml.tail] = elem
	// Calculate new tail position.
	ml.tail = ml.next(ml.tail)
	ml.count++
}

// PushFront prepends an element to the front of the queue.
func (ml *MoveList) PushFront(elem Move) {
	ml.growIfFull()

	// Calculate new head position.
	ml.head = ml.prev(ml.head)
	ml.buf[ml.head] = elem
	ml.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (ml *MoveList) PopFront() Move {
	if ml.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	ret := ml.buf[ml.head]
	ml.buf[ml.head] = MoveNone
	// Calculate new head position.
	ml.head = ml.next(ml.head)
	ml.count--

	ml.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (ml *MoveList) PopBack() Move {
	if ml.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}

	// Calculate new tail position
	ml.tail = ml.prev(ml.tail)

	// Remove value at tail.
	ret := ml.buf[ml.tail]
	ml.buf[ml.tail] = MoveNone
	ml.count--

	ml.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (ml *MoveList) Front() Move {
	if ml.count <= 0 {
		panic("deque: Front() called when empty")
	}
	return ml.buf[ml.head]
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (ml *MoveList) Back() Move {
	if ml.count <= 0 {
		panic("deque: Back() called when empty")
	}
	return ml.buf[ml.prev(ml.tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue.  This method accepts only non-negative index values.  At(0)
// refers to the first element and is the same as Front().  At(Len()-1) refers
// to the last element and is the same as Back().  If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Deque to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the deque, but may be read from any place within the deque.  Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end.  All the log entries
// in the buffer must be readable without altering the buffer contents.
func (ml *MoveList) At(i int) Move {
	if i < 0 || i >= ml.count {
		panic("deque: At() called with index out of range")
	}
	// bitwise modulus
	return ml.buf[(ml.head+i)&(len(ml.buf)-1)]
}

// Set puts the element at index i in the queue. Set shares the same purpose
// than At() but perform the opposite operation. The index i is the same
// index defined by At(). If the index is invalid, the call panics.
func (ml *MoveList) Set(i int, elem Move) {
	if i < 0 || i >= ml.count {
		panic("MoveList: Set() called with index out of range")
	}
	// bitwise modulus
	ml.buf[(ml.head+i)&(len(ml.buf)-1)] = elem
}

// Copy takes the indices of two elements and copies the element at index i
// into index j. Copy is a shortcut for q.Set(j) = q.At(i). The indices i and j
// are the same indices defined by At(). If one of the indices is invalid, the
// call panics.
func (ml *MoveList) Copy(i int, j int) {
	if i < 0 || i >= ml.count || j < 0 || j >= ml.count {
		panic("MoveList: Copy() called with index out of range")
	}
	// bitwise modulus
	ml.buf[(ml.head+j)&(len(ml.buf)-1)] = ml.buf[(ml.head+i)&(len(ml.buf)-1)]
}


// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse. The queue will not be resized smaller as long as items are
// only added.  Only when items are removed is the queue subject to getting
// resized smaller.
func (ml *MoveList) Clear() {
	// bitwise modulus
	modBits := len(ml.buf) - 1
	for h := ml.head; h != ml.tail; h = (h + 1) & modBits {
		ml.buf[h] = MoveNone
	}
	ml.head = 0
	ml.tail = 0
	ml.count = 0
}

// Rotate rotates the deque n steps front-to-back.  If n is negative, rotates
// back-to-front.  Having Deque provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (ml *MoveList) Rotate(n int) {
	if ml.count <= 1 {
		return
	}
	// Rotating a multiple of q.count is same as no rotation.
	n %= ml.count
	if n == 0 {
		return
	}

	modBits := len(ml.buf) - 1
	// If no empty space in buffer, only move head and tail indexes.
	if ml.head == ml.tail {
		// Calculate new head and tail using bitwise modulus.
		ml.head = (ml.head + n) & modBits
		ml.tail = (ml.tail + n) & modBits
		return
	}

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail using bitwise modulus.
			ml.head = (ml.head - 1) & modBits
			ml.tail = (ml.tail - 1) & modBits
			// Put tail value at head and remove value at tail.
			ml.buf[ml.head] = ml.buf[ml.tail]
			ml.buf[ml.tail] = MoveNone
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head value at tail and remove value at head.
		ml.buf[ml.tail] = ml.buf[ml.head]
		ml.buf[ml.head] = MoveNone
		// Calculate new head and tail using bitwise modulus.
		ml.head = (ml.head + 1) & modBits
		ml.tail = (ml.tail + 1) & modBits
	}
}

// SetMinCapacity sets a minimum capacity of 2^minCapacityExp.  If the value of
// the minimum capacity is less than or equal to the minimum allowed, then
// capacity is set to the minimum allowed.  This may be called at anytime to
// set a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing when the
// number of stored items changes frequently across a wide range.
func (ml *MoveList) SetMinCapacity(minCapacityExp uint) {
	if 1<<minCapacityExp > minCapacity {
		ml.minCap = 1 << minCapacityExp
	} else {
		ml.minCap = minCapacity
	}
}

// Less for MoveList sorts elements in descending order (so not less but more)
func (ml *MoveList) Less(i, j int) bool {
	return ml.buf[(ml.head+i)&(len(ml.buf)-1)] > ml.buf[(ml.head+j)&(len(ml.buf)-1)]
}

// Swap swaps the elements with indexes i and j.
func (ml *MoveList) Swap(i, j int) {
	l := len(ml.buf) - 1
	tmp := ml.buf[(ml.head+i)&l]
	ml.buf[(ml.head+i)&l] = ml.buf[(ml.head+j)&l]
	ml.buf[(ml.head+j)&l] = tmp
}

// String returns a string representation of a move list
func (ml *MoveList) String() string {
	var os strings.Builder
	size := ml.Len()
	os.WriteString(fmt.Sprintf("MoveList: [%d] { ", size))
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(", ")
		}
		m := ml.At(i)
		os.WriteString(m.String())
	}
	os.WriteString(" }")
	return os.String()
}

// StringUci returns a string with a sapce seperated list
// of all moves i the list in UCI protocol format
func (ml *MoveList) StringUci() string {
	var os strings.Builder
	size := ml.Len()
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := ml.At(i)
		os.WriteString(m.StringUci())
	}
	return os.String()
}

// prev returns the previous buffer position wrapping around buffer.
func (ml *MoveList) prev(i int) int {
	return (i - 1) & (len(ml.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (ml *MoveList) next(i int) int {
	return (i + 1) & (len(ml.buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (ml *MoveList) growIfFull() {
	if len(ml.buf) == 0 {
		if ml.minCap == 0 {
			ml.minCap = minCapacity
		}
		ml.buf = make([]Move, ml.minCap)
		return
	}
	if ml.count == len(ml.buf) {
		ml.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (ml *MoveList) shrinkIfExcess() {
	if len(ml.buf) > ml.minCap && (ml.count<<2) == len(ml.buf) {
		ml.resize()
	}
}

// resize resizes the deque to fit exactly twice its current contents.  This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (ml *MoveList) resize() {
	newBuf := make([]Move, ml.count<<1)
	if ml.tail > ml.head {
		copy(newBuf, ml.buf[ml.head:ml.tail])
	} else {
		n := copy(newBuf, ml.buf[ml.head:])
		copy(newBuf[n:], ml.buf[:ml.tail])
	}

	ml.head = 0
	ml.tail = ml.count
	ml.buf = newBuf
}



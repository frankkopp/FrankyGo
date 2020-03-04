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

package types

import (
	"fmt"
	"strings"

	"github.com/gammazero/deque"
)

// MoveList a list of moves based on a deque data structure
type MoveList struct {
	deque.Deque
}

func (ml MoveList) String() string {
	var os strings.Builder
	size := ml.Len()
	os.WriteString(fmt.Sprintf("MoveList: [%d] { ", size))
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(", ")
		}
		m := ml.At(i)
		os.WriteString(m.(Move).String())
	}
	os.WriteString(" }")
	return os.String()
	return ""
}

// StringUci returns a string with a sapce seperated list
// of all moves i the list in UCI protocol format
func (ml MoveList) StringUci() string {
	var os strings.Builder
	size := ml.Len()
	for i := 0; i < size; i++ {
		if i > 0 {
			os.WriteString(" ")
		}
		m := ml.At(i)
		os.WriteString(m.(Move).StringUci())
	}
	return os.String()
}





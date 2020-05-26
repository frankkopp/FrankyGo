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

// Package history provides data structures and functionality to manage
// history driven move tables (e.g. history counter, counter moves, etc.)
package history

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	. "github.com/frankkopp/FrankyGo/pkg/types"
)

var out = message.NewPrinter(language.German)

// History is a data structure updated during search to provide the move
// generator with valuable information for move sorting.
type History struct {
	HistoryCount [2][64][64]int64
	CounterMoves [64][64]Move
}

func (h History) String() string {
	sb := strings.Builder{}
	for sf := SqA1; sf < SqNone; sf++ {
		for st := SqA1; st < SqNone; st++ {
			sb.WriteString(out.Sprintf("Move=%s%s: ", sf.String(), st.String()))
			for c := White; c <= 1; c++ {
				count := h.HistoryCount[c][sf][st]
				sb.WriteString(out.Sprintf("%s=%-7d ", c.String(), count))
			}
			m := h.CounterMoves[sf][st]
			sb.WriteString(out.Sprintf("cm=%s\n", m.StringUci()))
		}
	}
	return sb.String()
}

// NewHistory creates a new History instance.
func NewHistory() *History {
	return &History{}
}

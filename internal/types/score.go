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

package types

import (
	"fmt"
)

// Score is a small struct for mid game and end game values
type Score struct {
	MidGameValue int
	EndGameValue int
}

// Add adds the corresponding parts of the given score to the
// calling score
func (s *Score) Add(a Score) {
	s.MidGameValue += a.MidGameValue
	s.EndGameValue += a.EndGameValue
}

// Sub subtracts the corresponding parts of the given score from the
// calling score
func (s *Score) Sub(a Score) {
	s.MidGameValue -= a.MidGameValue
	s.EndGameValue -= a.EndGameValue
}

// ValueFromScore adds up the mid and end games scores after multiplying
// them with the game phase factor
func (s *Score) ValueFromScore(gpf float64) Value {
	return Value(float64(s.MidGameValue)*gpf) + Value(float64(s.EndGameValue)*(1.0-gpf))
}

func (s *Score) String() string {
	return fmt.Sprintf("{ mid:%d end:%d }", s.MidGameValue, s.EndGameValue)
}

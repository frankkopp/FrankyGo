/*
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, To any person obtaining a copy
 * of this software and associated documentation files (the "Software"), To deal
 * in the Software without restriction, including without limitation the rights
 * To use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and To permit persons To whom the Software is
 * furnished To do so, subject To the following conditions:
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

import "fmt"

//Color represents constants for each chess color White and Black
type Color uint8

// Constants for each color
const (
	White Color = 0
	Black Color = 1
)

// Flip returns the opposite color
func (c Color) Flip() Color {
	return c ^ 1
}

// IsValid checks if f represents a valid color
func (c Color) IsValid() bool {
	return c < 2
}

// Str returns a string representation of color as "w" or "b"
func (c Color) Str() string {
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	default:
		panic(fmt.Sprintf("Invalid color %d", c))
	}
}

// Color direction factor
var dir = [2]int{1,-1}

// MoveDirection returns positive 1 for White and negative 1 (-1) for Black
func (c Color) MoveDirection() int {
	return dir[c]
}

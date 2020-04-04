// +build !debug

/*
 * FrankyGo - UCI chess engine in GO for learning purposes
 *
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

// Package assert is a helper to allow assert tests in a more standardized
// and simple manner. Using it makes it clear the this is an assertion
// used in non production setting.
package assert

import (
	"fmt"
)

func init()  {
	fmt.Println("RELEASE MODE")
}

// DEBUG if this is set to "true" asserts are evaluated
const DEBUG = false

// Assert runs the provided function and throws
// panic with the given message if the test evaluates to false.
// Unfortunately GO still executes parameters (e.g. value.String()
// of calls to this even if the function is a null function when
// DEBUG is set to false. So it is necessary to also have a if assert.DEBUG {}
// wrapper around calls to this to really avoid any run time
// impact. The GO compiler will then eliminate the whole statement
// if DEBUG as a const is set to false.
// Example:
//  if assert.DEBUG {
//	  assert.Assert(value > 0, "Error message if test fails %s", value.String())
//  }
func Assert(test bool, msg string, a ...interface{}) {}

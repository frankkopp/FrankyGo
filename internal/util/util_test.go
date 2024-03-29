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

package util

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)


var logTest *logging.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, Abs(-5))
	assert.Equal(t, 5, Abs(5))
	assert.Equal(t, int64(5), Abs64(int64(-5)))
	assert.Equal(t, int64(5), Abs64(int64(5)))
}

func TestMinMax(t *testing.T) {
	assert.Equal(t, -5, Min(-5, -3))
	assert.Equal(t, -3, Max(-5, -3))
	assert.Equal(t, int64(-5), Min64(int64(-5), int64(-3)))
	assert.Equal(t, int64(-3), Max64(int64(-5), int64(-3)))
}

var tmp, result int64
var index int64

func BenchmarkMax64(b *testing.B) {
	for index = -int64(b.N); index < int64(b.N); index++ {
		tmp = Max64(index, index+2)
	}
	result = tmp
}

func BenchmarkMin64(b *testing.B) {
	for index = -int64(b.N); index < int64(b.N); index++ {
		tmp = Min64(index, index+2)
	}
	result = tmp
}


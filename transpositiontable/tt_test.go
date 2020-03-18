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

package transpositiontable

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/franky_logging"
	"github.com/frankkopp/FrankyGo/types"
)

var logTest = franky_logging.GetLog("test")

func TestEntrySize(t *testing.T) {
	e := TtEntry{}
	assert.Equal(t, 16, unsafe.Sizeof(e))
	logTest.Debugf("Size of Entry %d bytes", unsafe.Sizeof(e))
}

func TestNew(t *testing.T) {

	tt := New(2)
	assert.Equal(t, uint64(131_072), tt.maxNumberOfEntries)
	assert.Equal(t, 131_072, cap(tt.data))

	tt = New(64)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = New(100)
	assert.Equal(t, uint64(4_194_304), tt.maxNumberOfEntries)
	assert.Equal(t, 4_194_304, cap(tt.data))

	tt = New(4_096)
	assert.Equal(t, uint64(268_435_456), tt.maxNumberOfEntries)
	assert.Equal(t, 268_435_456, cap(tt.data))

	tt = New(35_000)
	assert.Equal(t, uint64(2_147_483_648), tt.maxNumberOfEntries)
	assert.Equal(t, 2_147_483_648, cap(tt.data))
	assert.Equal(t, 32_768*types.MB, tt.sizeInByte)
}

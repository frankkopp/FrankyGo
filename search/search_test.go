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

package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/logging"
)

var logTest = logging.GetLog("test")

func TestWaitWhileSearching(t *testing.T) {
	search := NewSearch()
	start := time.Now()
	// FIXME: Prototype
	search.Start()
	logTest.Debug("Search started...waiting to finish")
	search.WaitWhileSearching()
	logTest.Debug("Search finished")
	elapsed := time.Since(start)
	out.Printf("Time %d ms\n", elapsed.Milliseconds())
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(2_000))
}

func TestIsSearching(t *testing.T) {
	search := NewSearch()
	start := time.Now()
	// FIXME: Prototype
	search.Start()
	logTest.Debug("Check searching in 1 sec")
	time.Sleep(time.Second)
	assert.True(t, search.IsSearching())
	logTest.Debugf("Is searching = %v", search.IsSearching())
	search.WaitWhileSearching()
	elapsed := time.Since(start)
	out.Printf("Time %d ms\n", elapsed.Milliseconds())
	assert.False(t, search.IsSearching())
	logTest.Debugf("Is searching = %v", search.IsSearching())
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(2_000))
}

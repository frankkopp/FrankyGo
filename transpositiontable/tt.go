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
	"math"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/franky_logging"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/types"
)

var out = message.NewPrinter(language.German)
var log = franky_logging.GetLog("tt")

// TtEntry struct is the data structure for each entry in the transposition
// table. Each entry has 16-bytes (128-bits)
type TtEntry struct {
	Key        position.Key    // 64-bit Zobrist Key
	Move       types.Move      // 32-bit Move and Value
	Depth      int8            // 7-bit 0-127 0b01111111
	Age        int8            // 3-bit 0-7   0b00000111 0=used 1=generated, not used, >1 older generation
	Type       types.ValueType // 2-bit None, Exact, Alpha (upper), Beta (lower)
	MateThreat bool            // 1-bit
}

const (
	// TtEntrySize is the size in bytes for each TtEntry
	TtEntrySize = 16 // 16 bytes

	// MaxSizeInMB maximal memory usage of tt
	MaxSizeInMB = 65_536
)

// TtTable is the actual transposition table
// object holding data and state.
type TtTable struct {
	data               []TtEntry
	sizeInByte         uint64
	hashKeyMask        uint64
	maxNumberOfEntries uint64
	numberOfEntries    uint64

	// statistics
	numberOfPuts       uint64
	numberOfCollisions uint64
	numberOfOverwrites uint64
	numberOfUpdates    uint64
	numberOfProbes     uint64
	numberOfHits       uint64
	numberOfMisses     uint64
}

// New creates a new TtTable with the given number of bytes
// as a maximum of memory usage. Actual size will be determined
// by the number of elements which need to be a power of 2
func New(sizeInMByte int) *TtTable {
	tt := TtTable{
		data:               nil,
		sizeInByte:         0,
		hashKeyMask:        0,
		maxNumberOfEntries: 0,
		numberOfEntries:    0,
	}

	tt.Resize(sizeInMByte)

	return &tt
}

// Resize resizes the tt table. All entries will be cleared.
func (tt *TtTable) Resize(sizeInMByte int) {
	if sizeInMByte > MaxSizeInMB {
		log.Error(out.Sprintf("Requested size for TT of %d MB reduced to max of %d MB", sizeInMByte, MaxSizeInMB))
		sizeInMByte = MaxSizeInMB
	}

	// calculate the maximum power of 2 of entries fitting into the given size in MB
	tt.sizeInByte = uint64(sizeInMByte) * types.MB
	tt.maxNumberOfEntries = 1 << uint64(math.Floor(math.Log2(float64(tt.sizeInByte/TtEntrySize))))
	tt.hashKeyMask = tt.maxNumberOfEntries - 1 // --> 0x0001111....111

	// if TT is resized to 0 we cant have any entries.
	if tt.sizeInByte == 0 {
		tt.maxNumberOfEntries = 0
	}

	// calculate the real memory usage
	tt.sizeInByte = tt.maxNumberOfEntries * TtEntrySize

	// Create new slice/array - garbage collections takes care of cleanup
	tt.data = make([]TtEntry, tt.maxNumberOfEntries, tt.maxNumberOfEntries)

	log.Info(out.Sprintf("TT Size %d MByte, Capacity %d entries (size=%dByte) (Requested were %d MBytes)",
		tt.sizeInByte/types.MB, tt.maxNumberOfEntries, unsafe.Sizeof(TtEntry{}), sizeInMByte))
}

// GetEntry returns a pointer to the corresponding tt entry
// or nil if it was not found
func (tt *TtTable) GetEntry(key position.Key) *TtEntry {
	e := &tt.data[tt.hash(key)]
	if e.Key == key {
		return e
	}
	return nil
}

// Probe returns a pointer to the corresponding tt entry
// or nil if it was not found. Decreases TtEntry.Age by 1
func (tt *TtTable) Probe(key position.Key) *TtEntry {
	e := &tt.data[tt.hash(key)]
	if e.Key == key {
		e.Age--
		if e.Age < 0 {
			e.Age = 0
		}
		return e
	}
	return nil
}

// Clear clears all entries of the tt
func (tt *TtTable) Clear() {
	// Create new slice/array - garbage collections takes care of cleanup
	tt.data = make([]TtEntry, tt.maxNumberOfEntries, tt.maxNumberOfEntries)
	tt.numberOfPuts = 0
	tt.numberOfEntries = 0
	tt.numberOfHits = 0
	tt.numberOfUpdates = 0
	tt.numberOfMisses = 0
	tt.numberOfCollisions = 0
	tt.numberOfOverwrites = 0
	tt.numberOfProbes = 0
}

// Hashfull returns how full the transposition table is in permill as per UCI
func (tt *TtTable) Hashfull() int {
	if tt.maxNumberOfEntries == 0 {
		return 0
	}
	return int((1000 * tt.numberOfEntries) / tt.maxNumberOfEntries)
}

// hash generates the internal hash key for the data array
func (tt *TtTable) hash(key position.Key) uint64 {
	return uint64(key) & tt.hashKeyMask
}

// ages each entry in the tt
// Creates a number of go routines with each a certain slice of data
// to process
func (tt *TtTable) ageEntries() {
	numberOfGoroutines := uint64(8)
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(int(numberOfGoroutines))
	slice := tt.maxNumberOfEntries / numberOfGoroutines
	for i := uint64(0); i < numberOfGoroutines; i++ {
		go func(i uint64) {
			defer wg.Done()
			start := i * slice
			end := start + slice
			if i == numberOfGoroutines-1 {
				end = tt.maxNumberOfEntries
			}
			for n := start; n < end; n++ {
				if tt.data[n].Key != 0 {
					tt.data[n].Age++
				}
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	log.Debug(out.Sprintf("Aged %d entries in %d ms\n", len(tt.data), elapsed.Milliseconds()))
}

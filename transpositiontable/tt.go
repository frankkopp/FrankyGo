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

// Package transpositiontable implements a transposition table (cache)
// data structure and functionality for a chess engine search.
// The TtTable class is not thread safe and needs to be synchronized
// externally if used from multiple threads. Is especially relevant
// for Resize and Clear which should not be called in parallel
// while searching.
package transpositiontable

import (
	"math"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/assert"
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/util"
)

var out = message.NewPrinter(language.German)
var log = logging.GetLog("tt")

// TtEntry struct is the data structure for each entry in the transposition
// table. Each entry has 16-bytes (128-bits)
type TtEntry struct {
	Key        position.Key // 64-bit Zobrist Key
	Move       Move         // 32-bit Move and Value
	Depth      int8         // 7-bit 0-127 0b01111111
	Age        int8         // 3-bit 0-7   0b00000111 0=used 1=generated, not used, >1 older generation
	Type       ValueType    // 2-bit None, Exact, Alpha (upper), Beta (lower)
	MateThreat bool         // 1-bit
}

const (
	// TtEntrySize is the size in bytes for each TtEntry
	TtEntrySize = 16 // 16 bytes

	// MaxSizeInMB maximal memory usage of tt
	MaxSizeInMB = 65_536
)

// TtTable is the actual transposition table
// object holding data and state.
// Create with NewTtTable()
type TtTable struct {
	data               []TtEntry
	sizeInByte         uint64
	hashKeyMask        uint64
	maxNumberOfEntries uint64
	numberOfEntries    uint64
	Stats              TtStats
}

// TtStats holds statistical data on tt usage
type TtStats struct {
	numberOfPuts       uint64
	numberOfCollisions uint64
	numberOfOverwrites uint64
	numberOfUpdates    uint64
	numberOfProbes     uint64
	numberOfHits       uint64
	numberOfMisses     uint64
}

// NewTtTable creates a new TtTable with the given number of bytes
// as a maximum of memory usage. Actual size will be determined
// by the number of elements fitting into this size which need
// to be a power of 2 for efficient hashing/addressing via bit
// masks
func NewTtTable(sizeInMByte int) *TtTable {
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
// The TtTable class is not thread safe and needs to be synchronized
// externally if used from multiple threads. Is especially relevant
// for Resize and Clear which should not be called in parallel
// while searching.
func (tt *TtTable) Resize(sizeInMByte int) {
	if sizeInMByte > MaxSizeInMB {
		log.Error(out.Sprintf("Requested size for TT of %d MB reduced to max of %d MB", sizeInMByte, MaxSizeInMB))
		sizeInMByte = MaxSizeInMB
	}

	// calculate the maximum power of 2 of entries fitting into the given size in MB
	tt.sizeInByte = uint64(sizeInMByte) * MB
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
		tt.sizeInByte/MB, tt.maxNumberOfEntries, unsafe.Sizeof(TtEntry{}), sizeInMByte))
	log.Debug(util.MemStat())
}

// GetEntry returns a pointer to the corresponding tt entry.
// The entry could be an empty entry with Key==0.
// Does not change statistics.
func (tt *TtTable) GetEntry(key position.Key) *TtEntry {
	return &tt.data[tt.hash(key)]
}

// Probe returns a pointer to the corresponding tt entry
// or nil if it was not found. Decreases TtEntry.Age by 1
func (tt *TtTable) Probe(key position.Key) *TtEntry {
	tt.Stats.numberOfProbes++
	e := &tt.data[tt.hash(key)]
	if e.Key == key {
		e.Age--
		if e.Age < 0 {
			e.Age = 0
		}
		tt.Stats.numberOfHits++
		return e
	}
	tt.Stats.numberOfMisses++
	return nil
}

// Put an TtEntry into the tt. Encodes value into the move.
func (tt *TtTable) Put(key position.Key, move Move, value Value, depth int8, valueType ValueType, mateThreat bool, forced bool) {
	if assert.DEBUG {
		assert.Assert(depth >= 0, "TT:put Depth must be > 0")
	}
	// if the size of the TT = 0 we
	// do not store anything
	if tt.maxNumberOfEntries == 0 {
		return
	}

	tt.Stats.numberOfPuts++
	// read the entries for this hash
	entryDataPtr := tt.GetEntry(key)
	// encode value into the move
	valueMove := move.SetValue(value)

	// NewTtTable entry
	if entryDataPtr.Key == 0 {
		tt.numberOfEntries++
		entryDataPtr.Key = key
		entryDataPtr.Move = valueMove
		entryDataPtr.Depth = depth
		entryDataPtr.Age = 1
		entryDataPtr.Type = valueType
		entryDataPtr.MateThreat = mateThreat
		return
	}

	// Same hash but different position
	if entryDataPtr.Key != key {
		tt.Stats.numberOfCollisions++
		// overwrite if
		// - the new entry's depth is higher
		// - the new entry's depth is same and the previous entry is old (is aged)
		if depth > entryDataPtr.Depth ||
			(depth == entryDataPtr.Depth && (forced || entryDataPtr.Age > 1)) {
			tt.Stats.numberOfOverwrites++
			entryDataPtr.Key = key
			entryDataPtr.Move = valueMove
			entryDataPtr.Depth = depth
			entryDataPtr.Age = 1
			entryDataPtr.Type = valueType
			entryDataPtr.MateThreat = mateThreat
		}
		return
	}

	// Same hash and same position -> update entry?
	if entryDataPtr.Key == key {
		tt.Stats.numberOfUpdates++
		// we always update as the stored moved can't be any good otherwise
		// we would have found this during the search in a previous probe
		// and we would not have come to store it again
		entryDataPtr.Key = key
		entryDataPtr.Move = valueMove
		entryDataPtr.Depth = depth
		entryDataPtr.Age = 1
		entryDataPtr.Type = valueType
		entryDataPtr.MateThreat = mateThreat
		return
	}

	if assert.DEBUG {
		assert.Assert(tt.Stats.numberOfPuts == (tt.numberOfEntries+tt.Stats.numberOfCollisions+tt.Stats.numberOfUpdates),
			"TT:put - stat values do not match")
	}
}

// Clear clears all entries of the tt
// The TtTable class is not thread safe and needs to be synchronized
// externally if used from multiple threads. Is especially relevant
// for Resize and Clear which should not be called in parallel
// while searching.
func (tt *TtTable) Clear() {
	// Create new slice/array - garbage collections takes care of cleanup
	tt.data = make([]TtEntry, tt.maxNumberOfEntries, tt.maxNumberOfEntries)
	tt.numberOfEntries = 0
	tt.Stats = TtStats{}
}

// Hashfull returns how full the transposition table is in permill as per UCI
func (tt *TtTable) Hashfull() int {
	if tt.maxNumberOfEntries == 0 {
		return 0
	}
	return int((1000 * tt.numberOfEntries) / tt.maxNumberOfEntries)
}

// String returns a string representation of this TtTable instance
func (tt *TtTable) String() string {
	return out.Sprintf("TT: size %d MB max entries %d of size %d Bytes entries %d (%d) puts %d "+
		"updates %d collisions %d overwrites %d probes %d hits %d (%d) misses %d (%d)",
		tt.sizeInByte/MB, tt.maxNumberOfEntries, unsafe.Sizeof(TtEntry{}), tt.numberOfEntries, tt.Hashfull(),
		tt.Stats.numberOfPuts, tt.Stats.numberOfUpdates, tt.Stats.numberOfCollisions, tt.Stats.numberOfOverwrites, tt.Stats.numberOfProbes,
		tt.Stats.numberOfHits, (tt.Stats.numberOfHits*100)/(1+tt.Stats.numberOfProbes),
		tt.Stats.numberOfMisses, (tt.Stats.numberOfMisses*100)/(1+tt.Stats.numberOfProbes))
}

// Len returns the number of non empty entries in the tt
func (tt *TtTable) Len() uint64 {
	return tt.numberOfEntries
}

// ///////////////////////////////////////////////////////////
// Private
// ///////////////////////////////////////////////////////////

// hash generates the internal hash key for the data array
func (tt *TtTable) hash(key position.Key) uint64 {
	return uint64(key) & tt.hashKeyMask
}

// AgeEntries ages each entry in the tt
// Creates a number of go routines with processes each
// a certain slice of data to process
func (tt *TtTable) AgeEntries() {
	startTime := time.Now()
	if tt.numberOfEntries > 0 {
		numberOfGoroutines := uint64(32) // arbitrary - uses up to 32 threads
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
	}
	elapsed := time.Since(startTime)
	log.Debug(out.Sprintf("Aged %d entries of %d in %d ms\n", tt.numberOfEntries, len(tt.data), elapsed.Milliseconds()))
}

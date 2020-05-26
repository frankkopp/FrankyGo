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

package evaluator

import (
	"math"
	"unsafe"

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/internal/config"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	. "github.com/frankkopp/FrankyGo/pkg/types"
)

const (
	// MaxSizeInMB maximal memory usage of pawnCache
	MaxSizeInMB = 1_024

	// EntrySize is the size in bytes for each pawn cache entry
	EntrySize = 16 // 16 bytes
)

type pawnCache struct {
	log                *logging.Logger
	data               []cacheEntry
	sizeInByte         uint64
	maxNumberOfEntries uint64
	hashKeyMask        uint64
	capacity           uint64
	entries            uint64
	hits               uint64
	misses             uint64
	replace            uint64
}

type cacheEntry struct {
	pawnKey Key
	score   Score
}

func newPawnCache() *pawnCache {
	pc := &pawnCache{
		log: myLogging.GetLog(),
	}
	pc.resize(config.Settings.Eval.PawnCacheSize)
	return pc
}

func (pc *pawnCache) resize(sizeInMByte int) {
	if sizeInMByte > MaxSizeInMB {
		pc.log.Error(out.Sprintf("Requested size for Pawn Cache of %d MB reduced to max of %d MB", sizeInMByte, MaxSizeInMB))
		sizeInMByte = MaxSizeInMB
	}

	// calculate the maximum power of 2 of entries fitting into the given size in MB
	pc.sizeInByte = uint64(sizeInMByte) * MB
	pc.maxNumberOfEntries = 1 << uint64(math.Floor(math.Log2(float64(pc.sizeInByte/EntrySize))))
	pc.hashKeyMask = pc.maxNumberOfEntries - 1 // --> 0x0001111....111

	// if TT is resized to 0 we cant have any entries.
	if pc.sizeInByte == 0 {
		pc.maxNumberOfEntries = 0
	}

	// calculate the real memory usage
	pc.sizeInByte = pc.maxNumberOfEntries * EntrySize

	// Create new slice/array - garbage collections takes care of cleanup
	pc.data = make([]cacheEntry, pc.maxNumberOfEntries)

	pc.log.Info(out.Sprintf("PawnCache Size %d MByte, Capacity %d entries (size=%dByte) (Requested were %d MBytes)",
		pc.sizeInByte/MB, pc.maxNumberOfEntries, unsafe.Sizeof(cacheEntry{}), sizeInMByte))
}

// GetEntry returns a pointer to the corresponding entry.
// Given key is checked against the entry's key. When
// equal pointer to entry will be returned. Otherwise
// nil will be returned.
func (pc *pawnCache) getEntry(key Key) *cacheEntry {
	e := &pc.data[pc.hash(key)]
	if e.pawnKey == key {
		pc.hits++
		return e
	}
	pc.misses++
	return nil
}

// putEntry stores a Score for a pawn structure represented by the
// pawn zobrist key in the cache.
func (pc *pawnCache) put(key Key, score *Score) {
	e := &pc.data[pc.hash(key)]
	if e.pawnKey == 0 {
		pc.entries++
		e.pawnKey = key
		e.score.MidGameValue = score.MidGameValue
		e.score.EndGameValue = score.EndGameValue
		return
	}
	// update - should not happen at all
	if e.pawnKey == key {
		pc.log.Warningf("Update to pawn cache entry - should not happen. Missing a read to cache?")
	}
	// replace
	pc.replace++
	e.pawnKey = key
	e.score.MidGameValue = score.MidGameValue
	e.score.EndGameValue = score.EndGameValue
}

// Clear clears all entries of the pawn cache
func (pc *pawnCache) clear() {
	// Create new slice/array - garbage collections takes care of cleanup
	pc.data = make([]cacheEntry, pc.maxNumberOfEntries)
	pc.entries = 0
	pc.hits = 0
	pc.misses = 0
	pc.replace = 0
}

// len returns the number of non empty entries in the cache
func (pc *pawnCache) len() uint64 {
	return pc.entries
}

// ///////////////////////////////////////////////////////////
// Private
// ///////////////////////////////////////////////////////////

// hash generates the internal hash key for the data array
func (pc *pawnCache) hash(key Key) uint64 {
	return uint64(key) & pc.hashKeyMask
}

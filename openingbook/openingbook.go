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

// Package openingbook
// The OpeningBook reads game databases of different formats into an internal
// data structure. It can then be queried for a book move on a certain position.
//
// Supported formats are currently:
//
// BookFormat::SIMPLE for files storing a game per line with from-square and
// to-square notation
//
// BookFormat::SAN for files with lines of moves in SAN notation
//
// BookFormat::PGN for PGN formatted games<br/>
package openingbook

import (
	"bufio"
	"encoding/gob"
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/franky_logging"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/types"
)

var out = message.NewPrinter(language.German)
var log = franky_logging.GetLog("openingbook")

// setting to use multiple goroutines or not - useful for debugging
const parallel = true

// BookFormat represent the supported book formats defined as constants
type BookFormat uint8

// Supported book formats
const (
	Simple BookFormat = iota
	San    BookFormat = iota
	Pgn    BookFormat = iota
)

// Successor  represents a tuple of a move
// and a zobrist key of the position the
// move leads to
type Successor struct {
	Move      uint32
	NextEntry uint64
}

// BookEntry represents a data structure for a move in the opening book
// data structure. It describes exactly one position defined by a zobrist
// key and has links to other entries representing moves and successor
// positions
type BookEntry struct {
	ZobristKey uint64
	Counter    int
	Moves      []Successor
}

// Book represents a structure for chess opening books which can
// be read from different file formats into an internal data structure.
type Book struct {
	bookMap     map[uint64]BookEntry
	rootEntry   uint64
	initialized bool
}

// to test found moves against positions
var bookLock sync.Mutex

//noinspection GoUnhandledErrorResult
func (b *Book) Initialize(bookPath string, bookFormat BookFormat, useCache bool, recreateCache bool) error {
	if b.initialized {
		return nil
	}

	log.Info("Initializing Opening Book")

	startTotal := time.Now()

	// check file path
	if _, err := os.Stat(bookPath); err != nil {
		log.Errorf("File \"%s\" does not exist\n", bookPath)
		return err
	}

	// if cache enabled check if we have a cache file and load from cache
	if useCache && !recreateCache {
		startReading := time.Now()
		hasCache, err := b.loadFromCache(bookPath)
		elapsedReading := time.Since(startReading)
		if err != nil {
			log.Warningf("Cache could not be loaded. Reading original data from \"%s\"", bookPath)
		}
		if hasCache {
			log.Infof("Finished reading cache from file in: %d ms\n", elapsedReading.Milliseconds())
			log.Infof("Book from cache file contains %d entries\n", len(b.bookMap))
			return nil
		} // else no cache file just load the data from original file
	}

	// read book from file
	log.Infof("Reading opening book file: %s\n", bookPath)
	startReading := time.Now()
	lines, err := readFile(bookPath)
	if err != nil {
		log.Errorf("File \"%s\" could not be read: %s\n", bookPath, err)
		return err
	}
	elapsedReading := time.Since(startReading)
	log.Infof("Finished reading %d lines from file in: %d ms\n", len(*lines), elapsedReading.Milliseconds())

	// add root position
	startPosition := position.New()
	b.bookMap = make(map[uint64]BookEntry)
	b.rootEntry = uint64(startPosition.ZobristKey())
	b.bookMap[uint64(startPosition.ZobristKey())] = BookEntry{ZobristKey: uint64(startPosition.ZobristKey()), Counter: 0, Moves: []Successor{}}

	// process lines
	if parallel {
		log.Infof("Processing %d lines in parallel with format: %v\n", len(*lines), bookFormat)
	} else {
		log.Infof("Processing %d lines sequential with format: %v\n", len(*lines), bookFormat)
	}
	startProcessing := time.Now()
	err = b.process(lines, bookFormat)
	if err != nil {
		log.Errorf("Error while processing: %s\n", err)
		return err
	}
	elapsedProcessing := time.Since(startProcessing)
	log.Infof("Finished processing %d lines in: %d ms\n", len(*lines), elapsedProcessing.Milliseconds())

	// finished
	elapsedTotal := time.Since(startTotal)

	log.Infof("Book contains %d entries\n", len(b.bookMap))
	log.Infof("Total initialization time : %d ms\n", elapsedTotal.Milliseconds())

	// saving to cache
	if useCache {
		log.Infof("Saving to cache...")
		startSave := time.Now()
		cacheFile, nBytes, err := b.saveToCache(bookPath)
		if err != nil {
			log.Errorf("Error while saving to cache: %s\n", err)
		}
		elapsedSave := time.Since(startSave)
		bytes := out.Sprintf("%d", nBytes/1_024)
		log.Infof("Saved %s kB to cache %s in %d ms\n", bytes, cacheFile, elapsedSave.Milliseconds())
	}

	b.initialized = true
	return nil
}

// NumberOfEntries returns the number of entries in the opening book
func (b *Book) NumberOfEntries() int {
	return len(b.bookMap)
}

// GetEntry returns a copy of the entry with the corresponding key
func (b *Book) GetEntry(key position.Key) (BookEntry, bool) {
	entryPtr, ok := b.bookMap[uint64(key)]
	if ok {
		return entryPtr, true
	} else {
		return entryPtr, false
	}
}

// Reset resets the opening book so it can/must be initialized again
func (b *Book) Reset() {
	b.bookMap = map[uint64]BookEntry{}
	b.rootEntry = 0
	b.initialized = false
}

// /////////////////////////////////////////////////
// Private
// /////////////////////////////////////////////////

// reads a complete file into a slice of strings
func readFile(bookPath string) (*[]string, error) {
	f, err := os.Open(bookPath)
	if err != nil {
		log.Errorf("File \"%s\" could not be read; %s\n", bookPath, err)
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Errorf("File \"%s\" could not be closed: %s\n", bookPath, err)
		}
	}()
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	err = s.Err()
	if err != nil {
		log.Errorf("Error while reading file \"%s\": %s\n", bookPath, err)
		return nil, err
	}
	return &lines, nil
}

// sends all lines to the correct processing depending on format
func (b *Book) process(lines *[]string, format BookFormat) error {
	switch format {
	case Simple:
		b.processSimple(lines)
	case San:
		b.processSan(lines)
	case Pgn:
		b.processPgn(lines)
	}
	return nil
}

// processes all lines of Simple format
// uses goroutines in parallel if enabled
func (b *Book) processSimple(lines *[]string) {
	if parallel {
		sliceLength := len(*lines)
		var wg sync.WaitGroup
		wg.Add(sliceLength)
		for _, line := range *lines {
			go func(line string) {
				defer wg.Done()
				b.processSimpleLine(line)
			}(line)
		}
		wg.Wait()
	} else {
		for _, line := range *lines {
			b.processSimpleLine(line)
		}
	}
}

// regular expressions for detecting moves
var regexSimpleUciMove = regexp.MustCompile("([a-h][1-8][a-h][1-8])")

// processes one line of simple format and adds each move to book
func (b *Book) processSimpleLine(line string) {
	line = strings.TrimSpace(line)

	// find Uci moves
	matches := regexSimpleUciMove.FindAllString(line, -1)

	// skip lines without matches
	if len(matches) == 0 {
		return
	}

	// start with root position
	pos := position.New()

	// increase counter for root position
	bookLock.Lock()
	e, found := b.bookMap[b.rootEntry]
	if found {
		e.Counter++
		b.bookMap[b.rootEntry] = e
	} else {
		panic("root entry of book map not found")
	}
	bookLock.Unlock()

	// move gen to check moves
	// movegen is not thread safe therefore we create a new instance for every line
	var mg = movegen.New()

	// add all matches to book
	for _, moveString := range matches {
		err := b.processSingleMove(moveString, &mg, &pos)
		// stop processing further matches when we had an error as it
		// would probably be fruitless as position will be wrong
		if err != nil {
			break
		}
	}
}

// processes all lines of Simple format
// uses goroutines in parallel if enabled
func (b *Book) processSan(lines *[]string) {
	if parallel {
		sliceLength := len(*lines)
		var wg sync.WaitGroup
		wg.Add(sliceLength)
		for _, line := range *lines {
			go func(line string) {
				defer wg.Done()
				b.processSanLine(line)
			}(line)
		}
		wg.Wait()
	} else {
		for _, line := range *lines {
			b.processSanLine(line)
		}
	}
}

var regexResult = regexp.MustCompile("((1-0)|(0-1)|(1/2-1/2)|(\\*))$")

// Processes PGN formatted file. PGN file have additional
// metadata and spreads its move section over several lines.
// We ignore the metadata for now and only look at the
// move section.
func (b *Book) processPgn(lines *[]string) {
	// bundle games in slices of lines by finding result patterns
	var gamesSlices [][]string

	// find slices for each game
	startSlicing := time.Now()
	start := 0
	// end := len(*lines)
	for i, l := range *lines {
		l = strings.TrimSpace(l)
		if regexResult.MatchString(l) {
			end := i + 1
			gamesSlices = append(gamesSlices, (*lines)[start:end])
			start = end
		}
	}
	elapsedReading := time.Since(startSlicing)
	log.Infof("Finished finding %d games from file in: %d ms\n", len(gamesSlices), elapsedReading.Milliseconds())

	// process each game
	// uses goroutines in parallel if enabled
	startProcessing := time.Now()
	if parallel {
		noOfSlices := len(gamesSlices)
		var wg sync.WaitGroup
		wg.Add(noOfSlices)
		for _, gs := range gamesSlices {
			go func(gs []string) {
				defer wg.Done()
				b.processPgnGame(gs)
			}(gs)
		}
		wg.Wait()
	} else {
		for _, gs := range gamesSlices {
			b.processPgnGame(gs)
		}

	}
	elapsedProcessing := time.Since(startProcessing)
	log.Infof("Finished processing %d games from file in: %d ms\n", len(gamesSlices), elapsedProcessing.Milliseconds())
}

var regexTrailingComments = regexp.MustCompile(";.*$")
var regexTagPairs = regexp.MustCompile("\\[\\w+ +\".*?\"\\]")
var regexNagAnnotation = regexp.MustCompile("(\\$\\d{1,3})") // no NAG annotation supported
var regexBracketComments = regexp.MustCompile("{[^{}]*}")    // bracket comments
var regexReservedSymbols = regexp.MustCompile("<[^<>]*>")    // reserved symbols < >
var regexRavVariants = regexp.MustCompile("\\([^()]*\\)")    // RAV variant comments < >

// processes one game comprising of several input lines
func (b *Book) processPgnGame(gameSlice []string) {
	// build a cleaned up string of the move part of the PGN
	var moveLine strings.Builder

	// cleanup lines and concatenate move lines
	for _, l := range gameSlice {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "%") { // skip comment lines
			continue
		}
		// remove unnecessary parts and lines / order is important - comments last
		l = regexTagPairs.ReplaceAllString(l, "")
		l = regexResult.ReplaceAllString(l, "")
		l = regexTrailingComments.ReplaceAllString(l, "")
		l = strings.TrimSpace(l)
		// after cleanup skip now empty lines
		if len(l) == 0 {
			continue
		}
		// add the rest to the moveLine
		moveLine.WriteString(" ")
		moveLine.WriteString(l)
	}
	line := moveLine.String()

	// clean up move section of PGN
	line = regexNagAnnotation.ReplaceAllString(line, " ")
	line = regexBracketComments.ReplaceAllString(line, " ")
	line = regexReservedSymbols.ReplaceAllString(line, " ")
	// RAV variation comments can be nested - therefore loop
	for regexRavVariants.MatchString(line) {
		line = regexRavVariants.ReplaceAllString(line, " ")
	}

	// process as SAN line
	b.processSanLine(line)
}

// regular expressions for handling SAN/UCI input lines
var regexSanLineStart = regexp.MustCompile("^\\d+\\. ?")
var regexSanLineCleanUpNumbers = regexp.MustCompile("(\\d+\\.{1,3} ?)")
var regexSanLineCleanUpResults = regexp.MustCompile("(1/2|1|0)-(1/2|1|0)")
var regexWhiteSpace = regexp.MustCompile("\\s+")

// processes one line of SAN format
func (b *Book) processSanLine(line string) {
	line = strings.TrimSpace(line)

	// check if line starts valid
	found := regexSanLineStart.MatchString(line)
	if !found {
		return
	}

	/*
	 Iterate over all tokens, ignore move numbers and results
	 Example:
	 1. f4 d5 2. Nf3 Nf6 3. e3 g6 4. b3 Bg7 5. Bb2 O-O 6. Be2 c5 7. O-O Nc6 8. Ne5 Qc7 1/2-1/2
	 1. f4 d5 2. Nf3 Nf6 3. e3 Bg4 4. Be2 e6 5. O-O Bd6 6. b3 O-O 7. Bb2 c5 1/2-1/2
	*/

	// remove unnecessary parts
	line = regexSanLineCleanUpNumbers.ReplaceAllString(line, "")
	line = regexSanLineCleanUpResults.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	// split at every whitespace and iterate through items
	moveStrings := regexWhiteSpace.Split(line, -1)
	// skip lines without matches
	if len(moveStrings) == 0 {
		return
	}

	// start with root position
	pos := position.New()

	// increase counter for root position
	bookLock.Lock()
	e, found := b.bookMap[b.rootEntry]
	if found {
		e.Counter++
		b.bookMap[b.rootEntry] = e
	} else {
		panic("root entry of book map not found")
	}
	bookLock.Unlock()

	// move gen to check moves
	// movegen is not thread safe therefore we create a new instance for every line
	var mg = movegen.New()

	for _, moveString := range moveStrings {
		err := b.processSingleMove(moveString, &mg, &pos)
		// stop processing further matches when we had an error as it
		// would probably be fruitless as position will be wrong
		if err != nil {
			log.Warningf("Move not valid %s on %s", moveString, pos.StringFen())
			break
		}
	}
}

var regexUciMove = regexp.MustCompile("([a-h][1-8][a-h][1-8])([NBRQnbrq])?")
var regexSanMove = regexp.MustCompile("([NBRQK])?([a-h])?([1-8])?x?([a-h][1-8]|O-O-O|O-O)(=?([NBRQ]))?([!?+#]*)?")

// Process a single move as a string in either UCI or SAN format.
// Uses pattern matching to distinguish format
func (b *Book) processSingleMove(s string, mgPtr *movegen.Movegen, posPtr *position.Position) error {
	// find move in the current position or stop processing
	var move = types.MoveNone
	if regexUciMove.MatchString(s) {
		move = mgPtr.GetMoveFromUci(posPtr, s)
	} else if regexSanMove.MatchString(s) {
		move = mgPtr.GetMoveFromSan(posPtr, s)
	}
	// if move is invalid return stop processing further matches
	if !move.IsValid() {
		return errors.New("Invalid move " + s)
	}
	// execute move on position and store the keys for the positions
	curPosKey := uint64(posPtr.ZobristKey())
	posPtr.DoMove(move)
	nextPosKey := uint64(posPtr.ZobristKey())
	// add the move
	b.addToBook(curPosKey, nextPosKey, uint32(move))
	// no error
	return nil
}

// adds a move to the book
// this function is thread save to be used in parallel
func (b *Book) addToBook(curPosKey uint64, nextPosKey uint64, move uint32) {
	// out.Printf("Add %s to position %s\n", move.StringUci(), p.StringFen())

	// mutex to synchronize parallel access
	bookLock.Lock()
	defer bookLock.Unlock()

	// find the current position's entry
	currentPosEntry, found := b.bookMap[curPosKey]
	if !found {
		log.Error("Could not find current position in book.")
		return
	}

	// create or update book entry
	nextPosEntry, found := b.bookMap[nextPosKey]
	if found { // entry already exists - update
		nextPosEntry.Counter++
		b.bookMap[nextPosKey] = nextPosEntry
		return
	} else { // new entry
		b.bookMap[nextPosKey] = BookEntry{
			ZobristKey: nextPosKey,
			Counter:    1,
			Moves:      nil}
		nextPosEntry = b.bookMap[nextPosKey]
		// add move and link to child position entry to current entry
		currentPosEntry.Moves = append(currentPosEntry.Moves, Successor{move, nextPosEntry.ZobristKey})
		b.bookMap[curPosKey] = currentPosEntry
	}
}

func (b *Book) loadFromCache(bookPath string) (bool, error) {
	// determine cache file name
	cachePath := bookPath + ".cache"

	// Read cache file
	// Open a RO file
	decodeFile, err := os.Open(cachePath)
	if err != nil {
		return false, err
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	bookLock.Lock()
	err = decoder.Decode(&b.bookMap)
	if err != nil {
		return false, err
	}
	bookLock.Unlock()

	// set root entry key
	p := position.New()
	b.rootEntry = uint64(p.ZobristKey())

	// no error
	return true, nil
}

func (b *Book) saveToCache(bookPath string) (string, int64, error) {
	// determine cache file name
	cachePath := bookPath + ".cache"

	// Create a file for IO
	encodeFile, err := os.Create(cachePath)
	if err != nil {
		return cachePath, 0, err
	}
	// create binary encoder
	enc := gob.NewEncoder(encodeFile)

	// encode bookMap
	bookLock.Lock()
	if err = enc.Encode(b.bookMap); err != nil {
		panic(err)
	}
	bookLock.Unlock()

	// close and return
	err = encodeFile.Close()
	if err != nil {
		return cachePath, 0, err
	}

	// no error
	fileInfo, _ := os.Stat(cachePath)
	return cachePath, fileInfo.Size(), nil
}

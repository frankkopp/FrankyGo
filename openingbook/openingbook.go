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
// <p/>
// Supported formats are currently:<br/>
// BookFormat::SIMPLE for files storing a game per line with from-square and
// to-square notation<br/>
// BookFormat::SAN for files with lines of moves in SAN notation<br/>
// BookFormat::PGN for PGN formatted games<br/>
// <p/>
// TODO: As reading these formats can be slow the OpeningBook keeps a cache file where
//  it stores the serialized data of the internal book.
//
package openingbook

import (
	"bufio"
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

// BookFormat represent the supported book formats defined as constants
type BookFormat uint8

// Supported book formats
const (
	Simple BookFormat = iota
	San    BookFormat = iota
	Pgn    BookFormat = iota
)

// Pair represents a tuple of two values
// It is used in the book entry to group a move and a pointer to
// the position the move leads to
// p := Pair{"Hello", false}
type successor struct {
	move      types.Move
	nextEntry *bookEntry
}

// BookEntry represents a data structure for a move in the opening book
// data structure. It describes exactly one position defined by a zobrist
// key and has links to other entries representing moves and successor
// positions
type bookEntry struct {
	zobristKey position.Key
	counter    int
	moves      []successor
}

// Book represents a structure for chess opening books which can
// be read from different file formats into an internal data structure.
type Book struct {
	bookMap     map[position.Key]*bookEntry
	rootEntry   *bookEntry
	initialized bool
}

// to test found moves against positions
var bookLock sync.Mutex

//noinspection GoUnhandledErrorResult
func (b *Book) Initialize(bookPath string, bookFormat BookFormat) error {
	if b.initialized {
		return nil
	}
	startTotal := time.Now()

	log.Info("Initializing Opening Book.")

	// check file path
	if _, err := os.Stat(bookPath); err != nil {
		log.Errorf("File \"%s\" does not exist\n", bookPath)
		return err
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
	b.bookMap = make(map[position.Key]*bookEntry)
	b.rootEntry = &bookEntry{zobristKey: startPosition.ZobristKey(), counter: 0, moves: []successor{}}
	b.bookMap[startPosition.ZobristKey()] = b.rootEntry

	// process lines
	log.Infof("Processing %d lines with format: %v\n", len(*lines), bookFormat)
	startProcessing := time.Now()
	err = b.process(lines, bookFormat)
	if err != nil {
		log.Errorf("Error while processing: %s\n", err)
		return err
	}
	elapsedProcessing := time.Since(startProcessing)
	log.Infof("Finished processing %d lines in: %d ms\n", len(*lines), elapsedProcessing.Milliseconds())

	log.Infof("Book contains %d entries\n", len(b.bookMap))

	// finished
	elapsedTotal := time.Since(startTotal)
	log.Infof("Total initialization time : %d ms\n", elapsedTotal.Milliseconds())

	b.initialized = true
	return nil
}

// NumberOfEntries returns the number of entries in the opening book
func (b *Book) NumberOfEntries() int {
	return len(b.bookMap)
}

// GetEntry returns a pointer to the entry with the corresponding key
func (b *Book) GetEntry(key position.Key) (*bookEntry, bool) {
	entryPtr, ok := b.bookMap[key]
	if ok {
		return entryPtr, true
	} else {
		return nil, false
	}
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
		panic("not yet implemented")
	}
	return nil
}

// processes all lines of Simple format
func (b *Book) processSimple(lines *[]string) {
	// for _, line := range *lines {
	// 	b.processSimpleLine(line)
	// }
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
}

// regular expressions for detecting moves
var regexUciMove = regexp.MustCompile("([a-h][1-8][a-h][1-8])")

// processes one line of simple format and adds each move to book
func (b *Book) processSimpleLine(line string) {
	line = strings.TrimSpace(line)

	// check if line starts with a move - otherwise skip
	matches := regexUciMove.FindAllString(line, -1)

	// skip lines without matches
	if len(matches) == 0 {
		return
	}

	// start with root position
	pos := position.New()

	// increase counter for root position
	bookLock.Lock()
	b.rootEntry.counter++
	bookLock.Unlock()

	// move gen to check moves
	// movegen is not thread safe therefore we create a new instance for every line
	var mg = movegen.New()

	// add all matches to book
	for _, m := range matches {
		// find move in the current position or stop processing
		move := mg.GetMoveFromUci(&pos, m)
		if !move.IsValid() {
			// we got an invalid move and stop processing further matches
			break
		}
		// execute move on position and store the keys for the positions
		curPosKey := pos.ZobristKey()
		pos.DoMove(move)
		nextPosKey := pos.ZobristKey()
		// add the move
		b.addToBook(curPosKey, nextPosKey, move)
	}
}

// processes all lines of Simple format
func (b *Book) processSan(lines *[]string) {
	// for _, line := range *lines {
	// 	b.processSanLine(line)
	// }
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
}

// regular expressions for handling input lines
var regexSanLineStart = regexp.MustCompile("^\\d+\\. ?")
var regexSanLineCleanUpNumbers = regexp.MustCompile("(\\d+\\. ?)")
var regexSanLineCleanUpResults = regexp.MustCompile("(1/2|1|0)-(1/2|1|0)")
var regexWhiteSpace = regexp.MustCompile("\\s+")

// processes one line of simple format and adds each move to book
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
	sans := regexWhiteSpace.Split(line, -1)
	// skip lines without matches
	if len(sans) == 0 {
		return
	}

	// start with root position
	pos := position.New()

	// increase counter for root position
	bookLock.Lock()
	b.rootEntry.counter++
	bookLock.Unlock()

	// move gen to check moves
	// movegen is not thread safe therefore we create a new instance for every line
	var mg = movegen.New()

	for _, s := range sans {
		// find move in the current position or stop processing
		move := mg.GetMoveFromSan(&pos, s)
		if !move.IsValid() {
			// we got an invalid move and stop processing further matches
			break
		}
		// execute move on position and store the keys for the positions
		curPosKey := pos.ZobristKey()
		pos.DoMove(move)
		nextPosKey := pos.ZobristKey()
		// add the move
		b.addToBook(curPosKey, nextPosKey, move)
	}

}

// adds a move to the book
// this function is thread save to be used in parallel
func (b *Book) addToBook(curPosKey position.Key, nextPosKey position.Key, move types.Move) {
	// out.Printf("Add %s to position %s\n", move.StringUci(), p.StringFen())

	// mutex to synchronize parallel access
	bookLock.Lock()
	defer bookLock.Unlock()

	// find the current position's entry
	currentPosEntryPtr, found := b.bookMap[curPosKey]
	if !found {
		log.Panic("Could not find current position in book.")
		return
	}

	// create or update book entry
	nextPosEntryPtr, found := b.bookMap[nextPosKey]
	if found { // entry already exists - update
		nextPosEntryPtr.counter++
		return
	} else { // new entry
		b.bookMap[nextPosKey] = &bookEntry{
			zobristKey: nextPosKey,
			counter:    1,
			moves:      nil}
		nextPosEntryPtr = b.bookMap[nextPosKey]
		// add move and link to child position entry to current entry
		currentPosEntryPtr.moves = append(currentPosEntryPtr.moves, successor{move, nextPosEntryPtr})
	}

}

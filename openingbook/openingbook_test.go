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

package openingbook

import (
	"os"
	"path"
	"runtime"
	"testing"

	logging2 "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/util"
)

var logTest *logging2.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	out.Println("Test Main Setup Tests ====================")
	config.Setup()
	logTest = logging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func Test_readingFile(t *testing.T) {
	b := NewBook()
	lines, err := b.readFile(config.Settings.Search.BookPath+"/superbook.pgn")
	assert.NoError(t, err, "Reading file threw error: %s", err)
	assert.Equal(t, 2_620_133, len(*lines))
}

func Test_readingNonExistingFile(t *testing.T) {
	b := NewBook()
	_, err := b.readFile(config.Settings.Search.BookPath+"/abc.txt")
	assert.Error(t, err, "Reading file should throw error: %s", err)
}

func Test_processingEmpty(t *testing.T) {
	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/empty.txt", "", Simple, false, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 1, book.NumberOfEntries())

	startPos := position.NewPosition()
	entry, ok := book.GetEntry(startPos.ZobristKey())
	assert.True(t, ok)
	assert.EqualValues(t, entry.ZobristKey, startPos.ZobristKey())

	entry, ok = book.GetEntry(position.Key(1234))
	assert.False(t, ok)
	assert.True(t, entry.ZobristKey == 0)
}

func Test_processingSimpleSmall(t *testing.T) {

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/book_smalltest.txt", "", Simple, false, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 11_196, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 10, len(entry.Moves))

	// get next entry from the first found entry
	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 10, len(entry.Moves))

	// for _, p := range entry.moves {
	// 	out.Printf("%s ==> %#v (%d)\n",p.move.StringUci(), p.nextEntry.zobristKey, p.nextEntry.counter)
	// }
}

func Test_processingSimple(t *testing.T) {
	t.SkipNow()

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/book.txt", "", Simple, false, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 273_578, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 15, len(entry.Moves))
	assert.Equal(t, 61_217, entry.Counter)

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 11, len(entry.Moves))
	assert.Equal(t, 24_350, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter)
	}
}

func Test_processingSANSmall(t *testing.T) {
	logTest.Info("Starting SAN small test")

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/book_graham.txt", "", San, false, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 1_256, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 8, len(entry.Moves))
	assert.Equal(t, 149, entry.Counter)

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 8, len(entry.Moves))
	assert.Equal(t, 94, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter)
	}
}

func Test_processingPGNSmall(t *testing.T) {
	logTest.Info("Starting PGN small test")

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/pgn_test.pgn", "", Pgn, false, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 1_428, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 2, len(entry.Moves))
	assert.Equal(t, 18, entry.Counter)

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 4, len(entry.Moves))
	assert.Equal(t, 12, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter)
	}
}

func Test_processingPGNLarge(t *testing.T) {
	t.SkipNow()

	logTest.Info("Starting PGN large test")

	book := NewBook()
	logTest.Debugf("Memory statistics: %s", util.MemStat())
	err := book.Initialize(config.Settings.Search.BookPath+"/superbook.pgn", "", Pgn, false, false)
	util.GcWithStats()

	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, 4_821_316, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 20, len(entry.Moves))
	assert.Equal(t, 190_775, entry.Counter)

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.EqualValues(t, entry.ZobristKey, pos.ZobristKey())
	assert.Equal(t, 18, len(entry.Moves))
	assert.Equal(t, 89_615, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter)
	}
}

func Test_processingPGNCacheSmall(t *testing.T) {
	logTest.Info("Starting PGN cache test")

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/pgn_test.pgn", "", Pgn, true, true)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	numberOfEntries := book.NumberOfEntries()
	assert.Equal(t, 1_428, numberOfEntries)

	book.Reset()
	assert.Equal(t, 0, book.NumberOfEntries())

	err = book.Initialize(config.Settings.Search.BookPath+"/pgn_test.pgn", "", Pgn, true, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)
	assert.Equal(t, numberOfEntries, book.NumberOfEntries())

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.Equal(t, entry.ZobristKey, uint64(pos.ZobristKey()))
	assert.Equal(t, 2, len(entry.Moves))
	assert.Equal(t, 18, entry.Counter)

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, entry.ZobristKey, uint64(pos.ZobristKey()))
	assert.Equal(t, 4, len(entry.Moves))
	assert.Equal(t, 12, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter )
	}
}

func Test_processingPGNCacheLarge(t *testing.T) {
	t.SkipNow()

	logTest.Info("Starting PGN large cache test")

	book := NewBook()
	err := book.Initialize(config.Settings.Search.BookPath+"/superbook.pgn", "", Pgn, true, true)
	assert.NoError(t, err, "Initialize book threw error: %s", err)

	book.Reset()

	err = book.Initialize(config.Settings.Search.BookPath+"/superbook.pgn", "", Pgn, true, false)
	assert.NoError(t, err, "Initialize book threw error: %s", err)

	// get root entry
	pos := position.NewPosition()
	entry, found := book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, book.rootEntry, entry.ZobristKey)
	assert.Equal(t, entry.ZobristKey, uint64(pos.ZobristKey()))
	assert.Equal(t, 20, len(entry.Moves))
	assert.Equal(t, 190_775, entry.Counter)

	for _, p := range entry.Moves {
		ne, _ := book.GetEntry(position.Key(p.NextEntry))
		out.Printf("%s ==> %#v (%d)\n",Move(p.Move).StringUci(), ne.ZobristKey, ne.Counter )
	}

	pos.DoMove(CreateMove(SqE2, SqE4, Normal, PtNone))
	entry, found = book.GetEntry(pos.ZobristKey())
	assert.True(t, found)
	assert.NotNil(t, entry)
	assert.Equal(t, entry.ZobristKey, uint64(pos.ZobristKey()))
	assert.Equal(t, 18, len(entry.Moves))
	assert.Equal(t, 89_615, entry.Counter)

}

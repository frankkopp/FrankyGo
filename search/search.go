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

package search

import (
	"context"
	golog "log"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/config"
	myLogging "github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/openingbook"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/transpositiontable"
	. "github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/uciInterface"
	"github.com/frankkopp/FrankyGo/util"
)

var out = message.NewPrinter(language.German)
var log = myLogging.GetLog()

// Search represents the data structure for a chess engine search
//  Create new instance with NewSearch()
type Search struct {
	log *logging.Logger

	uciHandlerPtr  uciInterface.UciDriver
	initSemaphore  *semaphore.Weighted
	isRunning      *semaphore.Weighted
	timerWaitGroup sync.WaitGroup

	book *openingbook.Book
	tt   *transpositiontable.TtTable

	// previous search
	lastSearchResult *Result

	// current search
	stopFlag        bool
	startTime       time.Time
	hasResult       bool
	currentPosition *position.Position
	searchLimits    *Limits
	timeLimit       time.Duration
	extraTime       time.Duration
	nodesVisited    int64
	curDepth        int
	curExtraDepth   int
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// NewSearch creates a new Search instance. If the given
// uci handler is nil all output will be sent to Stdout
func NewSearch() *Search {
	s := &Search{
		log:              getSearchLog(),
		uciHandlerPtr:    nil,
		initSemaphore:    semaphore.NewWeighted(int64(1)),
		isRunning:        semaphore.NewWeighted(int64(1)),
		book:             nil,
		tt:               nil,
		lastSearchResult: nil,
		stopFlag:         false,
		startTime:        time.Time{},
		hasResult:        false,
		currentPosition:  nil,
		searchLimits:     nil,
		timeLimit:        0,
		extraTime:        0,
		nodesVisited:     0,
		curDepth:         0,
		curExtraDepth:    0,
	}
	return s
}

// NewGame stops any running searches and resets the search state
// to be ready for a different game. Any caches or states will be reset.
func (s *Search) NewGame() {
	s.StopSearch()
	s.tt.Clear()
}

// StartSearch starts the search with on the given position with
// the given search limits. Search can be stopped with StopSearch().
// Search status can be checked with IsSearching()
// This takes a copy of the position and the search limits
func (s *Search) StartSearch(p position.Position, sl Limits) {
	// acquire init phase lock
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
	// set position and searchLimits for search
	s.currentPosition = &p
	s.searchLimits = &sl
	// run search
	go s.run(&p, &sl)
	// wait until search is running and initialization
	// is done before returning to caller
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
}

// StopSearch stops a running search as quickly as possible.
// The search stops gracefully and a result will be sent to UCI.
// This will wait for the search to be stopped before returning.
func (s *Search) StopSearch() {
	s.stopFlag = true
	s.WaitWhileSearching()
}

// IsSearching checks if search is running
func (s *Search) IsSearching() bool {
	if !s.isRunning.TryAcquire(1) {
		return true
	}
	s.isRunning.Release(1)
	return false
}

// WaitWhileSearching checks if search is running and blocks until
// search has stopped
func (s *Search) WaitWhileSearching() {
	// get and release semaphore. Will block if search is running
	_ = s.isRunning.Acquire(context.TODO(), 1)
	s.isRunning.Release(1)
}

// SetUciHandler sets the UCI handler to communicate with the
// UCI user interface. If not set output will be sent to Stdout.
func (s *Search) SetUciHandler(uciHandler uciInterface.UciDriver) {
	s.uciHandlerPtr = uciHandler
}

// GetUciHandlerPtr returns the current UciHandler or nil if none is set.
func (s *Search) GetUciHandlerPtr() uciInterface.UciDriver {
	return s.uciHandlerPtr
}

// IsReady signals the uciHandler that the search is ready.
// This is part if the UCI protocol to make sure a chess
// engine is initialized and ready to receive commands.
// When called this will initialize the search which might
// take a while. When finished this will call the uciHandler
// set in SetUciHandler to send "readyok" to the UCI user interface.
func (s *Search) IsReady() {
	s.initialize()
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendReadyOk()
	} else {
		s.log.Debug("uci >> readyok")
	}
}

// ClearHash clears the transposition table.
// Is ignored with a warning while searching.
func (s *Search) ClearHash() {
	if s.IsSearching() {
		msg := "Can't clear hash while searching."
		s.uciHandlerPtr.SendInfoString(msg)
		log.Warning(msg)
		return
	}
	s.tt.Clear()
	s.uciHandlerPtr.SendInfoString("Hash cleared")
}

// ResizeCache resizes and clears the transposition table.
// Is ignored with a warning while searching.
func (s *Search) ResizeCache() {
	if s.IsSearching() {
		msg := "Can't resize hash while searching."
		s.uciHandlerPtr.SendInfoString(msg)
		log.Warning(msg)
		return
	}
	// just remove the tt pointer and re-initialize
	s.tt = nil
	s.initialize()
	// good point in time to let the garbage collector do its work
	util.GcWithStats()
	s.uciHandlerPtr.SendInfoString("Hash resized")
}


// //////////////////////////////////////////////////////
// // Private
// //////////////////////////////////////////////////////

// run is called by StartSearch() in a separate go-routine
// It runs the actual search until a search limit is reached
// or the search has been stopped by StopSearch()
func (s *Search) run(position *position.Position, sl *Limits) {
	// check if there is already a search running
	// and if not grab the isRunning semaphore
	if !s.isRunning.TryAcquire(1) {
		s.log.Error("Search already running")
		return
	}
	// release the running semaphore after the search has ended
	defer func() {
		s.isRunning.Release(1)
	}()

	// start search timer
	s.startTime = time.Now()

	// init new search run
	s.initialize()
	s.hasResult = false

	// setup and report search limits
	s.setupSearchLimits(position, sl)
	if s.searchLimits.TimeControl {
		s.startTimer()
	}

	// check for opening book move when we have a time controlled game
	bookMove := MoveNone
	if s.book != nil && config.Settings.Search.UseBook && sl.TimeControl && len(s.searchLimits.Moves) == 0 {
		bookEntry, found := s.book.GetEntry(position.ZobristKey())
		if found && len(bookEntry.Moves) > 0 {
			// choose move - random for now
			rand.Seed(int64(time.Now().Nanosecond()))
			bookMove = Move(bookEntry.Moves[rand.Intn(len(bookEntry.Moves))].Move)
			s.log.Debug("Opening Book: Choosing book move: ", bookMove.StringUci())
		}
	} else {
		s.log.Debug("Opening Book: Not using book")
	}

	// age TT entries
	if s.tt != nil {
		s.log.Debugf("Transposition Table: Using TT (%s)", s.tt.String())
		s.tt.AgeEntries()
	} else {
		s.log.Debug("Transposition Table: Not using TT")
	}

	// Initialize ply based data
	// TODO

	// release the init phase lock to signal the calling go routine
	// waiting in StartSearch() to return
	s.initSemaphore.Release(1)

	// Start the actual search with iteration deepening
	var searchResult *Result
	if bookMove == MoveNone {
		// no book move --> do search
		searchResult = s.iterativeDeepening(position)
	} else {
		// create result based on book move
		searchResult = &Result{BestMove: bookMove, BookMove: true}
	}

	// If we arrive here and the search is not stopped it means that the search
	// was finished before it has been stopped by stopSearchFlag or ponderhit,
	// We wait here until search has completed.
	if !s.stopFlag && (s.searchLimits.Ponder || s.searchLimits.Infinite) {
		s.log.Debug("Search finished before stopped or ponderhit! Waiting for stop/ponderhit to send result")
		// relaxed busy wait
		for !s.stopFlag && (s.searchLimits.Ponder || s.searchLimits.Infinite) {
			time.Sleep(5 * time.Millisecond)
		}
	}

	// update search result with search time
	searchResult.SearchTime = time.Since(s.startTime)

	// send final search info update
	// TODO

	// At the end of a search we send the result in any case even if
	// searched has been stopped. Best move is the best move so far.
	s.sendResult(searchResult)

	// save result until overwritten by the next search
	s.lastSearchResult = searchResult
	s.hasResult = true

	// print stats to log
	s.log.Info(out.Sprintf("Search finished after %d ms ", searchResult.SearchTime.Milliseconds()))
	s.log.Info(out.Sprintf("Search depth was %d(%d) with %d nodes visited. NPS = %d nps",
		s.curDepth, s.curExtraDepth, s.nodesVisited,
		(s.nodesVisited*time.Second.Nanoseconds())/(1+searchResult.SearchTime.Nanoseconds())))

	// print result to log
	s.log.Infof("Search result: %s", searchResult.String())

	// Clean up
	// make sure timer stops as this could potentially still be running
	// when search finished without any stop signal/limit
	s.stopFlag = true
}

func (s *Search) iterativeDeepening(p *position.Position) *Result {
	// FIXME: prototype/DUMMY
	for !s.stopConditions() {
		s.nodesVisited++
		if s.nodesVisited%100 == 0 {
			s.log.Info("Simulating search...")
		}
		time.Sleep(5 * time.Millisecond)
	}
	mg := movegen.NewMoveGen()
	moves := mg.GenerateLegalMoves(p, movegen.GenAll)
	rand.Seed(int64(time.Now().Nanosecond()))
	bestMove := moves.At(rand.Intn(moves.Len()))
	result := &Result{
		BestMove:    bestMove,
		PonderMove:  MoveNone,
		SearchTime:  0,
		SearchDepth: 0,
		ExtraDepth:  0,
	}
	// FIXME: prototype/DUMMY
	return result
}

// Initialize sets up opening book, transposition table
// and other potentially time consuming setup tasks
// This can be called several times without doing
// initialization again.
func (s *Search) initialize() {
	// init opening book
	if config.Settings.Search.UseBook {
		if s.book == nil {
			s.book = openingbook.NewBook()
			bookPath := config.Settings.Search.BookPath
			bookFile := config.Settings.Search.BookFile
			bookFormat, found := openingbook.FormatFromString[config.Settings.Search.BookFormat]
			if !found {
				s.log.Warningf("Book format invalid %s", config.Settings.Search.BookFormat)
				s.book = nil
			}
			err := s.book.Initialize(bookPath, bookFile, bookFormat, true, false)
			if err != nil {
				s.log.Warningf("Book could not be initialized: %s (%s)", bookPath, err)
				s.book = nil
			}
		}
	} else {
		s.log.Info("Opening book is disabled in configuration")
	}

	// init transposition table
	if config.Settings.Search.UseTT {
		if s.tt == nil {
			sizeInMByte := config.Settings.Search.TTSize
			if sizeInMByte == 0 {
				sizeInMByte = 64
			}
			s.tt = transpositiontable.NewTtTable(sizeInMByte)
		}
	} else {
		s.log.Info("Transposition Table is disabled in configuration")
	}
}

// stopConditions checks if stopFlag is set or if nodesVisited have
// reached a potential maximum set in the search limits.
func (s *Search) stopConditions() bool {
	if s.stopFlag {
		return true
	}
	if s.searchLimits.Nodes > 0 && s.nodesVisited >= s.searchLimits.Nodes {
		s.stopFlag = true
	}
	return s.stopFlag
}

// setupSearchLimits reports to log.debug on search limits for the search
// and sets up time control
func (s *Search) setupSearchLimits(position *position.Position, sl *Limits) {
	if sl.Infinite {
		s.log.Debug("Search mode: Infinite")
	}
	if sl.Ponder {
		s.log.Debug("Search mode: Ponder")
	}
	if sl.Mate > 0 {
		s.log.Debug("Search mode: Search for mate in %s", sl.Mate)
	}
	if sl.TimeControl {
		s.timeLimit = s.setupTimeControl(position, sl)
		s.extraTime = 0
		if sl.MoveTime > 0 {
			s.log.Debugf("Search mode: Time controlled: Time per move %s ms",
				out.Sprintf("%d", sl.MoveTime.Milliseconds()))
		} else {
			s.log.Debug(out.Sprintf("Search mode: Time controlled: White = %d ms (inc %d ms) Black = %d ms (inc %d ms) Moves to go: %d",
				sl.WhiteTime.Milliseconds(), sl.WhiteInc.Milliseconds(),
				sl.BlackTime.Milliseconds(), sl.BlackInc.Milliseconds(),
				sl.MovesToGo))
			s.log.Debug(out.Sprintf("Search mode: Time limit     : %d ms", s.timeLimit.Milliseconds()))
		}
	} else {
		s.log.Debug("Search mode: No time control")
	}
	if sl.Depth > 0 {
		s.log.Debugf("Search mode: Depth limited  : %d", sl.Depth)
	}
	if sl.Nodes > 0 {
		s.log.Debugf(out.Sprintf("Search mode: Nodes limited  : %d", sl.Nodes))
	}
	if sl.Moves.Len() > 0 {
		s.log.Debugf(out.Sprintf("Search mode: Moves limited  : %s", sl.Moves.StringUci()))
	}
}

// setupTimeControl sets up time control according to the given search limits
// and returns a limit on the duration for the current search
func (s *Search) setupTimeControl(p *position.Position, sl *Limits) time.Duration {
	if sl.MoveTime > 0 { // mode time per move
		return sl.MoveTime
	} else { // remaining time - estimated time per move
		// moves left
		movesLeft := int64(sl.MovesToGo)
		if movesLeft == 0 { // default
			// we estimate minimum 10 more moves in final game phases
			// in early game phases this grows up to 40
			movesLeft = int64(10 + (30 * (p.GamePhase() / GamePhaseMax)))
		}
		// time left for current player
		var timeLeft time.Duration
		switch p.NextPlayer() {
		case White:
			timeLeft = sl.WhiteTime + time.Duration(movesLeft*sl.WhiteInc.Nanoseconds())
		case Black:
			timeLeft = sl.BlackTime + time.Duration(movesLeft*sl.BlackInc.Nanoseconds())
		}
		// estimate time per move
		timeLimit := time.Duration(timeLeft.Nanoseconds() / movesLeft)
		// account for code runtime
		if timeLimit.Milliseconds() < 100 {
			// limits for very short available time reduced by another 20%
			timeLimit = time.Duration(int64(0.8 * float64(timeLimit.Nanoseconds())))
		} else {
			// reduced by 10%
			timeLimit = time.Duration(int64(0.9 * float64(timeLimit.Nanoseconds())))
		}
		return timeLimit
	}
}

// addExtraTime certain situations might call for a extension or reduction
// of the given time limit for the search. This function add/subtracts
// a portion (%) of the current time limit.
//  Example:
//  f = 1.0 --> no change in search time
//  f = 0.9 --> reduction by 10%
//  f = 1.1 --> extension by 10%
func (s *Search) addExtraTime(f float64) {
	if s.searchLimits.TimeControl && s.searchLimits.MoveTime == 0 {
		duration := time.Duration(int64((f - 1.0) * float64(s.timeLimit.Nanoseconds())))
		s.extraTime += duration
		s.log.Debugf(out.Sprintf("Time added/reduced by %d ms to %d ms",
			duration.Milliseconds(), (s.timeLimit + s.extraTime).Milliseconds()))
	}
}

// startTimer starts a go routine which regularly checks the elapsed time against
// the time limit and extra time given. If time limit is reached this will set
// the stopFlag to true and terminate the go routine
func (s *Search) startTimer() {
	go func() {
		s.log.Debugf("Timer started with time limit of %d ms", s.timeLimit.Milliseconds())
		// relaxed busy wait
		// as timeLimit changes due to extra times we can't set a fixed timeout
		for time.Since(s.startTime) < s.timeLimit+s.extraTime && !s.stopFlag {
			time.Sleep(5 * time.Millisecond)
		}
		if s.stopFlag {
			s.log.Debugf("Timer stopped early after wall time: %d ms (time limit %d ms and extra time %d)",
				time.Since(s.startTime).Milliseconds(), s.timeLimit.Milliseconds(), s.extraTime.Milliseconds())
		} else {
			s.log.Debugf("Timer stops search after wall time: %d ms (time limit %d ms and extra time %d)",
				time.Since(s.startTime).Milliseconds(), s.timeLimit.Milliseconds(), s.extraTime.Milliseconds())
			s.stopFlag = true
		}
	}()
}

// sends the search result to the uci handler if a handler is available
func (s *Search) sendResult(searchResult *Result) {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendResult(searchResult.BestMove, searchResult.PonderMove)
	}
}

// getSearchLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
// for usage in the search itself
func getSearchLog() *logging.Logger {
	searchLog := logging.MustGetLogger("search")
	standardFormat := logging.MustStringFormatter(`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`)
	backend1 := logging.NewLogBackend(os.Stdout, "", golog.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, standardFormat)
	searchBackEnd := logging.AddModuleLevel(backend1Formatter)
	searchBackEnd.SetLevel(logging.Level(config.SearchLogLevel), "")
	searchLog.SetBackend(searchBackEnd)
	return searchLog
}

// //////////////////////////////////////////////////////
// Getter and Setter
// //////////////////////////////////////////////////////

// LastSearchResult returns a copy of the last search result
func (s *Search) LastSearchResult() Result {
	return *s.lastSearchResult
}



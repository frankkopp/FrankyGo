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
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/openingbook"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/transpositiontable"
	. "github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/uciInterface"
)

var out = message.NewPrinter(language.German)
var log = logging.GetSearchLog()

// Search represents the data structure for a chess engine search
//  Create new instance with NewSearch()
type Search struct {
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

// NewGame resets the search to be ready for a different game.
// Any caches or states will be reset.
func (s *Search) NewGame() {
	// TODO: NewGame
}

// StartSearch starts the search with on the given position with
// the given search limits. Search can be stopped with StopSearch().
// Search status can be checked with IsSearching()
// This takes a copy of the position and the search limits
func (s *Search) StartSearch(p position.Position, sl Limits) {
	// acquire init phase lock
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
	// set searchLimits for instance
	s.searchLimits = &sl
	// position for this search
	s.currentPosition = &p
	// run search
	go s.run(&p, &sl)
	// wait until search is running and initialization
	// is done before returning
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
}

// StopSearch stops a running search as quickly as possible.
// The search stops gracefully and a result will be sent to
// UCI.
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
// Currently this does nothing apart from immediately send
// the ok signal to the uciHandler which in turn send "readyok"
// to the UCI user interface.
// In the future this might be used to make the UCI user interface
// wait until the search has finished initializing.
func (s *Search) IsReady() {
	s.initialize()
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendReadyOk()
	} else {
		log.Debug("uci >> readyok")
	}
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
		log.Error("Search already running")
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

	// check opening book when we have a time controlled game
	bookMove := MoveNone
	if s.book != nil && sl.TimeControl && len(s.searchLimits.Moves) == 0 {
		bookEntry, found := s.book.GetEntry(position.ZobristKey())
		if found && len(bookEntry.Moves) > 0 {
			// choose move - random for now
			rand.Seed(int64(time.Now().Nanosecond()))
			bookMove = Move(bookEntry.Moves[rand.Intn(len(bookEntry.Moves))].Move)
			log.Debug("Opening Book: Choosing book move: ", bookMove.StringUci())
		}
	} else {
		log.Debug("Opening Book: Not using book")
	}

	// age TT entries
	if s.tt != nil {
		log.Debugf("Transposition Table: Using TT (%s)", s.tt.String())
		s.tt.AgeEntries()
	} else {
		log.Debug("Transposition Table: Not using TT")
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
		log.Debug("Search finished before stopped or ponderhit! Waiting for stop/ponderhit to send result")
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
	log.Info(out.Sprintf("Search finished after %d ms ", searchResult.SearchTime.Milliseconds()))
	log.Info(out.Sprintf("Search depth was %d(%d) with %d nodes visited. NPS = %d nps",
		s.curDepth, s.curExtraDepth, s.nodesVisited,
		(s.nodesVisited*time.Second.Nanoseconds())/(1+searchResult.SearchTime.Nanoseconds())))

	// print result to log
	log.Infof("Search result: %s", searchResult.String())

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
			log.Info("Simulating search...")
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
// initialization again
func (s *Search) initialize() {
	// init opening book
	if config.Settings.Search.UseBook {
		if s.book == nil {
			s.book = openingbook.NewBook()
			bookPath := "../books/book.txt" // TODO config option
			err := s.book.Initialize(bookPath, openingbook.Simple, true, false)
			if err != nil {
				log.Warningf("Book could not be initialized: %s", bookPath)
				s.book = nil
			}
		}
	} else {
		log.Info("Opening book is disabled in configuration")
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
		log.Info("Transposition Table is disabled in configuration")
	}
}

func (s *Search) stopConditions() bool {
	if s.stopFlag {
		return true
	}
	if s.searchLimits.Nodes > 0 && s.nodesVisited >= s.searchLimits.Nodes {
		s.stopFlag = true
	}
	return s.stopFlag
}

func (s *Search) setupSearchLimits(position *position.Position, sl *Limits) {
	if sl.Infinite {
		log.Debug("Search mode: Infinite")
	}
	if sl.Ponder {
		log.Debug("Search mode: Ponder")
	}
	if sl.Mate > 0 {
		log.Debug("Search mode: Search for mate in %s", sl.Mate)
	}
	if sl.TimeControl {
		s.timeLimit = s.setupTimeControl(position, sl)
		s.extraTime = 0
		if sl.MoveTime > 0 {
			log.Debugf("Search mode: Time controlled: Time per move %s ms",
				out.Sprintf("%d", sl.MoveTime.Milliseconds()))
		} else {
			log.Debug(out.Sprintf("Search mode: Time controlled: White = %d ms (inc %d ms) Black = %d ms (inc %d ms) Moves to go: %d",
				sl.WhiteTime.Milliseconds(), sl.WhiteInc.Milliseconds(),
				sl.BlackTime.Milliseconds(), sl.BlackInc.Milliseconds(),
				sl.MovesToGo))
			log.Debug(out.Sprintf("Search mode: Time limit     : %d ms", s.timeLimit.Milliseconds()))
		}
	} else {
		log.Debug("Search mode: No time control")
	}
	if sl.Depth > 0 {
		log.Debugf("Search mode: Depth limited  : %d", sl.Depth)
	}
	if sl.Nodes > 0 {
		log.Debugf(out.Sprintf("Search mode: Nodes limited  : %d", sl.Nodes))
	}
	if sl.Moves.Len() > 0 {
		log.Debugf(out.Sprintf("Search mode: Moves limited  : %s", sl.Moves.StringUci()))
	}
}

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

func (s *Search) addExtraTime(f float64) {
	if s.searchLimits.TimeControl && s.searchLimits.MoveTime == 0 {
		duration := time.Duration(int64(f * float64(s.timeLimit.Nanoseconds())))
		s.extraTime += duration
		log.Debugf(out.Sprintf("Time added/reduced by %d ms to %d ms",
			duration.Milliseconds(), (s.timeLimit + s.extraTime).Milliseconds()))
	}
}

func (s *Search) startTimer() {
	go func() {
		log.Debugf("Timer started with time limit of %d ms", s.timeLimit.Milliseconds())
		// relaxed busy wait
		// as timeLimit changes due to extra times we can't set a fixed timeout
		for time.Since(s.startTime) < s.timeLimit+s.extraTime && !s.stopFlag {
			time.Sleep(5 * time.Millisecond)
		}
		if s.stopFlag {
			log.Debugf("Timer stopped early after wall time: %d ms (time limit %d ms and extra time %d)",
				time.Since(s.startTime).Milliseconds(), s.timeLimit.Milliseconds(), s.extraTime.Milliseconds())
		} else {
			log.Debugf("Timer stops search after wall time: %d ms (time limit %d ms and extra time %d)",
				time.Since(s.startTime).Milliseconds(), s.timeLimit.Milliseconds(), s.extraTime.Milliseconds())
			s.stopFlag = true
		}
	}()
}

func (s *Search) sendResult(searchResult *Result) {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendResult(searchResult.BestMove, searchResult.PonderMove)
	}
}

// //////////////////////////////////////////////////////
// Getter and Setter
// //////////////////////////////////////////////////////

// LastSearchResult returns a copy of the last search result
func (s *Search) LastSearchResult() Result {
	return *s.lastSearchResult
}

//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package search

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/evaluator"
	"github.com/frankkopp/FrankyGo/internal/history"
	myLogging "github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/movegen"
	"github.com/frankkopp/FrankyGo/internal/moveslice"
	"github.com/frankkopp/FrankyGo/internal/openingbook"
	"github.com/frankkopp/FrankyGo/internal/position"
	"github.com/frankkopp/FrankyGo/internal/transpositiontable"
	. "github.com/frankkopp/FrankyGo/internal/types"
	"github.com/frankkopp/FrankyGo/internal/uciInterface"
	"github.com/frankkopp/FrankyGo/internal/util"
)

var out = message.NewPrinter(language.German)

// Search represents the data structure for a chess engine search
//  Create new instance with NewSearch()
type Search struct {
	log  *logging.Logger
	slog *logging.Logger

	uciHandlerPtr uciInterface.UciDriver
	initSemaphore *semaphore.Weighted
	isRunning     *semaphore.Weighted

	book *openingbook.Book
	tt   *transpositiontable.TtTable
	eval *evaluator.Evaluator

	// history heuristics
	history *history.History

	// previous search
	lastSearchResult *Result

	// current search state
	stopFlag          bool
	startTime         time.Time
	hasResult         bool
	currentPosition   *position.Position
	searchLimits      *Limits
	timeLimit         time.Duration
	extraTime         time.Duration
	nodesVisited      uint64
	mg                []*movegen.Movegen
	pv                []*moveslice.MoveSlice
	rootMoves         *moveslice.MoveSlice
	hadBookMove       bool
	lastUciUpdateTime time.Time
	statistics        Statistics
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// NewSearch creates a new Search instance. If the given
// uci handler is nil all output will be sent to Stdout.
func NewSearch() *Search {
	s := &Search{
		log:               myLogging.GetLog(),
		slog:              getSearchTraceLog(),
		uciHandlerPtr:     nil,
		initSemaphore:     semaphore.NewWeighted(int64(1)),
		isRunning:         semaphore.NewWeighted(int64(1)),
		book:              nil,
		tt:                nil,
		eval:              evaluator.NewEvaluator(),
		history:           history.NewHistory(),
		lastSearchResult:  nil,
		stopFlag:          false,
		startTime:         time.Time{},
		hasResult:         false,
		currentPosition:   nil,
		searchLimits:      nil,
		timeLimit:         0,
		extraTime:         0,
		nodesVisited:      0,
		mg:                nil,
		pv:                nil,
		rootMoves:         nil,
		hadBookMove:       false,
		lastUciUpdateTime: time.Time{},
		statistics:        Statistics{},
	}
	return s
}

// NewGame stops any running searches and resets the search state
// to be ready for a different game. Any caches or states will be reset.
func (s *Search) NewGame() {
	s.StopSearch()
	if s.tt != nil {
		s.tt.Clear()
		s.history = history.NewHistory()
	}
}

// StartSearch starts the search on the given position with
// the given search limits. Search can be stopped with StopSearch().
// Search status can be checked with IsSearching()
// This takes a copy of the position and the search limits.
func (s *Search) StartSearch(p position.Position, sl Limits) {
	// acquire init phase lock
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
	// set position and searchLimits into the current search state
	s.currentPosition = &p
	s.searchLimits = &sl
	// run search
	go s.run(&p, &sl)
	// wait until search is running and initialization
	// is done before returning to caller
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
	s.initSemaphore.Release(1)
}

// StopSearch stops a running search as quickly as possible.
// The search stops gracefully and a result will be sent to UCI.
// This will wait for the search to be stopped before returning.
func (s *Search) StopSearch() {
	s.stopFlag = true
	s.WaitWhileSearching()
}

// PonderHit is called by the UCI user interface when the engine has
// been instructed to ponder before. The engine usually is in search
// mode without time control when this is sent. This command will
// then activate time control without interrupting the running search.
// If no search is running this has no effect.
func (s *Search) PonderHit() {
	if s.IsSearching() && s.searchLimits.Ponder {
		s.log.Debug("Ponderhit during search - activating time control")
		s.startTimer()
		return
	}
	s.log.Warning("Ponderhit received while not pondering")
}

// IsSearching checks if search is running.
func (s *Search) IsSearching() bool {
	if !s.isRunning.TryAcquire(1) {
		return true
	}
	s.isRunning.Release(1)
	return false
}

// WaitWhileSearching checks if search is running and blocks until
// search has stopped.
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
		s.sendInfoStringToUci(msg)
		s.log.Warning(msg)
		return
	}
	if s.tt != nil {
		s.tt.Clear()
		s.sendInfoStringToUci("Hash cleared")
	}
}

// ResizeCache resizes and clears the transposition table.
// Is ignored with a warning while searching.
func (s *Search) ResizeCache() {
	if s.IsSearching() {
		msg := "Can't resize hash while searching."
		s.uciHandlerPtr.SendInfoString(msg)
		s.log.Warning(msg)
		return
	}
	// just remove the tt pointer and re-initialize
	s.tt = nil
	s.initialize()
	// good point in time to let the garbage collector do its work
	s.log.Debug(util.GcWithStats())
	if s.tt != nil {
		s.uciHandlerPtr.SendInfoString(out.Sprintf("Hash resized: %s", s.tt.String()))
	}
}

// //////////////////////////////////////////////////////
// // Private
// //////////////////////////////////////////////////////

// run is called by StartSearch() in a separate goroutine
// It runs the actual search until a search limit is reached
// or the search has been stopped by StopSearch().
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

	s.log.Infof("Searching: %s", position.StringFen())

	// init new search run
	s.stopFlag = false
	s.hasResult = false
	s.timeLimit = 0
	s.extraTime = 0
	s.nodesVisited = 0
	// s.history = history.NewHistory()
	s.statistics = Statistics{}
	s.lastUciUpdateTime = s.startTime
	s.initialize()

	// setup and report search limits
	s.setupSearchLimits(position, sl)

	// when not pondering and search is time controlled start timer
	if s.searchLimits.TimeControl && !s.searchLimits.Ponder {
		s.startTimer()
	}

	// check for opening book move when we have a time controlled game
	bookMove := MoveNone
	if s.book != nil && config.Settings.Search.UseBook && sl.TimeControl {
		bookEntry, found := s.book.GetEntry(position.ZobristKey())
		if found && len(bookEntry.Moves) > 0 {
			// choose move - random for now
			rand.Seed(int64(time.Now().Nanosecond()))
			bookMove = Move(bookEntry.Moves[rand.Intn(len(bookEntry.Moves))].Move)
			s.log.Debug("Opening Book: Choosing book move: ", bookMove.StringUci())
		}
	} else {
		s.log.Info("Opening Book: Not using book")
	}

	// age TT entries
	if s.tt != nil {
		s.log.Infof("Transposition Table: Using TT (%s)", s.tt.String())
		s.tt.AgeEntries()
	} else {
		s.log.Info("Transposition Table: Not using TT")
	}

	// Initialize ply based data
	s.mg = make([]*movegen.Movegen, 0, MaxDepth+1)
	s.pv = make([]*moveslice.MoveSlice, 0, MaxDepth+1)
	for i := 0; i <= MaxDepth; i++ {
		newMoveGen := movegen.NewMoveGen()
		if config.Settings.Search.UseHistoryCounter || config.Settings.Search.UseCounterMoves {
			newMoveGen.SetHistoryData(s.history)
		}
		s.mg = append(s.mg, newMoveGen)
		s.pv = append(s.pv, moveslice.NewMoveSlice(MaxDepth+1))
	}

	s.log.Infof("Search using: PVS=%t ASP=%t MTDf=%t",
		config.Settings.Search.UsePVS,
		config.Settings.Search.UseAspiration,
		config.Settings.Search.UseMTDf)

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
		s.hadBookMove = true
	}

	// If we arrive here during Ponder of Infinite mode and the search is not
	// stopped it means that the search was finished before it has been stopped
	// by stopSearchFlag or ponderhit,
	// We wait here until search has completed.
	if (s.searchLimits.Ponder || s.searchLimits.Infinite) && !s.stopFlag {
		s.log.Debug("Search finished before stopped or ponderhit! Waiting for stop/ponderhit to send result")
		// relaxed busy wait
		for !s.stopFlag && (s.searchLimits.Ponder || s.searchLimits.Infinite) {
			time.Sleep(5 * time.Millisecond)
		}
	}

	// update search result with search time and pv
	searchResult.SearchTime = time.Since(s.startTime)
	searchResult.Pv = *s.pv[0]

	// print stats to log
	s.log.Info(out.Sprintf("Search finished after %s", searchResult.SearchTime))
	s.log.Info(out.Sprintf("Search depth was %d(%d) with %d nodes visited. NPS = %d nps",
		s.statistics.CurrentSearchDepth, s.statistics.CurrentExtraSearchDepth, s.nodesVisited,
		util.Nps(s.nodesVisited, searchResult.SearchTime)))
	s.log.Debugf("Search stats: %s", s.statistics.String())
	// s.log.Debugf("History stats: %s", s.history.String())

	// print result to log
	s.log.Infof("Search result: %s", searchResult.String())

	// save result until overwritten by the next search
	s.lastSearchResult = searchResult
	s.hasResult = true

	// Clean up
	// make sure timer stops as this could potentially still be running
	// when search finished without any stop signal/limit
	s.stopFlag = true

	// At the end of a search we send the result in any case even if
	// searched has been stopped.
	s.sendResult(searchResult)
}

// Iterative Deepening:
// It works as follows: the program starts with a one ply search,
// then increments the search depth and does another search. This
// process is repeated until the time allocated for the search is
// exhausted. In case of an unfinished search, the current best
// move is returned by the search. The current best move is
// guaranteed to by at least as good as the best move from the
// last finished iteration as we sorted the root moves before
// the start of the new iteration and started with this best
// move. This way, also the results from the partial search
// can be accepted
//
//  IDEA in case of a severe drop of the score it is wise
//   to allocate some more time, as the first alternative is often
//   a bad capture, delaying the loss instead of preventing it
func (s *Search) iterativeDeepening(position *position.Position) *Result {

	// prepare search result
	var result *Result

	// check repetition and 50 moves
	if s.checkDrawRepAnd50(position, 2) {
		msg := "Search called on DRAW by Repetition or 50-moves-rule"
		s.sendInfoStringToUci(msg)
		s.log.Warning(msg)
		result = &Result{BestValue: ValueDraw}
		return result
	}

	// generate all legal root moves
	s.rootMoves = s.mg[0].GenerateLegalMoves(position, movegen.GenAll)

	// check if there are legal moves - if not it's mate or stalemate
	if s.rootMoves.Len() == 0 {
		if position.HasCheck() {
			s.statistics.Checkmates++
			msg := "Search called on a mate position"
			s.sendInfoStringToUci(msg)
			s.log.Warning(msg)
			result = &Result{BestValue: -ValueCheckMate}
		} else {
			s.statistics.Stalemates++
			msg := "Search called on a stalemate position"
			s.sendInfoStringToUci(msg)
			s.log.Warning(msg)
			result = &Result{BestValue: ValueDraw}
		}
		return result
	}

	// add some extra time for the move after the last book move
	// hadBookMove move will be true at his point if we ever had
	// a book move.
	if s.hadBookMove && s.searchLimits.TimeControl && s.searchLimits.MoveTime == 0 {
		s.log.Debugf(out.Sprintf("First non-book move to search. Adding extra time: Before: %d ms After: %s ms",
			s.timeLimit.Milliseconds(), 2*s.timeLimit.Milliseconds()))
		s.addExtraTime(2.0)
		s.hadBookMove = false
	}

	// prepare max depth from search limits
	maxDepth := MaxDepth
	if s.searchLimits.Depth > 0 {
		maxDepth = s.searchLimits.Depth
	}

	// Max window search in preparation for aspiration window search
	// not needed yet
	alpha := ValueMin
	beta := ValueMax
	bestValue := ValueNA

	// ###########################################
	// ### BEGIN Iterative Deepening
	for iterationDepth := 0; iterationDepth < maxDepth; {
		iterationDepth++

		// update search counter
		s.nodesVisited++

		s.statistics.CurrentIterationDepth = iterationDepth
		s.statistics.CurrentSearchDepth = s.statistics.CurrentIterationDepth
		if s.statistics.CurrentExtraSearchDepth < s.statistics.CurrentIterationDepth {
			s.statistics.CurrentExtraSearchDepth = s.statistics.CurrentIterationDepth
		}

		// ###########################################
		// Start actual alpha beta search
		switch {
		// ASPIRATION SEARCH
		case config.Settings.Search.UseAspiration && iterationDepth > 3:
			bestValue = s.aspirationSearch(position, iterationDepth, bestValue)
		// MTD(f) SEARCH
		case config.Settings.Search.UseMTDf && iterationDepth > 3:
			bestValue = s.mtdf(position, iterationDepth, bestValue)
		// PVS SEARCH (or pure ALPHA BETA when PVS deactivated)
		default:
			bestValue = s.rootSearch(position, iterationDepth, alpha, beta)
		}
		// ###########################################

		// check if we need to stop
		// doing this after the first iteration ensures that
		// we have done at least one complete search and have
		// a pv (best) move
		// If we only have one move to play also stop the search
		if !s.stopConditions() && s.rootMoves.Len() > 1 {
			// sort root moves for the next iteration
			s.rootMoves.Sort()
			s.statistics.CurrentBestRootMove = s.pv[0].At(0)
			s.statistics.CurrentBestRootMoveValue = s.pv[0].At(0).ValueOf()
			// update UCI GUI
			s.sendIterationEndInfoToUci()
		} else {
			break
		}
	}
	// ### END OF Iterative Deepening
	// ###########################################

	// update searchResult here
	// best move is pv[0][0] - we need to make sure this array entry exists at this time
	// best value is pv[0][0].valueOf
	result = &Result{
		BestMove:    s.pv[0].At(0).MoveOf(),
		BestValue:   s.pv[0].At(0).ValueOf(),
		PonderMove:  MoveNone,
		SearchTime:  0,
		SearchDepth: s.statistics.CurrentIterationDepth,
		ExtraDepth:  s.statistics.CurrentExtraSearchDepth,
		BookMove:    false,
	}

	// see if we have a move we could ponder on
	if s.pv[0].Len() > 1 {
		result.PonderMove = s.pv[0].At(1).MoveOf()
	} else {
		// we do not have a ponder move in the pv list
		// so let's check the TT
		if config.Settings.Search.UseTT {
			position.DoMove(result.BestMove)
			ttEntry := s.tt.Probe(position.ZobristKey())
			if ttEntry != nil { // tt hit
				s.statistics.TTHit++
				result.PonderMove = ttEntry.Move
				s.log.Debugf(out.Sprintf("Using ponder move from hash: %s", result.PonderMove.StringUci()))
			}
		}
	}

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
// and sets up time control.
func (s *Search) setupSearchLimits(position *position.Position, sl *Limits) {
	if sl.Infinite {
		s.log.Info("Search mode: Infinite")
	}
	if sl.Ponder {
		s.log.Info("Search mode: Ponder")
	}
	if sl.Mate > 0 {
		s.log.Infof("Search mode: Search for mate in %d", sl.Mate)
	}
	if sl.TimeControl {
		s.timeLimit = s.setupTimeControl(position, sl)
		s.extraTime = 0
		if sl.MoveTime > 0 {
			s.log.Infof("Search mode: Time controlled: Time per move %s", sl.MoveTime)
		} else {
			s.log.Info(out.Sprintf("Search mode: Time controlled: White = %s (inc %s) Black = %s (inc %s) Moves to go: %d",
				sl.WhiteTime, sl.WhiteInc, sl.BlackTime, sl.BlackInc, sl.MovesToGo))
			s.log.Info(out.Sprintf("Search mode: Time limit     : %s", s.timeLimit))
		}
		if sl.Ponder {
			s.log.Info("Search mode: Ponder - time control postponed until ponderhit received")
		}
	} else {
		s.log.Info("Search mode: No time control")
	}
	if sl.Depth > 0 {
		s.log.Debugf("Search mode: Depth limited  : %d", sl.Depth)
	}
	if sl.Nodes > 0 {
		s.log.Infof(out.Sprintf("Search mode: Nodes limited  : %d", sl.Nodes))
	}
	if sl.Moves.Len() > 0 {
		s.log.Infof(out.Sprintf("Search mode: Moves limited  : %s", sl.Moves.StringUci()))
	}
}

// setupTimeControl sets up time control according to the given search limits
// and returns a limit on the duration for the current search.
func (s *Search) setupTimeControl(p *position.Position, sl *Limits) time.Duration {
	if sl.MoveTime > 0 { // mode time per move
		// we need a little room for executing the code
		duration := sl.MoveTime - (20 * time.Millisecond)
		if duration < 0 {
			s.log.Warningf("Very short move time: %s. ", sl.MoveTime)
			return sl.MoveTime
		}
		return duration
	} else { // remaining time - estimated time per move
		// moves left
		movesLeft := int64(sl.MovesToGo)
		if movesLeft == 0 { // default
			// we estimate minimum 15 more moves in final game phases
			// in early game phases this grows up to 40
			movesLeft = int64(15 + (25 * p.GamePhaseFactor()))
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
		// account for runtime of our code
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
		s.log.Debugf(out.Sprintf("Time added/reduced by %s to %s ",
			duration, s.timeLimit+s.extraTime))
	}
}

// startTimer starts a go routine which regularly checks the elapsed time against
// the time limit and extra time given. If time limit is reached this will set
// the stopFlag to true and terminate itself.
func (s *Search) startTimer() {
	go func() {
		timerStart := time.Now()
		s.log.Debugf("Timer started with time limit of %s", s.timeLimit)
		// as timeLimit changes due to extra times we can't set a fixed timeout
		// so we do a relaxed busy wait
		for time.Since(timerStart) < s.timeLimit+s.extraTime && !s.stopFlag {
			time.Sleep(5 * time.Millisecond)
		}
		if s.stopFlag {
			s.log.Debugf("Timer stopped early after wall time: %s (time limit %s and extra time %s)",
				time.Since(timerStart), s.timeLimit, s.extraTime)
		} else {
			s.log.Debugf("Timer stops search after wall time: %s (time limit %s and extra time %s)",
				time.Since(timerStart), s.timeLimit, s.extraTime)
			s.stopFlag = true
		}
	}()
}

// checks repetitions and 50-moves rule. Returns true if the position
// has repeated itself at least the given number of times.
func (s *Search) checkDrawRepAnd50(p *position.Position, i int) bool {
	if p.CheckRepetitions(i) || p.HalfMoveClock() >= 100 {
		return true
	}
	return false
}

// sends the search result to the uci handler if a handler is available.
func (s *Search) sendResult(searchResult *Result) {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendResult(searchResult.BestMove, searchResult.PonderMove)
	}
}

// sends an info string to the uci handler if a handler is available.
func (s *Search) sendInfoStringToUci(msg string) {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendInfoString(msg)
	}
}

// send UCI information about search - could be called each 500ms or so.
func (s *Search) sendSearchUpdateToUci() {
	// also do a regular search update here
	if time.Since(s.lastUciUpdateTime) > time.Second {
		s.lastUciUpdateTime = time.Now()
		hashfull := 0
		if s.tt != nil {
			hashfull = s.tt.Hashfull()
		}
		if s.uciHandlerPtr != nil {
			s.uciHandlerPtr.SendSearchUpdate(
				s.statistics.CurrentSearchDepth,
				s.statistics.CurrentExtraSearchDepth,
				s.nodesVisited,
				s.getNps(),
				time.Since(s.startTime),
				hashfull)
			s.uciHandlerPtr.SendCurrentRootMove(s.statistics.CurrentRootMove, s.statistics.CurrentRootMoveIndex)
			s.uciHandlerPtr.SendCurrentLine(s.statistics.CurrentVariation)
		} else {
			s.log.Infof(out.Sprintf("depth %d seldepth %d value %s nodes %d nps %d time %d hashful %d",
				s.statistics.CurrentSearchDepth,
				s.statistics.CurrentExtraSearchDepth,
				s.statistics.CurrentBestRootMoveValue.String(),
				s.nodesVisited,
				s.getNps(),
				time.Since(s.startTime).Milliseconds(),
				hashfull))
		}
	}
}

// send UCI information after each depth iteration.
func (s *Search) sendIterationEndInfoToUci() {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendIterationEndInfo(
			s.statistics.CurrentSearchDepth,
			s.statistics.CurrentExtraSearchDepth,
			s.statistics.CurrentBestRootMoveValue,
			s.nodesVisited,
			s.getNps(),
			time.Since(s.startTime),
			*s.pv[0])
	} else {
		s.log.Infof(out.Sprintf("depth %d seldepth %d value %s nodes %d nps %d time %d pv %s",
			s.statistics.CurrentSearchDepth,
			s.statistics.CurrentExtraSearchDepth,
			s.statistics.CurrentBestRootMoveValue.String(),
			s.nodesVisited,
			s.getNps(),
			time.Since(s.startTime).Milliseconds(),
			s.pv[0].StringUci()))
	}
}

func (s *Search) sendAspirationResearchInfo(bound string) {
	if s.uciHandlerPtr != nil {
		s.uciHandlerPtr.SendAspirationResearchInfo(
			s.statistics.CurrentSearchDepth,
			s.statistics.CurrentExtraSearchDepth,
			s.statistics.CurrentBestRootMoveValue,
			bound,
			s.nodesVisited,
			s.getNps(),
			time.Since(s.startTime),
			*s.pv[0])
	} else {
		s.log.Infof(out.Sprintf("depth %d seldepth %d value %s %s nodes %d nps %d time %d pv %s",
			s.statistics.CurrentSearchDepth,
			s.statistics.CurrentExtraSearchDepth,
			s.statistics.CurrentBestRootMoveValue.String(),
			bound,
			s.nodesVisited,
			s.getNps(),
			time.Since(s.startTime).Milliseconds(),
			s.pv[0].StringUci()))
	}
}

// helper to calculate current nps relative to s.startTime.
// limits the value to 15M to avoid very small times
// returning unrealistic values.
func (s *Search) getNps() uint64 {
	nps := util.Nps(s.nodesVisited, time.Since(s.startTime)+100)
	if nps > 15_000_000 { // sanity value for very short times
		nps = 0
	}
	return nps
}

// //////////////////////////////////////////////////////
// Getter and Setter
// //////////////////////////////////////////////////////

// LastSearchResult returns a copy of the last search result.
func (s *Search) LastSearchResult() Result {
	return *s.lastSearchResult
}

// NodesVisited returns the number of visited nodes in the last search.
func (s *Search) NodesVisited() uint64 {
	return s.nodesVisited
}

// Statistics returns a pointer to the search statistics of the last search.
func (s *Search) Statistics() *Statistics {
	return &s.statistics
}

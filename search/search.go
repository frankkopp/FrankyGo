package search

import (
	"context"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/uciInterface"
)

var out = message.NewPrinter(language.German)
var log = logging.GetLog("search")

type Search struct {
	uciHandlerPtr uciInterface.UciDriver

	initSemaphore *semaphore.Weighted
	isRunning     *semaphore.Weighted

	stopFlag  bool
	hasResult bool
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

// NewSearch creates a new Search instance. If the given
// uci handler is nil all output will be sent to Stdout
func NewSearch() *Search {
	s := &Search{
		uciHandlerPtr: nil,
		stopFlag:      false,
		initSemaphore: semaphore.NewWeighted(int64(1)),
		isRunning:     semaphore.NewWeighted(int64(1)),
		hasResult:     false,
	}
	return s
}

// NewGame resets the search to be ready for a different game.
// ANy caches or states will be reset.
func (s *Search) NewGame() {
	// TODO: NewGame
}

// StartSearch starts the search with on the given position with
// the given search limits. Search can be stopped with StopSearch().
// Search status can be checked with IsSearching()
// This takes a copy of the position and the search limits
func (s *Search) StartSearch(p position.Position, sl SearchLimits) {
	// acquire init phase lock
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
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
func (s *Search) run(position *position.Position, searchLimits *SearchLimits) {
	// check if there is already a search running
	if !s.isRunning.TryAcquire(1) {
		log.Error("Search already running")
		return
	}
	// release the running lock after the search has ended
	defer func() {
		s.isRunning.Release(1)
	}()

	// setup search limits
	// TODO setup search limits

	// release the init phase lock to signal the call waiting in
	// StartSearch() to return
	s.initSemaphore.Release(1)

	// FIXME: prototype/DUMMY
	i := 0
	for !s.stopFlag && i < 5 {
		// simulate cpu intense calculation
		f := 10000000.0
		for f > 1 {
			f /= 1.0000001
		}
		log.Info("Searching...")
		i++
	}
	// TODO: search DUMMY
}

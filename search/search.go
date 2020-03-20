package search

import (
	"context"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
)

var out = message.NewPrinter(language.German)
var log = logging.GetLog("search")

type Search struct {
	stopFlag bool
	initSemaphore *semaphore.Weighted
	isRunning *semaphore.Weighted
}

// //////////////////////////////////////////////////////
// // Public
// //////////////////////////////////////////////////////

func NewSearch() Search {
	log.Debug("Initializing search")
	s := Search{
		stopFlag:  false,
		initSemaphore: semaphore.NewWeighted(int64(1)),
		isRunning: semaphore.NewWeighted(int64(1)),
	}
	// TODO init
	return s
}

func (s *Search) NewGame() {
	log.Debug("New game")
}

func (s *Search) Start() {
	// acquire init phase lock
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
	// run search
	// FIXME: prototype
	go s.run()
	// wait until search is running and initialization
	// is done before returning
	_ = s.initSemaphore.Acquire(context.TODO(), 1)
}

func(s *Search) Stop() {
	s.stopFlag = true
	s.WaitWhileSearching()
}

// IsSearching checks if search is running
func (s * Search) IsSearching() bool {
	if !s.isRunning.TryAcquire(1) {
		return true
	}
	s.isRunning.Release(1)
	return false
}

// WaitWhileSearching checks if search is running and blocks until
// search has stopped
func (s * Search) WaitWhileSearching() {
	// get and release semaphore. Will block if search is running
	_ = s.isRunning.Acquire(context.TODO(), 1)
	s.isRunning.Release(1)
}

// //////////////////////////////////////////////////////
// // Private
// //////////////////////////////////////////////////////

func (s *Search) run() {
	// check if there is already a search running
	if !s.isRunning.TryAcquire(1) {
		log.Error("Search already running")
		return
	}
	// release the running lock after the search has ended
	defer func() {
		s.isRunning.Release(1)
	}()

	// release the init phase lock to signal the call waiting in
	// Start() to return
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


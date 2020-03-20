package search

import (

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
)

var out = message.NewPrinter(language.German)
var log = logging.GetLog("search")

// Search not implemented yet
type Search struct {
	stopFlag bool
}

func New() Search {
	log.Debug("Initializing search")
	s := Search{}
	// TODO init
	return s
}

func (u *Search) NewGame() {
	log.Info("New game")
}

func (s *Search) Start() {
	go s.run()
}

func(s *Search) Stop() {
	s.stopFlag = true
}

func (s *Search) run() {
	log.Info("Search started.")
	defer log.Info("Search stopped.")
	for !s.stopFlag {
		// simulate cpu intense calculation
		f := 100000000.0
		for f > 1 {
			f /= 1.00000001
		}
		log.Info("Still searching...")
	}
	return
}


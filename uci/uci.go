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

package uci

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/search"
)

var out = message.NewPrinter(language.German)
var log = logging.GetLog("ucihandler")
var uciLog = logging.GetUciLog()

// UciHandler handles all communication with the chess ui via UCI
// and controls options and search. Create an instance with New()
type UciHandler struct {
	InIo       *bufio.Scanner
	OutIo      *bufio.Writer
	myMoveGen  movegen.Movegen
	mySearch   search.Search
	myPosition position.Position
}

// New creates a new UciHandler instance
func New() UciHandler {
	u := UciHandler{}
	u.InIo = bufio.NewScanner(os.Stdin)
	u.OutIo = bufio.NewWriter(os.Stdout)
	u.mySearch = search.New()
	u.myMoveGen = movegen.New()
	return u
}

// Loop starts the main loop to receive commands through
// input stream (pipe or user)
func (u *UciHandler) Loop() {
	u.loop()
}

func (u *UciHandler) loop() {
	// infinite loop until "quit" command are aborted
	for {
		log.Info("Waiting for command:")

		// read from stdin or other in stream
		for u.InIo.Scan() {

			// get cmd line
			cmd := u.InIo.Text()
			strings.ToLower(cmd)
			log.Infof("Received command: %s", cmd)
			uciLog.Infof("<< %s", cmd)

			// find command and execute by calling command function
			tokens := strings.Split(cmd, " ")
			strings.TrimSpace(tokens[0])
			switch tokens[0] {
			case "quit":
				return
			case "uci":
				u.uciCommand()
			case "isready":
				u.isReadyCommand()
			case "setoption":
				u.setOptionCommand(tokens)
			case "ucinewgame":
				u.uciNewGameCommand()
			case "position":
				u.positionCommand(tokens)
			case "go":
				u.goCommand(tokens)
			case "stop":
				u.stopCommand()
			case "ponderhit":
				u.ponderHitCommand()
			case "register":
				u.registerCommand()
			case "debug":
				u.debugCommand()
			case "noop":
			default:
				log.Infof("Error: Unknown command: %s", cmd)
			}
			log.Infof("Processed command: %s", cmd)
		}

	}
}

func (u *UciHandler) ponderHitCommand() {
	msg := "Command 'ponderhit' not yet implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) stopCommand() {
	u.mySearch.Stop()
}

func (u *UciHandler) goCommand(tokens []string) {
	log.Info("Search starting...")
	u.mySearch.Start()
	log.Info("...started")
}

func (u *UciHandler) positionCommand(tokens []string) {
	// position [fen <fenstring> | startpos ]  moves <move1> .... <movei>
	//	set up the position described in fenstring on the internal board and
	//	play the moves on the internal chess board.
	//	if the game was played  from the start position the string "startpos" will be sent
	//	Note: no "new" command is needed. However, if this position is from a different game than
	//	the last position sent to the engine, the GUI should have sent a "ucinewgame" in between.
	// position startpos  moves   e2e4
	// position fen 	 8/8/...8 moves e2e4

	// build initial position
	// for i, t := range tokens[1:] {
	//
	// }

	//
	// index := 1
	// if index < len(tokens) {
	// 	switch tokens[index] {
	// 	case "startpos":
	// 		index = 2
	// 		u.myPosition = position.New()
	// 	case "fen":
	// 		index = 2
	// 		if index < len(tokens) {
	// 			u.myPosition = position.NewFen(tokens[index])
	// 			break
	// 		}
	// 		fallthrough
	// 	default:
	// 		msg := out.Sprintf("Command 'position' malformed. %s", tokens)
	// 		u.sendInfoString(msg)
	// 		log.Warning(msg)
	// 	}
	// } else { // we except a lonely position command for startpos
	// 	u.myPosition = position.New()
	// }
	// log.Debugf("New position: %s", u.myPosition.StringFen())
	//
	// // get moves from command and execute them on position
	// index++
	// if len(tokens) > 3 && tokens[3] == "moves" {
	// 	for i := 4; i < len(tokens); i++ {
	// 		move := u.myMoveGen.GetMoveFromUci(&u.myPosition, tokens[i])
	// 		if move.IsValid() {
	// 			u.myPosition.DoMove(move)
	// 			log.Debugf("Do move: %s", move.StringUci())
	// 		} else {
	// 			msg := out.Sprintf("Command 'position' malformed. Invalid moves. %s", tokens)
	// 			u.sendInfoString(msg)
	// 			log.Warning(msg)
	// 		}
	// 	}
	// }
}

func (u *UciHandler) uciNewGameCommand() {
	u.mySearch.Stop()
	u.myPosition = position.New()
	u.mySearch.NewGame()
}

func (u *UciHandler) setOptionCommand(tokens []string) {
	msg := "Command 'setoption' not yet implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) isReadyCommand() {
	u.send("readyok")
}

func (u *UciHandler) uciCommand() {
	u.send("id name FrankyGo")
	u.send("id author Frank Kopp, Germany")
	u.send("uciok")
}

func (u *UciHandler) debugCommand() {
	msg := "Command 'debug' not implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) registerCommand() {
	msg := "Command 'register' not implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) sendInfoString(s string) {
	u.send(out.Sprintf("info string %s", s))
	_, _ = u.OutIo.WriteString(s + "\n")
	_ = u.OutIo.Flush()
}

func (u *UciHandler) send(s string) {
	uciLog.Infof(">> %s", s)
	_, _ = u.OutIo.WriteString(s + "\n")
	_ = u.OutIo.Flush()
}

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

// Package uci contains the UciHandler data structure and functionality to
// handle the UCI protocol communication between the Chess User Interface
// and the chess engine.
package uci

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
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
// and controls options and search.
// Create an instance with NewUciHandler()
type UciHandler struct {
	InIo       *bufio.Scanner
	OutIo      *bufio.Writer
	myMoveGen  movegen.Movegen
	mySearch   search.Search
	myPosition position.Position
	myPerft    movegen.Perft
}

// ///////////////////////////////////////////////////////////
// Public
// ///////////////////////////////////////////////////////////

// NewUciHandler creates a new UciHandler instance.
// Input / Output io can be replaced by changing the instance's
// InIo and OutIo members.
//  Example:
// 		u.InIo = bufio.NewScanner(os.Stdin)
//		u.OutIo = bufio.NewWriter(os.Stdout)
func NewUciHandler() UciHandler {
	u := UciHandler{}
	u.InIo = bufio.NewScanner(os.Stdin)
	u.OutIo = bufio.NewWriter(os.Stdout)
	u.mySearch = search.New()
	u.myMoveGen = movegen.NewMoveGen()
	u.myPerft = movegen.Perft{}
	return u
}

// Loop starts the main loop to receive commands through
// input stream (pipe or user)
func (u *UciHandler) Loop() {
	u.loop()
}

// ///////////////////////////////////////////////////////////
// Private
// ///////////////////////////////////////////////////////////

func (u *UciHandler) loop() {
	// infinite loop until "quit" command are aborted
	for {
		log.Debugf("Waiting for command:")

		// read from stdin or other in stream
		for u.InIo.Scan() {

			// get cmd line
			cmd := u.InIo.Text()
			strings.ToLower(cmd)
			log.Debugf("Received command: %s", cmd)
			uciLog.Infof("<< %s", cmd)

			// find command and execute by calling command function
			regexWhiteSpace := regexp.MustCompile("\\s+")
			tokens := regexWhiteSpace.Split(cmd, -1)
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
			case "perft":
				u.perftCommand(tokens)
			case "noop":
			default:
				log.Warningf("Error: Unknown command: %s", cmd)
			}
			log.Debugf("Processed command: %s", cmd)
			log.Debugf("Waiting for command:")
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
	u.myPerft.Stop()
}

func (u *UciHandler) perftCommand(tokens []string) {
	depth := 4 // default
	var err error = nil
	if len(tokens) > 1 {
		depth, err = strconv.Atoi(tokens[1])
		if err != nil {
			log.Warningf("Can't perft on depth='%s'", tokens[1])
		}
	}
	depth2 := depth
	if len(tokens) > 2 {
		tmp, err := strconv.Atoi(tokens[2])
		if err != nil {
			log.Warningf("Can't use second perft depth2='%s'", tokens[2])
		} else {
			depth2 = tmp
		}
	}
	go u.myPerft.StartPerftMulti(position.StartFen, depth, depth2, true)
}

func (u *UciHandler) goCommand(tokens []string) {
	log.Info("Search starting...")
	u.mySearch.Start()
	log.Info("...started")
}

func (u *UciHandler) positionCommand(tokens []string) {
	// build initial position
	fen := position.StartFen
	i := 1
	switch tokens[i] {
	case "startpos":
		i++
	case "fen":
		i++
		var fenb strings.Builder
		for i < len(tokens) && tokens[i] != "moves" {
			fenb.WriteString(tokens[i])
			fenb.WriteString(" ")
			i++
		}
		fen = strings.TrimSpace(fenb.String())
		if len(fen) > 0 {
			break
		}
		// fen empty
		fallthrough
	default:
		msg := out.Sprintf("Command 'position' malformed. %s", tokens)
		u.sendInfoString(msg)
		log.Warning(msg)
		return
	}
	u.myPosition = position.NewPositionFen(fen)

	// check for moves to make
	if i < len(tokens) {
		if tokens[i] == "moves" {
			i++
			for i < len(tokens) && tokens[i] != "moves" {
				move := u.myMoveGen.GetMoveFromUci(&u.myPosition, tokens[i])
				if move.IsValid() {
					u.myPosition.DoMove(move)
				} else {
					msg := out.Sprintf("Command 'position' malformed. Invalid move '%s' (%s)", move.String(), tokens)
					u.sendInfoString(msg)
					log.Warning(msg)
					return
				}
				i++
			}
		} else {
			msg := out.Sprintf("Command 'position' malformed moves. %s", tokens)
			u.sendInfoString(msg)
			log.Warning(msg)
			return
		}
	}
	log.Debugf("New position: %s", u.myPosition.StringFen())
}

func (u *UciHandler) uciNewGameCommand() {
	u.mySearch.Stop()
	u.myPosition = position.NewPosition()
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

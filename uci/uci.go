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
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/movegen"
	"github.com/frankkopp/FrankyGo/moveslice"
	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/search"
	. "github.com/frankkopp/FrankyGo/types"
	"github.com/frankkopp/FrankyGo/uciInterface"
	"github.com/frankkopp/FrankyGo/version"
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
	myMoveGen  *movegen.Movegen
	mySearch   *search.Search
	myPosition *position.Position
	myPerft    *movegen.Perft
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
func NewUciHandler() *UciHandler {
	u := &UciHandler{}
	u.InIo = bufio.NewScanner(os.Stdin)
	u.OutIo = bufio.NewWriter(os.Stdout)
	u.mySearch = search.NewSearch()
	u.myPosition = position.NewPosition()
	u.myMoveGen = movegen.NewMoveGen()
	u.myPerft = movegen.NewPerft()
	var uciDriver uciInterface.UciDriver
	uciDriver = u
	u.mySearch.SetUciHandler(uciDriver)
	return u
}

// Loop starts the main loop to receive commands through
// input stream (pipe or user)
func (u *UciHandler) Loop() {
	u.loop()
}

// Command handles a single line of UCI protocol aka command.
// Returns the uci response as string output.
// Mostly useful for debugging and unit testing.
func (u *UciHandler) Command(cmd string) string {
	tmp := u.OutIo
	buffer := new(bytes.Buffer)
	u.OutIo = bufio.NewWriter(buffer)
	strings.ToLower(cmd)
	u.handleReceivedCommand(cmd)
	u.OutIo = tmp
	return buffer.String()
}

// SendReadyOk tells the UciDriver to send the uci response "readyok" to the UCI user interface
func (u *UciHandler) SendReadyOk() {
	u.send("readyok")
}

// SendInfoString send a arbitrary string to the UCI user interface
func (u *UciHandler) SendInfoString(info string) {
	u.sendInfoString(info)
}

// SendIterationEndInfo sends information about the last search depth iteration to the UCI ui
func (u *UciHandler) SendIterationEndInfo(depth int, seldepth int, value Value, nodes uint64, nps uint64, time time.Duration, pv moveslice.MoveSlice) {
	// TODO
	panic("implement me")
}

// SendAspirationResearchInfo sends information about Aspiration researches to the UCI ui
func (u *UciHandler) SendAspirationResearchInfo(depth int, seldepth int, value Value, valueType ValueType, nodes uint64, nps uint64, time time.Duration, pv moveslice.MoveSlice) {
	// TODO
	panic("implement me")
}

// SendCurrentRootMove sends the currently searched root move to the UCI ui
func (u *UciHandler) SendCurrentRootMove(currMove Move, moveNumber int) {
	// TODO
	panic("implement me")
}

// SendSearchUpdate sends a periodically update about search stats to the UCI ui
func (u *UciHandler) SendSearchUpdate(depth int, seldepth int, nodes uint64, nps uint64, time time.Duration, hashfull int) {
	// TODO
	panic("implement me")
}

// SendCurrentLine sends a periodically update about the currently searched variation ti the UCI ui
func (u *UciHandler) SendCurrentLine(moveList moveslice.MoveSlice) {
	// TODO
	panic("implement me")
}

// SendResult send the search result to the UCI ui after the search has ended are has been stopped
func (u *UciHandler) SendResult(bestMove Move, ponderMove Move) {
	var resultStr strings.Builder
	resultStr.WriteString("bestmove ")
	resultStr.WriteString(bestMove.StringUci())
	if ponderMove != MoveNone {
		resultStr.WriteString(" ")
		resultStr.WriteString(ponderMove.StringUci())
	}
	u.send(resultStr.String())
}

// ///////////////////////////////////////////////////////////
// Private
// ///////////////////////////////////////////////////////////

func (u *UciHandler) loop() {
	// infinite loop until "quit" command is received
	for {
		log.Debugf("Waiting for command:")
		// read from stdin or other in stream
		for u.InIo.Scan() {
			if u.handleReceivedCommand(u.InIo.Text()) {
				// quit command received
				return
			}
			log.Debugf("Waiting for command:")
		}
	}
}

var regexWhiteSpace = regexp.MustCompile("\\s+")

func (u *UciHandler) handleReceivedCommand(cmd string) bool {
	log.Debugf("Received command: %s", cmd)
	uciLog.Infof("<< %s", cmd)
	// find command and execute by calling command function
	tokens := regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	switch tokens[0] {
	case "quit":
		return true
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
	return false
}

func (u *UciHandler) ponderHitCommand() {
	// TODO
	msg := "Command 'ponderhit' not yet implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) stopCommand() {
	u.mySearch.StopSearch()
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
	searchLimits, err := u.readSearchLimits(tokens)
	if err {
		return
	}
	// start the search
	u.mySearch.StartSearch(*u.myPosition, *searchLimits)
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
		// fen empty fall through to err msg
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
				move := u.myMoveGen.GetMoveFromUci(u.myPosition, tokens[i])
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
	u.mySearch.StopSearch()
	u.myPosition = position.NewPosition()
	u.mySearch.NewGame()
}

func (u *UciHandler) setOptionCommand(tokens []string) {
	// TODO
	msg := "Command 'setoption' not yet implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

func (u *UciHandler) isReadyCommand() {
	u.mySearch.IsReady()
}

func (u *UciHandler) uciCommand() {
	u.send("id name FrankyGo " + version.Version())
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

func (u *UciHandler) readSearchLimits(tokens []string) (*search.Limits, bool) {
	searchLimits := search.NewSearchLimits()
	i := 1
	for i < len(tokens) {
		var err error = nil
		switch tokens[i] {
		case "moves":
			i++
			for i < len(tokens) {
				move := u.myMoveGen.GetMoveFromUci(u.myPosition, tokens[i])
				if move.IsValid() {
					searchLimits.Moves.PushBack(move)
					i++
				} else {
					break
				}
			}
		case "infinite":
			i++
			searchLimits.Infinite = true
		case "ponder":
			i++
			searchLimits.Ponder = true
		case "depth":
			i++
			searchLimits.Depth, err = strconv.Atoi(tokens[i])
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Depth value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			i++
		case "nodes":
			i++
			searchLimits.Nodes, err = strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Nodes value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			i++
		case "mate":
			i++
			searchLimits.Mate, err = strconv.Atoi(tokens[i])
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Mate value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			i++
		case "moveTime":
			i++
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. MoveTime value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.MoveTime = time.Duration(parseInt * 1_000_000)
			searchLimits.TimeControl = true
			i++
		case "wtime":
			i++
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. WhiteTime value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.WhiteTime = time.Duration(parseInt * 1_000_000)
			searchLimits.TimeControl = true
			i++
		case "btime":
			i++
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Black value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.BlackTime = time.Duration(parseInt * 1_000_000)
			searchLimits.TimeControl = true
			i++
		case "winc":
			i++
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. WhiteInc value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.WhiteInc = time.Duration(parseInt * 1_000_000)
			i++
		case "binc":
			i++
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. BlackInc value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.BlackInc = time.Duration(parseInt * 1_000_000)
			i++
		case "movestogo":
			i++
			searchLimits.MovesToGo, err = strconv.Atoi(tokens[i])
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Movestogo value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			i++
		default:
			msg := out.Sprintf("UCI command go malformed. Invalid subcommand: %s", tokens[i])
			u.sendInfoString(msg)
			log.Warning(msg)
			return nil, true
		}
	}
	// sanity check / minimum settings
	if !(searchLimits.Infinite ||
		searchLimits.Ponder ||
		searchLimits.Depth > 0 ||
		searchLimits.Nodes > 0 ||
		searchLimits.Mate > 0 ||
		searchLimits.TimeControl) {

		msg := out.Sprintf("UCI command go malformed. No effective limits set %s", tokens)
		u.sendInfoString(msg)
		log.Warning(msg)
		return nil, true
	}
	// sanity check time control
	if searchLimits.TimeControl && searchLimits.MoveTime == 0 {
		if u.myPosition.NextPlayer() == White && searchLimits.WhiteTime == 0 {
			msg := out.Sprintf("UCI command go invalid. White to move but time for white is zero! %s", tokens)
			u.sendInfoString(msg)
			log.Warning(msg)
			return nil, true
		} else if u.myPosition.NextPlayer() == Black && searchLimits.BlackTime == 0 {
			msg := out.Sprintf("UCI command go invalid. Black to move but time for white is zero! %s", tokens)
			u.sendInfoString(msg)
			log.Warning(msg)
			return nil, true
		}
	}
	return searchLimits, false
}

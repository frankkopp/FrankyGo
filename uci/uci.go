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

// Package uci contains the UciHandler data structure and functionality to
// handle the UCI protocol communication between the Chess User Interface
// and the chess engine.
package uci

import (
	"bufio"
	"bytes"
	"fmt"
	golog "log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	logging2 "github.com/op/go-logging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/config"
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
var log = logging.GetLog()

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
	uciLog     *logging2.Logger
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
	u := &UciHandler{
		InIo:       bufio.NewScanner(os.Stdin),
		OutIo:      bufio.NewWriter(os.Stdout),
		myMoveGen:  movegen.NewMoveGen(),
		mySearch:   search.NewSearch(),
		myPosition: position.NewPosition(),
		myPerft:    movegen.NewPerft(),
		uciLog:     getUciLog(),
	}
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
	u.handleReceivedCommand(cmd)
	_ = u.OutIo.Flush()
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
	u.send(fmt.Sprintf("info depth %d seldepth %d multipv 1 score %s nodes %d nps %d time %d pv %s",
		depth, seldepth, value.String(), nodes, nps, time.Milliseconds(), pv.StringUci()))
}

// SendSearchUpdate sends a periodically update about search stats to the UCI ui
func (u *UciHandler) SendSearchUpdate(depth int, seldepth int, nodes uint64, nps uint64, time time.Duration, hashfull int) {
	u.send(fmt.Sprintf("info depth %d seldepth %d nodes %d nps %d time %d hashfull %d",
		depth, seldepth, nodes, nps, time.Milliseconds(), hashfull))
}

// SendAspirationResearchInfo sends information about Aspiration researches to the UCI ui
func (u *UciHandler) SendAspirationResearchInfo(depth int, seldepth int, value Value, valueType ValueType, nodes uint64, nps uint64, time time.Duration, pv moveslice.MoveSlice) {
	// TODO
	panic("implement me")
}

// SendCurrentRootMove sends the currently searched root move to the UCI ui
func (u *UciHandler) SendCurrentRootMove(currMove Move, moveNumber int) {
	u.send(fmt.Sprintf("info currmove %s currmovenumber %d", currMove.StringUci(), moveNumber))
}

// SendCurrentLine sends a periodically update about the currently searched variation ti the UCI ui
func (u *UciHandler) SendCurrentLine(moveList moveslice.MoveSlice) {
	u.send(fmt.Sprintf("info currline %s", moveList.StringUci()))
}

// SendResult send the search result to the UCI ui after the search has ended are has been stopped
func (u *UciHandler) SendResult(bestMove Move, ponderMove Move) {
	var resultStr strings.Builder
	resultStr.WriteString("bestmove ")
	resultStr.WriteString(bestMove.StringUci())
	if ponderMove != MoveNone {
		resultStr.WriteString(" ponder ")
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
	if len(cmd) == 0 {
		return false
	}
	log.Debugf("Received command: %s", cmd)
	u.uciLog.Infof("<< %s", cmd)
	// find command and execute by calling command function
	tokens := regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	switch tokens[0] {
	case "quit":
		return true
	case "uci":
		u.uciCommand()
	case "setoption":
		u.setOptionCommand(tokens)
	case "isready":
		u.isReadyCommand()
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

func (u *UciHandler) uciCommand() {
	u.send("id name FrankyGo " + version.Version())
	u.send("id author Frank Kopp, Germany")
	options := uciOptions.GetOptions()
	for _, o := range *options {
		u.send(o)
	}
	u.send("uciok")
}

// the set option command reads the option name and the optional value
// and checks if the uci option exists. If it does its new value will
// be stored and its handler function will be called
func (u *UciHandler) setOptionCommand(tokens []string) {
	name := ""
	value := ""
	if len(tokens) > 1 && tokens[1] == "name" {
		i := 2
		for i < len(tokens) && tokens[i] != "value" {
			name += tokens[i] + " "
			i++
		}
		name = strings.TrimSpace(name)
		if len(tokens) > i && tokens[i] == "value" && len(tokens) > i+1 {
			value += tokens[i+1]
		}
	} else {
		msg := "Command 'setoption' is malformed"
		u.sendInfoString(msg)
		log.Warning(msg)
		return
	}
	o, found := uciOptions[name]
	if found {
		o.CurrentValue = value
		o.HandlerFunc(u, o)
	} else {
		msg := out.Sprintf("Command 'setoption': No such option '%s'", name)
		u.sendInfoString(msg)
		log.Warning(msg)
		return
	}
}

// requests the isready status from the Search which in turn might
// initialize itself
func (u *UciHandler) isReadyCommand() {
	u.mySearch.IsReady()
}

func (u *UciHandler) ponderHitCommand() {
	// TODO
	msg := "Command 'ponderhit' not yet implemented"
	u.sendInfoString(msg)
	log.Warning(msg)
}

// sends a stop signal to search or perft
func (u *UciHandler) stopCommand() {
	u.mySearch.StopSearch()
	u.myPerft.Stop()
}

// starts a perft test with the given depth
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

// starts a search after reading in the search limits provided
func (u *UciHandler) goCommand(tokens []string) {
	searchLimits, err := u.readSearchLimits(tokens)
	if err {
		return
	}
	// start the search
	u.mySearch.StartSearch(*u.myPosition, *searchLimits)
}

// sets the current position as given by the uci command
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

// Signals the search to stop a running search and that a new game should
// be started. Usually this means resetting all search related data e.g.
// hash tables etc.
func (u *UciHandler) uciNewGameCommand() {
	u.myPosition = position.NewPosition()
	u.mySearch.NewGame()
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
	u.uciLog.Infof(">> %s", s)
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
			parseInt, err := strconv.ParseInt(tokens[i], 10, 64)
			if err != nil {
				msg := out.Sprintf("UCI command go malformed. Nodes value not an number: %s", tokens[i])
				u.sendInfoString(msg)
				log.Warning(msg)
				return nil, true
			}
			searchLimits.Nodes = uint64(parseInt)
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

// getUciLog returns an instance of a special Logger preconfigured for
// logging all UCI protocol communication to os.Stdout or file
// Format is very simple "time UCI <uci command>"
func getUciLog() *logging2.Logger {
	// create logger
	uciLog := logging2.MustGetLogger("UCI ")

	// Stdout backend
	uciFormat := logging2.MustStringFormatter(`%{time:15:04:05.000} UCI %{message}`)
	backend1 := logging2.NewLogBackend(os.Stdout, "", golog.Lmsgprefix)
	backend1Formatter := logging2.NewBackendFormatter(backend1, uciFormat)
	uciBackEnd1 := logging2.AddModuleLevel(backend1Formatter)
	uciBackEnd1.SetLevel(logging2.DEBUG, "")

	// File backend
	programName, _ := os.Executable()
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")
	var logPath string
	if filepath.IsAbs(config.Settings.Log.LogPath) {
		logPath = config.Settings.Log.LogPath
	} else {
		dir, _ := os.Getwd()
		logPath = dir + "/" + config.Settings.Log.LogPath
	}
	uciLogFilePath := logPath + "/" + exeName + "_ucilog.log"
	uciLogFilePath = filepath.Clean(uciLogFilePath)
	// out.Printf("Log %s \n", uciLogFilePath)

	// create file backend
	var err error
	uciLogFile, err := os.OpenFile(uciLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// we use either Stdout or file - if file is valid we use only file
	if err != nil {
		golog.Println("Logfile could not be created:", err)
		uciLog.SetBackend(uciBackEnd1)
	} else {
		backend2 := logging2.NewLogBackend(uciLogFile, "", golog.Lmsgprefix)
		backend2Formatter := logging2.NewBackendFormatter(backend2, uciFormat)
		uciBackEnd2 := logging2.AddModuleLevel(backend2Formatter)
		uciBackEnd2.SetLevel(logging2.DEBUG, "")
		// multi := logging2.SetBackend(uciBackEnd1, uciBackEnd2)
		uciLog.SetBackend(uciBackEnd2)
		uciLog.Infof("Log %s started at %s:", uciLogFile.Name(), time.Now().String())
	}

	return uciLog
}

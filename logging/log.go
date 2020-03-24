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

// Package logging is a helper for the "github.com/op/go-logging" package
// to reduce the lines of code within each go file to one line.
// The functions return Logger instances which are configured with
// the necessary backends and formatters.
package logging

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/config"
)

var out = message.NewPrinter(language.German)

var (
	standardLog *logging.Logger
	searchLog   *logging.Logger
	testLog     *logging.Logger
	uciLog      *logging.Logger
	uciLogFile  *os.File

	standardFormat = logging.MustStringFormatter(`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`)

	uciLogFilePath string
)

func init() {
	programName, _ := os.Executable()
	exePath := filepath.Dir(programName)
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")
	uciLogFilePath = exePath + "/../logs/" + exeName + "_ucilog.log"

	// global loggers
	standardLog = logging.MustGetLogger("standard")
	searchLog = logging.MustGetLogger("search")
	testLog = logging.MustGetLogger("test")
	uciLog = logging.MustGetLogger("UCI ")
}

// GetLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetLog() *logging.Logger {
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, standardFormat)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	level := logging.Level(config.LogLevel)
	standardBackEnd.SetLevel(level, "")
	standardLog.SetBackend(standardBackEnd)
	return standardLog
}

// GetSearchLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
// for usage in the search itself
func GetSearchLog() *logging.Logger {
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, standardFormat)
	searchBackEnd := logging.AddModuleLevel(backend1Formatter)
	level := logging.Level(config.LogLevel)
	searchBackEnd.SetLevel(level, "")
	searchLog.SetBackend(searchBackEnd)
	return searchLog
}

// GetTestLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetTestLog() *logging.Logger {
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	format := logging.MustStringFormatter(`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	standardBackEnd.SetLevel(logging.Level(config.TestLogLevel), "")
	testLog.SetBackend(standardBackEnd)
	return testLog
}

// GetUciLog returns an instance of a special Logger preconfigured for
// logging all UCI protocol communication to os.Stdout or file
// Format is very simple "time UCI <uci command>"
func GetUciLog() *logging.Logger {
	// Stdout backend
	uciFormat := logging.MustStringFormatter(`%{time:15:04:05.000} UCI %{message}`)
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, uciFormat)
	uciBackEnd1 := logging.AddModuleLevel(backend1Formatter)
	uciBackEnd1.SetLevel(logging.DEBUG, "")

	// File backend
	var err error
	uciLogFile, err = os.OpenFile(uciLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Logfile could not be created", err)
		uciLog.SetBackend(uciBackEnd1)
	} else {
		backend2 := logging.NewLogBackend(uciLogFile, "", log.Lmsgprefix)
		backend2Formatter := logging.NewBackendFormatter(backend2, uciFormat)
		uciBackEnd2 := logging.AddModuleLevel(backend2Formatter)
		uciBackEnd2.SetLevel(logging.DEBUG, "")
		multi := logging.SetBackend(uciBackEnd1, uciBackEnd2)
		uciLog.SetBackend(multi)
	}

	return uciLog
}

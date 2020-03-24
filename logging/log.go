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

	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/config"
)

var (
	standardLog *logging.Logger
	searchLog   *logging.Logger
	testLog     *logging.Logger
	uciLog      *logging.Logger
)

// GetLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetLog() *logging.Logger {
	if standardLog == nil {
		standardLog = logging.MustGetLogger("standard")
	}
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`,
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	level := logging.Level(config.LogLevel)
	if level == -1 {
		level = 5
	}
	standardBackEnd.SetLevel(level, "")
	standardLog.SetBackend(standardBackEnd)
	return standardLog
}

// GetSearchLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
// for usage in the search itself
func GetSearchLog() *logging.Logger {
	if searchLog == nil {
		searchLog = logging.MustGetLogger("search")
	}
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`,
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	searchBackEnd := logging.AddModuleLevel(backend1Formatter)
	level := logging.Level(config.LogLevel)
	if level == -1 {
		level = 5
	}
	searchBackEnd.SetLevel(level, "")
	standardLog.SetBackend(searchBackEnd)
	return searchLog
}

// GetUciLog returns an instance of a special Logger preconfigured for
// logging all UCI protocol communication to os.Stdout or file
// Format is very simple "time UCI <uci command>"
func GetUciLog() *logging.Logger {
	if uciLog == nil {
		uciLog = logging.MustGetLogger("UCI ")
	}
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} UCI %{message}`,
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	uciBackEnd := logging.AddModuleLevel(backend1Formatter)
	uciBackEnd.SetLevel(logging.DEBUG, "")
	uciLog.SetBackend(uciBackEnd)
	return uciLog
}

// GetTestLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetTestLog() *logging.Logger {
	if testLog == nil {
		testLog = logging.MustGetLogger("test")
	}
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-7.7s}:  %{message}`,
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	standardBackEnd.SetLevel(logging.Level(config.TestLogLevel), "")
	testLog.SetBackend(standardBackEnd)
	return testLog
}

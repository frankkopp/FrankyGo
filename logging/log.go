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

	"github.com/frankkopp/FrankyGo/config"
)

var (
	standardLog *logging.Logger
	testLog     *logging.Logger

	standardFormat = logging.MustStringFormatter(`%{time:15:04:05.000} %{shortpkg:-8.8s}:%{shortfile:-14.14s} %{level:-8.8s}:  %{message}`)
)

func init() {
	// global loggers
	standardLog = logging.MustGetLogger("standard")
}

// GetLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetLog() *logging.Logger {
	// Stdout backend
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, standardFormat)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	level := logging.Level(config.LogLevel)
	standardBackEnd.SetLevel(level, "")

	// File backend
	programName, _ := os.Executable()
	exeName := strings.TrimSuffix(filepath.Base(programName), ".exe")
	var logPath string
	if filepath.IsAbs(config.Settings.Log.LogPath) {
		logPath = config.Settings.Log.LogPath
	} else {
		executable, _ := os.Executable()
		dir := filepath.Dir(executable)
		logPath = dir + "/" + config.Settings.Log.LogPath
	}
	logFilePath := logPath + "/" + exeName + "_log.log"
	logFilePath = filepath.Clean(logFilePath)

	// create file backend
	var err error
	searchLogFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// buf := bufio.NewWriter(searchLogFile)

	// we use both  - Stdout or file
	if err != nil {
		log.Println("Logfile could not be created:", err)
		standardLog.SetBackend(standardBackEnd)
	} else {
		backend2 := logging.NewLogBackend(searchLogFile, "", log.Lmsgprefix)
		backend2Formatter := logging.NewBackendFormatter(backend2, standardFormat)
		searchBackEnd2 := logging.AddModuleLevel(backend2Formatter)
		searchBackEnd2.SetLevel(logging.DEBUG, "")
		multi := logging.SetBackend(standardBackEnd, searchBackEnd2)
		standardLog.SetBackend(multi)
	}

	return standardLog
}

// GetTestLog returns an instance of a standard Logger preconfigured with a
// os.Stdout backend and a "normal" logging format (e.g. time - file - level)
func GetTestLog() *logging.Logger {
	testLog = logging.MustGetLogger("test")
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	backend1Formatter := logging.NewBackendFormatter(backend1, standardFormat)
	standardBackEnd := logging.AddModuleLevel(backend1Formatter)
	standardBackEnd.SetLevel(logging.Level(config.TestLogLevel), "")
	testLog.SetBackend(standardBackEnd)
	return testLog
}

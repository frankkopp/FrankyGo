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

package main

import (
	"flag"
	"os"
	"runtime"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/config"
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/uci"
	"github.com/frankkopp/FrankyGo/version"
)

var out = message.NewPrinter(language.German)

func main() {

	// get command line args
	versionInfo := flag.Bool("version", false, "prints version and exits")
	configFile := flag.String("config", "../config/config.toml", "path to configuration settings file")
	logLvl := flag.String("loglvl", "", "standard log level\n(critical|error|warning|notice|info|debug)")
	searchlogLvl := flag.String("searchloglvl", "", "search log level\n(critical|error|warning|notice|info|debug)")
	flag.Parse()

	// print version info
	if *versionInfo {
		out.Printf("FrankyGo %s\n", version.Version())
		out.Println("Environment:")
		out.Printf("  Using GO version %s\n", runtime.Version())
		out.Printf("  Running %s using %s as a compiler\n", runtime.GOARCH, runtime.Compiler)
		out.Printf("  Number of CPU: %d\n", runtime.NumCPU())
		out.Printf("  Number of Goroutines: %d\n", runtime.NumGoroutine())
		cwd, _ := os.Getwd()
		out.Printf("  Working directory: %s\n", cwd)
		return
	}

	// set config file
	config.ConfigFile = *configFile

	// read config file
	config.Setup()

	// set log level from cmd line options overwriting config file or defaults
	if lvl, found := config.LogLevels[*logLvl]; found {
		config.LogLevel = lvl
	}
	if lvl, found := config.LogLevels[*searchlogLvl]; found {
		config.SearchLogLevel = lvl
	}

	// resetting log levels - must be done as most packages include a logger as a
	// global var and therefore even before main() is called. These loggers start
	// with the default log level and must be reset to the actual level required.
	logging.GetLog()
	logging.GetSearchLog()
	logging.GetUciLog()

	// starting the uci handler and waiting for communication with
	// the UCI user interface
	u := uci.NewUciHandler()
	u.Loop()
}

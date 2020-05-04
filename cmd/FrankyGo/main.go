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
	"fmt"
	"os"
	"runtime"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/internal/movegen"
	"github.com/frankkopp/FrankyGo/internal/position"
	"github.com/frankkopp/FrankyGo/internal/search"
	"github.com/frankkopp/FrankyGo/internal/testsuite"
	"github.com/frankkopp/FrankyGo/internal/uci"
	"github.com/frankkopp/FrankyGo/internal/util"
	"github.com/frankkopp/FrankyGo/internal/version"
)

var out = message.NewPrinter(language.German)

func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=localhost:8080 FrankyGo_Test.exe cpu.pprof

	// command line args
	versionInfo := flag.Bool("version", false, "prints version and exits")
	configFile := flag.String("config", "./config.toml", "path to configuration settings file")
	logLvl := flag.String("loglvl", "info", "standard log level\n(critical|error|warning|notice|info|debug)")
	searchlogLvl := flag.String("searchloglvl", "", "search log level\n(critical|error|warning|notice|info|debug)")
	logPath := flag.String("logpath", "../logs", "path where to write log files to")
	bookPath := flag.String("bookpath", "../assets/books", "path to opening book files")
	bookFile := flag.String("bookfile", "", "opening book file\nprovide path if file is not in same directory as executable\nPlease also provide bookFormat otherwise this will be ignored")
	bookFormat := flag.String("bookFormat", "", "format of opening book\n(Simple|San|Pgn)")
	testSuite := flag.String("testsuite", "", "path to file containing EPD tests or folder containing EPD files")
	testMovetime := flag.Int("testtime", 2000, "search time for each test position in milliseconds")
	testSearchdepth := flag.Int("testdepth", 0, "search depth limit for each test position")
	perft := flag.Int("perft", 0, "starts perft on the start position with the given depth\nuse -fen to provide a different position")
	fen := flag.String("fen", position.StartFen, "fen for perft and nps test")
	nps := flag.Int("nps", 0, "starts nodes per second test on the start position for given amount of seconds\nuse -fen to provide a different position")
	flag.Parse()

	// print version info and exit
	if *versionInfo {
		printVersionInfo()
		return
	}

	// set config file
	// this needs to be set before config.Setup() is called. Otherwise the default will be used.
	config.ConfFile = *configFile

	// read config file
	config.Setup()

	// After reading the configuration file and the defaults we can now overwrite
	// settings with command line options.

	// path to logfile
	if *logPath != "" {
		config.Settings.Log.LogPath = *logPath
	}

	// set log level from cmd line options overwriting config file or defaults
	if lvl, found := config.LogLevels[*logLvl]; found {
		config.LogLevel = lvl
	}
	if lvl, found := config.LogLevels[*searchlogLvl]; found {
		config.SearchLogLevel = lvl
	}

	// set book path if provided as cmd line option
	if *bookPath != "" {
		config.Settings.Search.BookPath = *bookPath
	}
	if *bookFile != "" && *bookFormat != "" {
		config.Settings.Search.BookFile = *bookFile
		config.Settings.Search.BookFormat = *bookFormat
	}

	// resetting log level auf standard log - required  as most packages include
	// the standard logger as a global var and therefore even before main() is
	// called. These loggers start with the default log level and must be reset
	// to the actual level required.
	logging.GetLog()

	// nps test
	if *nps != 0 {
		config.Settings.Search.UseBook = false
		s := search.NewSearch()
		p := position.NewPosition(*fen)
		sl := search.NewSearchLimits()
		sl.TimeControl = true
		sl.MoveTime = time.Duration(*nps * int(time.Second))
		s.StartSearch(*p, *sl)
		s.WaitWhileSearching()
		out.Println()
		out.Println("NPS : ", util.Nps(s.NodesVisited(), s.LastSearchResult().SearchTime))
		return
	}

	// perft
	if *perft != 0 {
		var perftTest movegen.Perft
		for i := 1; i <= *perft; i++ {
			perftTest.StartPerft(*fen, i, true)
		}
		return
	}

	// execute test suite if command line options are given
	if *testSuite != "" {
		name := *testSuite
		fi, err := os.Stat(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			testsuite.FeatureTests(name+"/", time.Duration(*testMovetime*int(time.Millisecond)), *testSearchdepth)
		case mode.IsRegular():
			ts, _ := testsuite.NewTestSuite(name, time.Duration(*testMovetime*1_000_000), *testSearchdepth)
			ts.RunTests()
		}
		return
	}

	// starting the uci handler and waiting for communication with
	// the UCI user interface
	u := uci.NewUciHandler()
	u.Loop()
}

func printVersionInfo() {
	out.Printf("FrankyGo %s\n", version.Version())
	out.Println("Environment:")
	out.Printf("  Using GO version %s\n", runtime.Version())
	out.Printf("  Running %s using %s as a compiler\n", runtime.GOARCH, runtime.Compiler)
	out.Printf("  Number of CPU: %d\n", runtime.NumCPU())
	out.Printf("  Number of Goroutines: %d\n", runtime.NumGoroutine())
	cwd, _ := os.Getwd()
	out.Printf("  Working directory: %s\n", cwd)
}

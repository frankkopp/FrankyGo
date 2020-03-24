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

package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// globally available config values
var (
	// ConfFile hold the path to the used config file
	ConfFile = "../config/config.toml"

	// LogLevel defines the general log level - can be overwritten by cmd line options or config file
	LogLevel = 5

	// SearchLogLevel defines the search log level - can be overwritten by cmd line options or config file
	SearchLogLevel = 5

	// TestLogLevel defines the test log level
	TestLogLevel = 5

	// Settings is the global configuration read in from file
	Settings conf

	initialized = false
)

type conf struct {
	Log    logConfiguration
	Search searchConfiguration
	Eval   evalConfiguration
}

// Setup reads configuration file and sets settings from this file or defaults
// to various aspects of the application. E.g. Search config, Eval config, etc.
func Setup() {
	if initialized {
		return
	}
	// read configuration file
	if _, err := toml.DecodeFile(ConfFile, &Settings); err != nil {
		log.Fatal(err)
		return
	}
	// setup log level - first check cmd line, then config file, finally leave defaults
	setupLogLvl()
	// setup search config after reading from configuration file if necessary
	setupSearch()
	// setup eval config after reading from configuration file if necessary
	setupEval()
	initialized = true
}



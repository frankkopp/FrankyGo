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
	"fmt"

	"github.com/BurntSushi/toml"
)

// globally available config values
var (
	// LogLevel defines the general log level set by default or given by the command line arguments
	LogLevel = 2

	// SearchLogLevel defines the search log level set by default or given by the command line arguments
	SearchLogLevel = 2

	// Settings is the global configuration read in from file
	Settings conf

	initialized = false
)

type conf struct {
	Log    logConfiguration
	Search searchConfiguration
	Eval   evalConfiguration
}

func Setup() {
	if initialized {
		return
	}

	// TODO command line options
	//  config file path
	//  log levels

	// read configuration file
	if _, err := toml.DecodeFile("../config/config.toml", &Settings); err != nil {
		fmt.Println(err)
	}

	// setup log level - first check cmd line, then config file, finally leave defaults
	setupLogLvl()

	// setup search config after reading from configuration file if necessary
	setupSearch()

	// setup eval config after reading from configuration file if necessary
	setupEval()

	initialized = true
}



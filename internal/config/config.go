//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// Package config holds globally available configuration variables
// which are either set by defaults, read from a config file or set
// by command line options.
package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/frankkopp/FrankyGo/internal/util"
)

// globally available config values.
var (
	// ConfFile hold the path to the used config file (relative to working directory)
	ConfFile = "./config.toml"

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

	path, _ := util.ResolveFile(ConfFile)
	if _, err := toml.DecodeFile(path, &Settings); err != nil {
		log.Println("Config file not found. Using defaults. (", err, ")")
	}

	// setup log level - first check cmd line, then config file, finally leave defaults
	setupLogLvl()
	// setup search config after reading from configuration file if necessary
	setupSearch()
	// setup eval config after reading from configuration file if necessary
	setupEval()
	initialized = true
}

// String() prints out the current configuration settings and values.
// This uses reflection to read variables and their values.
func (settings *conf) String() string {
	var c strings.Builder
	c.WriteString("Search Config:\n")
	s := reflect.ValueOf(&settings.Search).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		c.WriteString(fmt.Sprintf("%-2d: %-22s %-6s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface()))
	}
	c.WriteString("\nEvaluation Config:\n")
	s = reflect.ValueOf(&settings.Eval).Elem()
	typeOfT = s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		c.WriteString(fmt.Sprintf("%-2d: %-22s %-6s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface()))
	}
	return c.String()
}

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

type logConfiguration struct {
	LogLvl       string
	SearchLogLvl string
}

// sets defaults which might be overwritten by config file
func init() {
	Settings.Log.LogLvl = "debug"
	Settings.Log.SearchLogLvl = "debug"
}

// set defaults for configurations here in case a configuration
// is not available from the config file
func setupLogLvl() {
	// log level
	if Settings.Log.LogLvl != "" { // check config file
		LogLevel = LogLevels[Settings.Log.LogLvl]
	}
	// search log level
	if Settings.Log.SearchLogLvl != "" { // check config file
		SearchLogLevel = LogLevels[Settings.Log.SearchLogLvl]
	}
}

// LogLevels mapping of string representations of log levels to numerical values
var LogLevels = map[string]int{
	"off":      -1,
	"critical": 0,
	"error":    1,
	"warning":  2,
	"notice":   3,
	"info":     4,
	"debug":    5,
}

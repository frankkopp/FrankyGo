/*
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

package logging

import (
	"log"
	"os"

	"github.com/op/go-logging"
)

func GetLog (name string) *logging.Logger {
	l := logging.MustGetLogger(name)
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortpkg:-8s}:%{shortfile:-14s} %{level:.7s}:  %{message}`,
		// :%{shortfunc}
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.DEBUG, "")
	l.SetBackend(backend1Leveled)
	return l
}


func GetUciLog () *logging.Logger {
	l := logging.MustGetLogger("UCI ")
	backend1 := logging.NewLogBackend(os.Stdout, "", log.Lmsgprefix)
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} UCI %{message}`,
	)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.DEBUG, "")
	l.SetBackend(backend1Leveled)
	return l
}

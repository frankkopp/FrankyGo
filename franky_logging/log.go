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

package franky_logging

import (
	"os"

	. "github.com/op/go-logging"
)

func GetLog (name string) *Logger {
	log := MustGetLogger(name)
	backend1 := NewLogBackend(os.Stdout, "", 0)
	var format = MustStringFormatter(
		`%{time:15:04:05.000} %{shortfile}:%{shortfunc} %{level:7s}:  %{message}`,
	)
	backend1Formatter := NewBackendFormatter(backend1, format)
	backend1Leveled := AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(DEBUG, "")
	SetBackend(backend1Leveled)
	return log
}

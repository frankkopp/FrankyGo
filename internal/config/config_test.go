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

package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

// make tests run in the projects root directory.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestInit(t *testing.T) {
	Setup()
	fmt.Printf("LogLvl: %v\n", Settings.Log.LogLvl)
	fmt.Printf("SearchLogLvl: %v\n", Settings.Log.SearchLogLvl)
	fmt.Printf("LogLvl set: %v\n", LogLevel)
	fmt.Printf("SearchLogLvl set: %v\n", SearchLogLevel)
	fmt.Printf("UseTT: %v\n", Settings.Search.UseTT)
	fmt.Printf("TT Size: %v\n", Settings.Search.TTSize)
	fmt.Printf("UsePawnCache: %v\n", Settings.Eval.UsePawnCache)
	fmt.Printf("PawnCacheSize: %v\n", Settings.Eval.PawnCacheSize)
}

func Test(t *testing.T) {
	Setup()
	fmt.Println(Settings.String())
}

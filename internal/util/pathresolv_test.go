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

package util

import (
	"testing"
)

func TestResolveFile(t *testing.T) {

	// file := "D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/config.toml"
	// expected := filepath.Clean("D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/config.toml")
	// resolveFile, err := ResolveFile(file)
	// assert.EqualValues(t, expected, resolveFile)
	// assert.EqualValues(t, nil, err)
	//
	// file = "./config/config.toml"
	// expected = filepath.Clean("D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/config.toml")
	// resolveFile, err = ResolveFile(file)
	// assert.EqualValues(t, expected, resolveFile)
	// assert.EqualValues(t, nil, err)

}

func TestResolveCreateFolder(t *testing.T) {
	// file := "D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/"
	// expected := filepath.Clean("D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/")
	// resolvedFolder, err := ResolveCreateFolder(file)
	// assert.EqualValues(t, expected, resolvedFolder)
	// assert.EqualValues(t, nil, err)
	//
	// file = "./config/"
	// expected = filepath.Clean("D:/_DEV/go/src/github.com/frankkopp/FrankyGo/config/")
	// resolvedFolder, err = ResolveCreateFolder(file)
	// assert.EqualValues(t, expected, resolvedFolder)
	// assert.EqualValues(t, nil, err)


	// file = "./LICENSE"
	// expected = filepath.Clean(filepath.Join(os.TempDir(), "LICENSE"))
	// resolvedFolder, err = ResolveCreateFolder(file)
	// assert.EqualValues(t, expected, resolvedFolder)
	// assert.EqualValues(t, nil, err)
	//
	// // Cleanup
	// os.Remove(expected)
}

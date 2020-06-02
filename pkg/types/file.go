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

package types

// File represents a chess board file a-h
//  FileA    File = 0
//	FileB    File = 1
//	FileC    File = 2
//	FileD    File = 3
//	FileE    File = 4
//	FileF    File = 5
//	FileG    File = 6
//	FileH    File = 7
//	FileNone File = 8
//  FileLength    = FileNone
type File uint8

// File represents a chess board file a-h
//noinspection GoUnusedConst
const (
	FileA    File = 0
	FileB    File = 1
	FileC    File = 2
	FileD    File = 3
	FileE    File = 4
	FileF    File = 5
	FileG    File = 6
	FileH    File = 7
	FileNone File = 8
	FileLength    = FileNone
)

// IsValid checks if f represents a valid file
func (f File) IsValid() bool {
	return f < FileNone
}

const fileLabels string = "abcdefgh"

// String returns a string letter for the file (e.g. a - h)
// if f is not a valid file returns "-"
func (f File) String() string {
	if f > FileH {
		return "-"
	}
	return string(fileLabels[f])
}

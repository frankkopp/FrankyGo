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

package uci

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/position"
	"github.com/frankkopp/FrankyGo/types"
)

func TestUciCommand(t *testing.T) {
	cmds := []string{
		"uci",
	}
	_,result := sendUciCmds(cmds)
	out.Print(result)
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "uciok")
}

func TestIsreadyCmd(t *testing.T) {
	cmds := []string{
		"isready",
	}
	_,result := sendUciCmds(cmds)
	out.Print(result)
	assert.Contains(t, result, "readyok")
}

func TestPositionCmd(t *testing.T) {
	cmds := []string{ // startpos
		"position startpos",
	}
	uh, result := sendUciCmds(cmds)
	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	cmds = []string{ // position with fen no moves
		"position fen "+ position.StartFen,
	}
	uh, result = sendUciCmds(cmds)
	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	cmds = []string{ // missing fen
		"position fen",
	}
	uh, result = sendUciCmds(cmds)
	out.Print(result)
	assert.Contains(t, result, "Command 'position' malformed")

	cmds = []string{ // position with fen and moves
		"position fen "+ position.StartFen +"  moves     e2e4 e7e5 g1f3 b8c6",
	}
	uh, result = sendUciCmds(cmds)
	out.Print(result)
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

	cmds = []string{ // invalid moves
		"position fen "+ position.StartFen +"  moves e7e5 g1f3 b8c6",
	}
	uh, result = sendUciCmds(cmds)
	out.Print(result)
	assert.Contains(t, result, "Command 'position' malformed")

	cmds = []string{ // position with fen and moves
		"position startpos  moves  e2e4 e7e5 g1f3 b8c6",
	}
	uh, result = sendUciCmds(cmds)
	out.Print(result)
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

}


// //////////////////////////////////////
// Helper for tests
// //////////////////////////////////////

// Takes an array of commands as strings an sends it to the UCI loop
// Captures and returns the resulting response.
func sendUciCmds(cmds []string) (*UciHandler, string) {
	types.Init()
	uh := NewUciHandler()
	uh.InIo = bufio.NewScanner(strings.NewReader(cmdString(&cmds)))
	buffer := new(bytes.Buffer)
	uh.OutIo = bufio.NewWriter(buffer)
	uh.Loop()
	result := buffer.String()
	return &uh, result
}

// Creates a command string with newline between each command
// and adding quit as the last command to stop the loop
func cmdString(cmds *[]string) string {
	var os strings.Builder
	for _, c := range *cmds {
		os.WriteString(c+"\n")
	}
	os.WriteString("quit\n")
	return os.String()
}

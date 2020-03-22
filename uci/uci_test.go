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
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/position"
)

func TestNewUciHandler(t *testing.T) {
	u := NewUciHandler()
	assert.Same(t, u, u.mySearch.GetUciHandlerPtr())
}

func TestUciCommand(t *testing.T) {
	uh := NewUciHandler()
	cmds := []string{
		"uci",
	}
	result := sendUciCmds(uh, cmds)
	out.Print(result)
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "uciok")
}

func TestIsreadyCmd(t *testing.T) {
	uh := NewUciHandler()
	cmds := []string{
		"isready",
	}
	result := sendUciCmds(uh, cmds)
	out.Print(result)
	assert.Contains(t, result, "readyok")
}

func TestPositionCmd(t *testing.T) {
	uh := NewUciHandler()
	cmds := []string{ // startpos
		"position startpos",
	}
	result := sendUciCmds(uh, cmds)
	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	cmds = []string{ // position with fen no moves
		"position fen " + position.StartFen,
	}
	result = sendUciCmds(uh, cmds)

	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	cmds = []string{ // missing fen
		"position fen",
	}
	result = sendUciCmds(uh, cmds)
	out.Print(result)
	assert.Contains(t, result, "Command 'position' malformed")

	cmds = []string{ // position with fen and moves
		"position fen " + position.StartFen + "  moves     e2e4 e7e5 g1f3 b8c6",
	}
	result = sendUciCmds(uh, cmds)
	out.Print(result)
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

	cmds = []string{ // invalid moves
		"position fen " + position.StartFen + "  moves e7e5 g1f3 b8c6",
	}
	result = sendUciCmds(uh, cmds)
	out.Print(result)
	assert.Contains(t, result, "Command 'position' malformed")

	cmds = []string{ // position with fen and moves
		"position startpos  moves  e2e4 e7e5 g1f3 b8c6",
	}
	result = sendUciCmds(uh, cmds)
	out.Print(result)
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

}

func TestReadSearchLimits(t *testing.T) {
	var cmd string
	var tokens []string
	regexWhiteSpace := regexp.MustCompile("\\s+")
	uciHandler := NewUciHandler()

	cmd = "go infinite"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err := uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.True(t, sl.Infinite)
	assert.False(t, sl.TimeControl)

	cmd = "go infinite moves e2e4 d2d4"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.True(t, sl.Infinite)
	assert.EqualValues(t, "e2e4 d2d4", sl.Moves.StringUci())
	assert.False(t, sl.TimeControl)

	cmd = "go  moves e2e4 d2d4 infinite"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.True(t, sl.Infinite)
	assert.EqualValues(t, "e2e4 d2d4", sl.Moves.StringUci())
	assert.False(t, sl.TimeControl)

	cmd = "go ponder"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.True(t, sl.Ponder)
	assert.False(t, sl.TimeControl)

	cmd = "go depth 6"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 6, sl.Depth)
	assert.False(t, sl.TimeControl)

	cmd = "go nodes 10000000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 10_000_000, sl.Nodes)
	assert.False(t, sl.TimeControl)

	cmd = "go mate 4"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 4, sl.Mate)
	assert.False(t, sl.TimeControl)

	cmd = "go depth 6 mate 4"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 4, sl.Mate)
	assert.EqualValues(t, 6, sl.Depth)
	assert.False(t, sl.TimeControl)

	cmd = "go depth mate 4"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.True(t, err)

	cmd = "go movetime 5000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.True(t, err)

	cmd = "go moveTime 5000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 5000, sl.MoveTime.Milliseconds())
	assert.True(t, sl.TimeControl)

	cmd = "go moveTime 5000 mate 6"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 5000, sl.MoveTime.Milliseconds())
	assert.EqualValues(t, 6, sl.Mate)
	assert.True(t, sl.TimeControl)

	cmd = "go moveTime 5000 depth 6 nodes 1000000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 5000, sl.MoveTime.Milliseconds())
	assert.EqualValues(t, 6, sl.Depth)
	assert.EqualValues(t, 1_000_000, sl.Nodes)
	assert.True(t, sl.TimeControl)

	cmd = "go moveTime 5000 depth 6 nodex 1000000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.True(t, err)

	cmd = "go wtime 60000 btime 60000 depth 6 nodes 1000000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 60000, sl.WhiteTime.Milliseconds())
	assert.EqualValues(t, 60000, sl.BlackTime.Milliseconds())
	assert.EqualValues(t, 6, sl.Depth)
	assert.EqualValues(t, 1_000_000, sl.Nodes)
	assert.True(t, sl.TimeControl)

	cmd = "go wtime 60000 btime 60000 winc 2000 binc 2000 depth 6 nodes 1000000"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 60000, sl.WhiteTime.Milliseconds())
	assert.EqualValues(t, 60000, sl.BlackTime.Milliseconds())
	assert.EqualValues(t, 2000, sl.WhiteInc.Milliseconds())
	assert.EqualValues(t, 2000, sl.BlackInc.Milliseconds())
	assert.EqualValues(t, 6, sl.Depth)
	assert.EqualValues(t, 1_000_000, sl.Nodes)
	assert.True(t, sl.TimeControl)

	cmd = "go wtime 60000 btime 60000 winc 2000 binc 2000 depth 6 nodes 1000000 movestogo 20 moves e2e4 d2d4 g1f3"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.False(t, err)
	assert.EqualValues(t, 60000, sl.WhiteTime.Milliseconds())
	assert.EqualValues(t, 60000, sl.BlackTime.Milliseconds())
	assert.EqualValues(t, 2000, sl.WhiteInc.Milliseconds())
	assert.EqualValues(t, 2000, sl.BlackInc.Milliseconds())
	assert.EqualValues(t, 20, sl.MovesToGo)
	assert.EqualValues(t, 6, sl.Depth)
	assert.EqualValues(t, 1_000_000, sl.Nodes)
	assert.EqualValues(t, "e2e4 d2d4 g1f3", sl.Moves.StringUci())
	assert.True(t, sl.TimeControl)

	cmd = "go winc 2000 binc 2000 movestogo 20 moves e2e4 d2d4 g1f3"
	tokens = regexWhiteSpace.Split(cmd, -1)
	strings.TrimSpace(tokens[0])
	sl, err = uciHandler.readSearchLimits(tokens)
	assert.True(t, err)
}

func TestFullSearchProcess(t *testing.T) {
	uh := NewUciHandler()

	cmds := []string{ // startpos
		"uci",
		"isready",
		"position startpos",
	}
	result := sendUciCmds(uh, cmds)
	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	cmds = []string{ // startpos
		"go wtime 60000 btime 60000 winc 2000 binc 2000 depth 6 nodes 1000000 movestogo 20 moves e2e4 d2d4 g1f3",
	}
	result = sendUciCmds(uh, cmds)

	assert.True(t, uh.mySearch.IsSearching())
	time.Sleep(2 * time.Second)

	cmds = []string{ // startpos
		"stop",
	}
	result = sendUciCmds(uh, cmds)

	uh.mySearch.WaitWhileSearching()
	assert.False(t, uh.mySearch.IsSearching())
	time.Sleep(time.Second)

	cmds = []string{ // startpos
		"stop",
	}
	result = sendUciCmds(uh, cmds)

	out.Print(result)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())
}

// //////////////////////////////////////
// Helper for tests
// //////////////////////////////////////

// Takes an array of commands as strings an sends it to the UCI loop
// Captures and returns the resulting response.
func sendUciCmds(uh *UciHandler, cmds []string) string {
	uh.InIo = bufio.NewScanner(strings.NewReader(cmdString(&cmds)))
	buffer := new(bytes.Buffer)
	uh.OutIo = bufio.NewWriter(buffer)
	uh.Loop()
	result := buffer.String()
	return result
}

// Creates a command string with newline between each command
func cmdString(cmds *[]string) string {
	var os strings.Builder
	for _, c := range *cmds {
		os.WriteString(c + "\n")
	}
	os.WriteString("quit\n")
	return os.String()
}

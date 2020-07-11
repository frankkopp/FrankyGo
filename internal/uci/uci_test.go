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

package uci

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	logging2 "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/logging"
	"github.com/frankkopp/FrankyGo/pkg/position"
)

var logTest *logging2.Logger

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	out.Println("Test Main Setup Tests ====================")
	config.Setup()
	logTest = logging.GetTestLog()
	code := m.Run()
	os.Exit(code)
}

func TestNewUciHandler(t *testing.T) {
	u := NewUciHandler()
	assert.Same(t, u, u.mySearch.GetUciHandlerPtr())
}

func TestUciHandler_Loop(t *testing.T) {
	uh := NewUciHandler()
	uh.InIo = bufio.NewScanner(strings.NewReader("uci\nquit\n"))
	buffer := new(bytes.Buffer)
	uh.OutIo = bufio.NewWriter(buffer)
	uh.Loop()
	result := buffer.String()
	assert.Contains(t, result, "uciok")
}

func TestUciCommand(t *testing.T) {
	uh := NewUciHandler()
	result := uh.Command("uci")
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "Clear Hash")
	assert.Contains(t, result, "uciok")
}

func TestIsreadyCmd(t *testing.T) {
	uh := NewUciHandler()
	result := uh.Command("isready")
	assert.Contains(t, result, "readyok")
}

func TestClearHash(t *testing.T) {
	uh := NewUciHandler()
	result := uh.Command("isready")
	assert.Contains(t, result, "readyok")
	result = uh.Command("setoption name Clear Hash")
	assert.Contains(t, result, "Hash cleared")
}

func TestResizeHash(t *testing.T) {
	uh := NewUciHandler()
	result := uh.Command("isready")
	assert.Contains(t, result, "readyok")
	result = uh.Command("setoption name Hash value 512")
	assert.Contains(t, result, "Hash resized")
}

func TestPositionCmd(t *testing.T) {
	uh := NewUciHandler()
	result := uh.Command("position startpos")
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	result = uh.Command("position fen " + position.StartFen)
	assert.EqualValues(t, position.StartFen, uh.myPosition.StringFen())

	result = uh.Command("position fen")
	assert.Contains(t, result, "Command 'position' malformed")

	result = uh.Command("position fen " + position.StartFen + "  moves     e2e4 e7e5 g1f3 b8c6")
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

	result = uh.Command("position fen " + position.StartFen + "  moves e7e5 g1f3 b8c6")
	assert.Contains(t, result, "Command 'position' malformed")

	result = uh.Command("position startpos  moves  e2e4 e7e5 g1f3 b8c6")
	assert.EqualValues(t, "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", uh.myPosition.StringFen())

}

func TestReadSearchLimits(t *testing.T) {
	var cmd string
	var tokens []string
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

	result := uh.Command("uci")
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "uciok")

	result = uh.Command("isready")
	assert.Contains(t, result, "readyok")

	result = uh.Command("setoption name Hash value 512")
	assert.Contains(t, result, "Hash resized")

	result = uh.Command("setoption name OwnBook value false")

	uh.Command("position startpos moves e2e4 e7e5")
	assert.EqualValues(t, "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2", uh.myPosition.StringFen())

	result = uh.Command("go moveTime 5000")
	assert.True(t, uh.mySearch.IsSearching())
	time.Sleep(2 * time.Second)
	uh.mySearch.WaitWhileSearching()
	assert.False(t, uh.mySearch.LastSearchResult().BookMove)

	result = uh.Command("quit")
}

func TestBookMove(t *testing.T) {
	uh := NewUciHandler()

	result := uh.Command("uci")
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "uciok")

	result = uh.Command("isready")
	assert.Contains(t, result, "readyok")

	result = uh.Command("setoption name Use_Book value true")

	uh.Command("position startpos moves e2e4 e7e5")
	assert.EqualValues(t, "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2", uh.myPosition.StringFen())

	result = uh.Command("go wtime 60000 btime 60000")
	uh.mySearch.WaitWhileSearching()
	assert.True(t, uh.mySearch.LastSearchResult().BookMove)

	result = uh.Command("quit")
}

func TestInfiniteFinishedBeforeStop(t *testing.T) {
	uh := NewUciHandler()

	result := uh.Command("uci")
	assert.Contains(t, result, "id name FrankyGo")
	assert.Contains(t, result, "uciok")

	result = uh.Command("isready")
	assert.Contains(t, result, "readyok")

	uh.Command("position startpos moves e2e4 e7e5")
	assert.EqualValues(t, "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2", uh.myPosition.StringFen())

	result = uh.Command("go infinite")
	assert.True(t, uh.mySearch.IsSearching())

	time.Sleep(3 * time.Second)

	result = uh.Command("stop")
	uh.mySearch.WaitWhileSearching()
	assert.False(t, uh.mySearch.IsSearching())

	result = uh.Command("quit")
}

func TestDebug(t *testing.T) {
	t.SkipNow()

	uh := NewUciHandler()
	var result string

	result = uh.Command("uci")
	result = uh.Command("isready")
	result = uh.Command("position startpos moves a2a3 d7d5 g1f3 g8f6 g2g3 g7g6 f1g2 f8g7 d2d3 e8g8 e1g1 c7c6 b1c3 f8e8 e2e4 e7e5 f3d2 d5e4 d2e4 f6e4 c3e4 c8e6 c1e3 d8c7 a1c1 b8d7 e4g5 e6d5 g2h3 d7f6 c2c4 d5e6 g5e6 f7e6 h3g2 a8d8 d1b3 b7b6 h2h3 f6d7 b3c2 c6c5 g3g4 d7b8 c2a4 c7d7 a4d7 d8d7 g2e4 d7c7 b2b4 c5b4 a3b4 b8c6 e4c6 c7c6 c4c5 e8d8 f1d1 c6c7 d1d2 h7h6 h3h4 b6c5 b4c5 d8b8 g4g5 b8b4 c1c4 b4c4 d3c4 h6g5 h4g5 g7f8 d2d8 g8f7 d8a8 f8c5 e3c5 c7c5 a8a7 f7g8 a7a8 g8g7 a8a7 g7f8 g1f1 e5e4 f1e2 c5c4 a7a8 f8f7 a8a7 f7f8 a7a8 f8e7 a8a7 e7d6 a7g7 d6e5 g7g6 c4c2 e2e3 c2c3 e3d2 c3f3 d2e2 f3a3 g6g8 a3a2 e2e3 a2a3 e3e2 a3a2")
	result = uh.Command("go ponder wtime 20394 btime 22577 winc 1000 binc 1000")
	assert.True(t, uh.mySearch.IsSearching())
	uh.mySearch.WaitWhileSearching()
	time.Sleep(2 * time.Second)
	result = uh.Command("quit")
	_ = result
}

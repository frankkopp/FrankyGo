package uci

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var in *bufio.Scanner
var out *bufio.Writer

func Loop() {
	in = bufio.NewScanner(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	loop()
}

func loop() {

	// infinite loop until "quit" command are aborted
	for {
		log.Println("Waiting for command:")

		// read from stdin or other in stream
		for in.Scan() {

			// get cmd line
			cmd := in.Text()
			strings.ToLower(cmd)
			log.Printf("Received command %s:", cmd)

			// find command and execute by calling command function
			tokens := strings.Split(cmd, " ")
			strings.TrimSpace(tokens[0])
			switch tokens[0] {
			case "quit":
				return
			case "uci":
				uciCommand()
			case "isready":
				isReadyCommand()
			case "setoption":
				setOptionCommand(tokens)
			case "ucinewgame":
				uciNewGameCommand()
			case "position":
				positionCommand(tokens)
			case "go":
				goCommand(tokens)
			case "stop":
				stopCommand()
			case "ponderhit":
				ponderHitCommand()
			case "register":
				registerCommand()
			case "debug":
				debugCommand()
			case "noop":
			default:
				log.Printf("Error: Unknown command %s:", cmd)
			}
			log.Printf("Processed command %s:", cmd)
		}
	}
}

func debugCommand() {

}

func registerCommand() {

}

func ponderHitCommand() {

}

func stopCommand() {

}

func goCommand(tokens []string) {

}

func positionCommand(tokens []string) {

}

func uciNewGameCommand() {

}

func setOptionCommand(tokens []string) {

}

func isReadyCommand() {

}

func uciCommand() {
	send("id name FrankyGo")
	send("id author Frank Kopp, Germany")
	send("uciok")
}

func send(s string) {
	_, _ = out.WriteString(s + "\n")
	_ = out.Flush()
}

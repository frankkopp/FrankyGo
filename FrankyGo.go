package main

import (
	"fmt"

	"github.com/frankkopp/FrankyGo/uci"
	"github.com/frankkopp/FrankyGo/version"
)

func main() {
	fmt.Printf("FrankyGo %s", version.Version())
	u := uci.NewUciHandler()
	u.Loop()
}

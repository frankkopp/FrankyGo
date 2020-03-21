package main

import (
	"fmt"

	"github.com/frankkopp/FrankyGo/uci"
)

func main() {
	fmt.Println("FrankyGo v0.3 (dev started 20.3.2020)")
	u := uci.NewUciHandler()
	u.Loop()
}

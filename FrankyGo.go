package main

import (
	"fmt"

	"github.com/frankkopp/FrankyGo/types"
)
import "github.com/frankkopp/FrankyGo/uci"

func main() {
	fmt.Println("FrankyGo v0.3 (dev started 20.3.2020)")

	types.Init()
	u := uci.NewUciHandler()
	u.Loop()
}

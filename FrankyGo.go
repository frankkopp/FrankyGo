package main

import (
	"fmt"

	"github.com/frankkopp/FrankyGo/types"
)
import "github.com/frankkopp/FrankyGo/uci"

func main() {
	fmt.Println("FrankyGo")

	types.Init()
	u := uci.New()
	u.Loop()
}

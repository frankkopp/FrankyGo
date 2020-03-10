package search

import "log"

var stopFlag = false

func Start() {
	dosomething()
}

func Stop() {
	stopFlag = true
}

func dosomething() {
	log.Printf("Search started.")
	defer log.Printf("Search stopped.")
	for !stopFlag {
		// simulate cpu intense calculation
		f := 100000000.0
		for f > 1 {
			f /= 1.00000001
		}
		log.Printf("Still searching...")
	}
	return
}

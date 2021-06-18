package main

import "log"

func main() {
	for i := 0; i < 3; i++ {
		log.Printf("Hello, WASM %d!", i)
	}
}

package main

import (
	"fmt"
	"modules/extract"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "data path")
		os.Exit(2)
	}

	config := extract.Config{OsmPath: os.Args[1]}
	extract.Run(config)
}

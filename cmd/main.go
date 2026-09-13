package main

import (
	"fmt"
	"github/idbeholdv18/simple-memory/internal/pm"
	"os"
)

func main() {
	mem := pm.NewMemory(1024)
	err := mem.Write(1024, 0xAB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	println(mem.Read(0x10))
}

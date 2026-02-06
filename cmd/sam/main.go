package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
	printHelp()
	os.Exit(1)
}

func printHelp() {
	fmt.Println("Sam — your personal assistant")
	fmt.Println()
	fmt.Println("Usage: sam <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  (no capabilities registered yet)")
}

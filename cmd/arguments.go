package cmd

import (
	"flag"
	"fmt"
	"os"
)

func FileNameFromArguments(cmdName string) string {
	help := flag.Bool("h", false, "Display this help message and exit")

	flag.Parse()

	if *help {
		printHelp(cmdName)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		printHelp(cmdName)
		os.Exit(1)
	}

	return flag.Arg(0)
}

func printHelp(cmdName string) {
	fmt.Printf("Usage: %s [OPTIONS] FILE\n", cmdName)
	fmt.Println("Options:")
	fmt.Println("  -h    Display this help message and exit")
}

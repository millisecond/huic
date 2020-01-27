package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// set up the main commands

	getCommand := flag.NewFlagSet("get", flag.ExitOnError)

	benchmarkCommand := flag.NewFlagSet("benchmark", flag.ExitOnError)
	//protocolsPtr := benchmarkCommand.String("protocols", "all", "Protocols {huic|http|all}. (Default: all)")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("get or benchmark subcommand is required")
		os.Exit(1)
	}

	// Parse the flags for appropriate FlagSet
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "get":
		getCommand.Parse(os.Args[2:])
	case "benchmark":
		benchmarkCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if benchmarkCommand.Parsed() {
		fmt.Println("benchmark", URLOrExit(benchmarkCommand))
	}

	if getCommand.Parsed() {
		fmt.Println("get", URLOrExit(getCommand))
	}
}

func URLOrExit(flags *flag.FlagSet) string {
	url := flags.Arg(0)
	if len(url) != 0 {
		return url
	}
	fmt.Println("URL required")
	os.Exit(1)
	return ""
}

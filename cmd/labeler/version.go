package main

import (
	"flag"
	"fmt"
)

const version = "1.2.0"

func VersionCommand(args []string) error {
	fs := flag.NewFlagSet("version", flag.ExitOnError)
	fs.Parse(args)

	fmt.Println(version)

	return nil
}

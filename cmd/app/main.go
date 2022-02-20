package main

import (
	"fmt"
	"os"

	"Gartenschlaeger/github-labeler/pkg/colors"
)

type CommandHandler func([]string) error

var cmdsMap map[string]CommandHandler

func printUsage() {
	fmt.Printf("%vlabeler <command> [flags]%v\n", colors.Cyan, colors.Reset)
	os.Exit(1)
}

func printUnknownCommand(commandName string) {
	fmt.Printf("%vThe '%s' command is not a known command.%v\n", colors.Red, commandName, colors.Reset)
	os.Exit(1)
}

func init() {
	cmdsMap = make(map[string]CommandHandler)
	cmdsMap["version"] = VersionCommand
	cmdsMap["merge"] = MergeCommand
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	cmdName := os.Args[1]
	if h, f := cmdsMap[cmdName]; f {
		err := h(os.Args[2:])
		if err != nil {
			fmt.Printf("%v%v%v\n", colors.Red, err, colors.Reset)
		}
	} else {
		printUnknownCommand(cmdName)
	}
}

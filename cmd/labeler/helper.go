package main

import (
	"Gartenschlaeger/github-labeler/pkg/colors"
	"fmt"
	"os"
	"strings"
)

// hasOptionalFlag checks if a flag is passed
func hasOptionalBoolFlag(value *bool) bool {
	return value != nil && *value
}

// requireFlag checks if a required flag exists
func requireFlag(value *string, errorMsg string) {
	if value == nil || strings.TrimSpace(*value) == "" {
		fmt.Printf("%v%s%v\n", colors.Yellow, errorMsg, colors.Reset)
		os.Exit(1)
	}
}

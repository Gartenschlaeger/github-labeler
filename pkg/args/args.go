package args

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"Gartenschlaeger/github-labeler/pkg/cli"
)

// Arguments holds flag values
type Arguments struct {
	Token      *string
	Owner      *string
	Repository *string
	IsDryMode  bool
	SkipDelete bool
}

// hasOptionalFlag checks if a flag is passed
func hasOptionalBoolFlag(value *bool) bool {
	return value != nil && *value
}

// requireFlag checks if a required flag exists
func requireFlag(value *string, errorMsg string) {
	if value == nil || strings.TrimSpace(*value) == "" {
		fmt.Printf("%v%s%v\n", cli.Yellow, errorMsg, cli.Reset)
		os.Exit(1)
	}
}

// validateFlags validates for valid flags
func validateFlags(args *Arguments) {
	requireFlag(args.Token, "Token required. Use -t <token>")
	requireFlag(args.Owner, "Owner required. Use -o <owner>")
	requireFlag(args.Repository, "Repository required. Use -r <repository>")
}

// Parse parsed flags from command line input
func Parse() *Arguments {
	var args *Arguments = &Arguments{}
	args.Token = flag.String("t", os.Getenv("LABELER_TOKEN"), "Bearer token for Github API requests.")
	args.Owner = flag.String("o", os.Getenv("LABELER_OWNER"), "Github Owner")
	args.Repository = flag.String("r", os.Getenv("LABELER_REPO"), "Github repository name")

	dryMode := flag.Bool("dry-mode", false, "Enable dry mode")
	skipDelete := flag.Bool("skip-delete", false, "Skip deletion of unknown labels.")

	flag.Parse()

	args.IsDryMode = hasOptionalBoolFlag(dryMode)
	args.SkipDelete = hasOptionalBoolFlag(skipDelete)

	validateFlags(args)

	return args
}

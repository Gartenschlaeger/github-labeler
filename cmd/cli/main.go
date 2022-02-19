package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"Gartenschlaeger/github-labeler/pkg/githubapi"
)

var token, owner, repo *string

func checkFlagValue(value string, errorMsg string) {
	if strings.TrimSpace(value) == "" {
		fmt.Println(errorMsg)
		os.Exit(1)
	}
}

func parseFlags() {
	token = flag.String("t", os.Getenv("LABELER_TOKEN"), "Bearer token for Github API requests.")
	owner = flag.String("o", os.Getenv("LABELER_OWNER"), "Github Owner")
	repo = flag.String("r", os.Getenv("LABELER_REPO"), "Github repository name")
	flag.Parse()

	checkFlagValue(*token, "Token required. Use -t <token>")
	checkFlagValue(*owner, "Owner required. Use -o <owner>")
	checkFlagValue(*repo, "Repository required. Use -r <repository>")
}

func main() {
	parseFlags()

	githubapi.SetBearerToken(*token)

	labels := githubapi.GetLabelsForRepository(*owner, *repo)
	fmt.Println(labels)
}

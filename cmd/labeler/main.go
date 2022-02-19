package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"Gartenschlaeger/github-labeler/pkg/cli"
	"Gartenschlaeger/github-labeler/pkg/githubapi"
	"Gartenschlaeger/github-labeler/pkg/types"
)

var token, owner, repo *string

func checkFlagValue(value string, errorMsg string) {
	if strings.TrimSpace(value) == "" {
		fmt.Printf("%v%s%v\n", cli.Yellow, errorMsg, cli.Reset)
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

func readLabelsDefinitions() (*map[string]types.LabelDefinition, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fullPath := path.Join(homePath, ".config", "github-labeler", "labels.json")

	fileData, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	labelSet := []types.LabelDefinition{}
	json.Unmarshal(fileData, &labelSet)

	if len(labelSet) == 0 {
		return nil, errors.New("list of expected labels are empty")
	}

	m := make(map[string]types.LabelDefinition)
	for i := 0; i < len(labelSet); i++ {
		m[labelSet[i].Name] = labelSet[i]
	}

	return &m, nil
}

func convertGithubLabelsToMap(labels *[]githubapi.GithubLabelResponse) *map[string]githubapi.GithubLabelResponse {
	m := make(map[string]githubapi.GithubLabelResponse)
	for i := 0; i < len(*labels); i++ {
		l := (*labels)[i]
		m[l.Name] = l
	}

	return &m
}

func main() {
	parseFlags()

	githubapi.SetBearerToken(*token)

	definedLabels, err := readLabelsDefinitions()
	if err != nil {
		panic(err)
	}

	rl, err := githubapi.GetLabelsForRepository(*owner, *repo)
	if err != nil {
		panic(err)
	}

	repoLabels := convertGithubLabelsToMap(rl)

	// create
	for labelName, label := range *definedLabels {
		if _, found := (*repoLabels)[labelName]; !found {
			fmt.Printf("%vCREATE '%s'%v\n", cli.Green, label.Name, cli.Reset)
			githubapi.CreateLabel(*owner, *repo, &label)
		}
	}

	// update and delete
	for repoLabelName, repoLabel := range *repoLabels {
		if matchedLabel, found := (*definedLabels)[repoLabelName]; found {
			if matchedLabel.Name != repoLabel.Name || matchedLabel.Color != repoLabel.Color || matchedLabel.Description != repoLabel.Description {
				fmt.Printf("%vUPDATE '%s'%v\n", cli.Blue, repoLabel.Name, cli.Reset)
				githubapi.UpdateLabel(*owner, *repo, repoLabelName, &matchedLabel)
			}
		} else {
			fmt.Printf("%vDELETE '%s'%v\n", cli.Red, repoLabel.Name, cli.Reset)
			githubapi.DeleteLabel(*owner, *repo, repoLabelName)
		}
	}

}

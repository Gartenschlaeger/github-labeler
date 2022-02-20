package main

import (
	"fmt"

	"Gartenschlaeger/github-labeler/pkg/args"
	"Gartenschlaeger/github-labeler/pkg/cli"
	"Gartenschlaeger/github-labeler/pkg/config"
	"Gartenschlaeger/github-labeler/pkg/githubapi"
	"Gartenschlaeger/github-labeler/pkg/helper"
)

var arguments *args.Arguments

func mergeLabels() {
	definedLabels, err := config.ReadLabelDefinitions()
	if err != nil {
		panic(err)
	}

	repoLabels, err := githubapi.GetLabelsForRepository(*arguments.Owner, *arguments.Repository)
	if err != nil {
		panic(err)
	}

	repoLabelsMap := helper.ConvertGithubLabelsToMap(repoLabels)

	if arguments.IsDryMode {
		fmt.Printf("%vRunning in dry mode%v\n", cli.Yellow, cli.Reset)
	}

	// create
	for labelName, label := range *definedLabels {
		if _, found := (*repoLabelsMap)[labelName]; !found {
			fmt.Printf("%vCREATE '%s'%v\n", cli.Green, label.Name, cli.Reset)
			if !arguments.IsDryMode {
				githubapi.CreateLabel(*arguments.Owner, *arguments.Repository, &label)
			}
		}
	}

	// update and delete
	for repoLabelName, repoLabel := range *repoLabelsMap {
		if matchedLabel, found := (*definedLabels)[repoLabelName]; found {
			if matchedLabel.Name != repoLabel.Name || matchedLabel.Color != repoLabel.Color || matchedLabel.Description != repoLabel.Description {
				fmt.Printf("%vUPDATE '%s'%v\n", cli.Blue, repoLabel.Name, cli.Reset)

				if !arguments.IsDryMode {
					githubapi.UpdateLabel(*arguments.Owner, *arguments.Repository, repoLabelName, &matchedLabel)
				}
			}
		} else {
			if !arguments.SkipDelete {
				fmt.Printf("%vDELETE '%s'%v\n", cli.Red, repoLabel.Name, cli.Reset)

				if !arguments.IsDryMode {
					githubapi.DeleteLabel(*arguments.Owner, *arguments.Repository, repoLabelName)
				}
			}
		}
	}
}

func init() {
	arguments = args.Parse()

	githubapi.SetBearerToken(*arguments.Token)
}

func main() {
	mergeLabels()
}

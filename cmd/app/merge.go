package main

import (
	"flag"
	"fmt"
	"os"

	"Gartenschlaeger/github-labeler/pkg/colors"
	"Gartenschlaeger/github-labeler/pkg/githubapi"
	"Gartenschlaeger/github-labeler/pkg/labelconfig"
)

// mergCommandsFlagValues holds flag values
type mergCommandsFlagValues struct {
	token      *string
	owner      *string
	repository *string
	isDryMode  bool
	skipDelete bool
}

// convertGithubLabelsToMap converts the github api response to a map
func convertGithubLabelsToMap(labels *[]githubapi.GithubLabelResponse) *map[string]githubapi.GithubLabelResponse {
	m := make(map[string]githubapi.GithubLabelResponse)
	for i := 0; i < len(*labels); i++ {
		l := (*labels)[i]
		m[l.Name] = l
	}

	return &m
}

// mergeLabels starts the merge process
func mergeLabels(arguments *mergCommandsFlagValues) error {
	if arguments.isDryMode {
		fmt.Println("Running in dry mode")
	}

	lcp, err := labelconfig.GetDefaultPath()
	if err != nil {
		return err
	}

	fmt.Printf("Read label configuration from %s..\n", *lcp)

	definedLabels, err := labelconfig.ReadFromFile(*lcp)
	if err != nil {
		return err
	}

	fmt.Println("Read repository labels..")

	repoLabels, err := githubapi.GetLabelsForRepository(*arguments.owner, *arguments.repository)
	if err != nil {
		return err
	}

	repoLabelsMap := convertGithubLabelsToMap(repoLabels)

	fmt.Println("Start merge..")

	// create
	for labelName, label := range *definedLabels {
		if _, found := (*repoLabelsMap)[labelName]; !found {
			fmt.Printf("%vCreate label '%s'%v\n", colors.Green, label.Name, colors.Reset)
			if !arguments.isDryMode {
				githubapi.CreateLabel(*arguments.owner, *arguments.repository, label.Name, label.Color, label.Description)
			}
		}
	}

	// update and delete
	for repoLabelName, repoLabel := range *repoLabelsMap {
		if matchedLabel, found := (*definedLabels)[repoLabelName]; found {
			if matchedLabel.Name != repoLabel.Name || matchedLabel.Color != repoLabel.Color || matchedLabel.Description != repoLabel.Description {
				fmt.Printf("%vUpdate label '%s'%v\n", colors.Blue, repoLabel.Name, colors.Reset)

				if !arguments.isDryMode {
					githubapi.UpdateLabel(*arguments.owner, *arguments.repository, repoLabelName, matchedLabel.Name, matchedLabel.Color, matchedLabel.Description)
				}
			}
		} else {
			if !arguments.skipDelete {
				fmt.Printf("%vDelete label '%s'%v\n", colors.Red, repoLabel.Name, colors.Reset)

				if !arguments.isDryMode {
					githubapi.DeleteLabel(*arguments.owner, *arguments.repository, repoLabelName)
				}
			}
		}
	}

	fmt.Println("Merge finished")

	return nil
}

func MergeCommand(args []string) error {
	fv := mergCommandsFlagValues{}

	fs := flag.NewFlagSet("merge", flag.ExitOnError)
	fv.token = fs.String("t", os.Getenv("LABELER_TOKEN"), "Bearer token for Github API requests.")
	fv.owner = fs.String("o", os.Getenv("LABELER_OWNER"), "Github Owner")
	fv.repository = fs.String("r", os.Getenv("LABELER_REPO"), "Github repository name")
	dryMode := fs.Bool("dry-mode", false, "Enable dry mode")
	skipDelete := fs.Bool("skip-delete", false, "Skip deletion of unknown labels.")
	fs.Parse(args)

	fv.isDryMode = hasOptionalBoolFlag(dryMode)
	fv.skipDelete = hasOptionalBoolFlag(skipDelete)

	requireFlag(fv.token, "Token required. Use -t <token>")
	requireFlag(fv.owner, "Owner required. Use -o <owner>")
	requireFlag(fv.repository, "Repository required. Use -r <repository>")

	githubapi.SetBearerToken(*fv.token)

	return mergeLabels(&fv)
}

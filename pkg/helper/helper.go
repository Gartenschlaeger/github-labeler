package helper

import "Gartenschlaeger/github-labeler/pkg/githubapi"

func ConvertGithubLabelsToMap(labels *[]githubapi.GithubLabelResponse) *map[string]githubapi.GithubLabelResponse {
	m := make(map[string]githubapi.GithubLabelResponse)
	for i := 0; i < len(*labels); i++ {
		l := (*labels)[i]
		m[l.Name] = l
	}

	return &m
}

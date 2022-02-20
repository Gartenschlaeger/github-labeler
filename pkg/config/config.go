package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"Gartenschlaeger/github-labeler/pkg/types"
)

func ReadLabelDefinitions() (*map[string]types.LabelDefinition, error) {
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

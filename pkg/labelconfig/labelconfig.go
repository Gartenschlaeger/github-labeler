package labelconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type LabelDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func GetDefaultPath() (*string, error) {
	hp, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fp := path.Join(hp, ".config", "github-labeler", "labels.json")

	return &fp, nil
}

func ReadFromFile(path string) (*map[string]LabelDefinition, error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	labelSet := []LabelDefinition{}
	json.Unmarshal(fileData, &labelSet)

	if len(labelSet) == 0 {
		return nil, errors.New("list of expected labels are empty")
	}

	m := make(map[string]LabelDefinition)
	for i := 0; i < len(labelSet); i++ {
		m[labelSet[i].Name] = labelSet[i]
	}

	return &m, nil
}

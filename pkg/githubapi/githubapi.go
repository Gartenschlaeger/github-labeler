package githubapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"Gartenschlaeger/github-labeler/pkg/types"
)

const apiBaseUrl = "https://api.github.com"

var bearerToken = ""

func SetBearerToken(token string) {
	bearerToken = token
}

func isBearerTokenAvailable() bool {
	hasToken := strings.TrimSpace(bearerToken) != ""
	if !hasToken {
		log.Fatal("Token was not set. Call SetToken to set the api token before doing any request.")
		return false
	}

	return true
}

func getRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func postJsonRequest(url string, body []byte) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	bd := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPost, url, bd)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func patchJsonRequest(url string, body []byte) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	bd := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPatch, url, bd)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func deleteRequest(url string) (*http.Response, error) {
	if !isBearerTokenAvailable() {
		return nil, errors.New("bearer token missing")
	}

	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readAsBytes(res *http.Response) (*[]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

// https://docs.github.com/en/rest/reference/issues#list-labels-for-a-repository

func GetLabelsForRepository(owner string, repo string) (*[]GithubLabelResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/labels?page=1&per_page=100", apiBaseUrl, owner, repo)

	res, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	b, err := readAsBytes(res)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		labels := []GithubLabelResponse{}

		err := json.Unmarshal(*b, &labels)
		if err != nil {
			return nil, err
		}

		return &labels, nil
	} else {
		return nil, fmt.Errorf("labels request failed with unexpected status code %d", res.StatusCode)
	}
}

// https://docs.github.com/en/rest/reference/issues#delete-a-label

func DeleteLabel(owner string, repo string, labelName string) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/labels/%s", apiBaseUrl, owner, repo, labelName)

	res, err := deleteRequest(url)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusNoContent, nil
}

// https://docs.github.com/en/rest/reference/issues#create-a-label

func CreateLabel(owner string, repo string, label *types.LabelDefinition) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/labels", apiBaseUrl, owner, repo)

	reqObject := CreateLabelRequest{
		Name:        label.Name,
		Color:       label.Color,
		Description: label.Description,
	}

	reqData, err := json.Marshal(reqObject)
	if err != nil {
		return false, err
	}

	res, err := postJsonRequest(url, reqData)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusCreated, nil
}

// https://docs.github.com/en/rest/reference/issues#update-a-label

func UpdateLabel(owner string, repo string, oldName string, label *types.LabelDefinition) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/labels/%s", apiBaseUrl, owner, repo, oldName)

	reqObject := CreateLabelRequest{
		Name:        label.Name,
		Color:       label.Color,
		Description: label.Description,
	}

	reqData, err := json.Marshal(reqObject)
	if err != nil {
		return false, err
	}

	res, err := patchJsonRequest(url, reqData)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

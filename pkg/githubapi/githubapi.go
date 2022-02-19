package githubapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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

func getRequest(url string, token string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func deleteRequest(url string, token string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+token)

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
	if !isBearerTokenAvailable() {
		return nil, errors.New("cannot find bearer token for auth request")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/labels?page=1&per_page=100", apiBaseUrl, owner, repo)

	res, err := getRequest(url, bearerToken)
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

func DeleteLabel(owner string, repo string, labelName string) bool {
	if !isBearerTokenAvailable() {
		return false
	}

	url := fmt.Sprintf("%s/repos/%s/%s/labels/%s", apiBaseUrl, owner, repo, labelName)

	res, err := deleteRequest(url, bearerToken)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return res.StatusCode == http.StatusNoContent
}

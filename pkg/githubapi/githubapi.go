package githubapi

import (
	"encoding/json"
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

func doRequest(url string, token string) (*http.Response, error) {
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

func readAsBytes(res *http.Response) (*[]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

// https://docs.github.com/en/rest/reference/issues#delete-a-label

func DeleteLabel(owner string, repo string) {
	if !isBearerTokenAvailable() {
		return
	}

}

// https://docs.github.com/en/rest/reference/issues#list-labels-for-a-repository

func GetLabelsForRepository(owner string, repo string) GithubLabelsResponse {
	if !isBearerTokenAvailable() {
		return nil
	}

	url := fmt.Sprintf("%s/repos/%s/%s/labels?page=1&per_page=100", apiBaseUrl, owner, repo)

	res, err := doRequest(url, bearerToken)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	b, err := readAsBytes(res)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if res.StatusCode == 200 {
		labels := GithubLabelsResponse{}

		err := json.Unmarshal(*b, &labels)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		return labels
	} else {
		log.Fatal("Request responded with status code ", res.StatusCode)
		return nil
	}
}

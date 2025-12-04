package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	token        = ""
	org          = "CBE-Super-App"
	sourceBranch = "qa"
	targetBranch = "uat"
)

type Repo struct {
	Name string `json:"name"`
}

type Branch struct {
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

func apiRequest(method, url string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return client.Do(req)
}

func main() {
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100&page=%d", org, page)
		resp, _ := apiRequest("GET", url, nil)
		body, _ := ioutil.ReadAll(resp.Body)

		var repos []Repo
		json.Unmarshal(body, &repos)

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			fmt.Printf("\nProcessing %s...\n", repo.Name)

			// Check if qa exists
			qaURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", org, repo.Name, sourceBranch)
			qaResp, _ := apiRequest("GET", qaURL, nil)
			if qaResp.StatusCode != 200 {
				fmt.Println(" - qa branch missing, skipping.")
				continue
			}

			// Check if uat exists
			uatURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", org, repo.Name, targetBranch)
			uatResp, _ := apiRequest("GET", uatURL, nil)
			if uatResp.StatusCode == 200 {
				fmt.Println(" - uat already exists, skipping.")
				continue
			}

			// Get SHA for qa
			body, _ := ioutil.ReadAll(qaResp.Body)
			var branch Branch
			json.Unmarshal(body, &branch)
			sha := branch.Commit.SHA

			// Create uat
			createURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", org, repo.Name)
			payload := fmt.Sprintf(`{"ref":"refs/heads/%s","sha":"%s"}`, targetBranch, sha)

			createResp, _ := apiRequest("POST", createURL, []byte(payload))
			if createResp.StatusCode == 201 {
				fmt.Println(" - Created uat branch successfully.")
			} else {
				fmt.Printf(" - Failed to create uat branch: %s\n", createResp.Status)
			}
		}

		page++
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// https://docs.gitlab.com/ee/api/projects.html#list-all-projects
	projectsURL := fmt.Sprintf("%s/projects?per_page=1000", config.GitlabURL)
	data, err := APIRequest("GET", projectsURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Projects := Projects{}
	jsonErr := json.Unmarshal(data, &Projects)
	if jsonErr != nil {
		fmt.Println("Error parsing JSON for projects")
		os.Exit(1)
	}

	totalDeletedPipelines := 0
	totalDeletedJobs := 0
	totalDeletedSize := 0
	cutOff := float64(config.KeepDays * 24)

	for _, project := range Projects {
		if !project.JobsEnabled {
			continue
		}
		// https://docs.gitlab.com/ee/api/pipelines.html#list-project-pipelines
		pipelinesURL := fmt.Sprintf("%s/projects/%d/pipelines?per_page=1000", config.GitlabURL, project.ID)
		data, err := APIRequest("GET", pipelinesURL)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		Pipelines := Pipelines{}
		jsonErr := json.Unmarshal(data, &Pipelines)
		if jsonErr != nil {
			fmt.Println("Error parsing JSON for pipelines")
			continue
		}

		if len(Pipelines) <= config.MinPipelines {
			// fmt.Println(project.PathWithNamespace, "- no pipelines to delete")
			continue
		}

		PipelinesToDelete := Pipelines[config.MinPipelines:]

		fmt.Printf("-----\n%s - Examining %d pipeline(s)\n", project.PathWithNamespace, len(PipelinesToDelete))

		for _, pipeline := range PipelinesToDelete {

			if pipeline.Status == "running" || pipeline.Status == "pending" {
				fmt.Printf("- Skipping %s pipeline #%d\n", pipeline.Status, pipeline.ID)
				continue
			}

			// https://docs.gitlab.com/ee/api/jobs.html#list-pipeline-jobs
			jobsURL := fmt.Sprintf("%s/projects/%d/pipelines/%d/jobs", config.GitlabURL, project.ID, pipeline.ID)
			data, err := APIRequest("GET", jobsURL)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			Jobs := Jobs{}
			jsonErr := json.Unmarshal(data, &Jobs)
			if jsonErr != nil {
				fmt.Println("Error parsing JSON for pipelines")
				continue
			}

			canDelete := true

			for _, job := range Jobs {

				jobSize := 0

				if time.Since(job.CreatedAt).Hours() < cutOff {
					// fmt.Printf("- Skipping pipeline #%d as newer than %d days (%v)\n", pipeline.ID, config.KeepDays, job.CreatedAt)
					canDelete = false
					break
				}

				for _, artifact := range job.Artifacts {
					jobSize = jobSize + artifact.Size
				}

				if deleteFiles {
					// https://docs.gitlab.com/ee/api/jobs.html#erase-a-job
					eraseURL := fmt.Sprintf("%s/projects/%d/jobs/%d/erase", config.GitlabURL, project.ID, job.ID)
					_, err = APIRequest("POST", eraseURL)
					if err != nil {
						fmt.Printf("- Error deleting job #%d - %s\n", job.ID, err)
					} else {
						fmt.Printf("- Deleted %d artifact(s) - %dKb\n", len(job.Artifacts), bToKb(jobSize))
						totalDeletedSize = totalDeletedSize + jobSize
						totalDeletedJobs++
					}
				} else {
					fmt.Printf("- Would delete %d artifact(s) - %dKb\n", len(job.Artifacts), bToKb(jobSize))
					totalDeletedSize = totalDeletedSize + jobSize
					totalDeletedJobs++
				}

			}

			if canDelete {
				if deleteFiles {
					// https://docs.gitlab.com/ee/api/pipelines.html#delete-a-pipeline
					deleteURL := fmt.Sprintf("%s/projects/%d/pipelines/%d", config.GitlabURL, project.ID, pipeline.ID)
					_, err := APIRequest("DELETE", deleteURL)
					if err != nil {
						fmt.Printf("- Error deleting pipeline #%d: %s\n", pipeline.ID, err)
					} else {
						fmt.Printf("- Delete pipeline #%d\n", pipeline.ID)
						totalDeletedPipelines++
					}

				} else {
					fmt.Printf("- Would delete pipeline #%d\n", pipeline.ID)
					totalDeletedPipelines++
				}
			}
		}
	}

	if deleteFiles {
		fmt.Printf("\nReport (deleted)\n================\n")
	} else {
		fmt.Printf("\nReport (to delete / dry run)\n============================\n")
	}

	fmt.Printf("Minimum pipelines:  %d\n", config.MinPipelines)
	fmt.Printf("Keep for at least:  %d days\n", config.KeepDays)
	fmt.Printf("Deleted pipelines:  %d\n", totalDeletedPipelines)
	fmt.Printf("Deleted jobs:       %d\n", totalDeletedJobs)
	fmt.Printf("Total deleted size: %dKb\n", bToKb(totalDeletedSize))

}

// APIRequest is an authenticated request
// @param method eg: POST, GET, DELETE
// @param url The URL
func APIRequest(method, url string) ([]byte, error) {
	method = strings.ToUpper(method)

	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("PRIVATE-TOKEN", config.APIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}

	return nil, fmt.Errorf("%s %s returned a %v (%s)", method, url, resp.StatusCode, http.StatusText(resp.StatusCode))
}

// bToMb converts bytes to Kb
func bToKb(b int) int {
	return b / 1024
}

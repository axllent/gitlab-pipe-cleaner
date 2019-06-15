package main

import "time"

// Config struct
type Config struct {
	GitlabURL    string `json:"GitlabURL"`
	APIKey       string `json:"APIKey"`
	MinPipelines int    `json:"MinPipelines"`
	KeepDays     int    `json:"KeepDays"`
}

// Projects struct
type Projects []struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	CreatedAt         time.Time `json:"created_at"`
	JobsEnabled       bool      `json:"jobs_enabled"`
}

// Pipelines struct
type Pipelines []struct {
	ID     int    `json:"id"`
	Sha    string `json:"sha"`
	Ref    string `json:"ref"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// Jobs struct
type Jobs []struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at"`
	Duration      float64   `json:"duration"`
	ArtifactsFile struct {
		Filename string `json:"filename"`
		Size     int    `json:"size"`
	} `json:"artifacts_file"`
	Artifacts []struct {
		FileType   string `json:"file_type"`
		Size       int    `json:"size"`
		FileName   string `json:"filename"`
		FileFormat string `json:"file_format"`
	} `json:"artifacts"`
}

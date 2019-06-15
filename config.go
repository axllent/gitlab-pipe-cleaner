package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var (
	config      Config
	configFile  string
	deleteFiles bool
)

func init() {
	config = Config{}

	var configFile string

	defaultConfig := Home() + "/.config/gitlab-pipe-cleaner.json"

	flag.StringVar(&configFile, "c", defaultConfig, "Config file")
	flag.BoolVar(&deleteFiles, "delete", false, "Delete files (default dry run / report only)")

	// parse flags
	flag.Parse()

	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Use -c to specify the config location.")
		os.Exit(1)
	}

	err = json.Unmarshal(configJSON, &config)

	if config.GitlabURL == "" || config.APIKey == "" {
		fmt.Println("Error: no GitlabURL or APIKey found in your config")
		os.Exit(1)
	}

	// create the API (v4) URL
	config.GitlabURL = fmt.Sprintf("%s/api/v4", strings.TrimRight(config.GitlabURL, "/"))
	if !strings.HasPrefix(config.GitlabURL, "http") {
		fmt.Println("Error: GitlabURL must start with http or https")
		os.Exit(1)
	}
}

// Home returns the user's home directory
func Home() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

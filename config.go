package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/axllent/gitrel"
)

var (
	config      Config
	configFile  string
	deleteFiles bool
	showVersion bool
	update      bool
	version     = "dev"
)

func init() {
	config = Config{}

	var configFile string

	defaultConfig := Home() + "/.config/gitlab-pipe-cleaner.json"

	flag.StringVar(&configFile, "c", defaultConfig, "Config file")
	flag.BoolVar(&deleteFiles, "delete", false, "Delete files (default dry run / report only)")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&update, "update", false, "Update to latest version")

	// parse flags
	flag.Parse()

	if showVersion {
		fmt.Println(fmt.Sprintf("Version: %s", version))
		latest, _, _, err := gitrel.Latest("axllent/gitlab-pipe-cleaner", "gitlab-pipe-cleaner")
		if err == nil && latest != version {
			fmt.Println(fmt.Sprintf("Update available: %s\nRun `%s -update` to update.", latest, os.Args[0]))
		}
		os.Exit(0)
	}

	if update {
		rel, err := gitrel.Update("axllent/gitlab-pipe-cleaner", "gitlab-pipe-cleaner", version)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("Updated %s to version %s", os.Args[0], rel))
		os.Exit(0)
	}

	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Printf("\nUsage of %s\n", os.Args[0])
		flag.PrintDefaults()
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

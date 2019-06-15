# Gitlab Pipe Cleaner 

Gitlab Pipe Cleaner is a commandline tool that uses the Gitlab API to clean (prune) your servers CI pipelines and jobs, 
including all the build logs and artifacts.

It will run through **all projects** that your API key has access to, deleting pipelines based values on your configuration file ([example](#configuration-example)).

The tool is written in [Golang](https://golang.org/) so is cross-platform and efficient.


## Uage options

```
Usage of gitlab-pipe-cleaner:
  -c string
        Config file (default "~/.config/gitlab-pipe-cleaner.json")
  -delete
        Delete files (default dry run / report only)
  -update
        Update to latest version
  -v    Show version
```

Running gitlab-pipe-cleaner without the `-delete` option will not delete any data, but rather do a dry-run.


## Configuration example

Create a JSON confiburation file (edit with your own values).

This configuration file can be saved as `~/.config/gitlab-pipe-cleaner.json`, or specified with `gitlab-pipe-cleaner -c <path-to-your-configuration>`.

```json
{
  "GitlabURL": "https://git.example.com/",
  "APIKey": "XjsjQr2U1RkHAUpFyfF2",
  "MinPipelines": 10,
  "KeepDays": 7
}
```

- `GitlabURL` - The URL to your Gitlab server
- `APIKey` - Generate a personal access token with `api` access on `https://git.example.com/profile/personal_access_tokens`
- `MinPipelines` - The minimum number of latest pipelines to leave per project
- `KeepDays` - The minimum number of days to leave pipelines (regarless of `MinPipelines`)


## Installation

You can use of the pre-built binaries (see releases), or if you prefer to build it from source `go get github.com/axllent/gitlab-pipe-cleaner`

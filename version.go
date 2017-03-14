package main

import "github.com/rai-project/cmd"

var (
	// These fields are populated by govvv
	Version    = "0.2.5"
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func init() {
	cmd.Version = cmd.VersionInfo{
		Version:    Version,
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitBranch:  GitBranch,
		GitState:   GitState,
		GitSummary: GitSummary,
	}
}

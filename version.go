package main

import "github.com/rai-project/cmd"

var (
	// These fields are populated by govvv
	Version    string
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

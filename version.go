package main

import "github.com/rai-project/cmd"

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func init() {
	cmd.Version = cmd.VersionInfo{
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitBranch:  GitBranch,
		GitState:   GitState,
		GitSummary: GitSummary,
	}
}

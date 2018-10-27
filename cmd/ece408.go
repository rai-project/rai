// +build ece408ProjectMode

package cmd

import (
	"github.com/rai-project/client"
)

func validateEce408Options() error {
	return nil
}

func extraClientOptions(opts []client.Option) []client.Option {
	if ece408ProjectMode && submitionName != "" {
		return append(opts, client.SubmissionName(submitionName))
	}
	return opts
}

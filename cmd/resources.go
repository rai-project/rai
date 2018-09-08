//go:generate go get -v github.com/mjibson/esc
//go:generate go get -v golang.org/x/tools/cmd/stringer
//go:generate esc -o static_content.go -pkg cmd -private ../rai_config.yml

package cmd

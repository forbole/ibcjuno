package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

const (
	formatFlag = "format"
)

var (
	Version       = ""
	Commit        = ""
	versionFormat string
)

type versionInfo struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit" yaml:"commit"`
	Go      string `json:"go" yaml:"go"`
}

func getVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of IBCJuno",
		RunE: func(cmd *cobra.Command, args []string) error {
			verInfo := versionInfo{
				Version: Version,
				Commit:  Commit,
				Go:      fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
			}

			switch versionFormat {
			case "json":
				bz, err := json.Marshal(verInfo)
				if err != nil {
					return fmt.Errorf("error when marshaling the version: %s", err)
				}
				_, err = fmt.Println(string(bz))
				if err != nil {
					return err
				}

			default:
				bz, err := yaml.Marshal(&verInfo)
				if err != nil {
					return fmt.Errorf("error when marshaling the version: %s", err)
				}
				_, err = fmt.Println(string(bz))
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	versionCmd.Flags().StringVar(&versionFormat, formatFlag, "text", "Print the version in the given format (text | json)")

	return versionCmd
}

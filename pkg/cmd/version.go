package cmd

import (
	"github.com/asdfsx/k8s-device-plugin/pkg/version"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

type VersionCmd struct {
	command *cobra.Command
}

func NewVersionCmd() *VersionCmd {
	return &VersionCmd{
		command: &cobra.Command{
			Use:   "version",
			Short: "version: Print the version of k8s-device-plugin",
			Long: `version: Print the version of k8s-device-plugin

All software has versions. This is k8s-device-plugin's.`,
			Run: func(cmd *cobra.Command, args []string) {
				printWatcherVersion()
			},
		},
	}
}

func (cmd *VersionCmd) GetCommand() *cobra.Command {
	return cmd.command
}

func printWatcherVersion() {
	jww.FEEDBACK.Println(version.BuildVersionString())
}
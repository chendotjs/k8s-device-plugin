package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/chendotjs/k8s-device-plugin/pkg/version"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "version: Print the version of k8s-device-plugin",
		Long: `version: Print the version of k8s-device-plugin

All software has versions. This is k8s-device-plugin's.`,
		Run: func(cmd *cobra.Command, args []string) {
			printWatcherVersion()
		},
	}
}

func printWatcherVersion() {
	jww.FEEDBACK.Println(version.BuildVersionString())
}

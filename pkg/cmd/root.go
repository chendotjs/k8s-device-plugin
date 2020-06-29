package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "k8s-device-plugin",
		Short: "k8s-device-plugin: device-plugin for kubernetes",
		Long: `k8s-device-plugin: device-plugin for kubernetes

k8s-device-plugin is a device-plugin for kubernetes. Checking specific resources on the host,
then register them to the cluster.`,
		SilenceUsage: true,
	}

	command.AddCommand(NewDevicePluginCmd())
	command.AddCommand(NewVersionCmd())

	return command
}

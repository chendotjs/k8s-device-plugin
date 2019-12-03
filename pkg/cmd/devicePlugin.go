package cmd

import (
	"github.com/spf13/cobra"
	"github.com/asdfsx/k8s-device-plugin/pkg/utils"
	"github.com/asdfsx/k8s-device-plugin/pkg/plugins"
	"github.com/asdfsx/k8s-device-plugin/pkg/plugins/sample"
	"log"
	"syscall"
)

type DevicePluginCmd struct {
	command *cobra.Command
}

func NewDevicePluginCmd() *DevicePluginCmd {
	return &DevicePluginCmd{
		command: &cobra.Command{
			Use:   "device",
			Short: "device: device plugin for kubernetes",
			Long: `device: device-plugin for kubernetes

k8s-device-plugin is a device-plugin for kubernetes. Checking specific resouces on the host,
then regist them to the cluster.`,
			Run: func(cmd *cobra.Command, args []string) {
				log.Println("Starting OS watcher.")
				sigs := utils.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

				restart := true
				var devPlugin plugins.DevicePluginInterface
L:
				for {
					if restart {
						devPlugin = sample.NewSamplePlugin()
						if err := devPlugin.Start(); err != nil {
							log.Printf("Could not contact Kubelet, Cause %v. Retrying. Did you enable the device plugin feature gate?", err)
						} else {
							restart = false
						}
					}

					select{
					case s := <-sigs:
						switch s {
						case syscall.SIGHUP:
							log.Println("Received SIGHUP, restarting.")
							restart = true
						default:
							log.Printf("Received signal \"%v\", shutting down.", s)
							devPlugin.Stop()
							break L
						}
					}
				}
			},
		},
	}
}

func (cmd *DevicePluginCmd) GetCommand() *cobra.Command {
	return cmd.command
}

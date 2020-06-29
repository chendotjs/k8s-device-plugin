package cmd

import (
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"github.com/chendotjs/k8s-device-plugin/pkg/plugins"
	"github.com/chendotjs/k8s-device-plugin/pkg/utils"
)

func NewDevicePluginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "device",
		Short: "device: device plugin for kubernetes",
		Long: `device: device-plugin for kubernetes

k8s-device-plugin is a device-plugin for kubernetes. Checking specific resources on the host,
then register them to the cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
}

func run() {
	klog.Info("starting OS watcher")
	signals := utils.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	klog.Info("starting FS watcher")
	watcher, err := utils.NewFSWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		klog.Fatal("failed to create FS watcher: %v", err)
	}
	defer watcher.Close()

restart:
	plugin := plugins.NewSamplePlugin(plugins.NewDummyLinkManager("sample"))
	if err := plugin.Start(); err != nil {
		klog.Errorf("could not contact kubelet due to: %v. Retrying. Did you enable the device plugin feature gate?", err)
		goto restart
	}

	klog.Info("start device plugin successfully")

	for {
		select {
		// Detect a kubelet restart by watching for a newly created
		// 'pluginapi.KubeletSocket' file. When this occurs, restart this loop,
		// restarting all of the plugins in the process.
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				klog.Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
				goto restart
			}

		case err := <-watcher.Errors:
			klog.Errorf("inotify: %v", err)

		// Watch for any signals from the OS. On SIGHUP, restart this loop,
		// restarting all of the plugins in the process. On all other
		// signals, exit the loop and exit the program.
		case s := <-signals:
			switch s {
			case syscall.SIGHUP:
				klog.Info("received SIGHUP, restarting...")
				goto restart
			default:
				klog.Infof("Received signal \"%v\", shutting down...", s)
				plugin.Stop()
				os.Exit(0)
			}
		}
	}
}

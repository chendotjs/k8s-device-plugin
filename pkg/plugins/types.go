package plugins

import (
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type DevicePlugin interface {
	Start() error
	Stop() error
}

type ResourceManager interface {
	Devices() ([]*pluginapi.Device, error)
	CheckHealth(stop <-chan interface{}, devices []*pluginapi.Device, unhealthy chan<- *pluginapi.Device)
}

type SamplePlugin struct {
	resourceManager ResourceManager
	resourceName    string
	socket          string

	cachedDevices []*pluginapi.Device
	stop          chan interface{}
	health        chan *pluginapi.Device
	server        *grpc.Server
}

type DummyLinkManager struct {
	LinkPrefix string
}

var _ ResourceManager = &DummyLinkManager{}

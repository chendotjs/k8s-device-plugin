package plugins

import (
	"net"
	"strings"

	"github.com/vishvananda/netlink"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func NewDummyLinkManager(prefix string) ResourceManager {
	return &DummyLinkManager{
		LinkPrefix: prefix,
	}
}

func (m *DummyLinkManager) Devices() ([]*pluginapi.Device, error) {
	interfaces, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	var devices []*pluginapi.Device

	for _, iface := range interfaces {
		if iface.Type() != "dummy" {
			continue
		}
		if !strings.HasPrefix(iface.Attrs().Name, m.LinkPrefix) {
			continue
		}

		health := pluginapi.Unhealthy
		if isLinkUp(iface) {
			health = pluginapi.Healthy
		}

		devices = append(devices, &pluginapi.Device{
			ID:     iface.Attrs().Name,
			Health: health,
		})

	}
	return devices, nil
}

//TODO:
func (m *DummyLinkManager) CheckHealth(stop <-chan interface{}, devices []*pluginapi.Device, unhealthy chan<- *pluginapi.Device) {

}

func isLinkUp(link netlink.Link) bool {
	if link.Attrs().Flags&net.FlagUp != 0 {
		return true
	}
	return false
}

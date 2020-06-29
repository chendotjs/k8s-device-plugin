package plugins

import (
	"context"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"k8s.io/klog"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	SampleDpSocket     = "sample-dp.sock"
	SampleResourceName = "sample/dummy"
)

func NewSamplePlugin(resourceManager ResourceManager) *SamplePlugin {
	plugin := &SamplePlugin{
		resourceManager: resourceManager,
		resourceName:    SampleResourceName,
		socket:          pluginapi.DevicePluginPath + SampleDpSocket,
		stop:            make(chan interface{}),
		health:          make(chan *pluginapi.Device),
	}

	return plugin
}

func (m *SamplePlugin) initialize() {
	cachedDevices, err := m.resourceManager.Devices()
	if err != nil {
		klog.Fatal("error initialize: %v", err)
	}
	m.cachedDevices = cachedDevices
	m.server = grpc.NewServer([]grpc.ServerOption{}...)
	m.health = make(chan *pluginapi.Device)
	m.stop = make(chan interface{})
}

// dial establishes the gRPC communication with the registered device plugin.
func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	klog.Infof("start to dial %v", unixSocketPath)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return grpc.DialContext(ctx, unixSocketPath, grpc.WithBlock(), grpc.WithInsecure(),
		grpc.WithContextDialer(
			func(ctx context.Context, addr string) (net.Conn, error) {
				if deadline, ok := ctx.Deadline(); ok {
					return net.DialTimeout("unix", addr, time.Until(deadline))
				}
				return net.DialTimeout("unix", addr, 0)
			},
		),
	)
}

// Serve starts the gRPC server and register the device plugin to Kubelet
func (m *SamplePlugin) Start() error {
	m.initialize()

	err := m.serve()
	if err != nil {
		klog.Infof("Could not start device plugin: %s", err)
		return err
	}
	klog.Infof("starting to serve on %v", m.socket)

	err = m.Register(pluginapi.KubeletSocket, m.resourceName)
	if err != nil {
		klog.Infof("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	klog.Infof("registered device plugin for %s with kubelet", m.resourceName)

	return nil
}

// Stop stops the gRPC server
func (m *SamplePlugin) Stop() error {
	if m.server == nil {
		return nil
	}

	m.server.Stop()
	m.server = nil
	close(m.stop)

	return m.cleanup()
}

// Register registers the device plugin for the given resourceName with Kubelet.
func (m *SamplePlugin) Register(kubeletEndpoint, resourceName string) error {
	conn, err := dial(kubeletEndpoint, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: resourceName,
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// ListAndWatch lists devices and update that list according to the health status
func (m *SamplePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	klog.Infof("ListAndWatch: send all cached devices: %+v", m.cachedDevices)
	err := s.Send(&pluginapi.ListAndWatchResponse{Devices: m.cachedDevices})
	if err != nil {
		klog.Errorf("ListAndWatch: error sending devices: %v", err)
	}

	wait.PollImmediateInfinite(time.Second*10, func() (done bool, err error) {
		devices, err := m.resourceManager.Devices()
		if err != nil {
			klog.Errorf("ListAndWatch: error listing devices: %v", err)
			return false, nil
		}
		klog.Infof("ListAndWatch: send devices: %+v", devices)
		err = s.Send(&pluginapi.ListAndWatchResponse{Devices: devices})
		if err != nil {
			klog.Errorf("ListAndWatch: error sending devices: %v", err)
		}
		return false, nil
	})

	return nil
}

func (m *SamplePlugin) unhealthy(dev *pluginapi.Device) {
	m.health <- dev
}

// Allocate which return list of devices.
func (m *SamplePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	klog.Infof("allocate: recv request: %+v", *r)

	for _, req := range r.GetContainerRequests() {
		klog.Infof("allocate: assign device IDs: %v\n", req.DevicesIDs)
	}

	resp := make([]*pluginapi.ContainerAllocateResponse, len(r.GetContainerRequests()))

	for i, _ := range r.GetContainerRequests() {
		resp[i] = &pluginapi.ContainerAllocateResponse{
			Envs: map[string]string{"DummyLink": strings.Join(r.GetContainerRequests()[0].DevicesIDs, ",")},
		}
	}

	response := pluginapi.AllocateResponse{
		ContainerResponses: resp,
	}

	klog.Infof("allocate: send response: %+v", response)
	return &response, nil
}

func (m *SamplePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{
		PreStartRequired: false,
	}, nil
}

func (m *SamplePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (m *SamplePlugin) cleanup() error {
	if err := os.Remove(m.socket); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// Start starts the gRPC server of the device plugin
func (m *SamplePlugin) serve() error {
	err := m.cleanup()
	if err != nil {
		return err
	}

	sock, err := net.Listen("unix", m.socket)
	if err != nil {
		return err
	}

	m.server = grpc.NewServer([]grpc.ServerOption{}...)
	pluginapi.RegisterDevicePluginServer(m.server, m)

	go m.server.Serve(sock)

	// Wait for server to start by launching a blocking connection
	conn, err := dial(m.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

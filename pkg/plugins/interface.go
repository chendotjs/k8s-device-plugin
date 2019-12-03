package plugins

type DevicePluginInterface interface{
	Start() error
	Stop() error
}

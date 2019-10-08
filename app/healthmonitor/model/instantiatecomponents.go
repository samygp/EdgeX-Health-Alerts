package model

import (
	"fmt"

	"github.com/samygp/edgex-health-alerts/config"
)

//GetComponents instantiates all the components to be registered
//into the EdgeX metadata registry
func GetComponents() *EdgeXComponents {
	addressable := defaultAddressable()
	defaultDeviceProfile := newDeviceProfile()
	defaultDeviceService := newDeviceService(addressable.Name)
	defaultDevice := newDevice(defaultDeviceService.Name, defaultDeviceProfile.Name)

	return &EdgeXComponents{
		DeviceProfile:      newDeviceProfile(),
		DeviceService:      defaultDeviceService,
		DefaultAddressable: addressable,
		Device:             defaultDevice,
		ValueDescriptor:    newValueDescriptor("servicefailure"),
		ExportClients:      createExportClients(defaultDevice.Name),
	}
}

func defaultLabels() []string {
	return []string{"healthmonitor"}
}

func defaultAddressable() Addressable {
	return Addressable{
		Name:     fmt.Sprintf("%sAddressable", config.Config.App.Name),
		Method:   "POST",
		Protocol: "HTTP",
		Address:  "localhost",
		Port:     8000,
		Path:     "/",
	}
}

func newDevice(deviceServiceName, deviceProfileName string) Device {
	device := Device{
		Name:           "healthmonitor",
		AdminState:     "unlocked",
		OperatingState: "enabled",
		Labels:         defaultLabels(),
	}
	device.Service.Name = deviceServiceName
	device.Profile.Name = deviceProfileName
	device.Protocols.Protocol.Name = "default"
	return device
}

func newDeviceProfile() DeviceProfile {
	return DeviceProfile{
		Name:   "Health_monitor_profile",
		Labels: defaultLabels(),
	}
}

func newDeviceService(addressableName string) DeviceService {
	deviceService := DeviceService{
		Name:           "Health_Monitor_Device_Service",
		Description:    "Monitor health of EdgeX services and alert when services fail",
		Labels:         defaultLabels(),
		AdminState:     "unlocked",
		OperatingState: "enabled",
	}
	deviceService.Addressable.Name = addressableName
	return deviceService
}

func newValueDescriptor(name string) ValueDescriptor {
	return ValueDescriptor{
		Name:       name,
		Type:       "S",
		Formatting: "%s",
		Labels:     defaultLabels(),
	}
}

func createExportClients(deviceIdentifier string) []ExportClient {
	exportClients := make([]ExportClient, len(config.Config.ExportEndpoints))
	for i, exportEndpoint := range config.Config.ExportEndpoints {
		exportClients[i] = ExportClient{
			Name: fmt.Sprintf("%sExportClient", exportEndpoint.Name),
			Addressable: Addressable{
				Name:      exportEndpoint.Name,
				Protocol:  exportEndpoint.Protocol,
				Address:   exportEndpoint.Address,
				Port:      exportEndpoint.Port,
				Path:      exportEndpoint.Path,
				Method:    "POST",
				Publisher: "EdgeXExportPublisher",
			},
			Format:      "JSON",
			Enable:      true,
			Destination: "REST_ENDPOINT",
		}
		exportClients[i].Filter.DeviceIdentifiers = []string{deviceIdentifier}
	}
	return exportClients
}

package healthmonitor

import (
	"encoding/json"
	"fmt"

	ec "github.com/samygp/edgex-health-alerts/app/edgexconnector"
	"github.com/samygp/edgex-health-alerts/app/healthmonitor/model"
	"github.com/samygp/edgex-health-alerts/config"
	"github.com/samygp/edgex-health-alerts/fault"
	"github.com/samygp/edgex-health-alerts/log"
)

func (monitor *HealthMonitor) componentExists(clientID, path, componentName string) bool {
	namedPath := fmt.Sprintf("%s/name/%s", path, componentName)
	if err := monitor.connector.GetRequest(clientID, namedPath, false, nil); err != nil {
		if f, ok := err.(fault.Fault); !ok {
			log.Logger.Errorf("Error: %s", err.Error())
		} else if f.Status() != fault.NotFound {
			log.Logger.Errorf("Error: %s", f.Message())
		} else {
			log.Logger.Infof("%s not found on %s", componentName, clientID)
		}
		return false
	}
	return true
}

func (monitor *HealthMonitor) removeComponent(clientID, path, componentName string) {
	if monitor.componentExists(clientID, path, componentName) {
		if err := monitor.connector.DeleteRequest(clientID, fmt.Sprintf("%s/name/%s", path, componentName)); err != nil {
			log.Logger.Errorf("Error while deleting %s from %s: %s", componentName, clientID, err.Error())
			return
		}
	} else {
		log.Logger.Infof("Deleted %s from %s", componentName, clientID)
	}
}

func (monitor *HealthMonitor) registerComponent(clientID, path, componentName string, component interface{}) {
	if !monitor.componentExists(clientID, path, componentName) {
		log.Logger.Infof("Creating %s on %s", componentName, clientID)
		if err := monitor.connector.PostRequest(clientID, path, component); err != nil {
			log.Logger.Errorf("ERROR on registerComponent: %s", err.Error())
			return
		}
	}
	log.Logger.Infof("%s registered", componentName)
}

func (monitor *HealthMonitor) registerAllComponents() {
	c := &config.Config.EdgeXConnector
	components := model.GetComponents()
	buffer, _ := json.MarshalIndent(components, "", "  ")
	log.Logger.Debugf("Components: %s", buffer)
	monitor.components = components

	//register Device Profile
	monitor.registerComponent(ec.MetaData, c.MetaData.DeviceProfile, components.DeviceProfile.Name, components.DeviceProfile)

	//register default addressable
	monitor.registerComponent(ec.MetaData, c.MetaData.Addressable, components.DefaultAddressable.Name, components.DefaultAddressable)

	//register Device Service
	monitor.registerComponent(ec.MetaData, c.MetaData.DeviceService, components.DeviceService.Name, components.DeviceService)

	//register Device
	monitor.registerComponent(ec.MetaData, c.MetaData.Device, components.Device.Name, components.Device)

	//register Value Descriptor
	monitor.registerComponent(ec.CoreData, c.CoreData.ValueDescriptor, components.ValueDescriptor.Name, components.ValueDescriptor)

	//register export clients
	for _, client := range components.ExportClients {
		monitor.registerComponent(ec.ExportClient, c.ExportClient.Registration, client.Name, client)
	}
}

func (monitor *HealthMonitor) removeAllComponents() {
	c := &config.Config.EdgeXConnector

	//remove export clients
	for _, client := range monitor.components.ExportClients {
		monitor.removeComponent(ec.ExportClient, c.ExportClient.Registration, client.Name)
	}

	//remove Device
	monitor.removeComponent(ec.MetaData, c.MetaData.Device, monitor.components.Device.Name)

	//remove Device Service
	monitor.removeComponent(ec.MetaData, c.MetaData.DeviceService, monitor.components.DeviceService.Name)

	//remove default addressable
	monitor.removeComponent(ec.MetaData, c.MetaData.Addressable, monitor.components.DefaultAddressable.Name)

	//remove Device Profile
	monitor.removeComponent(ec.MetaData, c.MetaData.DeviceProfile, monitor.components.DeviceProfile.Name)

	//remove Value Descriptor
	monitor.removeComponent(ec.CoreData, c.CoreData.ValueDescriptor, monitor.components.ValueDescriptor.Name)

}

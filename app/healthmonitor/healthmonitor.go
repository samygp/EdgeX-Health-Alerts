package healthmonitor

import (
	ec "github.com/samygp/edgex-health-alerts/app/edgexconnector"
	"github.com/samygp/edgex-health-alerts/app/healthmonitor/model"
	"github.com/samygp/edgex-health-alerts/config"
	"github.com/samygp/edgex-health-alerts/log"
)

//HealthMonitor uses an EdgeXConnector to send messages
//to the EdgeX consul service, in order to receive health
//status for the running EdgeX services, and sends messages
//to subscribe the export client that will receive core data
//in order to report malfunctioning services
type HealthMonitor struct {
	connector  *ec.EdgeXConnector
	components *model.EdgeXComponents
}

//New instantiates a new health monitor
func New() *HealthMonitor {
	return &HealthMonitor{
		connector: ec.New(),
	}
}

//Start registers all the services required to export data to REST
//endpoints and send events to Core Data
func (monitor *HealthMonitor) Start() {
	log.Logger.Info("Starting...")
	monitor.registerAllComponents()
}

//Close removes all the subscribed components to EdgeX and stops monitoring
//health of the services in the EdgeX consul service
func (monitor *HealthMonitor) Close() {
	log.Logger.Info("Closing...")
	monitor.removeAllComponents()
}

func (monitor *HealthMonitor) queryConsul() model.ConsulStatusResponse {
	var result model.ConsulStatusResponse
	monitor.connector.GetRequest(ec.Consul, config.Config.EdgeXConnector.Consul.Health, true, &result)
	return result
}

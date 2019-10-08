package edgexconnector

import (
	"context"
	"fmt"
	"time"

	"github.com/samygp/edgex-health-alerts/config"
	"github.com/samygp/edgex-health-alerts/http"
)

const (
	//Consul reference to the clientID used for Consul API requests
	Consul = "consul"
	//CoreData reference to the clientID used for CoreData API requests
	CoreData = "coreData"
	//ExportClient reference to the clientID used for ExportClient API requests
	ExportClient = "exportClient"
	//MetaData reference to the clientID used for MetaData API requests
	MetaData = "metaData"
)

//EdgeXConnector handles sending periodical GET requests
//to EdgeX microservices, in order to monitor health of other
//services, via the Consul API, and to Post events using
//the device service API
type EdgeXConnector struct {
	client map[string]http.Client
}

//New instantiates a new EdgexConnector to handle REST requests
//to the EdgeX REST API
func New() *EdgeXConnector {
	edgeXConfig := &config.Config.EdgeXConnector
	debug := config.Config.App.Debug
	connector := &EdgeXConnector{
		client: make(map[string]http.Client),
	}

	connector.client[Consul] = http.New(fmt.Sprintf("%s%s", edgeXConfig.BaseURL, edgeXConfig.Consul.BasePath), nil, debug)
	connector.client[CoreData] = http.New(fmt.Sprintf("%s%s", edgeXConfig.BaseURL, edgeXConfig.CoreData.BasePath), nil, debug)
	connector.client[ExportClient] = http.New(fmt.Sprintf("%s%s", edgeXConfig.BaseURL, edgeXConfig.ExportClient.BasePath), nil, debug)
	connector.client[MetaData] = http.New(fmt.Sprintf("%s%s", edgeXConfig.BaseURL, edgeXConfig.MetaData.BasePath), nil, debug)

	return connector
}

//GetRequest generates a GET method request pointing to an EdgeX service's API
func (ec *EdgeXConnector) GetRequest(clientID, urlString string, withResponse bool, result interface{}) error {
	ctx, cancelFunction := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Config.EdgeXConnector.TimeoutMS))
	defer cancelFunction()
	if !withResponse {
		return ec.client[clientID].Get(ctx, urlString, http.WithStatusCode(http.StatusOK))
	}
	return ec.client[clientID].Get(ctx, urlString, http.WithResponse(result), http.WithStatusCode(http.StatusOK))
}

//PostRequest generates a POST method request pointing to an EdgeX service's API
func (ec *EdgeXConnector) PostRequest(clientID, urlString string, body interface{}) error {
	ctx, cancelFunction := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Config.EdgeXConnector.TimeoutMS))
	defer cancelFunction()
	return ec.client[clientID].Post(ctx, urlString, http.WithBody(body), http.WithStatusCode(http.StatusOK))
}

//DeleteRequest generates a POST method request pointing to an EdgeX service's API
func (ec *EdgeXConnector) DeleteRequest(clientID, urlString string) error {
	ctx, cancelFunction := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Config.EdgeXConnector.TimeoutMS))
	defer cancelFunction()
	return ec.client[clientID].Delete(ctx, urlString, http.WithStatusCode(http.StatusOK))
}

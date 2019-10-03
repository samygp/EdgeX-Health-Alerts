package healthmonitor

import (
	"context"
	"time"

	"github.com/samygp/Edgex-Health-Alerts/config"
	"github.com/samygp/Edgex-Health-Alerts/http"
)

//EdgeXConnector handles sending periodical GET requests
//to EdgeX microservices, in order to monitor health of other
//services, via the Consul API, and to Post events using
//the device service API
type EdgeXConnector struct {
	client http.Client
}

func (ec *EdgeXConnector) getRequest(urlString string, result interface{}) error {
	ctx, cancelFunction := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Config.EdgeXConnector.TimeoutMS))
	defer cancelFunction()
	return ec.client.Get(ctx, urlString, http.WithResponse(result), http.WithStatusCode(http.StatusOK))
}

func (ec *EdgeXConnector) postRequest(urlString string, body interface{}) error {
	ctx, cancelFunction := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Config.EdgeXConnector.TimeoutMS))
	defer cancelFunction()
	return ec.client.Post(ctx, urlString, http.WithBody(body), http.WithStatusCode(http.StatusOK))
}

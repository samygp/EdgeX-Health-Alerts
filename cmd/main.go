package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/samygp/edgex-health-alerts/app/healthmonitor"
	"github.com/samygp/edgex-health-alerts/config"
	"github.com/samygp/edgex-health-alerts/log"
)

func main() {
	config.Init()
	log.Init()

	defer log.Logger.Sync()

	buffer, err := json.MarshalIndent(config.Config, "", "  ")
	if err != nil {
		fmt.Printf("Unable to marshal config: %v\n", err)
	} else {
		log.Logger.Infof("Current config: %s\n", string(buffer))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	healthMonitor := healthmonitor.New()
	go healthMonitor.Start()
	<-stop
	healthMonitor.Close()
	log.Logger.Infof("Exit app %s", config.Config.App.Name)
}

package main

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/volteo/image-monitor/config"
	"bitbucket.org/volteo/image-monitor/log"
)

func main() {
	config.Init()
	log.Init()

	defer func() {
		if err := log.Logger.Sync(); err != nil {
			fmt.Printf("Unable to sync logger: %v\n", err)
		}
	}()

	buffer, err := json.Marshal(config.Config)
	if err != nil {
		fmt.Printf("Unable to marshal config: %v\n", err)
	} else {
		log.Logger.Infof("Current config: %s\n", string(buffer))
	}

	//TODO: add health check start
}

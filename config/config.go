package config

import (
	"github.com/jinzhu/configor"
)

// Config represents the configuration of the entire application.
var Config = struct {
	App struct {
		Debug bool   `json:"debug" default:"true"`
		Name  string `json:"name" default:"Health monitor"`
	} `json:"app"`
	ExportEndpoints []struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol" default:"HTTP"`
		Address  string `json:"address" default:"localhost"`
		Path     string `json:"path" default:"/"`
		Port     int    `json:"port" default:"8111"`
	} `json:"exportEndpoints"`
	EdgeXConnector struct {
		BaseURL string `json:"edgeXEndpoint" default:"http://localhost"`
		Consul  struct {
			BasePath string `json:"consulPath" default:":8500/v1"`
			Health   string `json:"health" default:"health"`
		} `json:"edgeXConnector"`
		ExportClient struct {
			BasePath     string `json:"exportClientPath" default:":48071/api/v1"`
			Registration string `json:"registration" default:"registration"`
		} `json:"exportClient"`
		MetaData struct {
			Addressable   string `json:"addressablePath" default:"addressable"`
			BasePath      string `json:"metaDataPath" default:":48081/api/v1"`
			Device        string `json:"devicePath" default:"device"`
			DeviceProfile string `json:"deviceProfilePath" default:"deviceprofile"`
			DeviceService string `json:"deviceServicePath" default:"deviceservice"`
		} `json:"metaData"`
		CoreData struct {
			BasePath        string `json:"coreDataPath" default:":48080/api/v1"`
			Event           string `json:"eventPath" default:"event"`
			ValueDescriptor string `json:"valuedescriptorPath" default:"valuedescriptor"`
		} `json:"coreData"`
		TimeoutMS int `json:"timeoutMS" default:"3000"`
	} `json:"edgeXConnector"`
	Logger struct {
		Level string `json:"level" default:"debug"`
	} `json:"logger"`
}{}

// Init config
func Init() {
	if err := configor.Load(&Config, "config.json"); err != nil {
		panic(err)
	}
}

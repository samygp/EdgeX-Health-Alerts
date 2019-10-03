package config

import (
	"github.com/jinzhu/configor"
)

// Config represents the configuration of the entire application.
var Config = struct {
	App struct {
		Name string `json:"name" default:"Image monitor"`
	} `json:"app"`
	EdgeXConnector struct {
		BaseURL string `json:"edgeXEndpoint" default:"http://localhost"`
		Consul  struct {
			BasePath string `json:"consulPath" default:":8500/v1"`
			Health   string `json:"health" default:":/health"`
		}
		ExportClient struct {
			BasePath     string `json:"exportClientPath" default:"48071/api/v1"`
			Registration string `json:"registration" default:"/registration"`
		}
		MetaData struct {
			Addressable   string `json:"addressablePath" default:"/addressable"`
			BasePath      string `json:"metaDataPath" default:":48081/api/v1"`
			Device        string `json:"devicePath" default:"/device"`
			DeviceProfile string `json:"deviceProfilePath" default:"/deviceprofile"`
			DeviceService string `json:"deviceServicePath" default:"/deviceservice"`
		}
		CoreData struct {
			BasePath        string `json:"coreDataPath" default:":48080/api/v1"`
			Event           string `json:"eventPath" default:"/event"`
			ValueDescriptor string `json:"valuedescriptorPath" default:"/valuedescriptor"`
		}
		MaxRetries int   `json:"maxretries" default:"3"`
		PollingMS  int64 `json:"pollingMS" default:"200"`
		TimeoutMS  int64 `json:"timeoutMS" default:"3000"`
	} `json:"edgeXConnector"`
	Logger struct {
		Level string `json:"level" default:"debug"`
	} `json:"logger"`
}{}

// Init config
func Init() {
	if err := configor.New(&configor.Config{ENVPrefix: "-"}).Load(&Config); err != nil {
		panic(err)
	}
}

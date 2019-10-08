package model

type consulStatus struct {
	Node        string      `json:"Node"`
	CheckID     string      `json:"CheckID"`
	Name        string      `json:"Name"`
	Status      string      `json:"Status"`
	Notes       string      `json:"Notes"`
	Output      string      `json:"Output"`
	ServiceID   string      `json:"ServiceID"`
	ServiceName string      `json:"ServiceName"`
	ServiceTags []string    `json:"ServiceTags"`
	Definition  interface{} `json:"Definition"`
	CreateIndex int32       `json:"CreateIndex"`
	ModifyIndex int32       `json:"ModifyIndex"`
}

//Addressable contains the payload to generate an Addressable entry
type Addressable struct {
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
	Address   string `json:"address"`
	Path      string `json:"path"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	Publisher string `json:"publisher"`
}

//ConsulStatusResponse holds the response objects when
//querying for the status of a Consul node's services
type ConsulStatusResponse []consulStatus

//Device contains the payload to generate a Device entry
type Device struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	AdminState     string   `json:"adminState"`
	OperatingState string   `json:"operatingState"`
	Labels         []string `json:"labels"`
	Service        struct {
		Name string `json:"name"`
	} `json:"service"`
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	Protocols struct {
		Protocol struct {
			Name string `json:"name"`
		} `json:"health monitor protocol"`
	} `json:"protocols"`
}

//DeviceProfile contains the payload to generate a Device Profile entry
type DeviceProfile struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
	Labels       []string `json:"labels"`
}

//DeviceService contains the payload to generate a DeviceService entry
type DeviceService struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Labels         []string `json:"labels"`
	AdminState     string   `json:"adminState"`
	OperatingState string   `json:"operatingState"`
	Addressable    struct {
		Name string `json:"name"`
	} `json:"addressable"`
}

//EdgeXResponseID represents the ID string that EdgeX Core API
//sends after receiving a POST/PUT request to uniquely identify
//a transaction or object
type EdgeXResponseID string

//EdgeXComponents contains a reference to instantiated
//EdgeX components to be registered
type EdgeXComponents struct {
	DeviceProfile      DeviceProfile
	DeviceService      DeviceService
	Device             Device
	ValueDescriptor    ValueDescriptor
	DefaultAddressable Addressable
	ExportClients      []ExportClient
}

//ExportClient contains the payload to generate a ExportClient entry
type ExportClient struct {
	Name        string      `json:"name"`
	Addressable Addressable `json:"addressable"`
	Filter      struct {
		DeviceIdentifiers []string `json:"deviceIdentifiers"`
	} `json:"filter"`
	Format      string `json:"format"`
	Enable      bool   `json:"enable"`
	Destination string `json:"destination"`
}

//ValueDescriptor contains the payload to generate a ValueDescriptor entry
type ValueDescriptor struct {
	Name         string   `json:"name" default:"servicefailure"`
	Description  string   `json:"description" default:"name of a failing service"`
	Type         string   `json:"type" default:"S"`
	DefaultValue string   `json:"defaultValue" default:""`
	Formatting   string   `json:"formatting" default:"%s"`
	Labels       []string `json:"labels" default:"['healthmonitor']"`
}

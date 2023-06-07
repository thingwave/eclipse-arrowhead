package datamodels

type SystemRequestDTO struct {
	SystemName         string             `json:"systemName"`
	Address            string             `json:"address"`
	Port               int                `json:"port"`
	AuthenticationInfo string             `json:"authenticationInfo,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
}

type ServiceRegistryEntryDTO struct {
	ServiceDefinition string             `json:"serviceDefinition"`
	ProviderSystem    SystemRequestDTO   `json:"providerSystem"`
	ServiceUri        string             `json:"serviceUri"`
	EndOfValidity     string             `json:"endOfValidity,omitempty"`
	Secure            string             `json:"secure"`
	Metadata          *map[string]string `json:"metadata,omitempty"`
	Version           int                `json:"version"`
	Interfaces        []string           `json:"interfaces"`
	CreatedAt         string             `json:"createdAt,omitempty"`
	UpdatedAt         string             `json:"updatedAt,omitempty"`
}

type ServiceRegistryEntryDTOIncomplete struct {
	ServiceDefinition string             `json:"serviceDefinition,omitempty"`
	ProviderSystem    *SystemRequestDTO  `json:"providerSystem,omitempty"`
	ServiceUri        string             `json:"serviceUri,omitempty"`
	EndOfValidity     string             `json:"endOfValidity,omitempty"`
	Secure            string             `json:"secure,omitempty"`
	Metadata          *map[string]string `json:"metadata,omitempty"`
	Version           *int               `json:"version,omitempty"`
	Interfaces        []string           `json:"interfaces,omitempty"`
}

type ServiceRegistryResponseDTO struct {
	Id                int64                         `json:"id"`
	ServiceDefinition ServiceDefinitionResponseDTO  `json:"serviceDefinition"`
	Provider          SystemResponseDTO             `json:"provider"`
	ServiceUri        string                        `json:"serviceUri"`
	EndOfValidity     string                        `json:"endOfValidity,omitempty"`
	Secure            string                        `json:"secure"`
	Metadata          *map[string]string            `json:"metadata,omitempty"`
	Version           int                           `json:"version"`
	Interfaces        []ServiceInterfaceResponseDTO `json:"interfaces"`
	CreatedAt         string                        `json:"createdAt"`
	UpdatedAt         string                        `json:"updatedAt"`
}

type ServiceDefinitionResponseDTO struct {
	Id                int64  `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type ServiceDefinitionResponseListDTO struct {
	Data  []ServiceDefinitionResponseDTO `json:"data"`
	Count int64                          `json:"count"`
}

type ServiceInterfaceListResponseDTO struct {
	Data  []ServiceInterfaceResponseDTO `json:"data"`
	Count int                           `json:"count"`
}

type ServiceInterfaceResponseDTO struct {
	Id            int64  `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type SystemResponseDTO struct {
	Id                 int64              `json:"id"`
	SystemName         string             `json:"systemName"`
	Address            string             `json:"address"`
	Port               int                `json:"port"`
	AuthenticationInfo string             `json:"authenticationInfo,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
	CreatedAt          string             `json:"createdAt"`
	UpdatedAt          string             `json:"updatedAt"`
}

type ServiceDefinitionDTO struct {
	Id                int64  `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type ServiceQueryResultListDTO struct {
	Results ServiceQueryResultDTO `json:"results"`
}

type ServiceQueryResultDTO struct {
	ServiceQueryData []ServiceRegistryResponseDTO `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}

type ServiceRegistryListResponseDTO struct {
	Data  []ServiceRegistryResponseDTO `json:"data"`
	Count int                          `json:"count"`
}

type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryResponseDTO `json:"serviceQueryData"`
}

type SystemListDTO struct {
	Data  []SystemResponseDTO `json:"data"`
	Count int                 `json:"count"`
}

type OrchestrationFormRequestDTO struct {
	RequesterSystem SystemRequestDTO `json:"requesterSystem"`
	RequesterCloud  CloudRequestDTO  `json:"requesterCloud"`
}

type CloudRequestDTO struct {
	Operator           string  `json:"operator"`
	Name               string  `json:"operator"`
	Secure             bool    `json:"secure"`
	Neighbor           bool    `json:"neighbor"`
	AuthenticationInfo string  `json:"authenticationInfo"`
	GatekeeperRelayIds []int64 `json:"gatekeeperRelayIds"`
	GatewayRelayIds    []int64 `json:"gatewayRelayIds"`
}

type OrchestrationResultDTO struct {
	Provider   SystemResponseDTO             `json:"provider"`
	Service    ServiceDefinitionResponseDTO  `json:"service"`
	Interfaces []ServiceInterfaceResponseDTO `json:"interfaces"`

	ServiceUri string             `json:"serviceUri"`
	Secure     string             `json:"secure"`
	Metadata   *map[string]string `json:"metadata,omitempty"`

	Version  int      `json:"version"`
	Warnings []string `json:"warnings"`
}

type OrchestrationResponseDTO struct {
	Response []OrchestrationResultDTO `json:"response"`
}

type StoreEntry struct {
	ID int `json:"id"`
}

type StoreEntryList struct {
	Count int          `json:"count"`
	Data  []StoreEntry `json:"data"`
}

type ServiceInterface struct {
	ID             int64  `json:"id"`
	Interface_name string `json:"interface_name"`
}

type PriorityList struct {
	PriorityMap map[string]int `json:"priorityMap"`
}

type ErrorMessageDTO struct {
	ErrorMessage  string `json:"errorMessage"`
	ErrorCode     int    `json:"errorCode"`
	ExceptionType string `json:"exceptionType"`
	Origin        string `json:"origin,omitempty"`
}

type InterfaceDTO struct {
	Id    int64  `json:"id"`
	Value string `json:"value"`
}

type ServiceDTO struct {
	Id    int64  `json:"id"`
	Value string `json:"value"`
}

type ServiceRequestForm struct {
	RequesterSystem    SystemRequestDTO       `json:"requesterSystem"`
	RequestedService   RequestedServiceDTO    `json:"requestedService"`
	PreferredProviders []PreferredProviderDTO `json:"preferredProviders,omitempty"`
}
type RequestedServiceDTO struct {
	ServiceDefinitionRequirement string             `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string           `json:"interfaceRequirements"`
	SecurityRequirements         []string           `json:"securityRequirements,omitempty"`
	MetadataRequirements         *map[string]string `json:"metadataRequirements,omitempty"`
	VersionRequirement           int64              `json:"versionRequirement,omitempty"`
	MaxVersionRequirement        int64              `json:"maxVersionRequirement,omitempty"`
	MinVersionRequirement        int64              `json:"minVersionRequirement,omitempty"`
}

type PreferredProviderDTO struct {
	ProviderCloud  ProviderCloudDTO  `json:"providerCloud"`
	ProviderSystem ProviderSystemDTO `json:"providerSystem"`
}

type ProviderCloudDTO struct {
	Operator string `json:"operator"`
	Name     string `json:"name"`
}
type ProviderSystemDTO struct {
	SystemName string `json:"systemName"`
	Address    string `json:"address"`
	Port       uint16 `json:"port"`
}

// /////////////////////////////////////////////////////////////////////////////
type SenMLEntry struct {
	Bn   *string  `json:"bn,omitempty"`
	Bt   *float64 `json:"bt,omitempty"`
	Bu   *string  `json:"bu,omitempty"`
	Bv   *float64 `json:"bv,omitempty"`
	Bs   *float64 `json:"bs,omitempty"`
	Bver *int     `json:"bver,omitempty"`

	N  *string  `json:"n,omitempty"`
	U  *string  `json:"u,omitempty"`
	V  *float64 `json:"v,omitempty"`
	Vs *string  `json:"vs,omitempty"`
	Vb *bool    `json:"vb,omitempty"`
	Vd *string  `json:"vd,omitempty"`

	S  *float64 `json:"s,omitempty"`
	T  *float64 `json:"t,omitempty"`
	Ut *float64 `json:"ut,omitempty"`
}

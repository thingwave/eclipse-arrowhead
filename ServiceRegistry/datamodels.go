package main

type JWT_hdr struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWT_payload struct {
	Iss string `json:"iss"`
	Exp int64  `json:"exp"`
	//Sub string `json:"sub"`
	//Aud string `json:"aud"`
	//Jti string `json:"jti"`
	Email string `json:"email"`
}

type ServiceQueryForm struct {
	ServiceDefinitionRequirement string             `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string           `json:"interfaceRequirements,omitempty"`
	SecurityRequirements         []string           `json:"securityRequirements,omitempty"`
	MetadataRequirements         *map[string]string `json:"metadataRequirements,omitempty"`
	VersionRequirement           *int               `json:"versionRequirement,omitempty"`
	MaxVersionRequirement        *int               `json:"maxVersionRequirement,omitempty"`
	MinVersionRequirement        *int               `json:"minVersionRequirement,omitempty"`
	PingProviders                *bool              `json:"pingProviders,omitempty"`
}

/*type ServiceDefinitionDTO struct {
	Id                int    `json:"id"`
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

type ServiceRegistryResponseDTO struct {
	Id                int                           `json:"id"`
	ServiceDefinition ServiceDefinitionResponseDTO  `json:"serviceDefinition"`
	Provider          SystemResponseDTO             `json:"provider"`
	ServiceUri        string                        `json:"serviceUri"`
	EndOfValidity     string                        `json:"endOfValidity"`
	Secure            string                        `json:"secure"`
	Metadata          *map[string]string            `json:"metadata,omitempty"`
	Version           int                           `json:"version"`
	Interfaces        []ServiceInterfaceResponseDTO `json:"interfaces"`
	CreatedAt         string                        `json:"createdAt"`
	UpdatedAt         string                        `json:"updatedAt"`
}

type ServiceDefinitionResponseDTO struct {
	Id                int    `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type ServiceInterfaceResponseDTO struct {
	Id            int    `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}*/

/*type ProviderDTO struct {
	Id                 int    `json:"id"`
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}*/

/*type SystemDTO struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}*/

/*type SystemRequestDTO struct {
	SystemName         string             `json:"systemName"`
	Address            string             `json:"address"`
	Port               int                `json:"port"`
	AuthenticationInfo string             `json:"authenticationInfo,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
}

type SystemResponseDTO struct {
	Id                 int                `json:"id"`
	SystemName         string             `json:"systemName"`
	Address            string             `json:"address"`
	Port               int                `json:"port"`
	AuthenticationInfo string             `json:"authenticationInfo,omitempty"`
	Metadata           *map[string]string `json:"metadata,omitempty"`
	CreatedAt          string             `json:"createdAt"`
	UpdatedAt          string             `json:"updatedAt"`
}

type ServiceQueryEntryDTO struct {
	Id                int                  `json:"id"`
	ServiceDefinition ServiceDefinitionDTO `json:"serviceDefinition"`
	Provider          SystemResponseDTO    `json:"provider"`
	ServiceUri        string               `json:"serviceUri"`
}

type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryResponseDTO `json:"serviceQueryData"`
}

type ServiceRegistryEntryDTO struct {
	ServiceDefinition string             `json:"serviceDefinition"`
	ProviderSystem    SystemRequestDTO   `json:"providerSystem"`
	ServiceUri        string             `json:"serviceUri"`
	EndOfValidity     string             `json:"endOfValidity"`
	Secure            string             `json:"secure"`
	Metadata          *map[string]string `json:"metadata,omitempty"`
	Version           int                `json:"version"`
	Interfaces        []string           `json:"interfaces"`
}*/

/*type InterfaceEntryDTO struct {
	Id            int    `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}*/

/*type ServiceRegistryEntryResp struct {
	Id                   int                  `json:"id"`
	ServiceDefinition    ServiceDefinitionDTO `json:"serviceDefinition"`
	Provider             ProviderDTO          `json:"provider"`
	ServiceUri           string               `json:"serviceUri"`
	EndOfValidity        string               `json:"endOfValidity"`
	Secure               string               `json:"secure"`
	MetadataRequirements *map[string]string   `json:"metadataRequirements,omitempty"`
	Version              int                  `json:"version"`
	Interfaces           []InterfaceEntryDTO  `json:"interfaces"`
	CreatedAt            string               `json:"createdAt"`
	UpdatedAt            string               `json:"updatedAt"`
}*/

type ServiceQueryFormDTO struct {
	ServiceDefinitionRequirement string             `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string           `json:"interfaceRequirements,omitempty"` // if specified at least one of the interfaces must match
	SecurityRequirements         []string           `json:"securityRequirements,omitempty"`  // if specified at least one of the types must match
	MetadataRequirements         *map[string]string `json:"metadataRequirements,omitempty"`  // if specified the whole content of the map must match
	VersionRequirement           *int               `json:"versionRequirement,omitempty"`    // if specified version must match
	MinVersionRequirement        *int               `json:"minVersionRequirement,omitempty"` // if specified version must be equals or higher; ignored if versionRequirement is specified
	MaxVersionRequirement        *int               `json:"maxVersionRequirement,omitempty"` // if specified version must be equals or lower; ignored if versionRequirement is specified
	PingProviders                bool               `json:"pingProviders,omitempty"`
}

/*type SystemListDTO struct {
	Data  []SystemResponseDTO `json:"data"`
	Count int                 `json:"count"`
}

type ErrorMessageDTO struct {
	ErrorMessage  string `json:"errorMessage"`
	ErrorCode     int    `json:"errorCode"`
	ExceptionType string `json:"exceptionType"`
	Origin        string `json:"origin,omitempty"`
}*/

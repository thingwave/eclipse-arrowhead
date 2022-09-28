package main

/*type JWT_hdr struct {
        Alg string `json:"alg"`
        Typ string `json:"typ"`
}

type JWT_payload struct {
  Iss string `json:"iss"`
  Exp int64 `json:"exp"`
  //Sub string `json:"sub"`
  //Aud string `json:"aud"`
  //Jti string `json:"jti"`
  Email string `json:"email"`
}

type MetadataRequirementsDTO struct {
  AdditionalProp1 *string `json:"additionalProp1,omitempty"`
  AdditionalProp2 *string `json:"additionalProp2,omitempty"`
  AdditionalProp3 *string `json:"additionalProp3,omitempty"`
}*/
/*
type ServiceQueryForm struct {
  ServiceDefinitionRequirement string `json:"serviceDefinitionRequirement"`
  InterfaceRequirements []string `json:"interfaceRequirements,omitempty"`
  SecurityRequirements []string `json:"securityRequirements,omitempty"`

  //MetadataRequirements *MetadataRequirementsDTO `json:"metadataRequirements,omitempty"`
  MetadataRequirements *map[string]string `json:"metadataRequirements,omitempty"`

  VersionRequirement *int `json:"versionRequirement,omitempty"`
  MaxVersionRequirement *int `json:"maxVersionRequirement,omitempty"`
  MinVersionRequirement *int `json:"minVersionRequirement,omitempty"`
  PingProviders *bool  `json:"pingProviders,omitempty"`
}

type ServiceDefinitionDTO struct {
  Id  int `json:"id"`
  ServiceDefinition string `json:"serviceDefinition"`
  CreatedAt string `json:"createdAt"`
  UpdatedAt string `json:"updatedAt"`
}

type ProviderDTO struct {
  Id  int `json:"id"`
  SystemName string `json:"systemName"`
  Address string `json:"address"`
  Port int `json:"port"`
  AuthenticationInfo string `json:"authenticationInfo"`
  CreatedAt string `json:"createdAt"`
  UpdatedAt string `json:"updatedAt"`
}

type SystemDTO struct {
  SystemName string `json:"systemName"`
  Address string `json:"address"`
  Port int `json:"port"`
  AuthenticationInfo string `json:"authenticationInfo"`
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
  Id int `json:"id"`
  ServiceDefinition ServiceDefinitionDTO `json:"serviceDefinition"`
  Provider ProviderDTO `json:"provider"`
  ServiceUri string  `json:"serviceUri"`
}

type ServiceQueryList struct {
  ServiceQueryData []ServiceQueryEntryDTO `json:"serviceQueryData"`
}

type ServiceRegistryEntryDTO struct {
  ServiceDefinition string `json:"serviceDefinition"`
  ProviderSystem SystemDTO `json:"providerSystem"`
  ServiceUri string `json:"serviceUri"`
  EndOfValidity string `json:"endOfValidity"`
  Secure string `json:"secure"`
  //Metadata *MetadataRequirementsDTO `json:"metadata"`
  Metadata *map[string]string `json:"metadata,omitempty"`
  Version int `json:"version"`
  Interfaces []string `json:"interfaces"`
}

type InterfaceEntryDTO struct {
  Id int `json:"id"`
  InterfaceName string `json:"interfaceName"`
  CreatedAt string `json:"createdAt"`
  UpdatedAt string `json:"updatedAt"`
}

type OrchestrationFormRequestDTO struct {
  RequesterSystem SystemDTO  `json:"requesterSystem"`
  RequesterCloud CloudRequestDTO  `json:"requesterCloud"`

}

type CloudRequestDTO struct {
  Operator string `json:"operator"`
  Name string `json:"operator"`
  Secure bool `json:"secure"`
  Neighbor bool `json:"neighbor"`
  AuthenticationInfo string `json:"authenticationInfo"`
  GatekeeperRelayIds []int64 `json:"gatekeeperRelayIds"`
  GatewayRelayIds []int64 `json:"gatewayRelayIds"`
}

type RequestedServiceDTO struct {
  ServiceDefinitionRequirement string`json:"serviceDefinitionRequirement"`
  InterfaceRequirements []string `json:"interfaceRequirements"`
}

type ServiceRequestForm struct {
  RequesterSystem SystemDTO  `json:"requesterSystem"`
  RequestedService RequestedServiceDTO  `json:"RequestedService"`
  //SecurityRequirements []string `json:"securityRequirements,omitempty"`
  //Metadata *MetadataRequirementsDTO `json:"metadata"`
  Metadata *map[string]string `json:"metadata,omitempty"`
  VersionRequirement *int `json:"versionRequirement,omitempty"`
  MaxVersionRequirement *int `json:"maxVersionRequirement,omitempty"`
  MinVersionRequirement *int `json:"minVersionRequirement,omitempty"`

  },
  "preferredProviders": [
    {
      "providerCloud": {
        "operator": "string",
        "name": "string"
      },
      "providerSystem": {
        "systemName": "string",
        "address": "string",
        "port": 0
      }
    }
  ],
  "orchestrationFlags": {
    "additionalProp1": true,
    "additionalProp2": true,
    "additionalProp3": true
  }
}
}

type OrchestrationResultDTO struct {
  Provider ProviderDTO `json:"provider"`

  ServiceUri string  `json:"serviceUri"`
  Secure string `json:"secure"`
  Metadata *map[string]string `json:"metadata,omitempty"`

  Version int `json:"version"`
  Warnings []string `json:"warnings"`
}

type OrchestrationResponseDTO struct {
  Response OrchestrationResultDTO  `json:"response"`
}

type StoreEntry struct {
  ID int  `json:"id"`
}

type StoreEntryList struct {
  Count int `json:"count"`
  Data []StoreEntry `json:"data"`
}

*/

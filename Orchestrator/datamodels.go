package main

/*
type MetadataRequirementsDTO struct {
  AdditionalProp1 *string `json:"additionalProp1,omitempty"`
  AdditionalProp2 *string `json:"additionalProp2,omitempty"`
  AdditionalProp3 *string `json:"additionalProp3,omitempty"`
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

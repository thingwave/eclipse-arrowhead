package main

import (
	dto "arrowhead.eu/common/datamodels"
)

type ErrorMessageDTO struct {
  ErrorMessage  string `json:"errorMessage"`
  ErrorCode     int    `json:"errorCode"`
  ExceptionType string `json:"exceptionType"`
  Origin        string `json:"origin,omitempty"`
}

/*type SystemResponseDTO struct {
  Id                 int64              `json:"id"`
  SystemName         string             `json:"systemName"`
  Address            string             `json:"address"`
  Port               int                `json:"port"`
  AuthenticationInfo string             `json:"authenticationInfo,omitempty"`
  Metadata           *map[string]string `json:"metadata,omitempty"`
  CreatedAt          string             `json:"createdAt"`
  UpdatedAt          string             `json:"updatedAt"`
}*/

type AuthorizationIntraCloudListResponseDTO struct {
  Data []AuthorizationIntraCloudResponseDTO  `json:"data"`
  Count int `json:"count"`
}

type AuthorizationIntraCloudResponseDTO struct {
  Id int64 `json:"id"`
  ConsumerSystem dto.SystemResponseDTO `json:"consumerSystem"`
  ProviderSystem dto.SystemResponseDTO `json:"providerSystem"`

  ServiceDefinition ServiceDefinitionResponseDTO `json:"serviceDefinition"`
  Interfaces []ServiceInterfaceResponseDTO `json:"interfaces"`

  CreatedAt string `json:"createdAt"`
  UpdatedAt string `json:"updatedAt"`
}

type ServiceDefinitionResponseDTO struct {
  Id                int64    `json:"id"`
  ServiceDefinition string `json:"serviceDefinition"`
  CreatedAt         string `json:"createdAt"`
  UpdatedAt         string `json:"updatedAt"`
}

type ServiceInterfaceResponseDTO struct {
  Id            int64   `json:"id"`
  InterfaceName string `json:"interfaceName"`
  CreatedAt     string `json:"createdAt"`
  UpdatedAt     string `json:"updatedAt"`
}

type AuthorizationIntraCloudRequestDTO struct {
  ConsumerId int64  `json:"consumerId"`
  ProviderIds []int64  `json:"providerIds"`
  ServiceDefinitionIds []int64 `json:"serviceDefinitionIds"`
  InterfaceIds []int64 `json:"interfaceIds"`
}

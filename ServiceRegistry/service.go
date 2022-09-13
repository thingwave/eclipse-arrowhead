package main

import(
	"errors"
	dto "arrowhead.eu/common/datamodels"
)

func validateRegistryEntry(req dto.ServiceRegistryEntryDTO) error {

	if req.ServiceDefinition == "" {
		return errors.New("serviceDefinition is missing!")
	}

	if req.ProviderSystem.SystemName == "" {
		return errors.New("providerSystem.systemName is missing!")
	}
	if req.ProviderSystem.Address == "" {
		return errors.New("providerSystem.Address is missing!")
	}
	if req.ProviderSystem.Port == 0 {
		return errors.New("providerSystem.Port is missing!")
	}
	if config.Server_ssl_enabled == true {
		if req.ProviderSystem.AuthenticationInfo == "" {
			return errors.New("providerSystem.AuthenticationInfo is missing!")
		}
	}

	return nil
}

/********************************************************************************
 * Copyright (c) 2022 Lulea University of Technology
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0.
 *
 * SPDX-License-Identifier: EPL-2.0
 *
 * Contributors:
 *   ThingWave AB - implementation
 *   Arrowhead Consortia - conceptualization
 ********************************************************************************/

package main

import (
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

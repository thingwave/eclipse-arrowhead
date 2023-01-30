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

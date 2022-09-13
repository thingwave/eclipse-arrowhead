# Eclipse Arrowhead Serviceregistry in Go
This is a an Arrowhead Service registry written in Go.


## Endpoints

### Client endpoint description


#### Echo
This endpoint returns a "Got it!" upon a GET request.

Status: 100% implemented

#### Query
Returns ServiceQueryList that fits the input specification. Mainly used by the Orchestrator.

Status: 60% implemented

#### Register
Registers a service. A provider is allowed to register only its own services. It means that provider system name and certificate common name must match for successful registration.

Status: 100% implemented

#### Unregister
Removes a registered service. A provider is allowed to unregister only its own services. It means that provider system name and certificate common name must match for successful unregistration.

Status: 100% implemented

### Private endpoints

#### Query System
POST /serviceregistry/query/system
This service can only be used by other core services, therefore is not part of the public API.

Status: 85% implemented - WIP

#### Query System by ID
GET /serviceregistry/system/{id}
This service can only be used by other core services, therefore is not part of the public API.

Status: 100% implemented

### Management endpoints

#### Get all entries
GET /serviceregistry/mgmt
Returns a list of Service Registry records. If page and item_per_page are not defined, returns all records.

Status: 90% implemented

#### Add an entry
POST /serviceregistry/mgmt
Creates service registry record and returns the newly created record.

Status: 90% implemented

#### Get an entry by ID
GET /serviceregistry/mgmt/{id}
Returns the Service Registry Entry specified by the ID path parameter.

Status: 0% implemented

#### Replace an entry by ID
Replace an entry by ID

Status: 0% implemented

#### PUT /serviceregistry/mgmt/{id}
PATCH /serviceregistry/mgmt/{id}

Status: 0% implemented

#### Delete an entry by ID
Delete an entry by ID

Status: 0% implemented

#### Get grouped view
GET /serviceregistry/mgmt/grouped
Returns all Service Registry Entries grouped for the purpose of the Management Tools' Service Registry view:
  - autoCompleteData
  - servicesGroupedByServiceDefinition
  - servicesGroupedBySystems

Status: 20% implemented

#### Get Service Registry Entries by Service Definition
GET /serviceregistry/mgmt/servicedef/{serviceDefinition}
Returns a list of Service Registry records specified by the serviceDefinition path parameter. If page and item_per_page are not defined, returns all records.

Status: 5% implemented

#### Get all services
GET /serviceregistry/mgmt/services
Get all services
Returns a list of Service Definition records. If page and item_per_page are not defined, returns all records.

Status: 100% implemented

### Add a service
POST /serviceregistry/mgmt/services
Creates service definition record and returns the newly created record.

Status: 10% implemented

#### Get a service by ID
GET /serviceregistry/mgmt/services/{id}
Returns the Service Definition record specified by the id path parameter.

Status: 100% implemented

#### Replace a service by ID
PUT /serviceregistry/mgmt/services/{id}
Updates and returns the modified Service Definition record specified by the ID path parameter.

Status: 100% implemented

#### Modify a service by ID
PATCH /serviceregistry/mgmt/services/{id}
Updates and returns the modified Service Definition record specified by the ID path parameter.

Status: 100% implemented

#### Delete a service by ID
DELETE /serviceregistry/mgmt/services/{id}
Removes the service definition record specified by the id path parameter.

Status: 100% implemented

#### Get all systems
GET /serviceregistry/mgmt/systems
Returns a list of System records. If page and item_per_page are not defined, it returns all records.

Status: 100% implemented

#### Add a system
POST /serviceregistry/mgmt/systems
Creates a System record and returns the newly created record.

Status: 100% implemented

#### Add a system
POST /serviceregistry/mgmt/systems
Creates a System record and returns the newly created record.

Status: 100% implemented

#### Get a system by ID
GET /serviceregistry/systems/{id}
Returns the System record specified by the ID path parameter.

Status: 100% implemented

#### Replace a system by ID
PUT /serviceregistry/mgmt/systems/{id}
Updates and returns the modified System record specified by the ID path parameter. Not defined fields are going to be updated to "null" value.

Status: 100% implemented

#### Modify a system by ID
PATCH /serviceregistry/mgmt/systems/{id}
Updates and returns the modified system record specified by the id path parameter. Not defined fields are going to be NOT updated.

Status: 100% implemented

#### Delete a system by ID
DELETE /serviceregistry/mgmt/systems/{id}
Removes the System record specified by the ID path parameter.

Status: 100% implemented


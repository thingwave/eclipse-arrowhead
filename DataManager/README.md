# Eclipse Arrowhead DataManager in Go
This is a an Arrowhead DataManager written in Go.


## Endpoints

### Client endpoint descriptions


#### Echo
This endpoint returns a "Got it!" upon a GET request.

Status: 100% implemented

#### Get Proxy system list
GET /datamanager/proxy
Returns a list of all systems that have at least one service endpoint.

Status: 100% implemented

#### Get Proxy service list
GET /datamanager/proxy/{systemName}

Status: 100% implemented

#### Get Historian system list
GET /datamanager/historian
Returns a list of all systems that have at least one service endpoint in the database.

Status: 100% implemented

#### Get Historian service list
GET /datamanager/historian/{systemName}

Status: 100% implemented

#### Fetch data from db
GET /datamanager/historian/{systemName}/{serviceName}
Returns sensor data from a service endpoint from the database.

Status: 100% implemented

### Private endpoint descriptions

### Management endpoint descriptions



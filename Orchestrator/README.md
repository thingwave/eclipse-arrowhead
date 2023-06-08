This is a an Arrowhead Orchestrator written in Go.

## Endpoints

### Client endpoint description


#### Echo
GET /orchestrator//echo
This endpoint returns a "Got it!" upon a GET request.

Status: 100% implemented

#### Orchestration
POST /orchestrator/orchestration

Status: 60% implemented

#### Start store Orchestration by ID
GET /orchestrator//orchestration/{id}

Status 40% implemented

### Private endpoint description
These services can only be used by other core services, therefore they are not part of the public API.

### Management endpoint Description
There endpoints are mainly used by the Management Tool and Cloud Administrators.

####  Get all Store Entries
GET /mgmt/store

Status 100% implemented

####  Add Store Entries
POST /mgmt/store

Status 10% implemented

#### Get Store Entry by ID
GET /mgmt/store/{id}

Status 10% implemented

#### Delete Store Entry by ID
DELETE /mgmt/store/{id}

Status 10% implemented

#### Get Entries by Consumer
POST /mgmt/store/all_by_consumer

Status 10% implemented

#### Get Top Priority Entries
GET /mgmt/store/all_top_priority

Status 10% implemented

#### Modify Priorities
POST /mgmt/store/modify_priorities

Status 10% implemented

## Bugs
1. "serviceUri": "", (missing in Orchstration response, do lockup!!)
# Eclipse Arrowhead Serviceregistry in Go
This is an Arrowhead Authorization system written in Go.


## Endpoints

### Client endpoint description


#### Echo
```
GET /authorization/echo
```
This endpoint returns a "Got it!" upon a GET request.

Status: 100% implemented

#### Get Public Key
```
GET /authorization/publickey
```
Returns the public key of the Authorization core service as a (Base64 encoded) text.

Status: 100% implemented

### Private endpoint description

#### Check an Intracloud rule
```
POST /authorization/intracloud/check
```

Status: 0% implemented

### Management Endpoint Description

#### Get all Intracloud rules
```
GET /authorization/mgmt/intracloud
```

Status: 100% implemented

#### Add Intracloud rules
```
POST /authorization/mgmt/intracloud
```

Status: 60% implemented

#### Get an Intracloud rule by ID
```
GET /authorization/mgmt/intracloud/{id}
```

Status: 100% implemented

#### Delete an Intracloud rule by ID
```
DELETE /authorization/mgmt/intracloud/{id}
```

Status: 100% implemented


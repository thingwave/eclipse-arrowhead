# Eclipse Arrowhead Certificate Authority in Go
This is a an Arrowhead Certificate Authority written in Go.


## Endpoints

### Client endpoint description

#### Echo
This endpoint returns a "Got it!" upon a GET request.

Status: 100% implemented

#### Check certificate validity
Returns whether the given certificate is valid or has been revoked. The client SHALL not trust a revoked certificate.

Status: 10% implemented

### Private endpoints

#### Sign CSR with the Cloud Certificate

Status: 10% implemented

#### Check trusted key

Status: 10% implemented

### Management endpoints

Status: 0% implemented

#### Get issued certificates

Status: 0% implemented

#### Revoke certificate

Status: 0% implemented

#### Get trusted keys

Status: 0% implemented

#### Add trusted key

Status: 0% implemented

#### Delete trusted key

Status: 0% implemented

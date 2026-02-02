#!/bin/bash

grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext localhost:9000 list

# protobuf timestamps require RFC3339 (ISO 8601)
grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "Applications": [
        {"date": "2026-01-30T17:11:47Z", "company": "DoorDash"}
        ]
    }'\
 localhost:9000 maxbear.maxhire.Applications/SetApplications

# List all applications
grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext localhost:9000 maxbear.maxhire.Applications/ListApplications

# List applications filtered by company name
grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "company": "DoorDash"
}' \
localhost:9000 maxbear.maxhire.Applications/ListApplications

grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "company": "DoorDash",
    "status": "PENDING"
}' \
localhost:9000 maxbear.maxhire.Applications/ListApplications


grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "start_date": "2025-10-03T00:00:00Z",
    "end_date": "2025-10-04T00:00:00Z"
}' \
localhost:9000 maxbear.maxhire.Applications/ListApplications
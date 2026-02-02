#!/bin/bash

grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext localhost:9000 list

protobuf timestamps require RFC3339 (ISO 8601)
grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "Applications": [
        {"date": "2026-01-30T17:11:47Z", "company": "DoorDash"}
        ]
    }'\
 localhost:9000 maxbear.maxhire.Applications/SetApplications

grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext localhost:9000 maxbear.maxhire.Applications/ListApplications

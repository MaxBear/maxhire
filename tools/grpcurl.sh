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
    "start_date": "2026-01-28T00:00:00Z",
    "end_date": "2026-01-29T00:00:00Z"
}' \
localhost:9000 maxbear.maxhire.Applications/ListApplications

# Set interviews for GitLab application
grpcurl -emit-defaults -import-path ./proto/applications/v1 -proto applications.proto -plaintext -d \
'{
    "date": "2026-01-28T02:50:10Z",
    "company": "GitLab",
    "interviews": [
        {
            "datetime": "2026-02-02T14:00:00Z",
            "interview_type": "RECRUITER_SCREEN"
        }
    ]
}' \
localhost:9000 maxbear.maxhire.Applications/SetInterviews
#!/bin/bash -eu
set -x
# Generate Go code for proto, exclude vendor directory
find ./proto -type d | \
  while read -r dir
  do    
    # Ignore directories with no proto files.
    ls ${dir}/*.proto > /dev/null 2>&1 || continue
    protoc --go-grpc_out=. --go_out=. ${dir}/*.proto
  done
#!/usr/bin/bash

# WORKROOT=$(pwd)
# cd ${WORKROOT}

# build
docker run --rm -i -v $(pwd)/go:/go dfinity/build-env go run main.go
if [ $? -ne 0 ];
then
    echo "fail to go build"
    exit 1
fi
echo "OK for go build"

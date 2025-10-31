#!/bin/bash

export ANAME="pool-dns"

rm -f ${ANAME}-*-a*

# ubuntu arm64
unset GOARM
export GOOS=linux
export GOARCH=arm64
echo "Building ${ANAME}-${GOOS}-${GOARCH}..."
go build -o ${ANAME}-${GOOS}-${GOARCH}

# for x86
unset GOARM
export GOOS=linux
export GOARCH=amd64
echo "Building ${ANAME}-${GOOS}-${GOARCH}..."
go build -o ${ANAME}-${GOOS}-${GOARCH}

# for x86
export GOOS=darwin
export GOARCH=arm64
echo "Building ${ANAME}-${GOOS}-${GOARCH}..."
go build -o ${ANAME}-${GOOS}-${GOARCH}


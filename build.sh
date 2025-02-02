#!/usr/bin/env bash

export TARGET="x86_64-apple-darwin20.4"
export OSXCROSS="/opt/osxcross"
export SDK_VERSION=11.3
export DARWIN="${OSXCROSS}/target"
export DARWIN_SDK="${DARWIN}/SDK/MacOSX${SDK_VERSION}.sdk"

export PATH="${DARWIN}/bin:${DARWIN_SDK}/bin:${PATH}"
export LDFLAGS="-L${DARWIN_SDK}/lib -mmacosx-version-min=10.10"
export CC="${TARGET}-clang"
export CXX="${TARGET}-clang++"

CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -tags static -ldflags "-s -w"

VERSION=`cat VERSION`
GIT_REV=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/aymerick/kowa/core.Version ${VERSION} -X github.com/aymerick/kowa/core.GitRev ${GIT_REV} -X github.com/aymerick/kowa/core.BuildDate ${BUILD_DATE}"

all: build

build:
	godep go build ${LDFLAGS}

build-from-scratch:
	go get github.com/tools/godep
	godep go build ${LDFLAGS} -o kowa-${VERSION}

VERSION=`cat VERSION`
GIT_REV=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/aymerick/kowa/core.Version ${VERSION} -X github.com/aymerick/kowa/core.GitRev ${GIT_REV} -X github.com/aymerick/kowa/core.BuildDate ${BUILD_DATE}"

# development default target
default: gen build

# development mode
build:
	godep go build ${LDFLAGS} -tags dev

# travis build
build-ci:
	godep go build ${LDFLAGS}

# to build in a docker container
build-from-scratch:
	go get github.com/tools/godep
	godep go build ${LDFLAGS} -o kowa-${VERSION}

test:
	godep go test ./...

gen:
	go-bindata -o core/bindata.go -pkg="core" locales/ mailers/templates/

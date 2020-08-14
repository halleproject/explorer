# VERSION               := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION 			  ?= manual
COMMIT                := $(shell git log -1 --format='%H')
COMMIT_SHA1           :=`git rev-parse HEAD`
TOOLS_DESTDIR         ?= bin/mintscan
#BUILD_FLAGS 		  := -ldflags "-X main.buildHash=${COMMIT} -X util.BuildVersion=${VERSION}"
BUILD_FLAGS           := -tags netgo -ldflags "-X github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.GitCommit=${COMMIT_HASH}  -X  github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.VERSION=0.1.0 -X github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.COMMIT_SHA1=${COMMIT_SHA1}   -X 'github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.BUILD_TIME=`date`'  -X  'github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.GO_VERSION=`go version`' -X 'github.com/cosmostation/mintscan-binance-dex-backend/mintscan/version.GIT_BRANCH=`git symbolic-ref --short -q HEAD`' "

build: go.sum
	@echo "building mintscan binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $(TOOLS_DESTDIR) .
	#GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) -o $(TOOLS_DESTDIR) .

install: go.sum
	@echo "installing mintscan binary..."
	@go install -mod=readonly $(BUILD_FLAGS) .

clean:
	@echo "cleaning mintscan binary..."
	rm -f $(TOOLS_DESTDIR) 2> /dev/null

HUBBLE_VER := "1.18.0"
FIDR_VERSION := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

all: 

proto:
	curl -s -L "https://github.com/farcasterxyz/hub-monorepo/archive/refs/tags/@farcaster/hubble@"${HUBBLE_VER}".tar.gz" \
	| tar -zxvf - -C . --strip-components 2 hub-monorepo--farcaster-hubble-${HUBBLE_VER}/protobufs/schemas/

farcaster-go: $(wildcard schemas/*.proto)
	protoc --proto_path=schemas --go_out=. \
	$(shell cd schemas; ls | xargs -I \{\} echo -n '--go_opt=M'{}=farcaster/" " '--go-grpc_opt=M'{}=farcaster/" " ) \
	--go-grpc_out=. \
	schemas/*.proto

local:
	@echo Building farma v${FIDR_VERSION}
	go build \
	-ldflags "-w -s" \
	-ldflags "-X github.com/vrypan/farma/config.FIDR_VERSION=${FIRD_VERSION}" \
	-o farma

releases:
	goreleaser release --clean


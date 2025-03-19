HUBBLE_VER := "1.19.1"
FARMA_VER := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

PROTO_FILES := $(wildcard schemas/*.proto)
SOURCES := $(wildcard utils/*.go config/*.go apiv2/*.go cmd/*.go fctools/*.go localdb/*.go models/*.go)
GREEN = \033[0;32m
NC = \033[0m

all: farma

.farcaster-built: $(PROTO_FILES)
	@echo -e "$(GREEN)Compiling .proto files...$(NC)"
	protoc --proto_path=schemas --go_out=. \
	$(shell cd schemas; ls | xargs -I \{\} echo -n '--go_opt=M'{}=farcaster/" " '--go-grpc_opt=M'{}=farcaster/" " ) \
	--go-grpc_out=. \
	schemas/*.proto
	@touch .farcaster-built

proto:
	@echo -e "$(GREEN)Downloading proto files (Hubble v$(HUBBLE_VER))...$(NC)"
	curl -s -L "https://github.com/farcasterxyz/hub-monorepo/archive/refs/tags/@farcaster/hubble@$(HUBBLE_VER).tar.gz" \
	| tar -zxvf - -C . --strip-components 2 "hub-monorepo--farcaster-hubble-$(HUBBLE_VER)/protobufs/schemas/"

clean:
	@echo -e "$(GREEN)Cleaning up...$(NC)"
	rm -f $(BINS) farcaster/*.pb.go farcaster/*.pb.gw.go .farcaster-built

.PHONY: all proto clean local release-notes tag tag-minor tag-major releases

farma: .farcaster-built $(SOURCES)
	@echo -e "$(GREEN)Building fcp ${FCP_VER} $(NC)"
	go build -o $@ -ldflags "-w -s -X github.com/vrypan/farma/config.FARMA_VERSION=${FARMA_VER}"

release-notes:
	# Autmatically generate release_notes.md
	./bin/generate_release_notes.sh

models/farma.pb.go:
	protoc --go_out=./models --go_opt=paths=source_relative --proto_path=./models models/farma.proto
tag:
	./bin/auto_increment_tag.sh patch

tag-minor:
	./bin/auto_increment_tag.sh minor

tag-major:
	./bin/auto_increment_tag.sh major

releases: models/farma.pb.go
	goreleaser release --clean

# Get version from git hash
git_hash := $(shell git rev-parse --short HEAD || echo 'development')

# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_hash}'

ifeq ($(GOOS),)
	GOOS := linux
endif

ifeq ($(GOARCH),)
	GOARCH := amd64
endif

ifeq ($(PREFIX),)
	PREFIX := /usr/local
endif

.PHONY:
build:
	@echo "Building binaries..."
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags=${linker_flags} -o=./bin/snxgo-${GOOS}-${GOARCH} ./cmd/snx-connnect/main.go
	ln -sf snxgo-${GOOS}-${GOARCH} bin/snxgo

clean:
	rm -rf ./bin

install:
	cp ./bin/snxgo ${PREFIX}/bin
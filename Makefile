# Get version from git hash
git_hash := $(shell git rev-parse --short HEAD || echo 'development')

# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_hash}'

ifeq ($(PREFIX),)
	PREFIX := /usr/local
endif

.PHONY:
build:
	@echo "Building binaries..."
	go build -ldflags=${linker_flags} -o=./bin/snxgo ./cmd/snx-connnect/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/snxgo-linux-amd64 ./cmd/snx-connnect/main.go

clean:
	rm -rf ./bin

install:
	cp ./bin/snxgo ${PREFIX}/bin
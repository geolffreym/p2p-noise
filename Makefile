# Small make tasks for go
.PHONY: test

# More tools:
# https://github.com/kisielk/godepgraph

USER=geolffreym
PACKAGE=p2p-noise
VERSION=0.1.0

INPUT=./examples/pingpong

BINARY=main
BINARY_WIN=${BINARY}-win
BINARY_OSX=${BINARY}-darwin
BINARY_LINUX=${BINARY}-linux

ARCH_64=amd64
ARCH_32=386

LINUX_64=${BINARY_LINUX}-${ARCH_64}
LINUX_32=${BINARY_LINUX}-${ARCH_32}
WIN_64=${BINARY_WIN}-${ARCH_64}
WIN_32=${BINARY_WIN}-${ARCH_32}
OSX_64=${BINARY_LINUX}-${ARCH_64}


# -count 1 idiomatic no cached testing
# -race test race condition for routines
# @ = dont echo the output
.PHONY: test ## run tests
test:
	@go test -v ./... -count 1 -race -covermode=atomic
	@echo "[OK] test finished"

# Could be compared using
# make benchmark > a.old
# make benchmark > b.new
# benchcmp a.old b.new
.PHONY: benchmark ## run benchmark tests
benchmark: 
	@perflock -governor=80% go test -run=^Benchmarck$  -benchtime 1s -bench=. -count=1
	@echo "[OK] benchmark finished"


# View standard output profiling:
# go tool pprof -top cpu.prof 

# For memory profiling type use:
# inuse_space	Display in-use memory size
# inuse_objects	Display in-use object counts
# alloc_space	Display allocated memory size
# alloc_objects	Display allocated object counts
# eg. go tool pprof --alloc_space -top prof.mem 

# For fancy visualization:
# Could use Graphviz (https://graphviz.org/download/)
# eg. go tool pprof -web bin/main-linux-amd64 cpu.prof
.PHONY: profiling ## run profiling tests
profiling: 
	@perflock -governor=80% go test -run=^Benchmarck$ -benchmem -benchtime 1s -bench=. -cpu 1,2,4,8 -count=1 -memprofile mem.prof -cpuprofile cpu.prof
	@echo "[OK] profiling finished"

.PHONY: coverage ## run tests coverage
coverage:
	@go test -v ./... -race -covermode=atomic -coverprofile coverage ./...
	@echo "[OK] coverage finished"

.PHONY: coverage-export ## run tests coverage export
coverage-export: coverage
	@go tool cover -html=coverage
	@echo "[OK] code test coverage finished"

# Allow to preview documentation.
# Please verify your GOPATH before run this command
.PHONY: preview-doc ## run local documentation server
preview-doc: 
	@godoc -http=localhost:6060 -links=true 

.PHONY: build ## compiles the command into and executable
build:
	@go build -v ./...

imports:
	@goimport

.PHONY: format ## automatically formats Go source cod
format: 
	@go fmt ./...
	@goimports -w .
	@echo "[OK] code format finished"

.PHONY: check ## examines Go source code and reports suspicious constructs
check:
	@go vet -v ./...
	@echo "[OK] code check finished"

.PHONY: clean ## removes generated files and clean go cache
clean:
	@go clean --cache ./... 
	@rm -f mem.prof
	@rm -f prof.mem
	@rm -rf bin
	@echo "[OK] cleaned"

.PHONY: compile-win ## compiles window exec
compile-win:
	@GOOS=windows GOARCH=amd64 go build -o bin/${WIN_64} ${INPUT}
	@GOOS=windows GOARCH=386 go build -o bin/${WIN_32} ${INPUT}

#Go1.15 deprecates 32-bit macOS builds	
# go build -x to show compilation details
#GOOS=darwin GOARCH=386 go build -o bin/main-mac-386 main.go
.PHONY: compile-mac ## compiles mac exec
compile-mac:
	@GOOS=darwin GOARCH=amd64 go build -o bin/${OSX_64} ${INPUT}

.PHONY: compile-linux ## compiles linux exec
compile-linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/${LINUX_64} ${INPUT}
	@GOOS=linux GOARCH=386 go build -o bin/${LINUX_32} ${INPUT}

.PHONY: compile ## compiles all os exec
compile: compile-linux compile-win compile-mac
	@echo "[OK] Compiling for every OS and Platform"

.PHONY: run ## compiles and runs the named main Go package
run: 
	@go run ${INPUT} $(filter-out $@,$(MAKECMDGOALS))

.PHONY: update-pkg-cache ## updated the package cache version
update-pkg-cache:
    GOPROXY=https://proxy.golang.org GO111MODULE=on \
    go get github.com/${USER}/${PACKAGE}@v${VERSION}

# https://go.dev/ref/mod#go-mod-vendor
.PHONY: vendorize ## lock dependencies
vendorize:
	@go mod vendor
	@echo "[OK]"

all: build test check-test-coverage code-check compile

.PHONY: help  ## display this message
help:
	@grep -E \
		'^.PHONY: .*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ".PHONY: |## "}; {printf "\033[36m%-19s\033[0m %s\n", $$2, $$3}'
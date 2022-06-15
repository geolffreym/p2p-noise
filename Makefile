# Small make tasks for go
.PHONY: test

# More tools:
# https://github.com/kisielk/godepgraph

USER=geolffreym
PACKAGE=p2p-noise
VERSION=0.1.0

INPUT=examples/listener.go

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
test:
	@go test -v ./... -count 1 -race -covermode=atomic
	@echo "[OK] test finished"

# Could be compared using
# make benchmark > a.old
# make benchmark > b.new
# benchcmp a.old b.new
benchmark: 
	@go test ./... -bench=. -benchtime 100000x -count 5
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

profiling: 
	@go test -bench=. -benchtime 100000x -run=^$ -cpuprofile=cpu.prof -memprofile=prof.mem
	@echo "[OK] profiling finished"

coverage:
	@go test -v ./... -race -covermode=atomic -coverprofile coverage ./...
	@echo "[OK] coverage finished"
	
coverage-export: coverage
	@go tool cover -html=coverage
	@echo "[OK] code test coverage finished"

build:
	@go build -v ./...

code-fmt: 
	@go fmt ./...
	@echo "[OK] code format finished"

code-check:
	@go vet -v ./...
	@echo "[OK] code check finished"

clean:
	@go clean --cache ./... 
	@rm -f mem.prof
	@rm -f prof.mem
	@rm -rf bin
	@echo "[OK] cleaned"

compile-win:
	@GOOS=windows GOARCH=amd64 go build -o bin/${WIN_64} ${INPUT}
	@GOOS=windows GOARCH=386 go build -o bin/${WIN_32} ${INPUT}

#Go1.15 deprecates 32-bit macOS builds	
#GOOS=darwin GOARCH=386 go build -o bin/main-mac-386 main.go
compile-mac:
	@GOOS=darwin GOARCH=amd64 go build -o bin/${OSX_64} ${INPUT}

compile-linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/${LINUX_64} ${INPUT}
	@GOOS=linux GOARCH=386 go build -o bin/${LINUX_32} ${INPUT}

compile: compile-linux compile-win compile-mac
	@echo "[OK] Compiling for every OS and Platform"
	
run: 
	@go run ${INPUT}

update-pkg-cache:
    GOPROXY=https://proxy.golang.org GO111MODULE=on \
    go get github.com/${USER}/${PACKAGE}@v${VERSION}

all: build test check-test-coverage code-check compile
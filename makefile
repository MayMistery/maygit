# Makefile for cross-compiling the mgit project

BINARY=mgit
VERSION=1.0.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
GOENV=CGO_ENABLED=0

# Build for all platforms
all: windows linux darwin

# Windows
windows:
	${GOENV} GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/${BINARY}-windows-amd64.exe
	${GOENV} GOOS=windows GOARCH=386 go build ${LDFLAGS} -o build/${BINARY}-windows-386.exe

# Linux
linux:
	${GOENV} GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/${BINARY}-linux-amd64
	${GOENV} GOOS=linux GOARCH=386 go build ${LDFLAGS} -o build/${BINARY}-linux-386
	${GOENV} GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o build/${BINARY}-linux-arm5
	${GOENV} GOOS=linux GOARCH=arm GOARM=6 go build ${LDFLAGS} -o build/${BINARY}-linux-arm6
	${GOENV} GOOS=linux GOARCH=arm GOARM=7 go build ${LDFLAGS} -o build/${BINARY}-linux-arm7
	${GOENV} GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o build/${BINARY}-linux-arm64

# macOS
darwin:
	${GOENV} GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/${BINARY}-darwin-amd64
	${GOENV} GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o build/${BINARY}-darwin-arm64

clean:
	rm -rf build/

.PHONY: windows linux darwin clean

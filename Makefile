# Go parameters

# CGO_ENABLED=0   -> Disable interoperate with C libraries -> speed up build time! Enable it, if dependencies use C libraries!
# GOOS=linux      -> compile to linux because scratch docker file is linux
# GOARCH=amd64    -> because, hmm, everthing works fine with 64 bit :)
# -a              -> force rebuilding of packages that are already up-to-date.
# -o gpio-test-x  -> force to build an executable gpio-test-x file (instead of default https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

BUILD_DIRECTORY=build
BINARY_NAME=$(BUILD_DIRECTORY)/ics_enhancer
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WIN=$(BINARY_NAME).exe

all: clean build-osx build-linux build-windows

build-osx:
	go build -a -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(BINARY_LINUX) -v

build-windows:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -o $(BINARY_WIN) -v

clean:
	go clean
	rm -rf $(BUILD_DIRECTORY)

deps:
	go get

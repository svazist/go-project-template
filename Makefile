# This how we want to name the binary output
BINARY=server

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildDate=${DATE}"

export GOOS=linux
export CGO_ENABLED=1
export GOARCH=amd64

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY} ./main.go

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install

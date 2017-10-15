include Makefile.conf

# Binary output name.
BINARY=./bin/$(shell basename `pwd`)

.DEFAULT: all
.DEFAULT_GOAL: $(BINARY)
.PHONY: clean

SDK_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")
SDK_UNIT_TEST_ONLY_PKGS=$(shell go list -tags ${UNIT_TEST_TAGS} ./... | grep -v "/vendor/")
SDK_GO_1_4=$(shell go version | grep "go1.4")
SDK_GO_1_5=$(shell go version | grep "go1.5")
SDK_GO_VERSION=$(shell go version | awk '''{print $$3}''' | tr -d '''\n''')

all: get-deps unit integration build

help:
	@echo "Please use \`make <target>\`, Available options for <target> are:"
	@echo "  build                   to build the project."
	@echo "  unit                    to run unit tests."
	@echo "  integration             to run integration tests."

${BINARY}:
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	for GOOS in $(OSs); do \
		for GOARCH in $(ARCHS); do \
        	env GOOS=$$GOOS GOARCH=$$GOARCH go build -a -v -o ${BINARY}-$$GOOS-$$GOARCH -installsuffix cgo --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"' .; \
		done; \
	done

build: clean \
    vet \
    ${BINARY}

# Cleaning the project, by deleting binaries.
clean:
	@echo "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	for GOOS in $(OSs); do \
        for GOARCH in $(ARCHS); do \
            TARGET_BINARY=${BINARY}-$$GOOS-$$GOARCH; \
            if [ -f $$TARGET_BINARY ] ; then rm -rf $$TARGET_BINARY ; fi \
        done; \
    done
	go clean

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet:
	go vet ./...

verify: get-deps-verify lint vet

lint:
	@echo "go lint SDK and vendor packages"
	@lint=`if [ \( -z "${SDK_GO_1_4}" \) -a \( -z "${SDK_GO_1_5}" \) ]; then  golint ./...; else echo "skipping golint"; fi`; \
	lint=`echo "$$lint" | grep -E -v -e ${LINTIGNOREDOT} -e ${LINTIGNOREDOC} -e ${LINTIGNORECONST} -e ${LINTIGNORESTUTTER} -e ${LINTIGNOREINFLECT} -e ${LINTIGNOREDEPS} -e ${LINTIGNOREINFLECTS3UPLOAD} -e ${LINTIGNOREPKGCOMMENT}`; \
	echo "$$lint"; \
	if [ "$$lint" != "" ] && [ "$$lint" != "skipping golint" ]; then exit 1; fi

get-deps: get-deps-tests get-deps-verify
	@echo "go get SDK dependencies"
	@go get -v $(SDK_ONLY_PKGS)

get-deps-tests:
	@echo "go get SDK testing dependencies"
	go get github.com/stretchr/testify

get-deps-verify:
	@echo "go get SDK verification utilities"
	@if [ \( -z "${SDK_GO_1_4}" \) -a \( -z "${SDK_GO_1_5}" \) ]; then  go get github.com/golang/lint/golint; else echo "skipped getting golint"; fi

# Unit tests
unit:
	go test ./... -v --cover

# Integration tests
integration: get-deps-tests build verify
    go test app/proxy_test.go -cover --tags=integration -v

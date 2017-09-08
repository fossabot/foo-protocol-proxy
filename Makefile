NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Binary output name.
BINARY=./bin/$(shell basename `pwd`)

.DEFAULT_GOAL: $(BINARY)
.PHONY: clean

${BINARY}:
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	go build -o ${BINARY} .

build: clean \
    vet \
    ${BINARY}

# Cleaning the project, by deleting binaries.
clean:
	@echo "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	if [ -f ${BINARY} ] ; then rm -rf ${BINARY} ; fi
	go clean

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet:
	go vet ./...
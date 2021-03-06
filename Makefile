# get name of directory containing this Makefile
# (stolen from https://stackoverflow.com/a/18137056)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
base_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

SERVICE ?= $(base_dir)
BUILDENV :=
BUILDENV += CGO_ENABLED=0
GIT_HASH := $(CIRCLE_SHA1)
ifeq ($(GIT_HASH),)
  GIT_HASH := $(shell git rev-parse HEAD)
endif
LINKFLAGS :=-s -X main.revision=$(GIT_HASH) -extldflags "-static"
TESTFLAGS := -v -cover
LINT_FLAGS :=--disable-all --enable=vet --enable=golint --enable=ineffassign --enable=goconst --enable=gofmt --enable=dupl --enable=gas
LINTER_EXE := gometalinter.v1
LINTER := $(GOPATH)/bin/$(LINTER_EXE)

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
join-with = $(subst $(SPACE),$1,$(strip $2))

LEXC :=
ifdef LINT_EXCLUDE
	LEXC := $(call join-with,|,$(LINT_EXCLUDE))
endif

.DEFAULT_GOAL := rebuild

.PHONY: install_packages
install_packages:
	go get -t -v ./...

.PHONY: update_packages
update_packages:
	go get -u -t -v ./...

$(LINTER):
	go get -u gopkg.in/alecthomas/$(LINTER_EXE)
	$(LINTER) --install

.PHONY: lint
lint: $(LINTER)
ifdef LEXC
	$(LINTER) --exclude '$(LEXC)' $(LINT_FLAGS) ./...
else
	$(LINTER) $(LINT_FLAGS) ./...
endif

.PHONY: clean
clean:
	rm -f $(SERVICE)

# builds our binary
$(SERVICE):
	$(BUILDENV) go build -o $(SERVICE) -a -ldflags '$(LINKFLAGS)' .

.PHONY: test
test:
	$(BUILDENV) go test $(TESTFLAGS) ./...

# remove any existing binary and build a new one
.PHONY: rebuild
rebuild: clean $(SERVICE);

.PHONY: default
default: rebuild ;

.PHONY: all
all: install_packages $(LINTER) lint test rebuild

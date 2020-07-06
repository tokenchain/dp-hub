#!/usr/bin/make -f
VERSION := v$(shell cat version.txt)              # $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
GPG_SIGNING_KEY = ''
export GO111MODULE = on
export COSMOS_SDK_TEST_KEYRING = n
define update_check
 sh update.sh
endef
# process build tags

build_tags = 
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif
ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = \
    -X github.com/cosmos/cosmos-sdk/version.Name=dpChain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=dpd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=dcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/tokenchain/ixo-blockchain/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
SHOWTIMECMD :=  date "+%Y/%m/%d H:%M:%S"

all: lint install
OS=linux


build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpd.exe ./cmd/dpd
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpcli.exe ./cmd/dpcli
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpfaucet.exe ./cmd/dpfaucet
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpd ./cmd/dpd
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpcli ./cmd/dpcli
	go build -mod=readonly $(BUILD_FLAGS) -o build/dpfaucet ./cmd/dpfaucet
endif
	$(update_check)

buildlinux: go.sum
	env GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dpd ./cmd/dpd
	env GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dpcli ./cmd/dpcli
	env GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dpfaucet ./cmd/dpfaucet

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/dpd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/dpcli
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/dpfaucet

sign-release:
	if test -n "$(GPG_SIGNING_KEY)"; then \
	  gpg --default-key $(GPG_SIGNING_KEY) -a \
	      -o SHA256SUMS.sign -b SHA256SUMS; \
	fi; 

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

.PHONY: all build install go.sum

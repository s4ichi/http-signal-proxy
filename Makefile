CURRENT := $(shell pwd)
BINARY = http-signal
PACKAGES = $(shell go list ./...)
BUILDDIR=./build
BINDIR=$(BUILDDIR)/bin
PKGDIR=$(BUILDDIR)/pkg
DISTDIR=$(PKGDIR)/dist

VERSION := $(shell git describe --tags --abbrev=0)

GOXOS := "darwin,linux"
GOXARCH := "386,amd64"
GOXOUTPUT := "$(PKGDIR)/$(VERSION)"

setup: dep
	go get github.com/motemen/gobump/cmd/gobump
	go get -u github.com/tcnksm/ghr
	go get github.com/Songmu/goxz/cmd/goxz

dep:
	@dep ensure -v

build:
	@go build -o $(BINDIR)/$(BINARY)

test:
	@go test -v -parallel=4 $(PACKAGES)

lint:
	@golint $(PACKAGES)

vet:
	@go vet $(PACKAGES)

coverage:
	@go test -v -race -cover -covermode=atomic -coverprofile=coverage.txt $(PACKAGES)

release: package
	ghr $(VERSION) $(GOXOUTPUT)

package:
	goxz -os=$(GOXOS) -arch=$(GOXARCH) -d=$(GOXOUTPUT) -pv=${VERSION}

.PHONY: dep build container push test lint vet coverage package release

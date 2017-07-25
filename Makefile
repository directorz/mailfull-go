GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS?=$(THIS_GOOS)
GOARCH?=$(THIS_GOARCH)
DIR_PKG=$(subst /src/github.com/directorz/mailfull-go,/pkg,$(PWD))
DIR_BUILD=build
DIR_RELEASE=release
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' version.go)))
GITTAG=$(shell git rev-parse --short HEAD)

.PHONY: build build-linux-amd64 build-linux-386 clean

default: build

$(DIR_BUILD)/bin/$(THIS_GOOS)_$(THIS_GOARCH)/dep:
	mkdir -p /tmp/go
	GOPATH=/tmp/go go get -d -v github.com/golang/dep
	GOPATH=/tmp/go go build -v -o $(DIR_BUILD)/bin/$(THIS_GOOS)_$(THIS_GOARCH)/dep github.com/golang/dep/cmd/dep
	rm -rf /tmp/go

dep: $(DIR_BUILD)/bin/$(THIS_GOOS)_$(THIS_GOARCH)/dep

installdeps: dep
	$(DIR_BUILD)/bin/$(THIS_GOOS)_$(THIS_GOARCH)/dep ensure -v

build:
	go build -v -i -ldflags "-X main.gittag=$(GITTAG)" -o $(DIR_BUILD)/mailfull_$(GOOS)_$(GOARCH)/mailfull cmd/mailfull/*.go

.build-docker:
	docker run --rm -v $(DIR_PKG):/go/pkg -v $(PWD):/go/src/github.com/directorz/mailfull-go -w /go/src/github.com/directorz/mailfull-go \
	-e GOOS=$(GOOS) -e GOARCH=$(GOARCH) golang:1.8.3 \
	go build -v -i -ldflags "-X main.gittag=$(GITTAG)" -o $(DIR_BUILD)/mailfull_$(GOOS)_$(GOARCH)/mailfull cmd/mailfull/*.go

build-linux-amd64:
	@$(MAKE) .build-docker GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) .build-docker GOOS=linux GOARCH=386

release: release-linux-amd64 release-linux-386

release-linux-amd64: build-linux-amd64
	@$(MAKE) release-doc release-targz GOOS=linux GOARCH=amd64

release-linux-386: build-linux-386
	@$(MAKE) release-doc release-targz GOOS=linux GOARCH=386

release-doc:
	cp -a README.md doc $(DIR_BUILD)/mailfull_$(GOOS)_$(GOARCH)

release-targz: dir-$(DIR_RELEASE)
	tar zcfp $(DIR_RELEASE)/mailfull_$(GOOS)_$(GOARCH).tar.gz -C $(DIR_BUILD) mailfull_$(GOOS)_$(GOARCH)

dir-$(DIR_RELEASE):
	mkdir -p $(DIR_RELEASE)

release-upload: release-linux-amd64 release-linux-386 release-github-token
	ghr -u directorz -r mailfull-go -t $(shell cat github_token) --replace --draft $(VERSION) $(DIR_RELEASE)

release-github-token: github_token
	@echo "file \"github_token\" is required"

clean:
	-rm -rf $(DIR_BUILD)
	-rm -rf $(DIR_RELEASE)

GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
DIR_BUILD=build
DIR_RELEASE=release
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' version.go)))

.PHONY: build build-linux-amd64 build-linux-386 clean

default: build

installdeps:
	glide install -v

build:
	go build -v -ldflags "-X main.gittag=`git rev-parse --short HEAD`" -o build/mailfull_$(GOOS)_$(GOARCH)/mailfull cmd/mailfull/mailfull.go

build-linux-amd64:
	docker run --rm -v $(PWD):/go/src/github.com/directorz/mailfull-go -w /go/src/github.com/directorz/mailfull-go \
	-e GOOS=linux -e GOARCH=amd64 golang:1.7.4 \
	go build -v -ldflags "-X main.gittag=`git rev-parse --short HEAD`" -o "build/mailfull_linux_amd64/mailfull" cmd/mailfull/mailfull.go

build-linux-386:
	docker run --rm -v $(PWD):/go/src/github.com/directorz/mailfull-go -w /go/src/github.com/directorz/mailfull-go \
	-e GOOS=linux -e GOARCH=386 golang:1.7.4 \
	go build -v -ldflags "-X main.gittag=`git rev-parse --short HEAD`" -o "build/mailfull_linux_386/mailfull" cmd/mailfull/mailfull.go

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

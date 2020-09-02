DATE=`date +%Y%m%d`
SHA=`git rev-parse HEAD | cut -c1-7`
NIGHTLY_VERSION="nightly-${DATE}${SHA}"

.PHONY: install clean

LDFLAGS=-ldflags "-X bitbucket.org/goreorto/benchwell/config.Version=${NIGHTLY_VERSION}"
RELEASELDFLAGS=-ldflags "-X bitbucket.org/goreorto/benchwell/config.Version=${VERSION} -s -w"

build:
	@go build -mod=vendor ${LDFLAGS} -o benchwell bitbucket.org/goreorto/benchwell

install:
	@go install -mod=vendor ${LDFLAGS} bitbucket.org/goreorto/benchwell

release:
	@go build -mod=vendor ${RELEASELDFLAGS} -o benchwell bitbucket.org/goreorto/benchwell
	@upx benchwell

assets: assets/data/*
	go generate assets/*.go

clean:
	rm -f benchwell

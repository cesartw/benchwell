DATE=`date +%Y%m%d`
SHA=`git rev-parse HEAD | cut -c1-7`
NIGHTLY_VERSION="nightly-${DATE}${SHA}"

.PHONY: install clean

LDFLAGS=-ldflags "-X bitbucket.org/goreorto/benchwell/cmd.version=${NIGHTLY_VERSION}"

build:
	@go build -mod=vendor ${LDFLAGS} -o benchwell bitbucket.org/goreorto/benchwell

install:
	@go install -mod=vendor ${LDFLAGS} bitbucket.org/goreorto/benchwell

assets: assets/data/*
	go generate assets/*.go

clean:
	rm -f benchwell

DATE=`date +%Y%m%d`
SHA=`git rev-parse HEAD | cut -c1-7`
NIGHTLY_VERSION="nightly-${DATE}${SHA}"

.PHONY: install clean

LDFLAGS=-ldflags "-X bitbucket.org/goreorto/sqlaid/cmd.version=${NIGHTLY_VERSION}"

build:
	@go build -mod=vendor ${LDFLAGS} -o sqlaid bitbucket.org/goreorto/sqlaid

install:
	@go install -mod=vendor ${LDFLAGS} bitbucket.org/goreorto/sqlaid

assets: assets/data/*
	go generate assets/*.go

clean:
	rm -f sqlaid

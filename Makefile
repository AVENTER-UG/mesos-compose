#Dockerfile vars

#vars
IMAGENAME=mesos-compose
REPO=avhost
TAG=`git describe --tags --abbrev=0`
BRANCH=`git rev-parse --abbrev-ref HEAD`
BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ`
IMAGEFULLNAME=${REPO}/${IMAGENAME}
LASTCOMMIT=$(shell git log -1 --pretty=short | tail -n 1 | tr -d " ")


.PHONY: help build all docs

help:
	    @echo "Makefile arguments:"
	    @echo ""
	    @echo "Makefile commands:"
	    @echo "build"
	    @echo "all"
			@echo "docs"
			@echo "publish"
			@echo ${TAG}

.DEFAULT_GOAL := all

ifeq (${BRANCH}, master) 
	BRANCH=latest
endif

ifneq ($(shell echo $(LASTCOMMIT) | grep -E '^v([0-9]+\.){0,2}(\*|[0-9]+)'),)
	BRANCH=${LASTCOMMIT}
else
	BRANCH=latest
endif

build:
	@echo ">>>> Build Docker"
	@docker build --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${TAG} .

build-bin:
	@echo ">>>> Build binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=${BUILDDATE} -X main.GitVersion=${TAG} -extldflags \"-static\"" .

push:
	@echo ">>>> Publish docker image"
	@docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} --build-arg VERSION_URL=${VERSION_URL} -t ${IMAGEFULLNAME}:latest .
	@docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} --build-arg VERSION_URL=${VERSION_URL} -t ${IMAGEFULLNAME}:${BRANCH} .

update-gomod:
	go get -u
	go mod tidy

docs:
	@echo ">>>> Build docs"
	$(MAKE) -C $@

sboom:
	syft dir:. > sbom.txt
	syft dir:. -o json > sbom.json

seccheck:
	grype --add-cpes-if-none .

go-fmt:
	@gofmt -w .

version:
	@echo ">>>> Generate version file"
	@echo "[{ \"version\":\"${TAG}\", \"builddate\":\"${BUILDDATE}\" }]" > .version.json
	@cat .version.json
	@echo "Saved under .version.json"

check: go-fmt sboom seccheck
all: check build version push

#Dockerfile vars

#vars
IMAGENAME=mesos-compose
REPO=avhost
TAG=`git describe --tags --abbrev=0`
BRANCH=`git rev-parse --abbrev-ref HEAD`
BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ`
IMAGEFULLNAME=${REPO}/${IMAGENAME}
IMAGEFULLNAMEPUB=avhost/${IMAGENAME}

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

build:
	@echo ">>>> Build docker image and publish it to private repo"
	@docker buildx build --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${BRANCH} .

build-bin:
	@echo ">>>> Build binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=${BUILDDATE} -X main.GitVersion=${TAG} -extldflags \"-static\"" .

publish:
	@echo ">>>> Publish docker image"
	@docker tag ${IMAGEFULLNAME}:${BRANCH} ${IMAGEFULLNAMEPUB}:${BRANCH}
	@docker tag ${IMAGEFULLNAME}:${BRANCH} ${IMAGEFULLNAMEPUB}:${TAG}
	@docker push ${IMAGEFULLNAMEPUB}:${BRANCH}
	@docker push ${IMAGEFULLNAMEPUB}:${TAG}

update-precommit:
	@virtualenv --no-site-packages ~/.virtualenv

update-gomod:
	go get -u
	go mod tidy

docs:
	@echo ">>>> Build docs"
	$(MAKE) -C $@

sboom:
	syft dir:. > sbom.txt
	syft dir:. -o json > sbom.json

version:
	@echo ">>>> Generate version file"
	@echo "[{ \"version\":\"${TAG}\", \"builddate\":\"${BUILDDATE}\" }]" > .version.json
	@cat .version.json
	@echo "Saved under .version.json"


all: build version

#Dockerfile vars

#vars
IMAGENAME=mesos-compose
REPO=${registry}
TAG=`git describe`
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
	@docker buildx build --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${BRANCH} --push .

publish:
	@echo ">>>> Publish docker image"
	@docker tag ${IMAGEFULLNAME}:${BRANCH} ${IMAGEFULLNAMEPUB}:${BRANCH}
	@docker push ${IMAGEFULLNAMEPUB}:${BRANCH}

update-precommit:
	@virtualenv --no-site-packages ~/.virtualenv

docs:
	@echo ">>>> Build docs"
	$(MAKE) -C $@

all: build

CONTAINER_BUILDER_DOCKER="docker"
CONTAINER_BUILDER_BUILDKIT="buildkit"
# Defaults
DEFAULT_CONTAINER_REGISTRY="ghcr.io"
DEFAULT_CONTAINER_REGISTRY_REPO="nationaloilwellvarco"
DEFAULT_GOLANG_BASE_IMAGE="golang:1.18.0-alpine3.15"
DEFAULT_HELLOWORLD_IMAGE_NAME="chirp-service"
DEFAULT_HELLOWORLD_IMAGE_VERSION="v1.0.0"
# Globals
CONTAINER_REGISTRY=$(DEFAULT_CONTAINER_REGISTRY)
CONTAINER_REGISTRY_REPO=$(DEFAULT_CONTAINER_REGISTRY_REPO)
CONTAINER_BUILDER=$(CONTAINER_BUILDER_DOCKER)
SHA=""
SHORT_SHA=""
OS_INFO=""
DOCKER_VERSION=""
GOLANG_VERSION=""
HELLOWORLD_IMAGE_NAME=$(DEFAULT_HELLOWORLD_IMAGE_NAME)
HELLOWORLD_IMAGE_VERSION=$(DEFAULT_HELLOWORLD_IMAGE_VERSION)
HELLOWORLD_IMAGE_FULL=""
 
# Check if local or action...
# This is janky but it does the job
ACTION=false
ifdef GITHUB_RUN_NUMBER #Check for env
    ACTION=true
endif
# set based on build enviroment
ifeq (ACTION, true)
    # Action
    SHA=${GITHUB_SHA}
    SHORT_SHA=$(shell git rev-parse --short=4 ${GITHUB_SHA})
else # Local
    SHA=$(shell git log -1 --format=%H)
    SHORT_SHA=$(shell git log -1 --pretty=format:%h)
endif
# set the ncs image name
# the tag always reflects the SDK revision
 
# check if there is a registry
ifdef  CONTAINER_REGISTRY
    HELLOWORLD_IMAGE_FULL="$(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_REPO)/$(HELLOWORLD_IMAGE_NAME):$(HELLOWORLD_IMAGE_VERSION)"
else
    HELLOWORLD_IMAGE_FULL="$(HELLOWORLD_IMAGE_NAME):$(HELLOWORLD_IMAGE_VERSION)"
endif
# Get the versions of required tools and enviroment
# Note: Some of these maybe should be pulled from the dockerfile instead of the host
# or have seperate variables.
OS_INFO=$(shell uname -a)
DOCKER_VERSION=$(shell docker --version 2>/dev/null | cut -d " " -f 3 | cut -d "," -f 1)
GO_VERSION=$(shell go version  2>/dev/null | cut -d " " -f 3)
 
.PHONY: about
about:
    @echo "[Git]"
    @echo "SHA: $(SHA)"
    @echo "SHORT_SHA: $(SHORT_SHA)"
    @echo "IN_ACTION: $(ACTION)"
    @echo ""
    @echo "[Enviroment]"
    @echo "OS_INFO: $(OS_INFO)"
    @echo "DOCKER_VERSION: $(DOCKER_VERSION)"
    @echo "GO_VERSION: $(GO_VERSION)"
    @echo ""
    @echo "[Container]"
    @echo "CONTAINER_REGISTRY: $(CONTAINER_REGISTRY)"
    @echo "CONTAINER_REGISTRY_REPO: $(CONTAINER_REGISTRY_REPO)"
    @echo "CONTAINER_BUILDER: $(CONTAINER_BUILDER)"
    @echo "HELLOWORLD_IMAGE_NAME: $(HELLOWORLD_IMAGE_NAME)"
    @echo "HELLOWORLD_IMAGE_VERSION $(HELLOWORLD_IMAGE_VERSION)"
    @echo "HELLOWORLD_IMAGE_FULL: $(HELLOWORLD_IMAGE_FULL)"
    @echo ""
 
.PHONY: usage
usage:
    @echo "##############################################################################"
    @echo "Usage"
    @echo "about - Logs meta info std out"
    @echo "gomod - Run go mod tidy and vendor"
    @echo "login - To login to the registry"
    @echo "build_amd64 - Builds the amd64 image"
    @echo "build_arm64 - Builds the arm64 image"
    @echo "pull - Pull image"
    @echo "push - Push image"
    @echo "run-docker - Run the server in a container"
    @echo "clean - Removes the images and any dangling images"
    @echo "##############################################################################"
 
.PHONY: mod
gomod:
#Perform Go Mod
    @echo "Go Mod"
    @cd src && go mod tidy && go mod vendor
 
.PHONY: build_amd64
build_amd64:
    @echo "Prepped for: $(HELLOWORLD_IMAGE_FULL)"
# Build for Docker
ifeq ($(CONTAINER_BUILDER),$(CONTAINER_BUILDER_DOCKER))
    @echo "Using Docker"
    @docker build . -f ./docker/Dockerfile.build -t $(HELLOWORLD_IMAGE_FULL) \
    --build-arg GO_ARCH="amd64" \
    --build-arg APP_GIT_COMMIT="$(SHORT_SHA)" \
    --build-arg APP_VERSION="$(DEFAULT_HELLOWORLD_IMAGE_VERSION)"
endif
# Build for Buildkit
ifeq ($(CONTAINER_BUILDER),$(CONTAINER_BUILDER_DOCKER))
    @echo "Using Buildkit"
    @docker buildx build . -f ./docker/Dockerfile.build -t $(HELLOWORLD_IMAGE_FULL) \
    --build-arg GO_ARCH="amd64" \
    --build-arg APP_GIT_COMMIT="$(SHORT_SHA)" \
    --build-arg APP_VERSION="$(DEFAULT_HELLOWORLD_IMAGE_VERSION)"
endif
 
.PHONY: build_arm64
build_arm64:
    @echo "Prepped for: $(HELLOWORLD_IMAGE_FULL)"
# Build for Docker
ifeq ($(CONTAINER_BUILDER),$(CONTAINER_BUILDER_DOCKER))
    @echo "Using Docker"
    @docker build . -f ./docker/Dockerfile.build -t $(HELLOWORLD_IMAGE_FULL) \
    --build-arg GO_ARCH="arm64" \
    --build-arg APP_GIT_COMMIT="$(SHORT_SHA)" \
    --build-arg APP_VERSION="$(DEFAULT_HELLOWORLD_IMAGE_VERSION)"
endif
# Build for Buildkit
ifeq ($(CONTAINER_BUILDER),$(CONTAINER_BUILDER_DOCKER))
    @echo "Using Buildkit"
    @docker buildx build . -f ./docker/Dockerfile.build -t $(HELLOWORLD_IMAGE_FULL) \
    --build-arg GO_ARCH="arm64" \
    --build-arg APP_GIT_COMMIT="$(SHORT_SHA)" \
    --build-arg APP_VERSION="$(DEFAULT_HELLOWORLD_IMAGE_VERSION)"
endif
 
.PHONY: login
login:
#Perform Docker Login
    @echo "Docker Login"
    @echo ${CR_PAT} | docker login ghcr.io -u USERNAME --password-stdin
 
.PHONY: push
push:
# Push the image to the registry
    @echo "Pushing image"
    @docker push $(HELLOWORLD_IMAGE_FULL)
 
.PHONY: pull
pull:
# pull the image from the registry
    @echo "Pulling image"
    @docker pull $(HELLOWORLD_IMAGE_FULL)
 
.PHONY: run-docker
run-docker:
    @echo "Running the server"
    @docker run -it --rm \
    $(HELLOWORLD_IMAGE_FULL)
   
 
.PHONY: clean
clean:
# Remove dangling images
# This might a little risky since it will delete ALL dangling images
# but you shouldn't have x number of <none>...
    @echo "Performing cleanup"
    # @yes | docker builder prune
    # @yes | docker image prune
    @docker rmi -f $(shell docker images -f "dangling=true" -q)
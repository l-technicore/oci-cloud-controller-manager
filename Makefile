# Copyright 2018 Oracle and/or its affiliates. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REGISTRY := iad.ocir.io/spinnaker
PKG := github.com/oracle/oci-cloud-controller-manager
BIN := oci-cloud-controller-manager
IMAGE := $(REGISTRY)/$(BIN)
DOCKER_REGISTRY_TENANCY := spinnaker
DOCKER_REGISTRY_USERNAME := spinnaker/ioana-madalina.patrichi@oracle.com
TEST_IMAGE ?= $(REGISTRY)/$(DOCKER_REGISTRY_TENANCY)/$(BIN)-test


BUILD := $(shell git describe --always --dirty)
# allow overriding for release versions
# Else just equal the build (git hash)
VERSION ?= ${BUILD}
GOOS ?= linux
ARCH ?= amd64

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

# Allows overriding where the CCM should look for the cloud provider config
# when running via make run-dev.
CLOUD_PROVIDER_CFG ?= $$(pwd)/cloud-provider.yaml

RETURN_CODE := $(shell sed --version >/dev/null 2>&1; echo $$?)
ifeq ($(RETURN_CODE),1)
    SED_INPLACE = -i ''
else
    SED_INPLACE = -i
endif

.PHONY: all
all: check test build

.PHONY: gofmt
gofmt:
	@./hack/check-gofmt.sh ${SRC_DIRS}

.PHONY: golint
golint:
	@./hack/check-golint.sh ${SRC_DIRS}

.PHONY: govet
govet:
	@./hack/check-govet.sh ${SRC_DIRS}

.PHONY: check
check: gofmt govet golint

.PHONY: build-dirs
build-dirs:
	@mkdir -p dist/

.PHONY: build
build: build-dirs manifests
	@GOOS=${GOOS} GOARCH=${ARCH} go build     \
	    -i                                    \
	    -o dist/oci-cloud-controller-manager  \
	    -installsuffix "static"               \
	    -ldflags "-X main.version=${VERSION} -X main.build=${BUILD}" \
	    ./cmd/oci-cloud-controller-manager

.PHONY: manifests
manifests: build-dirs
	@cp -a manifests/* dist
	@sed ${SED_INPLACE}                                            \
	    's#${IMAGE}:[0-9]\+.[0-9]\+.[0-9]\+#${IMAGE}:${VERSION}#g' \
	    dist/oci-cloud-controller-manager.yaml

.PHONY: test
test:
	@./hack/test.sh $(SRC_DIRS)

# Deploys the current version to a specified cluster.
# Requires binary, manifests, images to be built and pushed. Requires $KUBECONFIG set.
.PHONY: upgrade
upgrade:
	# Upgrade the current CCM to the specified version
	@./hack/deploy.sh deploy-build-version-ccm

# Deploys the current version to a specified cluster.
# Requires a 'dist/oci-cloud-controller-manager-rollback.yaml' manifest. Requires $KUBECONFIG set.
.PHONY: rollback
rollback:
	# Rollback the current CCM to the specified version
	@./hack/deploy.sh rollback-original-ccm

.PHONY: e2e
e2e:
	@./hack/test-e2e.sh

# Run the canary tests.
.PHONY: canary
canary:
	@./hack/test-canary.sh

# Validate the generated canary test image.
.PHONY: validate-canary
validate-canary:
	@./hack/validate-canary.sh

.PHONY: clean
clean:
	@rm -rf dist

.PHONY: deploy
deploy:
	kubectl -n kube-system set image ds/${BIN} ${BIN}=${IMAGE}:${VERSION}

.PHONY: image
image: build
	docker build -t ${IMAGE}:${VERSION} -f Dockerfile .
	docker build -t ${TEST_IMAGE}:${VERSION} -f Dockerfile.test .

.PHONY: push-dev
push-dev: image
	docker login -u '$(DOCKER_REGISTRY_USERNAME)' -p '$(DOCKER_REGISTRY_PASSWORD)' $(REGISTRY)
	docker push ${IMAGE}:${VERSION}
	docker push ${TEST_IMAGE}:${VERSION}

.PHONY: run-dev
run-dev: build
	@dist/oci-cloud-controller-manager          \
	    --kubeconfig=${KUBECONFIG}              \
	    --cloud-config=${CLOUD_PROVIDER_CFG}    \
	    --cluster-cidr=10.244.0.0/16            \
	    --leader-elect-resource-lock=configmaps \
	    --cloud-provider=oci                    \
	    -v=4

.PHONY: version
version:
	@echo ${VERSION}

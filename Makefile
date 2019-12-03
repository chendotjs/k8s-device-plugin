PACKAGE = github.com/asdfsx/k8s-device-plugin
REGISTRY_DOMAIN =
$(eval COMMIT_HASH=$(shell git rev-parse --short HEAD))
$(eval TAG=$(shell git tag -l --points-at HEAD))
$(eval BRANCH=$(shell git rev-parse --abbrev-ref HEAD))
$(eval BUILD_DATE=$(shell date +%FT%T%z))
LDFLAGS = -ldflags "-X ${PACKAGE}/pkg/version.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/pkg/version.BuildDate=${BUILD_DATE}"

ifdef TAG
        IMAGE_TAG=$(TAG)
else
        IMAGE_TAG=$(BRANCH)
endif

ifdef REGISTRY_DOMAIN
        REGISTRY=$(REGISTRY_DOMAIN)/
endif

.PHONY: build clean verify vendor fmt docker help

build: ## Only Build binarys
	cd bin && \
	go build -mod vendor ${LDFLAGS} ${PACKAGE}/cmd/k8s-device-plugin

clean: ## delete the build target
	rm bin/aios-operator

verify: ## veriy all dependencies
	go mod verify
	go mod tidy

vendor: verify ## put all dependencies into vendor
	go mod vendor

fmt: ## format code
	go fmt ./...

docker: ## build docker image
	docker build -t ${REGISTRY}asdfsx/k8s-device-plugin:$(IMAGE_TAG) .

push: docker
	docker push ${REGISTRY}asdfsx/k8s-device-plugin:$(IMAGE_TAG)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
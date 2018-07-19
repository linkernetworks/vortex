## Folder content generated files
BUILD_FOLDER = ./build

## command
GO           = go
GO_VENDOR    = govendor
MKDIR_P      = mkdir -p

################################################

.PHONY: all
all: build test

.PHONY: pre-build
pre-build:
	$(MAKE) govendor-sync

.PHONY: build
build: pre-build
	$(MAKE) src.build

.PHONY: test
test: build
	$(MAKE) src.test

.PHONY: check
check:
	$(MAKE) check-govendor

.PHONY: clean
clean:
	$(RM) -rf $(BUILD_FOLDER)

## vendor/ #####################################

.PHONY: govendor-sync
govendor-sync:
	$(GO_VENDOR) sync -v

## src/ ########################################

.PHONY: src.build
src.build:
	$(GO) build -v ./src/...
	$(MKDIR_P) $(BUILD_FOLDER)/src/cmd/vortex/
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/vortex/vortex -ldflags="-X github.com/linkernetworks/vortex/src/version.version=1.0.0" ./src/cmd/vortex/...

.PHONY: src.test
src.test:
	$(GO) test -v -race ./src/...

.PHONY: src.install
src.install:
	$(GO) install -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -race -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

.PHONY: src.test-coverage-minikube
src.test-coverage-minikube:
	sed -i.bak "s/localhost:9090/$$(minikube ip):30003/g; s/localhost:27017/$$(minikube ip):31717/g" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json

.PHONY: src.test-coverage-vagrant
src.test-coverage-vagrant:
	sed -i.bak "s/localhost:9090/172.17.8.100:30003/g; s/localhost:27017/172.17.8.100:31717/g" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json

## check build env #############################

.PHONY: check-govendor
check-govendor:
	$(info check govendor)
	@[ "`which $(GO_VENDOR)`" != "" ] || (echo "$(GO_VENDOR) is missing"; false)

## launch apps #############################

.PHONY: apps.init-helm
apps.init-helm:
	helm init
	kubectl create serviceaccount --namespace kube-system tiller
	kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
	kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

.PHONY: apps.launch-apps
apps.launch-apps:
	helm install --debug --wait --set global.environment=testing deploy/helm/foundation/mongodb
	helm install --debug --wait --set global.environment=testing deploy/helm/apps/prometheus/

.PHONY: apps.teardown
apps.teardown:
	helm ls --short | xargs -L1 helm delete

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --tag sdnvortex/vortex:latest --file ./dockerfiles/Dockerfile .

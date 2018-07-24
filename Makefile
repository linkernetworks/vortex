## vortex server version
SERVER_VERSION = v0.1.4
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

.PHONY: env-prepare
env-prepare:
	./prepare.sh prepare

.PHONY: env-k8s
env-k8s:
	./prepare.sh k8s

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
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/vortex/vortex -ldflags="-X github.com/linkernetworks/vortex/src/version.version=$(SERVER_VERSION)" ./src/cmd/vortex/...

.PHONY: src.test
src.test:
	$(GO) test -race ./src/...

.PHONY: src.install
src.install:
	$(GO) install -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	date
	$(GO) test -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	date
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html
	date

.PHONY: src.test-coverage-minikube
src.test-coverage-minikube:
	date
	sed -i.bak "s/localhost:9090/$$(minikube ip):30003/g; s/localhost:27017/$$(minikube ip):31717/g" config/testing.json
	date
	$(MAKE) src.test-coverage
	date
	mv config/testing.json.bak config/testing.json

.PHONY: src.test-coverage-vagrant
src.test-coverage-vagrant:
	sed -i.bak "s/localhost:9090/172.17.8.100:30003/g; s/localhost:27017/172.17.8.100:31717/g" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json

## check build env #############################

.PHONY: src.test-bats
src.test-bats:
	cd tests; \
	bats .;

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
	helm install --name vortex-foundation --debug --set global.environment=testing deploy/helm/foundation
	helm install --name vortex-apps --debug --set global.environment=testing deploy/helm/apps/

.PHONY: apps.teardown
apps.teardown:
	helm ls --short | xargs -i helm delete --purge {}

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --tag sdnvortex/vortex:$(SERVER_VERSION) --file ./dockerfiles/Dockerfile .

## git tag version ########################################

.PHONY: push.tag
push.tag:
	@echo "Current git tag version:"$(SERVER_VERSION)
	git tag $(SERVER_VERSION)
	git push --tags

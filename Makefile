## vortex server version
SERVER_VERSION = v0.1.7
## Folder content generated files
BUILD_FOLDER = ./build
PROJECT_URL  = github.com/linkernetworks/vortex
## command
GO           = go
GO_VENDOR    = govendor
MKDIR_P      = mkdir -p

## Random Alphanumeric String
SECRET_KEY   = $(shell cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

## UNAME
UNAME := $(shell uname)

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
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/vortex/vortex \
	-ldflags="-X $(PROJECT_URL)/src/version.version=$(SERVER_VERSION) -X $(PROJECT_URL)/src/server/backend.SecretKey=$(SECRET_KEY)" \
	./src/cmd/vortex/...

.PHONY: src.test
src.test:
	$(GO) test -v ./src/...

.PHONY: src.install
src.install:
	$(GO) install -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

.PHONY: src.test-coverage-minikube
src.test-coverage-minikube:
	sed -i.bak "s/localhost:9090/$$(minikube ip):30003/g; s/localhost:27017/$$(minikube ip):31717/g" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json
	kubectl exec mongo-0 -n vortex -- mongo vortex_test --eval 'printjson(db.dropDatabase())'

.PHONY: src.test-coverage-vagrant
src.test-coverage-vagrant:
	sed -i.bak "s/localhost:9090/172.17.8.100:30003/g; s/localhost:27017/172.17.8.100:31717/g" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json
	kubectl exec mongo-0 -n vortex -- mongo vortex_test --eval 'printjson(db.dropDatabase())'

## check build env #############################

.PHONY: src.test-bats
src.test-bats:
	sed -i.bak "s/localhost:9090/$$(minikube ip):30003/g; s/localhost:27017/$$(minikube ip):31717/g" config/testing.json
	./build/src/cmd/vortex/vortex -config config/testing.json -port 7890 &
	@cd tests;\
	./test.sh

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

.PHONY: apps.launch
apps.launch:
	yq -y .services deploy/helm/config/development.yaml | helm install --name vortex-services --debug --wait -f - deploy/helm/services
	yq -y .apps deploy/helm/config/development.yaml | helm install --name vortex-apps --debug --wait -f - --set vortex-server.controller.apiserverImageTag=$(SERVER_VERSION) deploy/helm/apps

.PHONY: apps.upgrade
apps.upgrade:
	yq -y .services deploy/helm/config/development.yaml | helm upgrade vortex-services --debug -f - deploy/helm/services
	yq -y .apps deploy/helm/config/development.yaml | helm  upgrade vortex-apps --debug -f - --set vortex-server.controller.apiserverImageTag=$(SERVER_VERSION) deploy/helm/apps

.PHONY: apps.teardown
apps.teardown:
ifeq ($(UNAME), Linux)
	helm ls --short | xargs -i helm delete --purge {}
else ifeq ($(UNAME), Darwin)
	helm ls --short | xargs helm delete --purge
endif

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --tag sdnvortex/vortex:$(SERVER_VERSION) .

## git tag version ########################################

.PHONY: push.tag
push.tag:
	@echo "Current git tag version:"$(SERVER_VERSION)
	git tag $(SERVER_VERSION)
	git push --tags

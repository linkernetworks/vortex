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
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/vortex/vortex ./src/cmd/vortex/...

.PHONY: src.test
src.test:
	$(GO) test -v -race ./src/...

.PHONY: src.install
src.install:
	$(GO) install -v ./src/...

.PHONY: src.test-prometheus
src.test-prometheus:
	kubectl apply -f deploy/kubernetes/apps/monitoring/ -R
	JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl get nodes -o jsonpath="$$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1; done
	#Verify kube-addon-manager.
	#kube-addon-manager is responsible for managing other kubernetes components, such as kube-dns, dashboard, storage-provisioner..
	JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n kube-system get pods -lcomponent=kube-addon-manager -o jsonpath="$$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for kube-addon-manager to be available"; kubectl get pods --all-namespaces; done
	# Wait for kube-dns to be ready.
	JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n kube-system get pods -lk8s-app=kube-dns -o jsonpath="$$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for kube-dns to be available"; kubectl get pods --all-namespaces; done
	#Chck the functrion of prometheus
	until curl --connect-timeout 1 -sL -w "%{http_code}\\n" http://`kubectl get service -n monitoring prometheus -o jsonpath="{.spec.clusterIP}"`:9090/api/v1/query?query=prometheus_build_info -o /dev/null | grep 200; do sleep 1; echo "wait the prometheus to be available"; done
	#Wait the cadvisor to be ready
	JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n monitoring get pods -lname=cadvisor -o jsonpath="$$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for cadvisor to be available"; kubectl get pods --all-namespaces; done

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -race -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

## check build env #############################

.PHONY: check-govendor
check-govendor:
	$(info check govendor)
	@[ "`which $(GO_VENDOR)`" != "" ] || (echo "$(GO_VENDOR) is missing"; false)

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --tag sdnvortex/vortex:latest --file ./dockerfiles/Dockerfile .

# Image URL to use all building/pushing image targets
IMG ?= 192.168.9.22:5000/hypersds-operator:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"
# Name prefix to generate the names of all resources. This value must be the same as 'namePrefix' defined in config/default/kustomization.yaml
NAME_PREFIX ?= hypersds-operator-
KUBE_VERSION ?= 1.19.8
KUBEBUILDER_VERSION ?= 2.3.1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests verify
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Clean all deployed resources
clean: uninstall
	kubectl delete namespace $(NAME_PREFIX)system
	kubectl delete clusterroles.rbac.authorization.k8s.io $(NAME_PREFIX)manager-role $(NAME_PREFIX)metrics-reader $(NAME_PREFIX)proxy-role
	kubectl delete clusterrolebinding.rbac.authorization.k8s.io $(NAME_PREFIX)manager-rolebinding $(NAME_PREFIX)proxy-rolebinding

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run go mod command against code
verify:
	go mod tidy
	go mod verify

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Run a local registry
registry:
ifeq (, $(shell docker ps | grep registry))
	docker run -d -p 5000:5000 --restart=always --name registry registry:2
else
	@echo "local registry is already exists"
endif

# Remove local registry
registry-clean:
	docker container stop registry && docker container rm -v registry

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# Run golangci-lint
lint:
	golangci-lint run ./... -v

# Verify go dependency and generated manifest files
static-test: generate fmt vet manifests verify
	git diff --exit-code

# Run unit test
unit-test:
	go test -v ./... -ginkgo.v -ginkgo.failFast

# Minikube related command
minikube-download:
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x minikube
	sudo mkdir -p /usr/local/bin/
	sudo install minikube /usr/local/bin/
	sudo snap install kubectl --classic
	sudo apt-get install conntrack

minikube-start:
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube start --driver=none --kubernetes-version=v$(KUBE_VERSION)
	sleep 3

minikube-clean:
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube delete

# Kubebuilder related command
kubebuilder-download:
	curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_linux_amd64.tar.gz | tar -xz -C /tmp/
	sudo mv /tmp/kubebuilder_$(KUBEBUILDER_VERSION)_linux_amd64 /usr/local/kubebuilder
	export PATH=$(PATH):/usr/local/kubebuilder/bin

e2e-deploy: registry docker-build docker-push deploy
	hack/e2e.sh bootstrap
	hack/e2e.sh update_cm_after_delete
	hack/e2e.sh update_secret_after_delete
	hack/e2e.sh add_host
	hack/e2e.sh add_disk
	hack/e2e.sh delete_cluster

e2e: e2e-deploy clean
	@echo "e2e completed"


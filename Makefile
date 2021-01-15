# Image URL to use all building/pushing image targets
IMG ?= 192.168.7.16:5000/hypersds-operator:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

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
clean:
	kustomize build config/default | kubectl delete -f -

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
	go test -v ./controllers/... -ginkgo.v -ginkgo.failFast

# Minikube related command
minikube-download:
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x minikube
	sudo mkdir -p /usr/local/bin/
	sudo install minikube /usr/local/bin/
	sudo snap install kubectl --classic
	sudo apt-get install conntrack

minikube-start:
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube start --driver=none --kubernetes-version=v1.19.7
	sleep 3

minikube-clean:
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube delete

# Kubebuilder related command
kubebuilder-download:
	curl -L https://go.kubebuilder.io/dl/2.3.1/linux/amd64 | tar -xz -C /tmp/
	sudo mv /tmp/kubebuilder_2.3.1_linux_amd64 /usr/local/kubebuilder
	export PATH=$(PATH):/usr/local/kubebuilder/bin

define wait-condition
	@cond="${1}"; \
	timeout="${2}"; \
	interval="${3}"; \
	n=0; \
	while [ $${n} -le $${timeout} ] ; do \
		n=`expr $$n + $$interval`; \
		echo "Waiting $$n seconds..."; \
		if [ -z "$${cond}" ]; then echo "Condition is met"; echo $${cond}; true; fi; \
		sleep $$interval; \
		kubectl get cephclusters.hypersds.tmax.io; \
	done; \

	@echo "Timeout"
	@false
endef

e2e-deploy: registry docker-build docker-push deploy
	@echo "deploying cr ..."
	kubectl apply -f config/samples/hypersds_v1alpha1_cephcluster.yaml
	$(call wait-condition, kubectl get cephclusters.hypersds.tmax.io | grep Completed, 18000, 30)

e2e: e2e-deploy clean
	@echo "e2e completed"


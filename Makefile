## Variables
ENVIRONMENT := develop
APIPORT := 8585
DBPORT := 5432
VERSION := 0.0.1
BASE_IMAGE := coffee-no-java
IMAGE_TAG := $(BASE_IMAGE):$(VERSION)
KIND_CLUSTER := coffee-cluser

## Setup
# For Mac Users with HomeBrew Package Manager
cli.setup.mac:
	brew update
	brew list kubectl || brew install kubectl
	brew list kind || brew install kind

# For Windows users on Choclatey Package Manager
cli.setup.windows:
	choco upgrade chocolatey
	choco install kind
	choco install kubernetes-cli


## Docker 
# Build the docker image
build:
	docker buildx build \
	--platform=linux/amd64 \
	--build-arg=BUILD_REF=$(ENVIRONMENT) \
	--build-arg=VERSION=$(VERSION) \
	--build-arg=PORT=$(APIPORT) \
	--build-arg=BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
	--tag=$(IMAGE_TAG) \
	.

## KiND Kubernetes 
# Create a new kind cluster
kind-up:
	kind create cluster \
		--image kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570 \
		--name $(KIND_CLUSTER) \
		--config k8s/kind/kind-config.yml 
	kubectl config set-context --current --namespace=coffee-shop

# Load docker into KiND environment
kind-load:
	kind load docker-image $(IMAGE_TAG) --name $(KIND_CLUSTER)

# Apply kubernetes manifests in k8s/ directory
# Deploy the application into kubernetes
kind-apply-dev:
	kustomize build k8s/base/database | kubectl apply -f -
	kubectl wait --namespace=coffee-shop --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build k8s/base/coffee-api | kubectl apply -f - 


# Production deployment, this increases cpu and memory
kind-apply-prod:
	kustomize build k8s/base/database | kubectl apply -f -
	kubectl wait --namespace=coffee-shop --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build k8s/kind/production | kubectl apply -f -


# Delete currently applied k8s manifests and objects. Needed before 
# running kind-apply again.
kind-delete:
	kubectl delete svc,deployment coffee-api

# Port forward kubernetes svc for localhost:8585/v1/liveness testing
kind-forward:
	kubectl port-forward -n coffee-shop svc/coffee-api $(APIPORT):$(APIPORT)

# Check that kubernetes objects with label app=coffee-api and app=database are up
kind-status:
	kubectl get deployments,pods,svc -l app=coffee-api
	kubectl get deployments,pods,svc -l app=database

# Starting at the last 100 logs, follow logging in realtime.
kind-logs:
	kubectl logs -l app=coffee-api -f --tail=100

# Delete the kind cluster
kind-down:
	kind delete cluster --name $(KIND_CLUSTER)
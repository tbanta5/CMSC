## Variables
ENVIRONMENT := develop
PORT := 8585 # Port needs to be manually adjusted in k8s/coffee-svc.yml and k8s/coffee-api.yml
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
	--build-arg=PORT=$(PORT) \
	--build-arg=BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
	--tag=$(IMAGE_TAG) \
	.

## KiND Kubernetes 
# Create a new kind cluster
kind-up:
	kind create cluster --name $(KIND_CLUSTER)

# Load docker into KiND environment
kind-load:
	kind load docker-image $(IMAGE_TAG) --name $(KIND_CLUSTER)

# Apply kubernetes manifests in k8s/ directory
kind-apply:
	kubectl apply -f k8s/

# Delete currently applied k8s manifests and objects. Needed before 
# running kind-apply again.
kind-delete:
	kubectl delete svc,deployment coffee-api

# Port forward kubernetes svc for localhost:8585/v1/liveness testing
kind-forward:
	kubectl port-forward svc/coffee-api $(PORT):$(PORT)

# Check that kubernetes objects with label app=coffee-api are running
kind-status:
	kubectl get deployments,pods,svc -l app=coffee-api

# Check logs of pods
kind-logs:
	kubectl logs -l app=appstore -f --tail=100

# Delete the kind cluster
kind-down:
	kind delete cluster --name $(KIND_CLUSTER)
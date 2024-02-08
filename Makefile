## Variables
ENVIRONMENT := develop
APIPORT := 8585 # Will also need manual changes in k8s manifests
DBPORT := 5432  # Will also need manual changes in k8s manifests
VERSION := 0.0.1
BASE_IMAGE := coffee-no-java
IMAGE_TAG := $(BASE_IMAGE):$(VERSION)
KIND_CLUSTER := coffee-cluser
KIND_IMAGE := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570 

## Setup
# For Mac Users with HomeBrew Package Manager
cli.setup.mac:
	brew update
	brew list kubectl || brew install kubectl
	brew list kind || brew install kind

# For Windows users with Choclatey Package Manager
# https://community.chocolatey.org/
# Must 'choco install make' in order to run the code.
# Choclatey commands will likely need to be run from an Administrative Console.
cli.setup.windows:
	choco upgrade chocolatey
	choco install kind
	choco install kustomize
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
# the k8s/kind/kind-config.yml file specifies
# host to container port mappings for easy ingress/egress
# of user and application data. 
# Similar to docker host:container port mappings with --expose or -p flags.
kind-up:
	kind create cluster \
		--image $(KIND_IMAGE) \
		--name $(KIND_CLUSTER) \
		--config k8s/kind/kind-config.yml 
	kubectl config set-context --current --namespace=coffee-shop

# Load docker image into KiND environment
kind-load:
	kind load docker-image $(IMAGE_TAG) --name $(KIND_CLUSTER)

# Apply kubernetes manifests in k8s/base directory
# Deploy the application and supporting k8s infrastructure into KiND.
kind-apply-dev:
	kustomize build k8s/base/database | kubectl apply -f -
	kubectl wait --namespace=coffee-shop --timeout=120s --for=condition=Available deployment/database
	kustomize build k8s/base/coffee-api | kubectl apply -f - 


# For production deployments. This increases cpu and memory usage of coffee-api.
# Kustomize allows for dynamic replacement of k8s/base manifest data 
# by "patching" the k8s/kind/production manifests into them.
kind-apply-prod:
	kustomize build k8s/base/database | kubectl apply -f -
	kubectl wait --namespace=coffee-shop --timeout=120s --for=condition=Available deployment/database
	kustomize build k8s/kind/production | kubectl apply -f -

# Delete currently applied k8s manifests and objects.
kind-delete:
	kubectl delete svc,deployment coffee-api
	kubectl delete svc,deployment database

# kind-restart: Rollout and restart new deployment of coffee-api
kind-restart:
	kubectl rollout restart deployment coffee-api

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
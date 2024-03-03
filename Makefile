## Variables
ENVIRONMENT := develop
APIPORT := 8585 # Will also need manual changes in k8s manifests
DBPORT := 5432  # Will also need manual changes in k8s manifests
VERSION := 0.0.1
ADMIN_PASSWD := "p@5fjaskdl45fadkfjl"
BUILD_DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
# ADMIN_PASSWD := 'password123$$' # In make dollar signs must be escaped. $$ = $ for this string.
DB_DSN := "postgres://postgres:pa55word123@localhost:5432/postgres?sslmode=disable"
BASE_IMAGE := coffee-no-java
IMAGE_TAG := $(BASE_IMAGE):$(VERSION)
KIND_CLUSTER := coffee-cluser
KIND_IMAGE := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570 

# ==============================================================================
# Environment Setup
#
#	Having brew installed will simplify the process of installing all the tooling.
#
#	Run this command to install brew on your machine. This works for Linux, Mac and Windows.
#	The script explains what it will do and then pauses before it does it.
#	$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
#
#	WINDOWS MACHINES
#	These are extra things you will most likely need to do after installing brew.
#   If you prefer, we have provided Choclatey Package Manager install steps for Windows users also.
#
# 	Run these three commands in your terminal to add Homebrew to your PATH:
# 	Replace <name> with your username.
#	$ echo '# Set PATH, MANPATH, etc., for Homebrew.' >> /home/<name>/.profile
#	$ echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/<name>/.profile
#	$ eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
#
# 	Install Homebrew's dependencies:
#	$ sudo apt-get install build-essential

# For Mac/Windows Users with HomeBrew Package Manager
setup.mac:
	brew update
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list kind || brew install kind

# For Windows users with Choclatey Package Manager
# https://community.chocolatey.org/
# Must 'choco install make' in order to run the code.
# Choclatey commands will likely need to be run from an Administrative Console.
setup.windows:
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
	--build-arg=DB_DSN=$(DB_DSN) \
	--build-arg=ADMIN_PASSWD=$(ADMIN_PASSWD) \
	--build-arg=BUILD_DATE=$(BUILD_DATE) \
	--tag=$(IMAGE_TAG) \
	.
# ==============================================================================
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

# Development Mode: Apply kubernetes manifests in k8s/base directory
# Deploy the application and supporting k8s infrastructure into KiND.
kind-apply:
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
# If kind-restart doesn't work, this will clear all objects.
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

# ==============================================================================
# Load Testing
#
# Optional tooling.
# To install run: 
# `[brew | choco] install hey`
#
# Average response time is reported to terminal.
# Note: This tests round-trip to the database and back
#       as hey doesn't provide support for sessions. This
# 		suggests that functionality is faster for users.
# 
# Issue HTTP GET request 10,000 times across 100 concurrent workloads.
load-test:
	hey -m GET -c 100 -n 10000 http://localhost:8585/coffee

# Load Test Example output: 
# 	Summary:
#   Total:	24.9957 secs
#   Slowest:	0.7804 secs
#   Fastest:	0.0052 secs
#   Average:	0.2493 secs
#   Requests/sec:	400.0696

#   Total data:	4800000 bytes
#   Size/request:	480 bytes

# Response time histogram:
#   0.005 [1]	|
#   0.083 [9]	|
#   0.160 [288]	|■■
#   0.238 [5111]|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
#   0.315 [4126]|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
#   0.393 [108]	|■
#   0.470 [239]	|■■
#   0.548 [64]	|■
#   0.625 [31]	|
#   0.703 [22]	|
#   0.780 [1]	|


# Latency distribution:
#   10% in 0.1986 secs
#   25% in 0.2028 secs
#   50% in 0.2101 secs
#   75% in 0.2953 secs
#   90% in 0.3010 secs
#   95% in 0.3087 secs
#   99% in 0.4886 secs

# Details (average, fastest, slowest):
#   DNS+dialup:	0.0001 secs, 0.0052 secs, 0.7804 secs
#   DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0037 secs
#   req write:	0.0000 secs, 0.0000 secs, 0.0033 secs
#   resp wait:	0.2492 secs, 0.0051 secs, 0.7734 secs
#   resp read:	0.0000 secs, 0.0000 secs, 0.0045 secs

# Status code distribution:
#   [200]	10000 responses
#  END SAMPLE OUTPUT 
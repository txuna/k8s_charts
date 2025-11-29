## VARIABLE
KIND_CLUSTER		= research
OS					= $(shell uname -s)
target				= local
DOCKER_BUILD_FLAG 	= build
REGISTRY			= localhost:5001


STAGE = $(target)
ifeq ($(OS), Darwin) # macOS
	ifeq ($(STAGE), local) # macOS and local
		DOCKER_BUILD_FLAG = build --output=type=docker
	else # macOS but not local
		DOCKER_BUILD_FLAG = buildx build --platform linux/amd64 --output=type=docker
	endif
endif

ifeq ($(STAGE), local)
	REGISTRY = localhost:5001
endif

include Makefile.image.mk Makefile.deploy.mk

CURRENT_CONTEXT=$(shell kubectl config current-context)


.PHONY: create-cluster
create-cluster:
	sh ./kind/create-cluster.sh $(KIND_CLUSTER)

# .PHONY: deploy-cilium
# deploy-cilium:
# 	@helm install cilium cilium/cilium --version 1.18.4 \
#    		--namespace kube-system \
# 		--set image.pullPolicy=IfNotPresent \
# 		--set ipam.mode=kubernetes

# remove-cilium:
# 	@helm uninstall cilium -n kube-system


.PHONY: helm-update
helm-update:
	@helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	@helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
	@helm repo add gitea-charts https://dl.gitea.com/charts/
	@helm repo add argo https://argoproj.github.io/argo-helm
	@helm repo add cilium https://helm.cilium.io/
	@helm repo add nats https://nats-io.github.io/k8s/helm/charts/
	@helm repo update
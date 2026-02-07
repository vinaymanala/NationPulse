# Makefile to run the nationpulse project with just one command locally
# Variables
CLUSTER_NAME=nationpulse-cluster
K8S_DIR=./k8s
SERVICES=nationpulse-bff nationpulse-ingestion nationpulse-reporting nationpulse-cronjob
run: start
	
.PHONY setup up down build load deploy restart logs clean
setup: ## Full first-time setup: Cluster, Infra, and Deploy
	@make up
	kind create cluster --name $(CLUSTER_NAME)
	@make build
	@make load
	@make deploy
	@echo "🚀 Setup complete. Run 'make tunnel' to access the API."

up: # Start Infrastructyre (kafka, redis, postgres) 
	docker compose up -d --build postgres redis kafka-1 kafka-2 kafka-3 kafka-init
# 	docker compose logs -f 

down: # Stop Infrastructure
	docker compose down

build: ## Build all Go service images
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		docker build -t $$service:latest ./$$service; \
	done

load: ## Load images into the kind cluster
	@for service in $(SERVICES); do \
		echo "Loading $$service into kind..."; \
		kind load docker-image $$service:latest --name $(CLUSTER_NAME); \
	done

deploy: ## Apply Kubernetes manifests
	kubectl apply -f $(K8S_DIR)/infra-config.yaml
	kubectl apply -f $(K8S_DIR)/

restart: ## Force a rollout restart of all deployments
	@for service in $(SERVICES); do \
		kubectl rollout restart deployment/$$service-service || true; \
	done

tunnel: ## Create a tunnel to the BFF service (Keep this running)
	kubectl port-forward svc/bff-service 8081:8081

logs: ## Follow logs for the BFF service
	kubectl logs -f deployment/bff-service deployments/kubepulse-service

clean: ## Nuke the cluster and the infra
	kind delete cluster --name $(CLUSTER_NAME)
	docker-compose down -v
APP_NAME=go-banking-api
IMAGE_NAME=go-banking-api:local
KIND_CLUSTER=go-banking
NAMESPACE=go-banking

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make docker-build        Build backend Docker image"
	@echo "  make kind-load           Load backend image into kind"
	@echo "  make k8s-apply           Apply local Kubernetes manifests"
	@echo "  make k8s-reset           Delete and recreate namespace"
	@echo "  make k8s-migrate         Run Kubernetes migration job"
	@echo "  make k8s-restart         Restart backend deployment"
	@echo "  make k8s-status          Show pods, services, jobs"
	@echo "  make k8s-logs            Show backend logs"
	@echo "  make k8s-migration-logs  Show migration job logs"
	@echo "  make k8s-forward-api     Port-forward backend API"
	@echo "  make k8s-forward-db      Port-forward PostgreSQL"
	@echo "  make k8s-clean           Delete namespace"

.PHONY: docker-build
docker-build:
	cd backend && docker build -t $(IMAGE_NAME) .

.PHONY: kind-load
kind-load:
	kind load docker-image $(IMAGE_NAME) --name $(KIND_CLUSTER)

.PHONY: k8s-apply
k8s-apply:
	kubectl apply -f k8s/local/namespace.yaml
	kubectl apply -f k8s/local/config-map.yaml
	kubectl apply -f k8s/local/secret.yaml
	kubectl apply -f k8s/local/postgres.yaml
	kubectl apply -f k8s/local/backend.yaml

.PHONY: k8s-reset
k8s-reset:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true
	@sleep 5
	kubectl apply -f k8s/local/namespace.yaml
	@sleep 2
	kubectl apply -f k8s/local/config-map.yaml
	kubectl apply -f k8s/local/secret.yaml
	kubectl apply -f k8s/local/postgres.yaml
	kubectl apply -f k8s/local/backend.yaml

.PHONY: k8s-migrate
k8s-migrate:
	kubectl delete job banking-migration -n $(NAMESPACE) --ignore-not-found=true
	kubectl apply -f k8s/local/migration-job.yaml

.PHONY: k8s-restart
k8s-restart:
	kubectl rollout restart deployment banking-api -n $(NAMESPACE)

.PHONY: k8s-status
k8s-status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
	kubectl get jobs -n $(NAMESPACE)

.PHONY: k8s-logs
k8s-logs:
	kubectl logs -n $(NAMESPACE) deploy/banking-api --tail=100 -f

.PHONY: k8s-migration-logs
k8s-migration-logs:
	kubectl logs -n $(NAMESPACE) job/banking-migration

.PHONY: k8s-forward-api
k8s-forward-api:
	kubectl port-forward -n $(NAMESPACE) svc/banking-api 8080:8080

.PHONY: k8s-forward-db
k8s-forward-db:
	kubectl port-forward -n $(NAMESPACE) svc/postgres 5435:5432

.PHONY: k8s-clean
k8s-clean:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true

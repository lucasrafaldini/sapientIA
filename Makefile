SHELL := /bin/zsh

.PHONY: roadmap-bootstrap project-v0.1 project-configure project-sync roadmap-all
.PHONY: build dev test clean install

# Roadmap e Project
roadmap-bootstrap:
	@./scripts/bootstrap_labels.sh
	@./scripts/bootstrap_milestone_v0_1.sh
	@./scripts/create_issues_v0_1.sh

project-v0.1:
	@./scripts/create_project_v0_1.sh

project-configure:
	@./scripts/configure_project_fields.sh

project-sync:
	@./scripts/sync_milestone_to_project.sh

roadmap-all: roadmap-bootstrap project-v0.1 project-configure
	@echo "Roadmap (v0.1) e Project configurados com campos personalizados."

# Build e desenvolvimento
build:
	@echo "🔨 Building Go CLI..."
	@go build -o bin/sapientia ./cmd/sapientia
	@echo "✅ Build completo: bin/sapientia"

install:
	@echo "📦 Installing sapientia..."
	@go install ./cmd/sapientia
	@echo "✅ Instalado: $(shell go env GOPATH)/bin/sapientia"

dev:
	@echo "🚀 Starting dev environment..."
	@docker compose up --build

dev-worker:
	@echo "🐍 Starting Python worker (dev mode)..."
	@cd worker && uvicorn api.main:app --reload --port 8001

test:
	@echo "🧪 Running tests..."
	@go test -v ./...
	@cd worker && pytest -v

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf runs/*
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@cd ui && rm -rf .next node_modules 2>/dev/null || true
	@echo "✅ Clean complete"

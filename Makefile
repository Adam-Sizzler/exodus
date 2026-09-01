# ==============================================================================
# Exodus Management Makefile
# ==============================================================================

.PHONY: help \
        openapi swagger \
        backend-build backend-test backend-lint \
        frontend-build frontend-typecheck frontend-lint frontend-format \
        contract-build contract-sync \
        release-patch release-minor release-major release-tag \
        docker-build docker-up docker-down docker-restart docker-logs

# Colors
CYAN  := $(shell printf '\033[36m')
GREEN := $(shell printf '\033[32m')
YELLOW:= $(shell printf '\033[33m')
RESET := $(shell printf '\033[0m')

# Default target: show help
help:
	@echo ""
	@echo "$(CYAN)Exodus Project Automation$(RESET)"
	@echo ""
	@echo "  $(YELLOW)API & Specs:$(RESET)"
	@echo "    $(GREEN)make openapi$(RESET)          - Generate OpenAPI/Swagger docs via swag"
	@echo ""
	@echo "  $(YELLOW)Backend (Go):$(RESET)"
	@echo "    $(GREEN)make backend-build$(RESET)    - Build backend Go binary"
	@echo "    $(GREEN)make backend-test$(RESET)     - Run all backend Go tests"
	@echo "    $(GREEN)make backend-lint$(RESET)     - Run go vet static analysis"
	@echo ""
	@echo "  $(YELLOW)Frontend (React/Vite):$(RESET)"
	@echo "    $(GREEN)make frontend-build$(RESET)   - Build frontend production bundle (dist/)"
	@echo "    $(GREEN)make frontend-typecheck$(RESET)- Run TypeScript typecheck without emit"
	@echo "    $(GREEN)make frontend-lint$(RESET)    - Run oxfmt and oxlint checks"
	@echo "    $(GREEN)make frontend-format$(RESET)  - Auto-format frontend code with oxfmt"
	@echo ""
	@echo "  $(YELLOW)Contracts:$(RESET)"
	@echo "    $(GREEN)make contract-build$(RESET)   - Build vendor/@exodus/backend-contract"
	@echo "    $(GREEN)make contract-sync$(RESET)    - Sync vendor lockfiles & deduplicate Zod 4"
	@echo ""
	@echo "  $(YELLOW)Releases & Tags:$(RESET)"
	@echo "    $(GREEN)make release-patch$(RESET)    - Bump patch tag (e.g. v26.9.1 -> v26.9.2) & push"
	@echo "    $(GREEN)make release-minor$(RESET)    - Bump minor tag (e.g. v26.9.1 -> v26.10.0) & push"
	@echo "    $(GREEN)make release-major$(RESET)    - Bump major tag (e.g. v26.9.1 -> v27.0.0) & push"
	@echo "    $(GREEN)make release-tag TAG=vX.Y.Z$(RESET) - Create and push specific custom tag"
	@echo ""
	@echo "  $(YELLOW)Docker:$(RESET)"
	@echo "    $(GREEN)make docker-build$(RESET)     - Build local Docker image with cache bust"
	@echo "    $(GREEN)make docker-up$(RESET)        - Start Docker Compose services"
	@echo "    $(GREEN)make docker-down$(RESET)      - Stop Docker Compose services"
	@echo "    $(GREEN)make docker-restart$(RESET)   - Restart exodus panel container"
	@echo "    $(GREEN)make docker-logs$(RESET)      - Follow live logs of exodus container"
	@echo ""

# ------------------------------------------------------------------------------
# OpenAPI / Swagger Generation
# ------------------------------------------------------------------------------
openapi:
	@echo "$(CYAN)Generating OpenAPI / Swagger documentation...$(RESET)"
	@cd backend && swag init --generalInfo main.go --dir . --output internal/httpapi/panelsettings/docs --outputTypes json,yaml --parseDependency --parseInternal
	@echo "$(GREEN)Swagger JSON & YAML generated at backend/internal/httpapi/panelsettings/docs/$(RESET)"

swagger: openapi

# ------------------------------------------------------------------------------
# Backend Commands
# ------------------------------------------------------------------------------
backend-build:
	@echo "$(CYAN)Building backend...$(RESET)"
	@cd backend && go build -o /tmp/exodus-test-bin .
	@rm -f /tmp/exodus-test-bin
	@echo "$(GREEN)Backend builds cleanly.$(RESET)"

backend-test:
	@echo "$(CYAN)Running backend tests...$(RESET)"
	@cd backend && go test ./...

backend-lint:
	@echo "$(CYAN)Linting backend with go vet...$(RESET)"
	@cd backend && go vet ./...

# ------------------------------------------------------------------------------
# Frontend Commands
# ------------------------------------------------------------------------------
frontend-build:
	@echo "$(CYAN)Building frontend production bundle...$(RESET)"
	@cd frontend && npm run cb
	@echo "$(GREEN)Frontend bundle built in frontend/dist/$(RESET)"

frontend-typecheck:
	@echo "$(CYAN)Checking frontend TypeScript types...$(RESET)"
	@cd frontend && npm run typecheck

frontend-lint:
	@echo "$(CYAN)Running frontend linter...$(RESET)"
	@cd frontend && npm run check

frontend-format:
	@echo "$(CYAN)Formatting frontend code...$(RESET)"
	@cd frontend && npm run format:fix

# ------------------------------------------------------------------------------
# Contracts Commands
# ------------------------------------------------------------------------------
contract-build:
	@echo "$(CYAN)Building backend contracts...$(RESET)"
	@npm --prefix frontend/vendor/@exodus/backend-contract run build
	@echo "$(GREEN)Contracts built successfully.$(RESET)"

contract-sync:
	@echo "$(CYAN)Synchronizing vendor lockfiles & deduplicating Zod 4...$(RESET)"
	@cd frontend && \
		rm -f package-lock.json && \
		find vendor/@exodus -name "node_modules" -exec rm -rf {} + && \
		find vendor/@exodus -name "package-lock.json" -exec rm -f {} + && \
		npm install --legacy-peer-deps && \
		npm ls zod
	@echo "$(GREEN)Zod 4 lockfile synchronized.$(RESET)"

# ------------------------------------------------------------------------------
# Release & Tag Automation
# ------------------------------------------------------------------------------
release-patch:
	@git fetch --tags --force 2>/dev/null || true
	@LATEST=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v26.9.1") && \
	CLEAN=$${LATEST#v} && \
	IFS="." read -r MAJOR MINOR PATCH <<< "$$CLEAN" && \
	NEXT_TAG="v$${MAJOR}.$${MINOR}.$$((PATCH + 1))" && \
	echo "$(CYAN)Current tag: $$LATEST -> New patch tag: $$NEXT_TAG$(RESET)" && \
	git tag "$$NEXT_TAG" -m "Release $$NEXT_TAG" && \
	git push origin "$$NEXT_TAG" && \
	echo "$(GREEN)Successfully created and pushed $$NEXT_TAG! GitHub Actions release started.$(RESET)"

release-minor:
	@git fetch --tags --force 2>/dev/null || true
	@LATEST=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v26.9.1") && \
	CLEAN=$${LATEST#v} && \
	IFS="." read -r MAJOR MINOR PATCH <<< "$$CLEAN" && \
	NEXT_TAG="v$${MAJOR}.$$((MINOR + 1)).0" && \
	echo "$(CYAN)Current tag: $$LATEST -> New minor tag: $$NEXT_TAG$(RESET)" && \
	git tag "$$NEXT_TAG" -m "Release $$NEXT_TAG" && \
	git push origin "$$NEXT_TAG" && \
	echo "$(GREEN)Successfully created and pushed $$NEXT_TAG! GitHub Actions release started.$(RESET)"

release-major:
	@git fetch --tags --force 2>/dev/null || true
	@LATEST=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v26.9.1") && \
	CLEAN=$${LATEST#v} && \
	IFS="." read -r MAJOR MINOR PATCH <<< "$$CLEAN" && \
	NEXT_TAG="v$$((MAJOR + 1)).0.0" && \
	echo "$(CYAN)Current tag: $$LATEST -> New major tag: $$NEXT_TAG$(RESET)" && \
	git tag "$$NEXT_TAG" -m "Release $$NEXT_TAG" && \
	git push origin "$$NEXT_TAG" && \
	echo "$(GREEN)Successfully created and pushed $$NEXT_TAG! GitHub Actions release started.$(RESET)"

release-tag:
	@if [ -z "$(TAG)" ]; then \
		echo "$(YELLOW)Error: Please specify TAG=vX.Y.Z (e.g. make release-tag TAG=v26.9.2)$(RESET)"; \
		exit 1; \
	fi
	@echo "$(CYAN)Creating and pushing tag $(TAG)...$(RESET)"
	@git tag "$(TAG)" -m "Release $(TAG)"
	@git push origin "$(TAG)"
	@echo "$(GREEN)Successfully pushed $(TAG)! GitHub Actions release started.$(RESET)"

# ------------------------------------------------------------------------------
# Docker Compose Helpers
# ------------------------------------------------------------------------------
docker-build:
	@echo "$(CYAN)Building Exodus Docker image...$(RESET)"
	@docker compose build --build-arg BUILD_BUST=$$(date +%s) exodus

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

docker-restart:
	@docker compose restart exodus

docker-logs:
	@docker compose logs -f exodus

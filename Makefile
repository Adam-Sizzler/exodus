# ==============================================================================
# Exodus Management Makefile
# ==============================================================================

.PHONY: help \
        openapi swagger \
        backend-build backend-test backend-lint \
        frontend-build frontend-typecheck frontend-lint frontend-format \
        contract-build contract-sync \
        release \
        deploy logs

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
	@echo "    $(GREEN)make openapi$(RESET)                 - Generate OpenAPI/Swagger docs (go generate ./...)"
	@echo ""
	@echo "  $(YELLOW)Backend (Go):$(RESET)"
	@echo "    $(GREEN)make backend-build$(RESET)           - Build backend Go binary"
	@echo "    $(GREEN)make backend-test$(RESET)            - Run backend Go tests"
	@echo "    $(GREEN)make backend-lint$(RESET)            - Run go vet static analysis"
	@echo ""
	@echo "  $(YELLOW)Frontend (React/Vite):$(RESET)"
	@echo "    $(GREEN)make frontend-build$(RESET)          - Build frontend production bundle (dist/)"
	@echo "    $(GREEN)make frontend-typecheck$(RESET)      - Run TypeScript typecheck"
	@echo "    $(GREEN)make frontend-lint$(RESET)           - Run oxfmt and oxlint checks"
	@echo "    $(GREEN)make frontend-format$(RESET)         - Auto-format frontend code with oxfmt"
	@echo ""
	@echo "  $(YELLOW)Contracts:$(RESET)"
	@echo "    $(GREEN)make contract-build$(RESET)          - Build vendor/@exodus/backend-contract"
	@echo "    $(GREEN)make contract-sync$(RESET)           - Sync vendor lockfiles & deduplicate Zod 4"
	@echo ""
	@echo "  $(YELLOW)Releases & GitHub Actions:$(RESET)"
	@echo "    $(GREEN)make release$(RESET)                 - Create & push release tag for today (vYY.M.D)"
	@echo "    $(GREEN)make release TAG=v26.9.1.1$(RESET)   - Create & push specific release tag"
	@echo ""
	@echo "  $(YELLOW)Local Deployment:$(RESET)"
	@echo "    $(GREEN)make deploy$(RESET)                  - Build with current branch/date & recreate exodus container"
	@echo "    $(GREEN)make deploy TAG=v26.9.1.1$(RESET)    - Deploy locally with specific version tag"
	@echo "    $(GREEN)make logs$(RESET)                    - Follow live logs of exodus container"
	@echo ""

# ------------------------------------------------------------------------------
# OpenAPI / Swagger Generation (via go generate)
# ------------------------------------------------------------------------------
openapi:
	@echo "$(CYAN)Generating OpenAPI / Swagger documentation via go generate...$(RESET)"
	@cd backend && go generate ./...
	@echo "$(GREEN)Swagger JSON & YAML generated successfully at backend/internal/httpapi/panelsettings/docs/$(RESET)"

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
release:
	@if [ -z "$(TAG)" ] && [ -z "$(v)" ]; then \
		echo "$(YELLOW)Error: Please specify TAG=vX.Y.Z (e.g. make release TAG=v26.9.2.1)$(RESET)"; \
		exit 1; \
	fi
	@TARGET_TAG=$$(if [ -n "$(TAG)" ]; then echo "$(TAG)"; else echo "$(v)"; fi) && \
	echo "$(CYAN)1/5 Syncing dev with origin...$(RESET)" && \
	git checkout dev && git pull origin dev && \
	echo "$(CYAN)2/5 Merging dev into main for release...$(RESET)" && \
	git checkout main && git pull origin main && \
	git merge dev -m "chore(release): merge dev into main for $$TARGET_TAG" && \
	echo "$(CYAN)3/5 Tagging $$TARGET_TAG on main...$(RESET)" && \
	git tag "$$TARGET_TAG" -m "Release $$TARGET_TAG" && \
	echo "$(CYAN)4/5 Pushing main and $$TARGET_TAG to origin...$(RESET)" && \
	git push origin main && \
	git push origin "$$TARGET_TAG" && \
	echo "$(CYAN)5/5 Returning to dev branch...$(RESET)" && \
	git checkout dev && \
	echo "$(GREEN)Successfully released $$TARGET_TAG on main! GitHub Actions Release workflow started.$(RESET)"

# ------------------------------------------------------------------------------
# Local Docker Deployment
# ------------------------------------------------------------------------------
deploy:
	@BRANCH=$$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "dev") && \
	COMMIT=$$(git rev-parse HEAD 2>/dev/null || echo "unknown") && \
	VERSION=$$(if [ -n "$(TAG)" ]; then \
		echo "$(TAG)"; \
	elif [ -n "$(v)" ]; then \
		echo "$(v)"; \
	elif [ -n "$(VERSION)" ]; then \
		echo "$(VERSION)"; \
	else \
		echo "v$$(date +%y).$$(date +%-m).$$(date +%-d)"; \
	fi) && \
	BUILD_TIME=$$(date -u +'%Y-%m-%dT%H:%M:%SZ') && \
	echo "$(CYAN)Building & deploying Exodus locally (Version: $$VERSION, Branch: $$BRANCH, Commit: $${COMMIT:0:8})...$(RESET)" && \
	docker compose build \
		--build-arg BUILD_BUST=$$(date +%s) \
		--build-arg BRANCH=$$BRANCH \
		--build-arg __EX_METADATA_GIT_BRANCH=$$BRANCH \
		--build-arg __EX_METADATA_VERSION=$$VERSION \
		--build-arg __EX_METADATA_GIT_BACKEND_COMMIT=$$COMMIT \
		--build-arg __EX_METADATA_GIT_FRONTEND_COMMIT=$$COMMIT \
		--build-arg __EX_METADATA_BUILD_TIME=$$BUILD_TIME \
		exodus && \
	docker compose up -d --no-deps --force-recreate exodus && \
	echo "$(GREEN)Exodus container deployed successfully!$(RESET)"

logs:
	@docker compose logs -f exodus

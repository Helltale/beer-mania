
# Variables
BACKEND_DIR = backend
TELEGRAM_DIR = telegram

lint-backend: ## Run golangci-lint on backend module
	@echo "Running golangci-lint on backend..."
	@cd $(BACKEND_DIR) && go mod download && go mod tidy && PATH=$$PATH:$$(go env GOPATH)/bin golangci-lint run --timeout=5m

lint-telegram: ## Run golangci-lint on telegram module
	@echo "Running golangci-lint on telegram..."
	@cd $(TELEGRAM_DIR) && go mod download && go mod tidy && PATH=$$PATH:$$(go env GOPATH)/bin golangci-lint run --timeout=5m

lint: lint-backend lint-telegram ## Run golangci-lint on all modules

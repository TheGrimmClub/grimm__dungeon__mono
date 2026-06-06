# TheGrimmClub — grimm monorepo
# Common developer tasks. Run `make help` for the list.

MODULE   := github.com/TheGrimmClub/grimm__dungeon__mono
BIN      := bin
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 windows/amd64

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build grimm and alchemist into ./bin
	@mkdir -p $(BIN)
	go build -o $(BIN)/grimm ./cmd/grimm
	go build -o $(BIN)/alchemist ./cmd/alchemist

.PHONY: run
run: ## Build and run grimm
	go run ./cmd/grimm

.PHONY: test
test: ## Run the test suite
	go test ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || \
		echo "golangci-lint not installed — skipping (CI runs it)"

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: check
check: vet test ## Run vet + tests (the quick gate)

.PHONY: cross
cross: ## Cross-compile grimm for all target platforms
	@mkdir -p $(BIN)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=""; \
		[ "$$os" = "windows" ] && ext=".exe"; \
		echo "building $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch go build -o $(BIN)/grimm-$$os-$$arch$$ext ./cmd/grimm; \
	done

.PHONY: clean
clean: ## Remove build output
	rm -rf $(BIN)

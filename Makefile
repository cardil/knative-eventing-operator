PROJECT_DIR = "./"
CLI_DIR     = "$(PROJECT_DIR)/cmd/manager"
BIN         = "$(PROJECT_DIR)/build/_output/bin/knative-eventing-operator"

GO           ?= go
GOLINT       ?= revive
RICHGO       ?= richgo

RESET          = \033[0m
make_std_color = \033[3$1m      # defined for 1 through 7
make_color     = \033[38;5;$1m  # defined for 1 through 255
BLUE = $(strip $(call make_color,44))
PINK = $(strip $(call make_color,210))
RED = $(strip $(call make_color,206))
GREEN = $(strip $(call make_color,120))
DGREEN = $(strip $(call make_color,106))
GRAY = $(strip $(call make_color,224))

.PHONY: default
default: binary

.PHONY: builddeps
builddeps:
	@$(GO) get github.com/kyoh86/richgo
	@$(GO) get github.com/mgechev/revive

.PHONY: clean
clean: builddeps
	@echo " $(GRAY)🛁 Cleaning$(RESET)"
	@rm -fv $(BIN)

.PHONY: check
check: builddeps
	@echo " $(PINK)🛂 Checking$(RESET)"
	$(GOLINT) -config revive.toml -formatter stylish -exclude ./vendor/... ./...

.PHONY: test
test: check
	@echo " $(GREEN)✔️ Testing$(RESET)"
	$(RICHGO) test -cover ./...

.PHONY: e2e
e2e: test
	@echo " $(DGREEN)🛫 E2E Testing$(RESET)"
	$(GO) test -v -count=1 -timeout=30m -tags e2e ./test/e2e

.PHONY: binary
binary: test
	@echo " $(BLUE)🔨 Building$(RESET)"
	$(RICHGO) build -o $(BIN) $(CLI_DIR)

.PHONY: run
run: binary
	@echo " $(RED)🏃 Running$(RESET)"
	$(BIN) $(args)

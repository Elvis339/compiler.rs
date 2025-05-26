REQUIRED_DEPS := go gotip

GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

BIN_DIR := bin
TRACE_DIR := traces

.PHONY: check-deps btree btreex all clean graph graphx codesign run

all: check-deps btree btreex memaccess memaccessx graph graphx codesign

check-deps:
	@for dep in $(REQUIRED_DEPS); do \
		if ! command -v $$dep > /dev/null 2>&1; then \
			echo "$(RED)$$dep is not installed or not in PATH$(NC)"; \
			exit 1; \
		else \
			echo "$(GREEN)✓ $$dep is installed$(NC)"; \
		fi; \
	done

btree:
	@echo "$(GREEN)Building btree without experimental flag...$(NC)"
	unset GOEXPERIMENT && go build -o $(BIN_DIR)/btree cmd/binarytrees/main.go

btreex:
	@echo "$(GREEN)Building btreex with experimental flag...$(NC)"
	GOEXPERIMENT=greenteagc gotip build -o $(BIN_DIR)/btreex cmd/binarytrees/main.go

memaccess:
	@echo "$(GREEN)Building memaccess without experimental flag...$(NC)"
	unset GOEXPERIMENT && go build -o $(BIN_DIR)/memaccess cmd/memaccess/memaccess.go

memaccessx:
	@echo "$(GREEN)Building memaccess with experimental flag...$(NC)"
	GOEXPERIMENT=greenteagc gotip build -o $(BIN_DIR)/memaccessx cmd/memaccess/memaccess.go

graph:
	@echo "$(GREEN)Building graph without experimental flag...$(NC)"
	unset GOEXPERIMENT && go build -o $(BIN_DIR)/graph ./cmd/graph

graphx:
	@echo "$(GREEN)Building graphx with experimental flag...$(NC)"
	GOEXPERIMENT=greenteagc gotip build -o $(BIN_DIR)/graphx ./cmd/graph

# Codesign all files in bin directory except .gitkeep
codesign:
	@if [ "$(shell uname -s)" = "Darwin" ]; then \
		echo "$(GREEN)Codesigning executables in $(BIN_DIR)/ directory...$(NC)"; \
		find $(BIN_DIR)/ -type f -not -name ".gitkeep" -exec codesign -s - -v -f --entitlements debug.plist {} \; ; \
	else \
		echo "$(GREEN)Skipping codesign (not on macOS)$(NC)"; \
	fi

# Run target with GODEBUG tracing enabled
# Usage: make run EXEC=<executable_name> [ARGS="<arguments>"]
run:
	@if [ -z "$(EXEC)" ]; then \
		echo "$(RED)Error: EXEC parameter is required. Usage: make run EXEC=<executable_name> [ARGS=\"<arguments>\"]$(NC)"; \
		exit 1; \
	fi
	@mkdir -p $(TRACE_DIR)
	@echo "$(GREEN)Running $(EXEC) with GC tracing enabled...$(NC)"
	@if [ -z "$(ARGS)" ]; then \
		GODEBUG=gctrace=1 $(BIN_DIR)/$(EXEC) 2>$(TRACE_DIR)/$(EXEC).gctrace; \
	else \
		echo "$(GREEN)Arguments: $(ARGS)$(NC)"; \
		GODEBUG=gctrace=1 $(BIN_DIR)/$(EXEC) $(ARGS) 2>$(TRACE_DIR)/$(EXEC).gctrace; \
	fi
	@echo "$(GREEN)GC trace saved to $(TRACE_DIR)/$(EXEC).gctrace$(NC)"

# Clean all files in bin/ except bin dir and .gitkeep
clean:
	@echo "$(GREEN)Cleaning executables in $(BIN_DIR)/ directory...$(NC)"
	@find $(BIN_DIR)/ -type f -not -name ".gitkeep" -delete
	@echo "$(GREEN)Cleaning trace files in $(TRACE_DIR)/ directory...$(NC)"
	@rm -rf $(TRACE_DIR)
	@echo "$(GREEN)Cleaning pprof's...$(NC)"
	@rm -rf *.pprof
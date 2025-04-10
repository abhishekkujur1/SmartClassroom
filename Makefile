PYTHON=python3
PIP=$(VENV_DIR)/bin/pip
VENV_DIR=mlapi/venv
GO_EXECUTABLE=$(shell command -v go 2> /dev/null)

# Default Target
all: setup_python setup_go

## Python Environment Setup
setup_python:
	@if [ ! -d "$(VENV_DIR)" ]; then \
		echo "Creating Python virtual environment..."; \
		$(PYTHON) -m venv $(VENV_DIR); \
		fi
	@echo "Installing Python dependencies..."
	$(PIP) install --upgrade pip
	$(PIP) install -r mlapi/requirements.txt

## Golang Environment Setup
setup_go:
	ifndef GO_EXECUTABLE
	@echo "Go not found! Please install Golang manually for your platform."
	@exit 1
else
	@echo "Go found at $(GO_EXECUTABLE). Setting up Go modules..."
	cd . && go mod tidy
	cd . && go build ./...
endif

## Run Python API
run_mlapi:
	@echo "Starting Python Face Detection API..."
	source $(VENV_DIR)/bin/activate && uvicorn mlapi.main:app --host 0.0.0.0 --port 8000

## Run Go App (Example)
run_go_server:
	cd cmd/server && go run main.go

## Clean All
clean:
	rm -rf $(VENV_DIR)
	go clean

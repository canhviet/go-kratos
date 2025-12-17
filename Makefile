# Project variables
APP_NAME := myapp

# Proto directories
API_DIR := api
THIRD_PARTY := third_party

# Default target
.PHONY: all
all: proto conf build

# ----------------------------------
# 🔥 Generate Protobuf (HTTP + GRPC)
# ----------------------------------
.PHONY: proto
proto:
	@protoc \
	  --proto_path=. \
	  --proto_path=$(API_DIR) \
	  --proto_path=$(THIRD_PARTY) \
	  --go_out=paths=source_relative:. \
	  --go-grpc_out=paths=source_relative:. \
	  --go-http_out=paths=source_relative:. \
	  $(shell find $(API_DIR) -name "*.proto")

# -----------------------
# 🔥 Generate conf.pb.go
# -----------------------
.PHONY: conf
conf:
	@protoc \
	  --proto_path=. \
	  --proto_path=$(THIRD_PARTY) \
	  --go_out=paths=source_relative:. \
	  internal/conf/conf.proto

# -----------------------
# 🔥 Build Application
# -----------------------
.PHONY: build
build:
	@go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# -----------------------
# 🔥 Run Application
# -----------------------
.PHONY: run
run:
	@go run ./cmd/$(APP_NAME)

# -----------------------
# 🧹 Clean generated files
# -----------------------
.PHONY: clean
clean:
	@rm -rf bin
	@find . -name "*.pb.go" -type f -delete

# -----------------------
# 🐳 Docker build
# -----------------------
.PHONY: docker
docker:
	@docker build -t $(APP_NAME):latest .

GREET_DIR = greet/proto
BLOG_DIR = blog/proto
KAFKA_DIR = kafka/proto
KAFKA_V2_DIR = kafka.schema.v2/proto

ifeq ($(OS), Windows_NT)
    OS = windows
    SHELL := powershell.exe
    .SHELLFLAGS := -NoProfile -Command
    # Extract package name from go.mod safely
    PACKAGE = $(shell Get-Content go.mod -head 1 | Foreach-Object { $$data = $$_ -split " "; "{0}" -f $$data[1]})
    BIN = server.exe
else
    UNAME := $(shell uname -s)
    ifeq ($(UNAME),Darwin)
       OS = macos
    else ifeq ($(UNAME),Linux)
       OS = linux
    endif
    PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
    BIN = server
endif

.PHONY: all generate clean build

all: generate build

build: generate
	go build -o bin/greet/server ./greet/server/
	go build -o bin/greet/client ./greet/client/
	go build -o bin/blog/server ./blog/server/
	go build -o bin/blog/client ./blog/client/
	go build -o bin/kafka/producer ./kafka/producer/
	go build -o bin/kafka/consumer ./kafka/consumer/
	go build -o bin/kafka.schema.v2/producer ./kafka.schema.v2/producer/
	go build -o bin/kafka.schema.v2/consumer ./kafka.schema.v2/consumer/
generate:
	protoc --proto_path=${GREET_DIR} --go_out=${GREET_DIR} --go_opt=paths=source_relative --go-grpc_out=${GREET_DIR} --go-grpc_opt=paths=source_relative ${GREET_DIR}/*.proto
	protoc --proto_path=${BLOG_DIR} --go_out=${BLOG_DIR} --go_opt=paths=source_relative --go-grpc_out=${BLOG_DIR} --go-grpc_opt=paths=source_relative ${BLOG_DIR}/*.proto
	protoc --proto_path=${KAFKA_DIR} --go_out=${KAFKA_DIR} --go_opt=paths=source_relative --go-grpc_out=${KAFKA_DIR} --go-grpc_opt=paths=source_relative ${KAFKA_DIR}/*.proto
	protoc --proto_path=${KAFKA_V2_DIR} --go_out=${KAFKA_V2_DIR} --go_opt=paths=source_relative --go-grpc_out=${KAFKA_V2_DIR} --go-grpc_opt=paths=source_relative ${KAFKA_V2_DIR}/*.proto

clean:
	@if [ "$(OS)" = "windows" ]; then \
		powershell -Command "Remove-Item ${GREET_DIR}/*.pb.go -ErrorAction SilentlyContinue"; \
		powershell -Command "Remove-Item ${BIN} -ErrorAction SilentlyContinue"; \
	else \
		rm -f ${GREET_DIR}/*.pb.go ${BIN}; \
	fi
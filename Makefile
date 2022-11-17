# Install:
# go install github.com/pressly/goose/v3/cmd/goose@latest
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
# go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
# go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

DB ?= $(shell bash -c 'read -p "DB Address [127.0.0.1:3306]: " db; echo $${db:-127.0.0.1:3306}')
DB_USERNAME ?= $(shell bash -c 'read -p "DB Username: " username; echo $$username')
DB_PASSWORD ?= $(shell bash -c 'read -p "DB Password: " pwd; echo $$pwd')

protoc:
	protoc -I ./pkg/grpc/schema/ \
    		--go_out=./pkg/grpc/api/proto/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./pkg/grpc/api/proto/ \
    		--go-grpc_opt=paths=source_relative \
    		./pkg/grpc/schema/*.proto

migrate:
	@cd ./migrations && goose mysql '${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB})/signaller' up

build:
	env GOOS=linux GOARCH=amd64	go build -o ./build/signaller ./cmd/signaller

deploy: build
	frsync ./build/signaller dev1:/home/ubuntu/signaller/app/app
	ssh ubuntu@dev1 "cd ~/signaller/ && docker compose build && docker compose down && docker compose up -d"

default: protoc

.PHONY: *
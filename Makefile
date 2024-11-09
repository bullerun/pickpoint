OUT_PATH:=$(CURDIR)/pkg
LOCAL_BIN:=$(CURDIR)/build/bin
.PHONY: test
createMocks:
	minimock -i github.com/IBM/sarama.SyncProducer -o .\client\internal\commands\mock\sync_producer_mock.go -n SyncProducerMock
	minimock -i .\client\internal\commands.AcceptOrderService -o .\client\internal\commands\mock\accept_order_service_mock.go -n AcceptOrderServiceMock
	minimock -i .\client\internal\commands.AcceptReturnService -o .\client\internal\commands\mock\accept_return_service_mock.go -n AcceptReturnServiceMock
	minimock -i .\client\internal\commands.IssueOrderService -o .\client\internal\commands\mock\issue_order_service_mock.go -n IssueOrderServiceMock
	minimock -i .\client\internal\commands.ListOrdersService -o .\client\internal\commands\mock\list_order_service_mock.go -n ListOrderServiceMock
	minimock -i .\client\internal\commands.ReturnOrderService -o .\client\internal\commands\mock\return_order_service_mock.go -n ReturnOrderServiceMock
test: createMocks
	go test ./internal/commands -cover

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations/prod postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" create rename_me sql

goose-up:
	goose -dir ./migrations/prod postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

goose-status:
	goose -dir ./migrations/prod postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" status

goose-down:
	goose -dir ./migrations/prod postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down-to $(v)

goose-up-test:
	goose -dir ./migrations/test postgres "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable" up

goose-down-test:
	goose -dir ./migrations/test  postgres "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable" down-to 0

# на wsl делал)
squawk-install:
	npm install -g squawk-cli

squawk:
	squawk ./migrations/*/*

gofakeit-install:
	go get github.com/brianvoe/gofakeit/v7

startTestAll: goose-down-test goose-up-test
	go test ./... -cover
testAll: startTestAll

create-fake-data:
	go run server/cmd/gofackit/gofackit.go

# ---------------------------------
# Запуск кодогенерации через protoc
# ---------------------------------

deps: .vendor-proto
	go get -u github.com/go-chi/chi/v5
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${OUT_PATH}
	protoc --proto_path api --proto_path build/vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		./api/order-service/v1/order_service.proto

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway build/vendor.protogen/grpc-ecosystem && \
 		cd build/vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p build/vendor.protogen/protoc-gen-openapiv2
		mv build/vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options build/vendor.protogen/protoc-gen-openapiv2
		rm -rf build/vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf build/vendor.protogen/protobuf &&\
		cd build/vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p build/vendor.protogen/google
		mv build/vendor.protogen/protobuf/src/google/protobuf build/vendor.protogen/google
		rm -rf build/vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis build/vendor.protogen/googleapis && \
 		cd build/vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  build/vendor.protogen/google
		mv build/vendor.protogen/googleapis/google/api build/vendor.protogen/google
		rm -rf build/vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate build/vendor.protogen/tmp && \
		cd build/vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p build/vendor.protogen/validate
		mv build/vendor.protogen/tmp/validate build/vendor.protogen/
		rm -rf build/vendor.protogen/tmp

build:
	go build -o $(CURDIR)/build/tmp/client.exe ./client/cmd/order_client
	go build -o $(CURDIR)/build/tmp/server.exe ./server/cmd/order_server
run:
	$(CURDIR)/build/tmp/server.exe & $(CURDIR)/build/tmp/client.exe

all: deps generate build run
.PHONY: build

.PHONY: run-prometheus
run-prometheus:
	prometheus --config.file ./config/prometheus.yaml
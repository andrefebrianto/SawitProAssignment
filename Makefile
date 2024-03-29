

.PHONY: clean all init generate generate_mocks generate_key

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated
	rm -rf cert

init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.tmp.out -v ./...
	cat coverage.tmp.out | grep -v ".mock.gen.go" > coverage.out
	go tool cover -html=coverage.out

generate: generated generate_mocks generate_key

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find . -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

validate_api_spec:
	openapi-spec-validator api.yml

generate_key:
	@echo "Gemerating RSA certificate..."
	mkdir cert || true
	openssl genrsa -out cert/id_rsa 4096
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub

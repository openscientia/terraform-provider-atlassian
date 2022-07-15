TEST?=$$(go list ./... | grep -v 'vendor')

default: install

build:
	go build -v ./...

install: build
	go install -v ./...

# See https://golangci-lint.run/
lint: 
	golangci-lint run

# See https://github.com/hashicorp/terraform-plugin-docs#usage
generate:
	tfplugindocs generate

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

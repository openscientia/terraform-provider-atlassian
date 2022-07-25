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
docs:
	tfplugindocs generate

gen:
	rm -f .github/labeler-issue-labels.yml
	rm -f .github/labeler-pr-labels.yml
	rm -f infrastructure/repository/labels-resource.tf
	go generate ./...

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

PKG 				 = internal/provider/...
ACCTEST_COUNT		?= 1
ACCTEST_TIMEOUT     ?= 180m
ACCTEST_PARALLELISM ?= 20

ifneq ($(origin TESTS), undefined)
	RUNARGS = -run='$(TESTS)'
endif

default: build

# See https://go.dev/ref/mod#go-install
build:
	go install -v

# See https://golangci-lint.run/
lint: 
	golangci-lint run

# See https://github.com/hashicorp/terraform-plugin-docs#usage
tfdocs:
	tfplugindocs generate

# See https://go.dev/blog/generate
gen:
	rm -f .github/labeler-issue-labels.yml
	rm -f .github/labeler-pr-labels.yml
	rm -f infrastructure/repository/labels-resource.tf
	go generate ./...

# See https://pkg.go.dev/cmd/go/internal/test
testacc:
	@if [ "$(TESTARGS)" = "-run=TestAccXXX" ]; then \
		echo ""; \
		echo "Error: Skipping example acceptance testing pattern. Update TESTS for the relevant *_test.go file."; \
		echo ""; \
		echo "For example if testing internal/provider/resource_jira_issue_type.go, use the test names in internal/provider/resource_jira_issue_type_test.go starting with TestAcc and up to the underscore:"; \
		echo "make testacc TESTS=TestAccJiraIssueType_"; \
		echo ""; \
		echo "See the contributing guide for more information: https://github.com/openscientia/terraform-provider-atlassian/blob/main/.github/contributing/acceptance-tests.md"; \
		exit 1; \
	fi
	TF_ACC=1 go test ./$(PKG) -v -count $(ACCTEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(RUNARGS) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

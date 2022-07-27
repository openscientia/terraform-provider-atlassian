package atlassian

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"atlassian": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ATLASSIAN_USERNAME"); v == "" {
		t.Fatal("ATLASSIAN_USERNAME must be set to run acceptance tests.")
	}

	if v := os.Getenv("ATLASSIAN_TOKEN"); v == "" {
		t.Fatal("ATLASSIAN_TOKEN must be set to run acceptance tests.")
	}

	if v := os.Getenv("ATLASSIAN_URL"); v == "" {
		t.Error("ATLASSIAN_URL must be set to run acceptance tests.")
	}
}

func TestProvider_InvalidUrlAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: `
					provider "atlassian" {
						url = " https://test.atlassian.net"
					}

					resource "atlassian_jira_issue_type" "test" {
						name = "test"
					}
				`,
				ExpectError: regexp.MustCompile(`Parsing URL\s.*\sfailed`),
			},
			{
				Config: `
					provider "atlassian" {
						url = "test.atlassian.net"
					}

					resource "atlassian_jira_issue_type" "test" {
						name = "test"
					}
				`,
				ExpectError: regexp.MustCompile(`contains no host`),
			},
			{
				Config: `
					provider "atlassian" {
						url = "http://test.atlassian.net"
					}

					resource "atlassian_jira_issue_type" "test" {
						name = "test"
					}
				`,
				ExpectError: regexp.MustCompile(`expected to use scheme from`),
			},
		},
	})
}

package atlassian

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"atlassian": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

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

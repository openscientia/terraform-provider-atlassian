package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	atlassian "github.com/openscientia/terraform-provider-atlassian/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/openscientia/atlassian",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), atlassian.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

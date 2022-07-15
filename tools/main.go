//go:build toolss

package main

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/hashicorp/go-changelog/cmd/changelog-build"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/rhysd/actionlint/cmd/actionlint"
)

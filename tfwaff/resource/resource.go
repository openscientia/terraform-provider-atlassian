package resource

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"

	"github.com/openscientia/terraform-provider-atlassian/tfwaff/utils"
)

//go:embed resource.tmpl
var resourceTmpl string

//go:embed resource_test.tmpl
var resourceTestTmpl string

type resourceTemplateData struct {
	ProviderLower          string
	ProviderSuffix         string
	ResourceCamel          string
	ResourceKebab          string
	ResourcePascal         string
	ResourceProse          string
	ResourceSnake          string
	ResourceSnakeFull      string
	ResourceTitle          string
	ResourceSuffix         string
	ResourceModelSuffix    string
	ResourceFilenamePrefix string
	ServiceLower           string
	ServiceTitle           string
}

func Create(provider, name string, force, dry_run bool) error {
	if !utils.IsPascalCase(name) {
		return fmt.Errorf("'name' must be in pascal case, e.g., FooBarBaz")
	}

	r := strings.SplitN(utils.GetSnakeCase(name), "_", 2)
	service := r[0]
	serviceTitle := utils.GetTitleCase(service)
	rSnake := r[1]
	rKebab := utils.GetKebabCase(name)
	rTitle := utils.GetTitleCase(strings.ReplaceAll(rSnake, "_", " "))
	rCamel := strings.ReplaceAll(strings.ToLower(rTitle[:1])+rTitle[1:], " ", "")
	rPascal := strings.ReplaceAll(rTitle, " ", "")

	rtd := resourceTemplateData{
		ProviderLower:          provider,
		ProviderSuffix:         "Provider",
		ResourceCamel:          rCamel,
		ResourceKebab:          rKebab,
		ResourcePascal:         rPascal,
		ResourceProse:          strings.ToLower(rTitle),
		ResourceSnake:          rSnake,
		ResourceTitle:          rTitle,
		ResourceSnakeFull:      strings.Join([]string{provider, service, rSnake}, "_"),
		ResourceSuffix:         "Resource",
		ResourceModelSuffix:    "ResourceModel",
		ResourceFilenamePrefix: "resource",
		ServiceLower:           service,
		ServiceTitle:           serviceTitle,
	}

	if !dry_run {
		filename := fmt.Sprintf("%s_%s_%s.go", rtd.ResourceFilenamePrefix, rtd.ServiceLower, rtd.ResourceSnake)
		if err := writeTemplate("new-resource-file", filename, resourceTmpl, force, rtd); err != nil {
			return fmt.Errorf("writing resource template: %w", err)
		}

		testFilename := fmt.Sprintf("%s_%s_%s_test.go", rtd.ResourceFilenamePrefix, rtd.ServiceLower, rtd.ResourceSnake)
		if err := writeTemplate("new-resource-test-file", testFilename, resourceTestTmpl, force, rtd); err != nil {
			return fmt.Errorf("writing resource test template: %w", err)
		}
	}

	fmt.Println("created new resource:", name)

	return nil
}

func writeTemplate(templateName, filename, tmpl string, force bool, td resourceTemplateData) error {
	if _, err := os.Stat(filename); !errors.Is(err, fs.ErrNotExist) && !force {
		return fmt.Errorf("file (%s) already exists and force is not set", filename)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file (%s): %s", filename, err)
	}

	tplate, err := template.New(templateName).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tplate.Execute(&buffer, td)
	if err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return fmt.Errorf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing file (%s): %s", filename, err)
	}

	return nil
}

package datasource

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/openscientia/terraform-provider-atlassian/tfwaff/utils"
	"github.com/spf13/cobra"
)

//go:embed datasource.tmpl
var datasourceTmpl string

//go:embed datasource_test.tmpl
var datasourceTestTmpl string

type dataSourceTemplateData struct {
	ProviderLower            string
	ProviderSuffix           string
	DataSourceCamel          string
	DataSourceKebab          string
	DataSourcePascal         string
	DataSourceProse          string
	DataSourceSnake          string
	DataSourceSnakeFull      string
	DataSourceTitle          string
	DataSourceSuffix         string
	DataSourceModelSuffix    string
	DataSourceFilenamePrefix string
	ServiceLower             string
	ServiceTitle             string
}

func Create(provider, name string, force, dry_run bool) error {
	if !utils.IsPascalCase(name) {
		return fmt.Errorf("'name' must be in pascal case, e.g., FooBarBaz")
	}

	d := strings.SplitN(utils.GetSnakeCase(name), "_", 2)
	service := d[0]
	serviceTitle := utils.GetTitleCase(service)
	dSnake := d[1]
	dKebab := utils.GetKebabCase(name)
	dTitle := utils.GetTitleCase(strings.ReplaceAll(dSnake, "_", " "))
	dCamel := strings.ReplaceAll(strings.ToLower(dTitle[:1])+dTitle[1:], " ", "")
	dPascal := strings.ReplaceAll(dTitle, " ", "")

	dstd := dataSourceTemplateData{
		ProviderLower:            provider,
		ProviderSuffix:           "Provider",
		DataSourceCamel:          dCamel,
		DataSourceKebab:          dKebab,
		DataSourcePascal:         dPascal,
		DataSourceProse:          strings.ToLower(dTitle),
		DataSourceSnake:          dSnake,
		DataSourceSnakeFull:      strings.Join([]string{provider, service, dSnake}, "_"),
		DataSourceTitle:          dTitle,
		DataSourceSuffix:         "DataSource",
		DataSourceModelSuffix:    "DataSourceModel",
		DataSourceFilenamePrefix: "data_source",
		ServiceLower:             service,
		ServiceTitle:             serviceTitle,
	}

	if !dry_run {
		filename := fmt.Sprintf("%s_%s_%s.go", dstd.DataSourceFilenamePrefix, dstd.ServiceLower, dstd.DataSourceSnake)
		if err := writeTemplate("new-datasource", filename, datasourceTmpl, force, dstd); err != nil {
			cobra.CheckErr(err)
			return fmt.Errorf("writing data source template: %w", err)
		}

		testFilename := fmt.Sprintf("%s_%s_%s_test.go", dstd.DataSourceFilenamePrefix, dstd.ServiceLower, dstd.DataSourceSnake)
		if err := writeTemplate("new-resource-test-file", testFilename, datasourceTestTmpl, force, dstd); err != nil {
			return fmt.Errorf("writing resource test template: %w", err)
		}
	}

	fmt.Println("created new data source:", name)

	return nil
}

func writeTemplate(templateName, filename, tmpl string, force bool, td dataSourceTemplateData) error {
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

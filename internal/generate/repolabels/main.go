//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	filename = `../../../infrastructure/repository/labels-resource.tf`
)

type Resource struct {
	Name string
}

type templateData struct {
	Resources []Resource
}

func main() {
	fmt.Printf("Generating %s\n", strings.TrimPrefix(filename, "../../../"))

	productUrls := []string{
		"https://developer.atlassian.com/cloud/jira/platform/swagger-v3.v3.json",
		"https://developer.atlassian.com/cloud/confluence/swagger.v3.json",
	}

	lbs := []Resource{}
	for _, v := range productUrls {
		lbs = append(lbs, getResources(v)...)
	}

	// Add custom labels
	lbs = append(lbs, getJiraCustomResources()...)

	td := templateData{}
	td.Resources = append(td.Resources, lbs...)

	sort.SliceStable(td.Resources, func(i, j int) bool {
		return td.Resources[i].Name < td.Resources[j].Name
	})

	writeTemplate(tmpl, "repolabeler", td)
}

func getResources(url string) []Resource {
	c := http.Client{Timeout: time.Duration(2) * time.Second}
	resp, getErr := c.Get(url)
	if getErr != nil {
		log.Fatalf("error calling url (%s): %s", url, getErr)
		return nil
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalf("error reading response body %s", readErr)
	}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		log.Fatalf("error unmarshalling JSON response %s", parseErr)
		return nil
	}
	resp.Body.Close()

	r := regexp.MustCompile(`\(apps\)`)
	r2 := regexp.MustCompile(`jira|confluence`)
	productName := r2.FindString(url)
	var tags = result["tags"].([]interface{})
	var resources []Resource
	replacer := strings.NewReplacer(" ", "", "-", "")
	for _, value := range tags {
		// Each value is an interface{} type, that is type asserted as map[string]interface{}
		// due to nested objects in the original JSON response
		m := value.(map[string]interface{})
		rawLabel := m["name"]
		ok := r.MatchString(rawLabel.(string))
		if ok {
			continue
		}
		rsc := Resource{}
		rsc.Name = productName + "/" + strings.ToLower(replacer.Replace(rawLabel.(string)))
		resources = append(resources, rsc)
	}

	return resources
}

func getJiraCustomResources() []Resource {
	var customResources []Resource
	resources := []string{
		"issuecustomfields",
		"issuefieldconfigurationitems",
		"issuefieldconfigurationschemes",
		"issuefieldconfigurationschememappings",
		"permissiongrants",
	}
	for _, r := range resources {
		customResources = append(customResources, Resource{
			Name: fmt.Sprintf("jira/%s", r),
		})
	}

	return customResources
}

func writeTemplate(body string, templateName string, td templateData) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", filename, err)
	}

	tp, err := template.New(templateName).Parse(body)
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tp.Execute(&buffer, td)
	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close()
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("error closing file (%s): %s", filename, err)
	}
}

var tmpl = `# Generated by internal/generate/repolabels/main.go; DO NOT EDIT.
variable "resource_labels" {
  default = [
    {{- range .Resources }}
    "{{ .Name }}",
    {{- end }}
  ]
  description = "Set of ATLASSIAN products' resource-specific labels."
  type        = set(string)
}

resource "github_issue_label" "resource" {
  for_each = var.resource_labels

  repository  = "terraform-provider-atlassian"
  name        = each.value
  color       = "5a4edd" # color: https://registry.terraform.io/
  description = "Issues and PRs that pertain to ${each.value} resources."
}
`

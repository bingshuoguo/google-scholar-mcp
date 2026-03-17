package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var requiredTools = []string{
	"search_google_scholar_key_words",
	"search_google_scholar_advanced",
	"get_author_info",
}

var requiredResources = []string{
	"scholar://server/overview",
	"scholar://server/tools",
	"scholar://server/config",
	"scholar://server/limitations",
}

var requiredResourceTemplates = []string{
	"scholar://search-guide/{topic}",
}

var requiredPrompts = []string{
	"scholar_literature_scan",
	"scholar_author_brief",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <server-binary>\n", os.Args[0])
		os.Exit(2)
	}

	serverBinary := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.Command(serverBinary, "stdio")
	client := mcp.NewClient(&mcp.Implementation{Name: "google-scholar-smoke-client", Version: "dev"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		fatalf("connect to stdio server: %v", err)
	}
	defer session.Close()

	res, err := session.ListTools(ctx, nil)
	if err != nil {
		fatalf("list tools: %v", err)
	}

	if len(res.Tools) == 0 {
		fatalf("list tools returned zero tools")
	}

	names := make([]string, 0, len(res.Tools))
	for _, tool := range res.Tools {
		if tool != nil {
			names = append(names, tool.Name)
		}
	}
	slices.Sort(names)

	missing := make([]string, 0, len(requiredTools))
	for _, required := range requiredTools {
		if !slices.Contains(names, required) {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		fatalf("missing required tools: %v", missing)
	}

	resources, err := session.ListResources(ctx, nil)
	if err != nil {
		fatalf("list resources: %v", err)
	}
	resourceNames := make([]string, 0, len(resources.Resources))
	for _, resource := range resources.Resources {
		if resource != nil {
			resourceNames = append(resourceNames, resource.URI)
		}
	}
	assertContains("resources", resourceNames, requiredResources)

	templates, err := session.ListResourceTemplates(ctx, nil)
	if err != nil {
		fatalf("list resource templates: %v", err)
	}
	templateNames := make([]string, 0, len(templates.ResourceTemplates))
	for _, template := range templates.ResourceTemplates {
		if template != nil {
			templateNames = append(templateNames, template.URITemplate)
		}
	}
	assertContains("resource templates", templateNames, requiredResourceTemplates)

	prompts, err := session.ListPrompts(ctx, nil)
	if err != nil {
		fatalf("list prompts: %v", err)
	}
	promptNames := make([]string, 0, len(prompts.Prompts))
	for _, prompt := range prompts.Prompts {
		if prompt != nil {
			promptNames = append(promptNames, prompt.Name)
		}
	}
	assertContains("prompts", promptNames, requiredPrompts)

	fmt.Println("Smoke test passed")
	fmt.Println("Tools:")
	for _, name := range names {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println("Resources:")
	for _, name := range resourceNames {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println("Resource templates:")
	for _, name := range templateNames {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println("Prompts:")
	for _, name := range promptNames {
		fmt.Printf("- %s\n", name)
	}
}

func fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func assertContains(kind string, actual []string, required []string) {
	slices.Sort(actual)
	missing := make([]string, 0, len(required))
	for _, item := range required {
		if !slices.Contains(actual, item) {
			missing = append(missing, item)
		}
	}
	if len(missing) > 0 {
		fatalf("missing required %s: %v", kind, missing)
	}
}

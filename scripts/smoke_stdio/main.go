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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <server-binary>\n", os.Args[0])
		os.Exit(2)
	}

	serverBinary := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.Command(serverBinary)
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

	fmt.Println("Smoke test passed")
	for _, name := range names {
		fmt.Printf("- %s\n", name)
	}
}

func fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

package main

import (
	"bytes"
	"runtime/debug"
	"testing"
)

func TestParseCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
		err  string
	}{
		{name: "default stdio", args: nil, want: commandStdio},
		{name: "stdio subcommand", args: []string{"stdio"}, want: commandStdio},
		{name: "long version flag", args: []string{"--version"}, want: commandVersion},
		{name: "single dash version", args: []string{"-version"}, want: commandVersion},
		{name: "version subcommand", args: []string{"version"}, want: commandVersion},
		{name: "help flag", args: []string{"--help"}, want: commandHelp},
		{name: "help subcommand", args: []string{"help"}, want: commandHelp},
		{name: "unknown command", args: []string{"serve"}, err: `unknown command "serve"`},
		{name: "extra version args", args: []string{"version", "extra"}, err: "version does not accept additional arguments"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseCommand(tc.args)
			if tc.err != "" {
				if err == nil || err.Error() != tc.err {
					t.Fatalf("parseCommand(%v) error = %v, want %q", tc.args, err, tc.err)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseCommand(%v) returned error: %v", tc.args, err)
			}
			if got != tc.want {
				t.Fatalf("parseCommand(%v) = %q, want %q", tc.args, got, tc.want)
			}
		})
	}
}

func TestRunWritesVersionToStdout(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--version"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run() exit code = %d, want 0", exitCode)
	}
	if stdout.Len() == 0 {
		t.Fatal("run() did not write version output")
	}
	if stderr.Len() != 0 {
		t.Fatalf("run() wrote unexpected stderr: %q", stderr.String())
	}
}

func TestResolveVersionPrefersInjectedVersion(t *testing.T) {
	t.Parallel()

	got := resolveVersion("v1.2.3")
	if got != "v1.2.3" {
		t.Fatalf("resolveVersion() = %q, want %q", got, "v1.2.3")
	}
}

func TestResolveVersionFallsBackToBuildInfo(t *testing.T) {
	original := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: "v0.0.0-test"}}, true
	}
	t.Cleanup(func() {
		readBuildInfo = original
	})

	got := resolveVersion("dev")
	if got != "v0.0.0-test" {
		t.Fatalf("resolveVersion() = %q, want %q", got, "v0.0.0-test")
	}
}

func TestResolveVersionUsesDevWhenBuildInfoMissing(t *testing.T) {
	original := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return nil, false
	}
	t.Cleanup(func() {
		readBuildInfo = original
	})

	got := resolveVersion("dev")
	if got != "dev" {
		t.Fatalf("resolveVersion() = %q, want %q", got, "dev")
	}
}

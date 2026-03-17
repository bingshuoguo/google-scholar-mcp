package main

import (
	"runtime/debug"
	"testing"
)

func TestWantsVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "long flag", args: []string{"--version"}, want: true},
		{name: "single dash", args: []string{"-version"}, want: true},
		{name: "subcommand", args: []string{"version"}, want: true},
		{name: "empty", args: nil, want: false},
		{name: "unknown flag", args: []string{"--help"}, want: false},
		{name: "extra args", args: []string{"version", "extra"}, want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := wantsVersion(tc.args); got != tc.want {
				t.Fatalf("wantsVersion(%v) = %v, want %v", tc.args, got, tc.want)
			}
		})
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

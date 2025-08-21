package gomod_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func mustReadGoMod(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("go.mod")
	require.NoError(t, err, "go.mod should exist at repository root")
	require.NotEmpty(t, data, "go.mod must not be empty")
	return data
}

func mustParseGoMod(t *testing.T, data []byte) *modfile.File {
	t.Helper()
	f, err := modfile.Parse("go.mod", data, nil)
	require.NoError(t, err, "failed to parse go.mod")
	require.NotNil(t, f)
	return f
}

func isPseudoVersion(v string) bool {
	// Pseudo version pattern: vX.Y.Z-yyyymmddhhmmss-abcdefabcdef
	re := regexp.MustCompile(`^v\d+\.\d+\.\d+-\d{14}-[0-9a-f]{12}$`)
	return re.MatchString(v)
}

func TestGoMod_Basics(t *testing.T) {
	data := mustReadGoMod(t)
	text := string(data)
	// Assert the exact lines present (from PR diff)
	require.Contains(t, text, "module github.com/projectdiscovery/katana", "module path must match")
	require.Contains(t, text, "go 1.24.0", "go directive should be pinned to 1.24.0 per diff")
	require.Contains(t, text, "toolchain go1.24.2", "toolchain directive should be pinned to go1.24.2 per diff")

	f := mustParseGoMod(t, data)
	require.NotNil(t, f.Module)
	assert.Equal(t, "github.com/projectdiscovery/katana", f.Module.Mod.Path, "unexpected module path")

	// Some versions of x/mod may normalize go version; allow both 1.24 and 1.24.0
	if f.Go != nil {
		got := f.Go.Version
		assert.Truef(t, got == "1.24.0" || got == "1.24",
			"unexpected go version in parsed file: %s", got)
	} else {
		t.Fatalf("parsed go.mod did not include Go version")
	}
}

func TestGoMod_NoReplaceDirectives(t *testing.T) {
	f := mustParseGoMod(t, mustReadGoMod(t))
	assert.Len(t, f.Replace, 0, "replace directives should be absent per diff")
}

func TestGoMod_NoDuplicateModuleRequirements(t *testing.T) {
	f := mustParseGoMod(t, mustReadGoMod(t))
	seen := make(map[string]struct{})
	for _, req := range f.Require {
		path := req.Mod.Path
		if _, ok := seen[path]; ok {
			t.Fatalf("duplicate requirement detected for module: %s", path)
		}
		seen[path] = struct{}{}
	}
}

func TestGoMod_DirectDependenciesPinned(t *testing.T) {
	// Validate a representative set of direct dependencies and versions from the diff.
	expected := map[string]string{
		"github.com/BishopFox/jsluice":       "v0.0.0-20240110145140-0ddfab153e06",
		"github.com/PuerkitoBio/goquery":     "v1.8.1",
		"github.com/go-rod/rod":              "v0.114.1",
		"github.com/json-iterator/go":        "v1.1.12",
		"github.com/logrusorgru/aurora":      "v2.0.3+incompatible",
		"github.com/mitchellh/mapstructure":  "v1.5.0",
		"github.com/pkg/errors":              "v0.9.1",
		"github.com/projectdiscovery/dsl":    "v0.5.0",
		"github.com/projectdiscovery/fastdialer": "v0.4.1",
		"github.com/projectdiscovery/goflags":    "v0.1.74",
		"github.com/projectdiscovery/gologger":   "v1.1.54",
		"github.com/projectdiscovery/hmap":       "v0.0.91",
		"github.com/projectdiscovery/mapcidr":    "v1.1.34",
		"github.com/projectdiscovery/ratelimit":  "v0.0.79",
		"github.com/projectdiscovery/retryablehttp-go": "v1.0.118",
		"github.com/projectdiscovery/utils":      "v0.4.21",
		"github.com/projectdiscovery/wappalyzergo": "v0.2.38",
		"github.com/remeh/sizedwaitgroup":       "v1.0.0",
		"github.com/rs/xid":                     "v1.5.0",
		"github.com/stoewer/go-strcase":         "v1.3.0",
		"github.com/stretchr/testify":           "v1.10.0",
		"github.com/valyala/fasttemplate":       "v1.2.2",
		"go.uber.org/multierr":                  "v1.11.0",
		"golang.org/x/net":                      "v0.42.0",
		"gopkg.in/yaml.v3":                      "v3.0.1",
	}

	f := mustParseGoMod(t, mustReadGoMod(t))

	// Build a map of direct requirements
	got := make(map[string]struct {
		version  string
		indirect bool
	})
	for _, req := range f.Require {
		got[req.Mod.Path] = struct {
			version  string
			indirect bool
		}{version: req.Mod.Version, indirect: req.Indirect}
	}

	for path, expVer := range expected {
		r, ok := got[path]
		require.Truef(t, ok, "expected direct module not found: %s", path)
		assert.Falsef(t, r.indirect, "module %s should be a direct requirement (not indirect)", path)
		assert.Equalf(t, expVer, r.version, "module %s has unexpected version", path)
	}
}

func TestGoMod_IndirectExamplesAndYamlv2Direct(t *testing.T) {
	// Validate an indirect example and ensure yaml.v2 is direct per diff.
	f := mustParseGoMod(t, mustReadGoMod(t))

	type rec struct {
		version  string
		indirect bool
	}
	got := make(map[string]rec)
	for _, req := range f.Require {
		got[req.Mod.Path] = rec{version: req.Mod.Version, indirect: req.Indirect}
	}

	// Example indirect dependency from diff
	if r, ok := got["github.com/Knetic/govaluate"]; ok {
		assert.True(t, r.indirect, "github.com/Knetic/govaluate should be marked // indirect")
		assert.Equal(t, "v3.0.0+incompatible", r.version)
	} else {
		t.Fatalf("expected to find github.com/Knetic/govaluate in go.mod")
	}

	// yaml v2 should be present as a direct requirement at v2.4.0
	if r, ok := got["gopkg.in/yaml.v2"]; ok {
		assert.False(t, r.indirect, "gopkg.in/yaml.v2 should be a direct requirement")
		assert.Equal(t, "v2.4.0", r.version, "unexpected version for gopkg.in/yaml.v2")
	} else {
		t.Fatalf("expected to find gopkg.in/yaml.v2 in go.mod")
	}
}

func TestGoMod_AllVersionsAreSemverOrPseudo(t *testing.T) {
	f := mustParseGoMod(t, mustReadGoMod(t))
	for _, req := range f.Require {
		v := req.Mod.Version
		// valid if semver or pseudo
		validSemver := semver.IsValid(v)
		validPseudo := isPseudoVersion(v)
		assert.Truef(t, validSemver || validPseudo, "invalid version format for %s: %s", req.Mod.Path, v)
	}
}

func TestGoMod_ToolchainSyntaxPresent(t *testing.T) {
	data := mustReadGoMod(t)
	text := string(data)
	// The exact directive is validated in Basics; here we ensure syntax more generally.
	re := regexp.MustCompile(`(?m)^toolchain\s+go\d+\.\d+(\.\d+)?$`)
	assert.Regexp(t, re, text, "toolchain directive should exist with go<major>.<minor>[.<patch>]")
}

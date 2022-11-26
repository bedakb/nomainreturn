package nomainreturn

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWithDefaultConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	analysistest.Run(t,
		filepath.Join(filepath.Dir(wd), "nomainreturn/testdata"),
		NewAnalyzer(DefaultConfig),
		"default")
}

func TestWithCustomConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	analysistest.Run(t,
		filepath.Join(filepath.Dir(wd), "nomainreturn/testdata"),
		NewAnalyzer(Config{
			AllowPackages: []string{"custom"},
		}),
		"custom")
}

func Test_isPkgAllowed_PkgIsAllowed(t *testing.T) {
	pkg := "foo"
	allowedPkgs := []string{"main", "bar", "foo"}
	allowed := isPkgAllowed(pkg, allowedPkgs)
	if !allowed {
		t.Errorf("isPkgAllowed() = %t; want = %t", allowed, true)
	}
}

func Test_isPkgAllowed_PkgIsNotAllowed(t *testing.T) {
	pkg := "bar"
	allowedPkgs := []string{"main", "server"}
	allowed := isPkgAllowed(pkg, allowedPkgs)
	if allowed {
		t.Errorf("isPkgAllowed() = %t; want = %t", allowed, false)
	}
}

package nomainreturn

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	t.Logf("wd is %+v", wd)

	analysistest.Run(t,
		filepath.Join(filepath.Dir(wd), "testdata"),
		Analyzer,
		"p")
}

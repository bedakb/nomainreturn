package main

import (
	"github.com/bedakb/nomainreturn/nomainreturn"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(nomainreturn.Analyzer)
}

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/zackarysantana/overstory/src/clientmux"
)

func main() {
	var (
		muxFunc string
		pkg     string
		outFile string
		baseVar string
	)
	flag.StringVar(&muxFunc, "mux", "", "<import/path>.<Func>")
	flag.StringVar(&pkg, "pkg", "api", "package name in generated file")
	flag.StringVar(&outFile, "o", "client_gen.go", "output file")
	flag.StringVar(&baseVar, "base", "baseURL", "parameter name for base URL")
	flag.Parse()

	parts := strings.Split(muxFunc, ".")
	if len(parts) < 2 {
		log.Fatal("-mux must be import/path.FuncName")
	}
	importPath := strings.Join(parts[:len(parts)-1], ".")
	funcName := parts[len(parts)-1]

	// 1. Build plugin containing the mux-builder.
	tmp := filepath.Join(os.TempDir(), "cmuxgen_plugin.so")
	cmd := exec.Command("go", "build",
		"-buildmode=plugin",
		"-tags=cmuxgen", //  ⇦ keep only the tag you need
		"-o", tmp, importPath)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("building plugin: %v", err)
	}

	// 2. Load plugin and call the mux factory.
	p, err := plugin.Open(tmp)
	if err != nil {
		log.Fatal(err)
	}

	sym, err := p.Lookup(funcName)
	if err != nil {
		log.Fatal(err)
	}

	buildMux, ok := sym.(func() interface{})
	if !ok {
		log.Fatalf("%s does not match func() interface{}", funcName)
	}

	mux, ok := buildMux().(*clientmux.ClientMux)
	if !ok {
		log.Fatalf("%s did not return *clientmux.ClientMux", funcName)
	}

	// 3. Generate client code.
	src, err := clientmux.Generate(mux, pkg, baseVar)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outFile, src, 0o644); err != nil {
		log.Fatal(err)
	}
}

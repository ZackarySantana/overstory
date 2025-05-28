// cmd/clientmux-gen/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/zackarysantana/overstory/src/clientmux"
	clientmuxgen "github.com/zackarysantana/overstory/src/clientmux-gen"
)

func main() {
	var (
		muxFunc string
		pkg     string
		outFile string
		baseVar string
	)
	flag.StringVar(&muxFunc, "mux", "", "<import/path>.<Func> (required)")
	flag.StringVar(&pkg, "pkg", "api", "package name in generated file")
	flag.StringVar(&outFile, "o", "client_gen.go", "output file")
	flag.StringVar(&baseVar, "base", "baseURL", "parameter name for base URL")
	flag.Parse()

	if muxFunc == "" {
		log.Fatal("-mux is required")
	}

	parts := strings.Split(muxFunc, ".")
	if len(parts) < 2 {
		log.Fatal("-mux must look like import/path.Func")
	}
	importPath := strings.Join(parts[:len(parts)-1], ".")
	funcName := parts[len(parts)-1]

	//----------------------------------------------------------------------
	// 1.  Build a tiny wrapper plugin that calls <importPath>.<funcName>.
	//----------------------------------------------------------------------
	tmpDir, err := os.MkdirTemp("", "cmuxgen-*")
	check(err)
	defer os.RemoveAll(tmpDir) // best-effort cleanup

	src := fmt.Sprintf(`
		package main
		import upkg "%s"
		// BuildMux is what the generator will reflect over.
		func BuildMux() interface{} { return upkg.%s() }
	`, importPath, funcName)

	wrapper := filepath.Join(tmpDir, "main.go")
	check(os.WriteFile(wrapper, []byte(src), 0o644))

	soFile := filepath.Join(tmpDir, "cmuxgen.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, wrapper)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	check(cmd.Run())

	//----------------------------------------------------------------------
	// 2.  Load the plugin and get the prepared *clientmux.ClientMux.
	//----------------------------------------------------------------------
	p, err := plugin.Open(soFile) // <--  FIX 1
	check(err)

	sym, err := p.Lookup("BuildMux") // <--  FIX 2
	check(err)

	buildMux, ok := sym.(func() interface{})
	if !ok {
		log.Fatal("BuildMux has wrong type signature")
	}

	mux, ok := buildMux().(*clientmux.ClientMux)
	if !ok {
		log.Fatalf("%s() did not return *clientmux.ClientMux", funcName)
	}

	//----------------------------------------------------------------------
	// 3.  Generate client code and write it out.
	//----------------------------------------------------------------------
	code, err := clientmuxgen.Generate(mux, pkg, importPath)
	check(err)
	check(os.MkdirAll(filepath.Dir(outFile), 0o755))
	check(os.WriteFile(outFile, code, 0o644))
}

// simple helper
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

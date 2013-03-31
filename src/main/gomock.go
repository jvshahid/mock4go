package main

import (
	"fmt"
	"go/build"
	"gomock"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func createTempDir() (string, error) {
	tmp := os.TempDir()
	name := path.Join(tmp, "gomock")
	return name, os.MkdirAll(name, os.ModePerm)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s packages\n", os.Args[0])
		os.Exit(1)
	}

	tmpDir, err := createTempDir()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create temporary directory. Error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("goroot: %s\n", build.Default.GOROOT)

	pkgs := make([]string, 0, len(os.Args[1:]))

	for _, packageName := range os.Args[1:] {
		pkg := gomock.InstrumentPackage(packageName, tmpDir)
		pkgs = append(pkgs, pkg.Name)
	}

	// run the tests
	args := append([]string{"", "test", "-v"}, pkgs...)
	goBinPath, err := exec.LookPath("go")
	proc, err := os.StartProcess(goBinPath, args, &os.ProcAttr{
		Env: []string{
			"GOPATH=" + tmpDir,
			"GOROOT=" + runtime.GOROOT(),
		},
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	proc.Wait()
}

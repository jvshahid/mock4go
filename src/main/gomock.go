package main

import (
	"fmt"
	"go/build"
	"gomock"
	"io"
	"os"
	"path"
)

func createTempDir() (string, error) {
	tmp := os.TempDir()
	name := path.Join(tmp, "gomock")
	return name, os.MkdirAll(name, os.ModePerm)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyPackage(pkg *build.Package, tmpDir string) error {
	// create a subdirectory

	dst := path.Join(tmpDir, pkg.Dir)
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range pkg.GoFiles {
		err := copyFile(path.Join(pkg.Dir, file), path.Join(dst, file))
		if err != nil {
			return err
		}
	}

	for _, file := range pkg.TestGoFiles {
		err := copyFile(path.Join(pkg.Dir, file), path.Join(dst, file))
		if err != nil {
			return err
		}
	}
	return nil
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

	for _, packageName := range os.Args[1:] {
		pkg, _ := gomock.GetPackage(packageName)

		err := copyPackage(pkg, tmpDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
		}

		// fmt.Printf("package %s contains: %s\n", pkg, strings.Join(files, ","))
		for _, file := range pkg.GoFiles {
			content, err := gomock.InstrumentFile(packageName + "/" + file)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Content of %s:\n%s\n", file, content)
		}
	}
}

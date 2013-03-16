package gomock

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
)

type Foo struct {
}

func (f *Foo) Bar() {
}

func GetFiles(packageName string) (goFiles []string, err error) {
	pkg, err := build.ImportDir(packageName, 0)
	goFiles = pkg.GoFiles
	return
}

func AddGoMockImport(f *ast.File) {
	importSpec := &ast.ImportSpec{
		Name: &ast.Ident{
			Name: "foo",
		},
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%#v", "lcs/handler/foo"),
		},
	}

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT, Specs: []ast.Spec{importSpec},
	}

	f.Decls = append([]ast.Decl{importDecl}, f.Decls...)
}

func InstrumentFunctions(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil {
				// fieldList := x.Recv.List[0]
				// name := fieldList.Names[0]
				// with receiver
			} else {
				// without receiver
			}
		}
		return true
	})
}

func InstrumentFile(fileName string) (string, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return "", err
	}
	AddGoMockImport(f)
	buf := bytes.NewBufferString("")
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

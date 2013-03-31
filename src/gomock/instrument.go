package gomock

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path"
	"strconv"
)

type Foo struct {
}

func (f *Foo) Bar() {
}

func GetPackage(packageName string) (*build.Package, error) {
	return build.ImportDir(packageName, 0)
}

func makeIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func AddGoMockImport(f *ast.File) {
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%#v", "gomock"),
		},
	}

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT, Specs: []ast.Spec{importSpec},
	}

	f.Decls = append([]ast.Decl{importDecl}, f.Decls...)
}

func functionName(f *ast.FuncDecl) ast.Expr {
	if f.Recv == nil {
		return &ast.Ident{
			Name: f.Name.Name,
		}
	}
	return &ast.BinaryExpr{
		X: &ast.ParenExpr{
			X: f.Recv.List[0].Type,
		},
		Op: token.PERIOD,
		Y: &ast.Ident{
			Name: f.Name.Name,
		},
	}
}

// construct the return statement that converts the interface{} type
// returned from gomock.FunctionCalled to the expected returned type
func functionReturnExprs(f *ast.FuncDecl) []ast.Expr {
	if f.Type.Results == nil {
		return nil
	}
	exprs := make([]ast.Expr, 0)
	for idx, param := range f.Type.Results.List {
		exprs = append(exprs, &ast.TypeAssertExpr{
			X: &ast.IndexExpr{
				X: &ast.Ident{
					Name: "values",
				},
				Index: &ast.Ident{
					Name: strconv.Itoa(idx),
				},
			},
			Type: param.Type,
		})
	}
	return exprs
}

// Programatically generate the following code:
//    if value, ok, err := gomock.FunctionCalled(myFunctionName, args); ok && err != nil {
//      return value[0].(Type1), value[1].(Type2)
//    }
// at the beginning of the given function declaration.
func instrumentFunction(f *ast.FuncDecl) {
	returnStmts := functionReturnExprs(f)
	returnValues := "_"
	if returnStmts != nil {
		returnValues = "values"
	}
	initStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: returnValues,
			},
			&ast.Ident{
				Name: "ok",
			},
			&ast.Ident{
				Name: "err",
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.Ident{
					Name: "gomock.FunctionCalled",
				},
				Args: []ast.Expr{
					functionName(f),
				},
			},
		},
	}
	condStmt := &ast.BinaryExpr{
		X: &ast.Ident{
			Name: "ok",
		},
		Op: token.LAND,
		Y: &ast.BinaryExpr{
			X: &ast.Ident{
				Name: "err",
			},
			Op: token.EQL,
			Y: &ast.Ident{
				Name: "nil",
			},
		},
	}
	bodyStmt := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: returnStmts,
			},
		},
	}
	stmt := &ast.IfStmt{
		Init: initStmt,
		Cond: condStmt,
		Body: bodyStmt,
	}
	body := f.Body
	stmts := body.List
	stmts = append([]ast.Stmt{stmt}, stmts...)
	body.List = stmts
	// return &ast.CallExpr{Fun: makeIdent("gomock.FunctionCalled"), Args: []ast.Expr{}}
}

func instrumentInterface(intrface *ast.InterfaceType) []ast.Decl {
	declarations := make([]ast.Decl, 0)
	declarations = append(declarations,
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: "MockInterface",
					},
					Type: &ast.StructType{
						Fields: &ast.FieldList{},
					},
				},
			},
		},
	)
	for _, fun := range intrface.Methods.List {
		fmt.Printf("I'm here")
		funType := fun.Type.(*ast.FuncType)
		// add a general declaration one per return value to guarantee
		// they are assigned the zero value, then return those variables
		stmts := make([]ast.Stmt, 0)
		returnVariables := make([]ast.Expr, 0)

		for idx, returnValue := range funType.Results.List {
			// add a declaration
			name := fmt.Sprintf("_temp%d", idx)
			returnVariables = append(returnVariables, &ast.Ident{
				Name: name,
			})
			stmts = append(stmts, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Type: returnValue.Type,
							Names: []*ast.Ident{
								&ast.Ident{
									Name: name,
								},
							},
						},
					},
				},
			})
		}

		stmts = append(stmts, &ast.ReturnStmt{
			Results: returnVariables,
		})

		newDecl := &ast.FuncDecl{
			Name: fun.Names[0],
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "recv",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "MockInterface",
							},
						},
					},
				},
			},
			Type: funType,
			Body: &ast.BlockStmt{
				List: stmts,
			},
		}
		instrumentFunction(newDecl)
		declarations = append(declarations, newDecl)
	}
	return declarations
}

func InstrumentFunctions(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			instrumentFunction(x)
			if x.Recv != nil {
				// fieldList := x.Recv.List[0]
				// name := fieldList.Names[0]
				// with receiver
			} else {
				// without receiver
			}
		case *ast.InterfaceType:
			if x.Incomplete {
				// TODO: what should we do here
				panic("incomplete interface type")
			}
			fmt.Printf("decls: %v\n", f.Decls)
			decls := instrumentInterface(x)
			f.Decls = append(f.Decls, decls...)
			fmt.Printf("decls: %v\n", f.Decls)
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
	InstrumentFunctions(f)
	buf := bytes.NewBufferString("")
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func InstrumentPackage(packageName string, tmpDir string) *build.Package {
	pkg, _ := GetPackage(packageName)

	fmt.Printf("dir: %s, goroot: %v, root: %s\n", pkg.Dir, pkg.Goroot, pkg.SrcRoot)

	if pkg.Goroot {
		return pkg
	}

	for _, importPackageName := range pkg.Imports {
		fmt.Printf("import: %s\n", importPackageName)
		InstrumentPackage(importPackageName, tmpDir)
	}
	for _, importPackageName := range pkg.TestImports {
		fmt.Printf("import: %s\n", importPackageName)
		InstrumentPackage(importPackageName, tmpDir)
	}
	InstrumentPackageRecur(pkg, tmpDir, make(map[string]bool))

	return pkg
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

func InstrumentPackageRecur(pkg *build.Package, tmpDir string, instrumented map[string]bool) {
	if instrumented[pkg.Name] {
		return
	}

	err := copyPackage(pkg, tmpDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

	// fmt.Printf("package %s contains: %s\n", pkg, strings.Join(files, ","))
	for _, file := range pkg.GoFiles {
		fileName := path.Join(tmpDir, pkg.Dir, file)
		content, err := InstrumentFile(fileName)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		file, err := os.OpenFile(fileName, os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		_, err = fmt.Fprintf(file, content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		file.Close()
	}

}

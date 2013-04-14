package api

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

func GetPackage(packageName string) (*build.Package, error) {
	return build.Default.Import(packageName, ".", 0)
}

func makeIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

const GoMockImport = "github.com/jvshahid/gomock"

func AddGoMockImport(f *ast.File) {
	importSpec := &ast.ImportSpec{
		Name: makeIdent("gomock"),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%#v", GoMockImport),
		},
	}

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT, Specs: []ast.Spec{importSpec},
	}

	f.Decls = append([]ast.Decl{importDecl}, f.Decls...)
}

func functionName(f *ast.FuncDecl) ast.Expr {
	if f.Recv == nil {
		return makeIdent(f.Name.Name)
	}
	return &ast.BinaryExpr{
		X: &ast.ParenExpr{
			X: f.Recv.List[0].Type,
		},
		Op: token.PERIOD,
		Y:  makeIdent(f.Name.Name),
	}
}

// construct the return statement that converts the interface{} type
// returned from gomock.FunctionCalled to the expected returned type
func functionReturnExprs(f *ast.FuncDecl, stmts []ast.Stmt) []ast.Stmt {
	if f.Type.Results == nil {
		return nil
	}

	for idx, variable := range f.Type.Results.List {
		value := &ast.IndexExpr{
			X:     makeIdent("values"),
			Index: makeIdent(strconv.Itoa(idx)),
		}
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  value,
				Op: token.NEQ,
				Y:  makeIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{makeIdent(fmt.Sprintf("_temp%d", idx))},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.TypeAssertExpr{
								X:    value,
								Type: variable.Type,
							},
						},
					},
				},
			},
		})
	}
	return stmts
}

// Programatically generate the following code:
//    if value, ok, err := gomock.FunctionCalled(myFunctionName, args); ok && err != nil {
//      return value[0].(Type1), value[1].(Type2)
//    }
// at the beginning of the given function declaration.
func instrumentFunction(f *ast.FuncDecl) bool {
	if f.Name.Name == "init" {
		return false
	}

	variableDeclaration, returnVariables := declareReturnValuesVariables(f.Type)

	returnStmts := functionReturnExprs(f, variableDeclaration)
	returnStmts = append(returnStmts, &ast.ReturnStmt{
		Results: returnVariables,
	})
	returnValues := "_"
	if len(returnStmts) > 1 {
		returnValues = "values"
	}

	functionCalledArgs := []ast.Expr{
		functionName(f),
	}

	if f.Recv != nil && len(f.Recv.List) > 0 {
		functionCalledArgs = append(functionCalledArgs, f.Recv.List[0].Names[0])
	}

	for _, arg := range f.Type.Params.List {
		functionCalledArgs = append(functionCalledArgs, arg.Names[0])
	}

	initStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			makeIdent(returnValues),
			makeIdent("ok"),
			makeIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  makeIdent("gomock.FunctionCalled"),
				Args: functionCalledArgs,
			},
		},
	}
	condStmt := &ast.BinaryExpr{
		X:  makeIdent("ok"),
		Op: token.LAND,
		Y: &ast.BinaryExpr{
			X:  makeIdent("err"),
			Op: token.EQL,
			Y:  makeIdent("nil"),
		},
	}
	bodyStmt := &ast.BlockStmt{
		List: returnStmts,
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
	return true
}

type Interface struct {
	pkg           string
	name          string
	interfaceType *ast.InterfaceType
	file          string
}

var interfaces = make(map[string]Interface) // map from pkg.InterfaceName to the Interface type

func declareReturnValuesVariables(funType *ast.FuncType) ([]ast.Stmt, []ast.Expr) {
	stmts := make([]ast.Stmt, 0)
	returnVariables := make([]ast.Expr, 0)

	if funType.Results != nil {
		for idx, returnValue := range funType.Results.List {
			// add a declaration
			name := fmt.Sprintf("_temp%d", idx)
			returnVariables = append(returnVariables, makeIdent(name))
			stmts = append(stmts, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Type:  returnValue.Type,
							Names: []*ast.Ident{makeIdent(name)},
						},
					},
				},
			})
		}
	}

	return stmts, returnVariables
}

func instrumentInterfaceFunction(interfaceName string,
	intrface *ast.InterfaceType, funName *ast.Ident, funType *ast.FuncType) ast.Decl {
	// add a general declaration one per return value to guarantee
	// they are assigned the zero value, then return those variables
	stmts, returnVariables := declareReturnValuesVariables(funType)

	// set a name to each function argument
	if funType.Params != nil {
		for idx, arg := range funType.Params.List {
			arg.Names = []*ast.Ident{makeIdent(fmt.Sprintf("arg%d", idx))}
		}
	}

	stmts = append(stmts, &ast.ReturnStmt{
		Results: returnVariables,
	})

	newDecl := &ast.FuncDecl{
		Name: funName,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{makeIdent("recv")},
					Type: &ast.StarExpr{
						X: makeIdent("Mock" + interfaceName),
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
	return newDecl
}

func instrumentInterface(name string, intrface *ast.InterfaceType) []ast.Decl {
	declarations := make([]ast.Decl, 0)

	structFunctions := make([]*ast.Field, 0)
	structEmbedded := make([]*ast.Field, 0)

	for _, fun := range intrface.Methods.List {
		switch x := fun.Type.(type) {
		case *ast.FuncType:
			structFunctions = append(structFunctions, fun)
		case *ast.Ident:
			fmt.Printf("x = %s, type = %v\n", x.Name, fun.Type)
			structEmbedded = append(structEmbedded, &ast.Field{
				Type: makeIdent("Mock" + x.Name),
			})
		}
	}

	declarations = append(declarations,
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: makeIdent("Mock" + name),
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: structEmbedded,
						},
					},
				},
			},
		},
	)

	for _, fun := range structFunctions {
		decl := instrumentInterfaceFunction(name, intrface, fun.Names[0], fun.Type.(*ast.FuncType))
		declarations = append(declarations, decl)
	}
	return declarations
}

func InstrumentFunctionsAndInterfaces(f *ast.File) bool {
	addGoMockImport := false

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if instrumentFunction(x) {
				addGoMockImport = true
			}
			if x.Recv != nil {
				// fieldList := x.Recv.List[0]
				// name := fieldList.Names[0]
				// with receiver
			} else {
				// without receiver
			}
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				typeSpec := x.Specs[0].(*ast.TypeSpec)
				if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					if interfaceType.Incomplete {
						// TODO: what should we do here
						panic("incomplete interface type")
					}
					decls := instrumentInterface(typeSpec.Name.Name, interfaceType)
					addGoMockImport = true
					f.Decls = append(f.Decls, decls...)
				}
			}
		}
		return true
	})

	return addGoMockImport
}

func InstrumentFile(fileName string) (string, error) {
	fmt.Fprintf(os.Stderr, "instrumenting file %s\n", fileName)
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return "", err
	}
	if InstrumentFunctionsAndInterfaces(f) {
		AddGoMockImport(f)
	}
	f.Comments = nil
	buf := bytes.NewBufferString("")
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func InstrumentPackage(packageName string, tmpDir string) *build.Package {
	if packageName == "C" {
		return nil
	}

	pkg, err := GetPackage(packageName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if pkg.Goroot {
		return pkg
	}

	fmt.Fprintf(os.Stderr, "instrumenting package %s\n", packageName)

	for _, importPackageName := range pkg.Imports {
		InstrumentPackage(importPackageName, tmpDir)
	}
	for _, importPackageName := range pkg.TestImports {
		InstrumentPackage(importPackageName, tmpDir)
	}
	InstrumentPackageRecur(pkg, tmpDir, make(map[string]bool))

	return pkg
}

func copyPackage(pkg *build.Package, tmpDir string) error {
	// create a subdirectory

	dst := path.Join(tmpDir, pkg.ImportPath)
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	filesLists := [][]string{
		pkg.GoFiles,
		pkg.TestGoFiles,
		pkg.CgoFiles,
		pkg.CFiles,
		pkg.HFiles,
		pkg.SFiles,
	}

	for _, list := range filesLists {
		for _, file := range list {
			err := copyFile(path.Join(pkg.Dir, file), path.Join(dst, file))
			if err != nil {
				return err
			}
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

	// copy only, don't instrument gomock
	if pkg.ImportPath == GoMockImport || pkg.ImportPath == "launchpad.net/gocheck" {
		return
	}

	// fmt.Printf("package %s contains: %s\n", pkg, strings.Join(files, ","))
	for _, file := range pkg.GoFiles {
		fileName := path.Join(tmpDir, pkg.ImportPath, file)
		content, err := InstrumentFile(fileName)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		_, err = file.Write([]byte(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		file.Close()
	}
}

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	fset := token.NewFileSet()

	// file := &ast.File{}

	f2, err := parser.ParseFile(fset, "whatever.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	// printer.Fprint(os.Stdout, fset, f2)
	var done bool
	astutil.Apply(f2, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.FuncDecl:
			if done {
				return true
			}
			done = true
			c.Replace(&ast.FuncDecl{
				Name: ast.NewIdent("foobar"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.SelectStmt{
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.CommClause{
										Comm: &ast.AssignStmt{
											Lhs: []ast.Expr{ast.NewIdent("v1"), ast.NewIdent("ok")},
											Tok: token.DEFINE,
											// Rhs: []ast.Expr{&ast.ChanType{Dir: ast.RECV, Value: ast.NewIdent("c1")}},
											Rhs: []ast.Expr{&ast.UnaryExpr{Op: token.ARROW, X: ast.NewIdent("c1")}},
										},
										Body: []ast.Stmt{
											&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("f1")}},
										},
									},
								},
							},
						},
					},
				},
				// Recv: &ast.FieldList{
				// 	List: []*ast.Field{
				// 		{
				// 			Names: []*ast.Ident{
				// 				ast.NewIdent("foo"),
				// 			},
				// 			Type: ast.NewIdent("int"),
				// 		},
				// 	},
				// },
			})
			// return false
		}
		fmt.Printf("%T: %v\n", n, n)

		return true
	})

	fmt.Println("Modified AST:")
	printer.Fprint(os.Stdout, fset, f2)
}

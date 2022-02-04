package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

func genSelectCall(count int, withDefault bool, withOk bool) (n ast.Node) {

	var (
		defaultFunction = ast.NewIdent("df")
	)

	name := fmt.Sprintf("Select%d", count)

	if withDefault {
		name += "Default"
	}
	if withOk {
		name += "OK"
	}

	body := &ast.BlockStmt{}
	params := &ast.FieldList{}
	typeParams := &ast.Field{
		Type: ast.NewIdent("any"),
	}

	for i := 0; i < count; i++ {
		isLast := i+1 == count
		selectStmt := &ast.SelectStmt{
			Body: &ast.BlockStmt{},
		}

		fxnArgs := []*ast.Field{
			{
				Names: []*ast.Ident{
					// ast.NewIdent(fmt.Sprintf("v%d", i)),
				},
				Type: ast.NewIdent(fmt.Sprintf("T%d", i)),
			},
		}
		if withOk {
			fxnArgs = append(fxnArgs, &ast.Field{
				Names: []*ast.Ident{
					// ast.NewIdent("ok"),
				},
				Type: ast.NewIdent("bool"),
			})
		}

		params.List = append(params.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("c%d", i))},
			Type: &ast.ChanType{
				Dir:   ast.RECV,
				Value: ast.NewIdent(fmt.Sprintf("T%d", i)),
			},
		},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("f%d", i))},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: fxnArgs,
					},
				},
			},
		)

		typeParams.Names = append(typeParams.Names, ast.NewIdent(fmt.Sprintf("T%d", i)))

		for j := 0; j <= i; j++ {
			val := ast.NewIdent(fmt.Sprintf("v%d", j))
			channel := ast.NewIdent(fmt.Sprintf("c%d", j))
			fxn := ast.NewIdent(fmt.Sprintf("f%d", j))
			boolean := ast.NewIdent("ok")

			args := []ast.Expr{val}
			if withOk {
				args = append(args, boolean)
			}

			clause := &ast.CommClause{
				Comm: &ast.AssignStmt{
					Lhs: args,
					Tok: token.DEFINE,
					Rhs: []ast.Expr{&ast.UnaryExpr{Op: token.ARROW, X: channel}},
				},
				Body: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun:  fxn,
						Args: args,
					}},
				},
			}

			if !isLast {
				clause.Body = append(clause.Body, &ast.ReturnStmt{})
			}

			selectStmt.Body.List = append(selectStmt.Body.List, clause)
		}
		if !isLast {
			selectStmt.Body.List = append(selectStmt.Body.List,
				&ast.CommClause{
					Comm: nil, // default
					Body: nil,
				},
			)
		} else if isLast && withDefault {
			selectStmt.Body.List = append(selectStmt.Body.List,
				&ast.CommClause{
					Comm: nil, // default
					Body: []ast.Stmt{
						&ast.ExprStmt{X: &ast.CallExpr{Fun: defaultFunction}},
					},
				},
			)
		}

		body.List = append(body.List, selectStmt)
	}

	if withDefault {
		params.List = append(params.List,
			&ast.Field{
				Names: []*ast.Ident{defaultFunction},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: nil,
					},
				},
			},
		)
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: &ast.FuncType{
			Params:     params,
			TypeParams: &ast.FieldList{List: []*ast.Field{typeParams}},
		},
		Body: body,
	}
}

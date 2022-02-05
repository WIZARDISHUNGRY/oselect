package main

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func genSelectCall(c *astutil.Cursor, count int, withDefault bool, withOk bool, isSend bool) *ast.FuncDecl {

	var (
		defaultFunction = ast.NewIdent("df")
	)

	if isSend && withOk {
		panic("isSend && withOk")
	}

	direction := ast.ChanDir(ast.RECV)

	name := "Recv"
	if isSend {
		direction = ast.SEND
		name = "Send"
	}
	name = fmt.Sprintf("%s%d", name, count)

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
				Type: ast.NewIdent(fmt.Sprintf("T%d", i)),
			},
		}
		if withOk {
			fxnArgs = append(fxnArgs, &ast.Field{
				Type: ast.NewIdent("bool"),
			})
		}

		var results []*ast.Field
		if isSend {
			results, fxnArgs = fxnArgs, results
		}

		params.List = append(params.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("c%d", i))},
			Type: &ast.ChanType{
				Dir:   direction,
				Value: ast.NewIdent(fmt.Sprintf("T%d", i)),
			},
		},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("f%d", i))},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: fxnArgs,
					},
					Results: &ast.FieldList{
						List: results,
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

			var clause *ast.CommClause

			if isSend {
				clause = &ast.CommClause{
					Comm: &ast.AssignStmt{
						Lhs: []ast.Expr{channel},
						Tok: token.ARROW,
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun: fxn,
						}},
					},
				}
			} else {
				clause = &ast.CommClause{
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
		// Pain in the ass!
		// Doc: &ast.CommentGroup{
		// 	List: []*ast.Comment{
		// 		{Text: fmt.Sprintf("// %s", name), Slash: token.NoPos},
		// 	},
		// },
		Body: body,
	}
}

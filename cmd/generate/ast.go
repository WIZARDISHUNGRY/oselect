package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// TODO: use parser.ParseExpr() to make this more readable

type functionType int

const (
	typeRecv = iota
	typeSend
	typeSelect
)

var names = map[functionType]string{
	typeRecv:   "Recv",
	typeSend:   "Send",
	typeSelect: "Select",
}

func mustParseExp(s string) ast.Expr {
	expr, err := parser.ParseExpr(s)
	if err != nil {
		panic(fmt.Sprintf("mustParseExp(%s): %s", s, err.Error()))
	}
	return expr
}

func genSelectCall(c *astutil.Cursor, count int, withDefault bool, withOk bool, fxnType functionType) *ast.FuncDecl {
	defaultFunction := ast.NewIdent("df")

	isSend := fxnType == typeSend

	if isSend && withOk {
		panic("isSend && withOk")
	}

	direction := ast.ChanDir(ast.RECV)

	name := names[fxnType]
	if isSend {
		direction = ast.SEND
	}
	if withOk && fxnType == typeRecv {
		name += "OK"
	}

	name = fmt.Sprintf("%s%d", name, count)

	if withDefault {
		name += "Default"
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

		if fxnType == typeSelect {
			typeExpr := mustParseExp(fmt.Sprintf("Param[T%d]", i))

			params.List = append(params.List, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("p%d", i))},
				Type:  typeExpr,
			})
		} else {

			params.List = append(params.List, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("c%d", i))},
				Type: &ast.ChanType{
					Dir:   direction,
					Value: ast.NewIdent(fmt.Sprintf("T%d", i)),
				},
			})

			if fxnType == typeRecv {
				params.List = append(params.List,
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
			} else {
				typeExpr := mustParseExp(fmt.Sprintf("T%d", i))

				params.List = append(params.List, &ast.Field{
					Names: []*ast.Ident{ast.NewIdent(fmt.Sprintf("v%d", i))},
					Type:  typeExpr,
				})
			}
		}

		typeParams.Names = append(typeParams.Names, ast.NewIdent(fmt.Sprintf("T%d", i)))

		for j := 0; j <= i; j++ {
			val := ast.NewIdent(fmt.Sprintf("v%d", j))
			channel := ast.NewIdent(fmt.Sprintf("c%d", j))
			fxn := ast.NewIdent(fmt.Sprintf("f%d", j))
			boolean := ast.NewIdent("ok")
			param := fmt.Sprintf("p%d", j)

			args := []ast.Expr{val}
			if withOk {
				args = append(args, boolean)
			}

			var clauses []*ast.CommClause

			if fxnType == typeSend {
				clauses = []*ast.CommClause{{
					Comm: &ast.AssignStmt{
						Lhs: []ast.Expr{channel},
						Tok: token.ARROW,
						Rhs: []ast.Expr{mustParseExp(fmt.Sprintf("v%d", j))},
					},
				}}
			} else if fxnType == typeRecv {
				clauses = []*ast.CommClause{{
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
				}}
			} else if fxnType == typeSelect {
				clauses = []*ast.CommClause{{
					Comm: &ast.AssignStmt{
						Lhs: []ast.Expr{mustParseExp(param + ".SendChan")},
						Tok: token.ARROW,
						Rhs: []ast.Expr{mustParseExp(param + ".SendValue")},
					},
				}, {
					Comm: &ast.AssignStmt{
						Lhs: args,
						Tok: token.DEFINE,
						Rhs: []ast.Expr{&ast.UnaryExpr{Op: token.ARROW, X: mustParseExp(param + ".RecvChan")}},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{X: &ast.CallExpr{
							Fun:  mustParseExp(param + ".RecvFunc"),
							Args: args,
						}},
					},
				},
				}
			}

			if !isLast {
				for _, clause := range clauses {
					clause.Body = append(clause.Body, &ast.ReturnStmt{})
				}
			}

			for _, clause := range clauses {
				selectStmt.Body.List = append(selectStmt.Body.List, clause)
			}

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

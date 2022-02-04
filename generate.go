package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

const MAX_CHANNELS = 10

func main() {
	fset := token.NewFileSet()

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
			c.Delete()
			for i := 1; i <= MAX_CHANNELS; i++ {
				c.InsertBefore(genSelectCall(i, false, false))
				c.InsertBefore(genSelectCall(i, true, false))
				c.InsertBefore(genSelectCall(i, false, true))
				c.InsertBefore(genSelectCall(i, true, true))
			}
			// return false
		}
		// fmt.Printf("%T: %v\n", n, n)

		return true
	})

	// fmt.Println("Modified AST:")
	printer.Fprint(os.Stdout, fset, f2)
}

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

const MAX_CHANNELS = 9

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Println(args)
		panic("takes zero or one args")
	}
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "template", template,
		parser.AllErrors|parser.ParseComments,
	)
	if err != nil {
		panic(err)
	}

	f.Comments = nil

	var done bool
	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.FuncDecl:
			if done {
				return true
			}
			done = true

			c.Delete()
			for i := 2; i <= MAX_CHANNELS; i++ {
				// Select
				c.InsertBefore(genSelectCall(c, i, false, true, typeSelect))
				c.InsertBefore(genSelectCall(c, i, true, true, typeSelect))
				// Recv
				c.InsertBefore(genSelectCall(c, i, false, false, typeRecv))
				c.InsertBefore(genSelectCall(c, i, true, false, typeRecv))
				c.InsertBefore(genSelectCall(c, i, false, true, typeRecv))
				c.InsertBefore(genSelectCall(c, i, true, true, typeRecv))
				// Send
				c.InsertBefore(genSelectCall(c, i, false, false, typeSend))
				c.InsertBefore(genSelectCall(c, i, true, false, typeSend))
			}
		}

		return true
	})

	var outStream io.Writer

	if len(args) > 0 {
		out := args[0]
		outFile, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()
		outStream = outFile
	} else {
		outStream = os.Stdout
	}

	fmt.Fprint(outStream, "// Code generated by a tool. DO NOT EDIT.\n\n") // easier than AST
	err = format.Node(outStream, fset, f)
	if err != nil {
		panic(err)
	}
}

const template = `
package oselect

func init() {
	panic("this should never be included")
}
`

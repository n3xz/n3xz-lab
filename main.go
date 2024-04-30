package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, "../infra/main.go", nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		
		mapOfAwsCalls := make(map[string]string)
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if x.X.(*ast.Ident).Name == "s3" {
					log.Println("Found pulumi code that uses the aws sdk")
					mapOfAwsCalls[x.Sel.Name] = x.X.(*ast.Ident).Name
				}
			}
			return true
		})
		fmt.Println(mapOfAwsCalls)

		return nil
	})
}

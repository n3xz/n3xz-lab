package iamAnalyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func AnalyzePermissions(path string) (map[string]int, error) {
	permissions := make(map[string]int)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", path, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		funcName, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		resourceType := funcName.X.(*ast.Ident).Name
		operation := funcName.Sel.Name

		permissions := []string{}

		if resourceType == "ec2" && operation == "CreateInstance" {
			permissions = append(permissions, "ec2:CreateInstance")
			permissions = append(permissions, "ec2:DescribeInstances")
		} else if resourceType == "s3" && operation == "CreateBucket" {
			permissions = append(permissions, "s3:CreateBucket")
		}
		return true
	})

	return permissions, nil
}
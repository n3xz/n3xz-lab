package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error{
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, "../static-site/main.go", nil, 0)
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

					tmpJSON0, err := json.Marshal(map[string]interface{}{
						"Version": "2012-10-17",
						"Statement": []map[string]interface{}{
							map[string]interface{}{
								"Action": []string{
									"ec2:Describe*",
								},
								"Effect":   "Allow",
								"Resource": "*",
							},
						},
					})
					if err != nil {
						log.Fatal(err)
					}

					_, err = iam.NewRolePolicy(ctx, "s3Policy", &iam.RolePolicyArgs{
						Policy: pulumi.String(tmpJSON0),
						Role:   pulumi.Any(aws_iam_role.Example.Arn),
						
				}
			}
			
			return true
		})
		return nil
	})
}

package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	roledeploy "github.com/the0xsec/pulumi-apps/iam-manager/iamrole"
)

type role struct {
	roleName string
	serviceName string
	permissionSet string
}

var roles = []role{
	{
		roleName:      "lambda-admin",
		serviceName:   "lambda.amazonaws.com",
		permissionSet: "arn:aws:iam::aws:policy/AWSLambda_FullAccess",
	},
	{
		roleName:      "s3-admin",
		serviceName:   "s3.amazonaws.com",
		permissionSet: "arn:aws:iam::aws:policy/AmazonS3FullAccess",
	},
	{
		roleName:      "ec2-admin",
		serviceName:   "ec2.amazonaws.com",
		permissionSet: "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
	},
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, role := range roles {
			_, err := roledeploy.CreateRole(ctx, role.roleName, role.serviceName, role.permissionSet)
			if err != nil {
				return err
			}

			ctx.Export(role.roleName, pulumi.String("Role created"))
		}
		
		return nil
	})
}

package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ec2Assume, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
            "Statement": []interface{}{
                map[string]interface{}{
                    "Effect": "Allow",
                    "Principal": map[string]interface{}{
                        "Service": "ec2.amazonaws.com",
                    },
                    "Action": "sts:AssumeRole",
                },
            },
		})
		if err != nil {
			return err
		}

		json0 := string(ec2Assume)
		_, err = iam.NewRole(ctx, "api-gateway-admin", &iam.RoleArgs{
			Name:             pulumi.String("api-gateway-admin"),
			AssumeRolePolicy: pulumi.String(json0),
			Tags: pulumi.StringMap{
				"iam-manager": pulumi.String("pulumi-stacked"),
			},
		})
		if err!= nil {
            return err
        }

		return nil
	})
}

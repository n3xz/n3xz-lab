package roledeploy

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRole(ctx *pulumi.Context, roleName string, serviceName string, permissionSet string) (*iam.Role, error) {
    assumeService, err := json.Marshal(map[string]interface{}{
        "Version": "2012-10-17",
        "Statement": []interface{}{
            map[string]interface{}{
                "Effect": "Allow",
                "Principal": map[string]interface{}{
                    "Service": serviceName,
                },
                "Action": "sts:AssumeRole",
            },
        },
    })
    if err != nil {
        return nil, err
    }

    role, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
        AssumeRolePolicy: pulumi.String(string(assumeService)),
    })
    if err != nil {
        return nil, err
    }

    _, err = iam.NewRolePolicyAttachment(ctx, roleName+"-policy", &iam.RolePolicyAttachmentArgs{
        PolicyArn: pulumi.String(permissionSet),
        Role:      role.Name,
    })
    if err != nil {
        return nil, err
    }

    return role, nil
}
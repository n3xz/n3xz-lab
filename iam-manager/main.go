package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-doppler/sdk/go/doppler"
)



func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		project := string(pulumi.String("iam-manager"))
		config := string(pulumi.String("dev"))

		credentials, err := doppler.GetSecrets(ctx, &doppler.GetSecretsArgs{
			Project: &project,
			Config:  &config,
		})
		if err != nil {
			return err
		}

		for key, value := range credentials.Map {
			if key != "AWS_SECRET_ACCESS_KEY" {
				ctx.Export(key, pulumi.String(value))
			} else {
				ctx.Export(key, pulumi.String("********"))
			}
		}

		return nil
	})
}

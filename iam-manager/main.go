package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-doppler/sdk/go/doppler"
)

func credGrabber(ctx *pulumi.Context) map[string]string  {

	project := string(pulumi.String("iam-manager"))
	config := string(pulumi.String("dev"))

	credentials, err := doppler.GetSecrets(ctx, &doppler.GetSecretsArgs{
		Project: &project,
		Config:  &config,
	})
	if err != nil {
		return nil
	}

	return credentials.Map
}


func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		creds := credGrabber(ctx)
		if creds == nil {
			return nil
		}

		provider, err := aws.NewProvider(ctx, "aws", &aws.ProviderArgs{
			Region: pulumi.String("us-west-2"),
			AccessKey: pulumi.String(creds["AWS_ACCESS_KEY_ID"]),
			SecretKey: pulumi.String(creds["AWS_SECRET_ACCESS_KEY"]),
		})
		if err != nil {
			return err
		}

		user, errp := iam.NewUser(ctx, "my-user", &iam.UserArgs{
			Path: pulumi.String("/"),
		}, pulumi.Provider(provider))
		if errp != nil {
			return errp
		}

		ctx.Export("user-name", user.Name)
		return nil
	})
}

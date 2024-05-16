package main

import (
	"io/ioutil"
	"net/http"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/vpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getPublicIP() (string, error) {
    resp, err := http.Get("https://api.ipify.org")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    ip, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(ip), nil
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		allowSSH, err := ec2.NewSecurityGroup(ctx, "allowSSH", &ec2.SecurityGroupArgs{
			Name:        pulumi.String("allow-ssh"),
			Description: pulumi.String("Allow SSH inbound traffic from my IP"),
			VpcId:       pulumi.String("vpc-076c3059fe1161101"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("allow-ssh"),
			},
		})
		errCheck(err)

		ip, err := getPublicIP()
		errCheck(err)

		ollama, err := ec2.NewInstance(ctx, "ollama", &ec2.InstanceArgs{
			Ami:                      pulumi.String("ami-03af5c83add1a2df7"),
			InstanceType:             pulumi.String("t2.micro"),
			AssociatePublicIpAddress: pulumi.Bool(false),
			VpcSecurityGroupIds: 	pulumi.StringArray{allowSSH.ID()},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("ollama"),
			},
		})
		errCheck(err)

		_, err = vpc.NewSecurityGroupIngressRule(ctx, "allow_ssh_ipv4", &vpc.SecurityGroupIngressRuleArgs{
			SecurityGroupId: allowSSH.ID(),
			CidrIpv4:        pulumi.String(ip + "/32"),
			FromPort:        pulumi.Int(22),
			ToPort:          pulumi.Int(22),
			IpProtocol:      pulumi.String("tcp"),
		})
		errCheck(err)

		_, err = vpc.NewSecurityGroupIngressRule(ctx, "allow_all_icmp_ipv4", &vpc.SecurityGroupIngressRuleArgs{
			SecurityGroupId: allowSSH.ID(),
			CidrIpv4:        pulumi.String(ip + "/32"),
			FromPort:        pulumi.Int(0),
			ToPort:          pulumi.Int(0),
			IpProtocol:      pulumi.String("icmp"),
		})
		errCheck(err)

		_, err = vpc.NewSecurityGroupEgressRule(ctx, "allowSSHOutbound", &vpc.SecurityGroupEgressRuleArgs{
			SecurityGroupId: allowSSH.ID(),
			CidrIpv4:        pulumi.String(ip + "/32"),
			FromPort:        pulumi.Int(22),
			ToPort:          pulumi.Int(22),
			IpProtocol:      pulumi.String("tcp"),
		})
		errCheck(err)

		ctx.Export("amiId", ollama.Ami)
		ctx.Export("instanceType", ollama.InstanceType)
		ctx.Export("publicIp", ollama.PublicIp)


		return nil
	})
}

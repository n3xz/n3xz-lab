package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func tags(name string) pulumi.StringMap {
	return pulumi.StringMap{
		"Name": pulumi.String(name),
	}
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		vpc, err := ec2.NewVpc(ctx, "my-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String("192.168.0.0/22"),
			Tags:      tags("my-vpc"),
		})
		errCheck(err)

		privateSub, err := ec2.NewSubnet(ctx, "my-private-subnet", &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String("192.168.0.0/24"),
			Tags:      tags("my-private-subnet"),
		})

		publicSub, err := ec2.NewSubnet(ctx, "my-public-subnet", &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String("192.168.1.0/24"),
			Tags:      tags("my-public-subnet"),
		})
		errCheck(err)

		routeTable, err := ec2.NewRouteTable(ctx, "my-route-table", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Tags: tags("my-route-table"),
		})
		errCheck(err)

		_, err = ec2.NewRouteTableAssociation(ctx, "my-public-subnet-association", &ec2.RouteTableAssociationArgs{
			SubnetId:     publicSub.ID(),
			RouteTableId: routeTable.ID(),
		})
		errCheck(err)

		_, err = ec2.NewRouteTableAssociation(ctx, "my-private-subnet-association", &ec2.RouteTableAssociationArgs{
			SubnetId:     privateSub.ID(),
			RouteTableId: routeTable.ID(),
		})
		errCheck(err)

		ig, err := ec2.NewInternetGateway(ctx, "my-ig", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags:  tags("my-ig"),
		})
		errCheck(err)

		_, err = ec2.NewRoute(ctx, "my-public-subnet-route", &ec2.RouteArgs{
			RouteTableId:         routeTable.ID(),
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			GatewayId:            ig.ID(),
		})
		errCheck(err)

		
		NetACL, err := ec2.NewNetworkAcl(ctx, "my-network-acl", &ec2.NetworkAclArgs{
			VpcId: vpc.ID(),
			Tags:  tags("my-network-acl"),
		})
		errCheck(err)

		_, err = ec2.NewNetworkAclAssociation(ctx, "my-public-subnet-acl-association", &ec2.NetworkAclAssociationArgs{
			SubnetId:     publicSub.ID(),
			NetworkAclId: NetACL.ID(),
		})
		errCheck(err)
	
		denyRules := []struct {
			name    string
			egress  bool
		}{
			{"my-public-subnet-inbound-acl-rule", false},
			{"my-public-subnet-outbound-acl-rule", true},
		}

		for _, rule := range denyRules {
			_, err = ec2.NewNetworkAclRule(ctx, rule.name, &ec2.NetworkAclRuleArgs{
				NetworkAclId: NetACL.ID(),
				CidrBlock:    publicSub.CidrBlock,
				Egress:       pulumi.Bool(rule.egress),
				Protocol:     pulumi.String("-1"),
				RuleAction:   pulumi.String("deny"),
				RuleNumber:   pulumi.Int(100),
			})
			errCheck(err)
		}

		allowRules := []struct {
			name    string
			egress  bool
		}{
			{"my-public-subnet-outbound-443", true},
			{"my-public-subnet-inbound-443", false},
		}

		for _, rule := range allowRules {
			_, err = ec2.NewNetworkAclRule(ctx, rule.name, &ec2.NetworkAclRuleArgs{
				NetworkAclId: NetACL.ID(),
				CidrBlock:    publicSub.CidrBlock,
				Egress:       pulumi.Bool(rule.egress),
				Protocol:     pulumi.String("6"),
				RuleAction:   pulumi.String("allow"),
				RuleNumber:   pulumi.Int(10),
				FromPort:     pulumi.Int(443),
				ToPort:       pulumi.Int(443),
			})
			errCheck(err)
		}

		

		webSg, err := ec2.NewSecurityGroup(ctx, "web-sg", &ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Name:  pulumi.String("web-sg"),
			Tags:  tags("web-sg"),
		})
		errCheck(err)

		sgRules := []struct {
			name string
		}{
			{"web-sg-inbound-443"},
			{"web-sg-outbound-443"},
		}

		for _, rule := range sgRules {

			var ruleType pulumi.String
			if rule.name == "web-sg-inbound-443" {
				ruleType = pulumi.String("ingress")
			} else {
				ruleType = pulumi.String("egress")
			}

			_, err = ec2.NewSecurityGroupRule(ctx, rule.name, &ec2.SecurityGroupRuleArgs{
				SecurityGroupId: webSg.ID(),
				Type:            ruleType,
				Protocol:        pulumi.String("-1"),
				CidrBlocks:      pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				FromPort:        pulumi.Int(443),
				ToPort:          pulumi.Int(443),
			})
			errCheck(err)
		}
		errCheck(err)
		return nil
	})
}

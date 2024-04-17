package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
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
		bucket, err := s3.NewBucket(ctx, "hive-net-infra", nil)
		errCheck(err)

		hiveVPC, err := ec2.NewVpc(ctx, "hive-net-infra", &ec2.VpcArgs{
			CidrBlock: pulumi.String("10.0.0.0/25"),
			Tags:      tags("hive-net-infra"),
		})
		errCheck(err)

		hivePubSubnet, err := ec2.NewSubnet(ctx, "hive-net-infra-PubSubNet",
			&ec2.SubnetArgs{
				VpcId:     hiveVPC.ID(),
				CidrBlock: pulumi.String("10.0.0.0/26"),
				Tags:      tags("hive-net-infra-PubSubNet"),
			})
		errCheck(err)

		hivePrivSubnet, err := ec2.NewSubnet(ctx, "hive-net-infra-PrivSubNet", &ec2.SubnetArgs{
			VpcId:     hiveVPC.ID(),
			CidrBlock: pulumi.String("10.0.0.64/26"),
			Tags:      tags("hive-net-infra-PrivSubNet"),
		})
		errCheck(err)

		hiveIGW, err := ec2.NewInternetGateway(ctx, "hive-net-infra", &ec2.InternetGatewayArgs{
			VpcId: hiveVPC.ID(),
			Tags:  tags("hive-net-infra"),
		})
		errCheck(err)

		routeTable, err := ec2.NewRouteTable(ctx, "hive-net-infra", &ec2.RouteTableArgs{
			VpcId: hiveVPC.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: hiveIGW.ID(),
				},
			},
			Tags: tags("hive-net-infra"),
		})
		errCheck(err)

		_, err = ec2.NewRouteTableAssociation(ctx, "hive-net-infra-PubSubNet", &ec2.RouteTableAssociationArgs{
			SubnetId:     hivePubSubnet.ID(),
			RouteTableId: routeTable.ID(),
		})
		errCheck(err)

		AccessControlList, err := ec2.NewNetworkAcl(ctx, "hive-net-infra", &ec2.NetworkAclArgs{
			VpcId: hiveVPC.ID(),
			Egress: ec2.NetworkAclEgressArray{
				&ec2.NetworkAclEgressArgs{
					Action:    pulumi.String("allow"),
					FromPort:  pulumi.Int(0),
					Protocol:  pulumi.String("-1"),
					RuleNo:    pulumi.Int(100),
					ToPort:    pulumi.Int(0),
					CidrBlock: pulumi.String("0.0.0.0/0"),
				},
			},
			Tags: tags("hive-net-infra"),
		})
		errCheck(err)

		ctx.Export("ACL", AccessControlList.ID())
		ctx.Export("routeTableID", routeTable.ID())
		ctx.Export("IGW", hiveIGW.ID())
		ctx.Export("privSubnetID", hivePrivSubnet.ID())
		ctx.Export("subnetID", hivePubSubnet.ID())
		ctx.Export("vpcID", hiveVPC.ID())
		ctx.Export("bucketName", bucket.ID())

		fmt.Println("You're all set! 🚀")

		return nil
	})
}

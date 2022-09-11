package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func k8s() {
	pulumi.Run( func(ctx *pulumi.Context) error {
		isDefault := true
		
		// Set to use default VPC
		vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Default: &isDefault})
		if err != nil {
			return err
		}

		// Set to use the above the VPC
		subnet, err := ec2.LookupSubnet(ctxm &ec2.LookupSubnetArgs{VpcId: &vpc.Id})
		if err != nil {
			return err
		}

		// Create a Kubernetes IAM Role
		eksRole, err := iam.NewRole(ctx, "kubernetes-iam-eksRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
			"Version": "2008-10-17",
		    "Statement": [{
		        "Sid": "",
		        "Effect": "Allow",
		        "Principal": {
		            "Service": "eks.amazonaws.com"
		        },
		        "Action": "sts:AssumeRole"
		    }]
		}`),
		})
		if err != nil {
			return err
		}
		
		// List of EKS Policies
		eksPolicies := []string {
			"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
		}

		// Assign the Polciees to the role we have created
		for i, eksPolicy := range eksPolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("rpa-%d", i), &iam.RolePolicyAttachmentArgs{
				PolicyArn: pulumi.String(eksPolicy),
				Role: eksRole.Name,
			})
			if err != nil {
				return err
			}
		}

		// Create the EC2 NodeGroup Role
		nodeGroupRole, err := iam.NewRole(ctx, "nodegroup-iam-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
		})
		if err != nil {
			return err``
		}

		nodeGroupPolicies := []string {
			"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
			"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
			"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
		}

		for i, nodeGroupPolicy := range nodeGroupPolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("ngpa-%d", i), &iam.RolePolicyAttachmentArgs{
				Role: nodeGroupRole.Name,
				PolicyArn: pulumi.String(nodeGroupPolicy),
			})
			if err != nil {
				return err
			}
		}



	}
	)
}

package main

import "testing"


func TestRoles(t *testing.T) {
	expectedRoles := []role{
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

	for i, expectedRole := range expectedRoles {
		if i >= len(roles) {
			t.Errorf("Missing role at index %d", i)
			continue
		}

		role := roles[i]
		if role.roleName != expectedRole.roleName {
			t.Errorf("Expected roleName %s, got %s", expectedRole.roleName, role.roleName)
		}
		if role.serviceName != expectedRole.serviceName {
			t.Errorf("Expected serviceName %s, got %s", expectedRole.serviceName, role.serviceName)
		}
		if role.permissionSet != expectedRole.permissionSet {
			t.Errorf("Expected permissionSet %s, got %s", expectedRole.permissionSet, role.permissionSet)
		}
	}

	if len(roles) > len(expectedRoles) {
		t.Errorf("Extra roles found")
	}
}
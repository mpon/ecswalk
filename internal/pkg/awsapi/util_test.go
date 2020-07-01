package awsapi

import (
	"testing"
)

func TestShortArn(t *testing.T) {
	arn := "arn:aws:ecs:region:account-id:cluster/cluster-name"
	expect := "cluster-name"
	result := ShortArn(arn)

	if result != expect {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}

func TestShortArnInvalidArguments(t *testing.T) {
	arn := "invalid"
	expect := "invalid"
	result := ShortArn(arn)

	if result != expect {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}

func TestShortDockerImage(t *testing.T) {
	image := "us.gcr.io/my-project/my-image:test"
	expectImage := "my-image"
	expectTag := "test"
	resultImage, resultTag := ShortDockerImage(image)

	if resultImage != expectImage {
		t.Fatalf("expect %s\nbut %s", expectImage, resultImage)
	}

	if resultTag != expectTag {
		t.Fatalf("expect %s\nbut %s", expectTag, resultTag)
	}
}

func TestShortDockerImageWithoutTag(t *testing.T) {
	image := "us.gcr.io/my-project/my-image"
	expectImage := "my-image"
	expectTag := ""
	resultImage, resultTag := ShortDockerImage(image)

	if resultImage != expectImage {
		t.Fatalf("expect %s\nbut %s", expectImage, resultImage)
	}

	if resultTag != expectTag {
		t.Fatalf("expect %s\nbut %s", expectTag, resultTag)
	}
}

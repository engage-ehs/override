package override

import (
	"errors"
	"testing"
)

func TestNorm(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"aws_s3bucket", "awss3bucket"},
		{"AWSS3Bucket", "awss3bucket"},
		{"awss3bucket", "awss3bucket"},
	}

	for _, c := range cases {
		if norm(c.in) != c.out {
			t.Errorf("norm(%s): want %s, got %s", c.in, c.out, norm(c.in))
		}
	}
}

type DeploySettings struct {
	AWSS3Bucket                  string
	VaultAddress                 string
	AWSDeploymentApplicationName string

	AWSSession string `canset:"no"`
}

func TestScan(t *testing.T) {
	args := []string{
		"aws_deployment_application_name=dark-knight",
		"aws_s3_bucket=my-codedeploy",
		"vault_address=http://vault.super.secret",
	}

	var cfg DeploySettings

	if err := Scan(args, &cfg); err != nil {
		t.Fatal(err)
	}

	hasValue := func(field, value string) {
		if field != value {
			t.Errorf("want %s, got %s", field, value)
		}
	}

	hasValue(cfg.AWSS3Bucket, "my-codedeploy")
	hasValue(cfg.AWSDeploymentApplicationName, "dark-knight")
	hasValue(cfg.VaultAddress, "http://vault.super.secret")
}

func TestCannotSet(t *testing.T) {
	var cfg DeploySettings

	err := Scan([]string{"aws_session=mysession"}, &cfg)
	want := cannotSetViolation("aws_session")

	if !errors.Is(err, want) {
		t.Fatalf("Did not prevent to change the session: %s", err)
	}
}

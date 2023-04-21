package infastructure

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AwsConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
}

func NewAwsConfig(AWSConfig AwsConfig) aws.Config {
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(AWSConfig.Region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					AWSConfig.AccessKeyID,
					AWSConfig.SecretAccessKey,
					"",
				),
			),
		),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: AWSConfig.Endpoint,
					}, nil
				},
			),
		),
	)
	if err != nil {
		panic(err)
	}
	return awsConfig
}

func AssumeRole(conf aws.Config, roleARn string) aws.Config {
	stsSvc := sts.NewFromConfig(conf)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, roleARn)

	conf.Credentials = aws.NewCredentialsCache(creds)
	return conf
}

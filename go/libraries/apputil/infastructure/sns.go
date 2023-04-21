package infastructure

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNSClient struct {
	*sns.Client
	topicArn string
}

func NewSNSClient(config aws.Config, topicArn string) *SNSClient {
	return &SNSClient{
		Client:   sns.NewFromConfig(config),
		topicArn: topicArn,
	}
}

func (snsClient *SNSClient) PublishMessage(ctx context.Context, message string) (*sns.PublishOutput, error) {
	result, err := snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: &snsClient.topicArn,
		Message:  &message,
	})
	return result, err
}

func (snsClient *SNSClient) PublishMessageWithAttributes(ctx context.Context, message string, messageAttributes map[string]types.MessageAttributeValue) (*sns.PublishOutput, error) {
	result, err := snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn:          &snsClient.topicArn,
		Message:           &message,
		MessageAttributes: messageAttributes,
	})
	return result, err
}

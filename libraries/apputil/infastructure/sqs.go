package infastructure

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log/slog"
)

type SQSClient struct {
	*sqs.Client
	queueUrl string
}

func NewSQSClient(config aws.Config, queueUrl string) *SQSClient {
	return &SQSClient{
		Client:   sqs.NewFromConfig(config),
		queueUrl: queueUrl,
	}
}

func (sqsClient *SQSClient) PollMessage(ctx context.Context, channel chan<- *types.Message) {
	receiveMessageOutput, err := sqsClient.ReceiveMessage(
		ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &sqsClient.queueUrl,
			AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
			WaitTimeSeconds:     20,
			MaxNumberOfMessages: 10,
		},
	)
	if err != nil {
		slog.Error(
			"failed to receive message",
			slog.Any("err", err),
			slog.String("queue_url", sqsClient.queueUrl),
		)
	}

	for _, message := range receiveMessageOutput.Messages {
		channel <- &message
	}
}

func (sqsClient *SQSClient) DeleteMessage(ctx context.Context, message types.Message) {
	_, err := sqsClient.Client.DeleteMessage(
		ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &sqsClient.queueUrl,
			ReceiptHandle: message.ReceiptHandle,
		},
	)
	if err != nil {
		slog.Error(
			"failed to delete message",
			slog.Any("err", err),
			slog.String("queue_url", sqsClient.queueUrl),
		)
	}
}

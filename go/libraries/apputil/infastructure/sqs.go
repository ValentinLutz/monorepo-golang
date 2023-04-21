package infastructure

import (
	"context"
	"monorepo/libraries/apputil/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	*sqs.Client
	queueUrl string
	logger   logging.Logger
}

func NewSQSClient(config aws.Config, queueUrl string, logger logging.Logger) *SQSClient {
	return &SQSClient{
		Client:   sqs.NewFromConfig(config),
		queueUrl: queueUrl,
		logger:   logger,
	}
}

func (sqsClient *SQSClient) PollMessage(ctx context.Context, channel chan<- *types.Message) {
	receiveMessageOutput, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &sqsClient.queueUrl,
		AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
		WaitTimeSeconds:     20,
		MaxNumberOfMessages: 10,
	})
	if err != nil {
		sqsClient.logger.Error().
			Err(err).
			Msgf("failed to receive message from '%v'", sqsClient.queueUrl)
	}

	for _, message := range receiveMessageOutput.Messages {
		channel <- &message
	}
}

func (sqsClient *SQSClient) DeleteMessage(ctx context.Context, message types.Message) {
	_, err := sqsClient.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &sqsClient.queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		sqsClient.logger.Error().
			Err(err).
			Msgf("failed to delete message from '%v'", sqsClient.queueUrl)
	}
}

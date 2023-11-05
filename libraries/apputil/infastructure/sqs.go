package infastructure

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSConfig struct {
	MaxNumberOfMessages int
	QueueURL            string
	WaitTime            time.Duration
}

type SQSClient struct {
	*sqs.Client
	config SQSConfig
	wg     sync.WaitGroup
}

type HandleFunc func(message types.Message) error

func NewSQSClient(awsConfig aws.Config, config SQSConfig) *SQSClient {
	return &SQSClient{
		Client: sqs.NewFromConfig(awsConfig),
		config: config,
	}
}

func (sqsClient *SQSClient) Start(handleFunc HandleFunc) {
	slog.Info(
		"starting sqs polling",
		slog.String("queue_url", sqsClient.config.QueueURL),
	)

	for {
		err := sqsClient.PollMessages(handleFunc)
		if err != nil {
			slog.Error(
				"failed to receive messages",
				slog.Any("err", err),
				slog.String("queue_url", sqsClient.config.QueueURL),
			)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (sqsClient *SQSClient) Stop() {
	slog.Info(
		"stopping sqs polling",
		slog.String("queue_url", sqsClient.config.QueueURL),
	)

	//wait for all messages to be processed
	sqsClient.wg.Wait()

	slog.Info(
		"stopped sqs polling",
		slog.String("queue_url", sqsClient.config.QueueURL),
	)
}

func (sqsClient *SQSClient) PollMessages(handleFunc HandleFunc) error {
	slog.Info("polling messages", slog.String("queue_url", sqsClient.config.QueueURL))
	receiveMessageOutput, err := sqsClient.ReceiveMessage(
		context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            &sqsClient.config.QueueURL,
			AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
			WaitTimeSeconds:     20,
			MaxNumberOfMessages: 10,
		},
	)
	if err != nil {
		return err
	}

	if len(receiveMessageOutput.Messages) != 0 {
		sqsClient.handleMessages(handleFunc, receiveMessageOutput.Messages)
	}
	return nil
}

func (sqsClient *SQSClient) handleMessages(handleFunc HandleFunc, messages []types.Message) {
	messageCount := len(messages)
	sqsClient.wg.Add(messageCount)

	for i := range messages {
		go func(message types.Message) {
			defer sqsClient.wg.Done()

			err := handleFunc(message)
			if err != nil {
				slog.Error(
					"failed to handle message",
					slog.Any("err", err),
					slog.String("queue_url", sqsClient.config.QueueURL),
				)
				return
			}
			sqsClient.deleteMessage(message)
		}(messages[i])
	}

	sqsClient.wg.Wait()
}

func (sqsClient *SQSClient) deleteMessage(message types.Message) {
	_, err := sqsClient.DeleteMessage(
		context.Background(), &sqs.DeleteMessageInput{
			QueueUrl:      &sqsClient.config.QueueURL,
			ReceiptHandle: message.ReceiptHandle,
		},
	)
	if err != nil {
		slog.Error(
			"failed to delete message",
			slog.Any("err", err),
			slog.String("queue_url", sqsClient.config.QueueURL),
			slog.String("message_id", *message.MessageId),
			slog.String("receipt_handle", *message.ReceiptHandle),
		)
	}
}

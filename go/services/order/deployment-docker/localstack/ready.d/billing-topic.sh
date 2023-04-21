#!/bin/bash

# Create payment-status SNS topic
awslocal sns create-topic \
  --region eu-central-1 \
  --name payment-status

# Create payment-status SQS queue
awslocal sqs create-queue \
  --region eu-central-1 \
  --queue-name order-service-payment-status \
  --attributes '{
    "RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:eu-central-1:000000000000:order-service-payment-status-dlq\",\"maxReceiveCount\":2}"
  }'

# Get payment-status SNS queue attributes
# awslocal sqs get-queue-attributes \
#   --attribute-name All \
#   --queue-url http://localhost:4566/000000000000/order-service-payment-status

# Create payment-status dead-letter SQS queue
awslocal sqs create-queue \
  --region eu-central-1 \
  --queue-name order-service-payment-status-dlq

# Get payment-status dead-letter SNS queue attributes
# awslocal sqs get-queue-attributes \
#   --attribute-name All \
#   --queue-url http://localhost:4566/000000000000/order-service-payment-status-dlq

# Subscribe payment-status SQS queue to payment-status SNS topic
awslocal sns subscribe \
  --region eu-central-1 \
  --protocol sqs \
  --topic-arn arn:aws:sns:eu-central-1:000000000000:payment-status \
  --notification-endpoint arn:aws:sqs:eu-central-1:000000000000:order-service-payment-status \
  --attributes RawMessageDelivery=true

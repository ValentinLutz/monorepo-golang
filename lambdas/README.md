# lambdas

For testing lambda functions locally, you need to install
the [AWS Lambda Runtime Interface Emulator](https://github.com/aws/aws-lambda-runtime-interface-emulator).

To trigger the lambda function, you can use the following command:
```bash
curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
```

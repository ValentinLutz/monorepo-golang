package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type fact struct {
	Text string `json:"text"`
}

type factEntity struct {
	//ID   int    `db:"fact_id"`
	Text string `db:"fact_text"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db := getDatabaseConnection()

	var factEntities []factEntity
	err := db.SelectContext(context.Background(), &factEntities, "SELECT fact_text FROM public.fact")
	if err != nil {
		logger.Error(
			"failed to select facts",
			slog.Any("err", err),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	rand.Intn(len(factEntities))
	randomFactEntity := factEntities[rand.Intn(len(factEntities))]
	randomFactBody, err := json.Marshal(fact{Text: randomFactEntity.Text})
	if err != nil {
		logger.Error(
			"failed to marshal random fact",
			slog.Any("err", err),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(randomFactBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

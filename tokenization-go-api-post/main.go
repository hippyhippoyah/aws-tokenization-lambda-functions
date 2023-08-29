package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type InputData struct {
	Data string `json:"data"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var inputData InputData
	if err := json.Unmarshal([]byte(event.Body), &inputData); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, err
	}

	// Initialize a new session using the default credentials and AWS Region
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal server error"}, err
	}

	// Create a new DynamoDB client
	dbSvc := dynamodb.New(sess)

	//UUID
	id := uuid.New()
	fmt.Println(id.String())
	var ref = id.String()

	// Define the input for the PutItem operation
	item := map[string]*dynamodb.AttributeValue{
		"reference": {
			S: aws.String(ref),
		},
		"data": {
			S: &inputData.Data,
		},
	}

	// Specify the name of the DynamoDB table
	tableName := "goDataTest"

	// Create the PutItem input
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &tableName,
	}

	// Perform the PutItem operation
	_, putErr := dbSvc.PutItem(input)
	if putErr != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to save data"}, putErr
	}

	responseBody := fmt.Sprintf("Reference Token: %s", ref)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: responseBody}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

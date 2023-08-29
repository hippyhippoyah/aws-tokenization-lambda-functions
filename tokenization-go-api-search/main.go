package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize an AWS session
	searchKey := event.QueryStringParameters["ref"]
	if searchKey == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing 'ref' parameter",
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Change this to your desired AWS region
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create a DynamoDB client
	svc := dynamodb.New(sess)
	fmt.Println("Connect DynamoDB")

	// Specify the table name
	tableName := "goDataTest"

	// Read data from DynamoDB
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"reference": {
				S: aws.String(searchKey),
			},
			// Add additional keys if you have a composite key
		},
	})
	fmt.Println("Item Search")

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Process the result and construct the response
	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Item not found",
		}, nil
	}

	// Assuming you have an attribute called "data" in your DynamoDB table
	itemData := *result.Item["data"].S

	fmt.Println("Failed")

	// You can do further processing here based on your use case

	// Return a successful response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       itemData,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

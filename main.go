package main

/*
https://api.backendless.com/<application-id>/<api-key>/<operation-specific-path>

*/

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response message
type Response struct {
	Message string `json:"msg"`
}

// Handler for output
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	r := Response{
		Message: "Hello from golang function",
	}
	rbytes, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(rbytes), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

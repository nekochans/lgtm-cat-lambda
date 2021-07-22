package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.S3Event) error {

	fmt.Println(event.Records[0].S3.Object)
	return nil
}

func main() {
	lambda.Start(Handler)
}

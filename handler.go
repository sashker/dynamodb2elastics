/*

	This file is a part of the dynamodb2elastics project.
	Licenced under Apache License 2.0
	https://www.apache.org/licenses/LICENSE-2.0

	Author: Alexander Mustafin <sashker@inbox.ru>
*/

package main

import (
	"context" //"encoding/base64"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
	//"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-lambda-go/events"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	//Global request identificator
	requestID string
)

//Main function which receives requests from AWS API Gateway and does a number of tasks
func handleRequest(ctx context.Context, event events.DynamoDBEvent) (err error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		requestID = "UNDEFINED"
	}
	requestID = lc.AwsRequestID

	//Define the request logger with common fields
	reqLogger := log.WithFields(logrus.Fields{"request_id": requestID, "function": "handleRequest"})
	reqLogger.Debug("Raw request: %#v", event)

	// Create an Elasticsearch client
	esClient, err := newESClient(awsRegion)
	if err != nil {
		log.Fatalf("Can't create client for ElasticSearch: ", err)
	}

	//New ES index if it doesn't exist
	err = createESIndex(ctx, esClient, esIndex)
	if err != nil {
		log.Fatalf("Something went wrong during index checking: %s", err)
	}

	for _, r := range event.Records {
		data := map[string]interface{}{}
		for k, v := range r.Change.NewImage {
			reqLogger.Debugf("key %s type %d", k, v.DataType())
			switch v.DataType() {
			case 5:
				data[k], err = v.Integer()
				if err != nil {
					reqLogger.Error("Can't convert value to integer")
					return err
				}
			case 7:
				data[k] = "NULL"
			case 8:
				data[k] = v.String()
			default:
				data[k] = v.String()
			}
		}

		var id string
		if id, ok = data[recID].(string); !ok {
			return errors.New("can't get the value for given ID")
		}

		reqLogger.Debugf("ES document id: %s", id)
		status, err := createESDocument(ctx, esClient, esIndex, id, data)
		if err != nil {
			reqLogger.Error(err)
			return err
		}
		reqLogger.Info(status.Result)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

/*

	This file is a part of the dynamodb2elastics project.
	Licenced under Apache License 2.0
	https://www.apache.org/licenses/LICENSE-2.0

	Author: Alexander Mustafin <sashker@inbox.ru>
*/

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log                       = logrus.New()
	cred                      *credentials.Credentials
	awsRegion, esURL, esIndex string
	ctx                       = context.Background()
	recID                     string
)

func init() {
	//Please keep in mind that DynamoDB and ES must be placed in the same region
	val, ok := os.LookupEnv("REGION")
	if !ok {
		//AWS Lambda exposes AWS_REGION automatically
		val, ok = os.LookupEnv("AWS_REGION")
		if ok && val != "" {
			awsRegion = val
		} else {
			log.Fatal("REGION or AWS_REGION variables are empty")
		}
	} else if val == "" {
		log.Fatal("REGION variable is empty")
	}
	awsRegion = val

	val, ok = os.LookupEnv("ES_URL")
	if !ok {
		log.Fatal("ES_URL variable is not set")
	} else if val == "" {
		log.Fatal("ES_URL variable is empty")
	}
	esURL = val

	val, ok = os.LookupEnv("ES_INDEX")
	if !ok {
		log.Fatal("ES_INDEX variable is not set")
	} else if val == "" {
		log.Fatal("ES_INDEX variable is empty")
	}
	esIndex = val

	val, ok = os.LookupEnv("RECORD_ID")
	if !ok {
		log.Fatal("RECORD_ID variable is not set")
	} else if val == "" {
		log.Fatal("RECORD_ID variable is empty")
	}
	recID = val

	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "DEVELOPMENT"
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.TraceLevel)
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.ErrorLevel)
	}

	//We use AWS credentials stored as environment variables
	cred = credentials.NewEnvCredentials()
	if _, err := cred.Get(); err != nil {
		log.Fatal(err)
	}

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
}

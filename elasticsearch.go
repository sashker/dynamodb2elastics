/*

	This file is a part of the dynamodb2elastics project.
	Licenced under Apache License 2.0
	https://www.apache.org/licenses/LICENSE-2.0

	Author: Alexander Mustafin <sashker@inbox.ru>
*/

package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws4 "github.com/olivere/elastic/aws/v4"
	"github.com/olivere/elastic/v7"
)

func newESClient(awsRegion string) (client *elastic.Client, err error) {
	signingClient := aws4.NewV4SigningClient(credentials.NewEnvCredentials(), awsRegion)

	// Create an Elasticsearch client
	client, err = elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.Errorf("Can't create client for ElasticSearch: ", err)
		return nil, err
	}

	return client, nil
}

func createESIndex(ctx context.Context, client *elastic.Client, indexName string) error {
	// Check if the index exists
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		log.Infof("ES index %s doesn't exist\n", indexName)

		createIndex, err := client.CreateIndex(indexName).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			log.Warning("ES index creation is not acknowledged")
		}
		log.Infof("ES index %s is created", indexName)
	}
	return nil
}

func createESDocument(ctx context.Context, client *elastic.Client, index string, id string, entry map[string]interface{}) (status *elastic.UpdateResponse, err error) {
	if entry == nil {
		return nil, errors.New("entry's data is nil")
	}

	//status, err = client.Update().Index(index).Doc(entry).Upsert(entry).Do(ctx)
	status, err = client.Update().Index(index).Id(id).Doc(entry).Upsert(entry).Do(ctx)
	if err != nil {
		return nil, err
	}

	return status, nil
}

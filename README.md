# dynamodb2elastics

This quick'n'dirty implementation of the lambda function which ships records from AWS DynamoDB -> AWS ElasticSearch database

It requires the DynamoDB Stream to be activated for the database which you're going to sync.
The IAM role for the lambda-function must have permissions to write data to an ES server.

You also have to create a trigger in the DynamoDB database which fires the lambda-function when new records arrive
to a stream.

Please, pay attention that here we use AWS authentication to send data to the ES instance.
If you want to use a standalone ES server or the other cloud provider - feel free to refactor the code
and change this part.

## Usage
* Compile binary for the lambda function (Linux)
`GOOS=linux GOARCH=amd64 go build .`
* Make a zip archive of that binary file
`zip dynamodb2elastics.zip dynamodb2elastics`
* Create a new lambda-function with Go1.x environment and upload the zip-acrhive using the console or S3
* Fill all necessary environment variables

## Environment variables
* REGION - overwrites the default AWS_REGION variable, so you may have a lambda-function launched in the different region
* ES_URL - https endpoint of an ElasticSearch server
* ES_INDEX - the index to which we send data
* RECORD_ID - an identificator of each record in a DynamoDB database, which we use as id of each record in the ES index
* STAGE - enables you to have different kinds of logging for different environments(e.g. dev and production)

## Considerations
* DynamoDB database and an ES server must be placed in the same region
* New index is created dynamically, meaning that maybe you have to create it manually first with all necessary fields and their types
* We use the Upsert ES method which updates records if something changes or creates new entries
* All kinds of event will be sent to an ES index (INSERT,MODIFY,REMOVE)
* !!!It doesn't convert all DynamoDB types to the proper ES representations!!!

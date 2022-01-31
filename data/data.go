package data

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crseat/patient_feedback/errs"
	"github.com/crseat/patient_feedback/logger"
	"github.com/m7shapan/njson"
	"io/ioutil"
	"time"
)

// Patient defines information we need about the patient.
type Patient struct {
	Name      string `njson:"entry.0.resource.name.0.given.0"`
	Doctor    string `njson:"entry.1.resource.name.0.family"`
	Diagnosis string `njson:"entry.3.resource.code.coding.0.name"`
	PatientId string `njson:"entry.0.resource.id"`
}

// Response defines the answers we will save from user.
type Response struct {
	Diagnosis        string
	PatientId        string
	RecommendNumber  int
	ExplainedWell    string
	DiagnosisFeeling string
}

var svc *dynamodb.DynamoDB

func GetPatient() (*Patient, *errs.AppError) {
	raw, err := ioutil.ReadFile("./data/patient-feedback-raw-data.json")
	if err != nil {
		appError := errs.NewAppError("Error while reading JSON file: " + err.Error())
		logger.ErrorLogger.Println(appError.Message)
	}
	patient := Patient{}
	err = njson.Unmarshal(raw, &patient)
	if err != nil {
		appError := errs.NewAppError("Error while unmarshalling JSON file: " + err.Error())
		logger.ErrorLogger.Println(appError.Message)
		return nil, appError
	}
	return &patient, nil
}

func SaveResponse(response Response) *errs.AppError {
	// Add each item to PatientData table:
	tableName := "Responses"

	av, err := dynamodbattribute.MarshalMap(response)
	itemInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(itemInput)
	if err != nil {
		appError := errs.NewAppError("Error while putting item: " + err.Error())
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}
	return nil

}

func OpenDatabase() *errs.AppError {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	var sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)

	return nil
}

func CreateTable() *errs.AppError {

	// Create table Movies
	tableName := "Responses"

	//Ensure table doesn't exist and then create it
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "ResourceNotFoundException" {
			//Define table attributes and key schema
			input := &dynamodb.CreateTableInput{
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("PatientId"),
						AttributeType: aws.String("S"),
					},
					{
						AttributeName: aws.String("Diagnosis"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PatientId"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("Diagnosis"),
						KeyType:       aws.String("RANGE"),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
				TableName: aws.String(tableName),
			}

			_, err := svc.CreateTable(input)
			if err != nil {
				appError := errs.NewAppError("Error while creating table: " + err.Error())
				logger.ErrorLogger.Println(appError.Message)
				return appError
			}
			tableStatus := "CREATING"
			logger.InfoLogger.Println("Waiting for table to be active...")
			for tableStatus != "ACTIVE" {
				time.Sleep(1 * time.Second)
				table, _ := svc.DescribeTable(&dynamodb.DescribeTableInput{
					TableName: aws.String(tableName),
				})
				tableStatus = *table.Table.TableStatus
			}
			logger.InfoLogger.Println(tableName + " Table Active")

		}
	}
	return nil
}

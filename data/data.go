package data

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crseat/patient_feedback/errs"
	"github.com/crseat/patient_feedback/logger"
	"github.com/m7shapan/njson"
	"io/ioutil"
	"strconv"
)

// Patient defines information we need about the patient.
type Patient struct {
	Name      string `njson:"entry.0.resource.name.0.given.0"`
	Doctor    string `njson:"entry.1.resource.name.0.family"`
	Diagnosis string `njson:"entry.3.resource.code.coding.0.name"`
}

type Response struct {
	RecommendNumber  int
	ExplainedWell    string
	DiagnosisFeeling string
}

func GetItems() (*Patient, *errs.AppError) {
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
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Add each item to PatientData table:
	tableName := "Responses"

	// Create item in table Responses
	av, err := dynamodbattribute.MarshalMap(response)
	if err != nil {
		appError := errs.NewAppError("Error while marshalling response: " + err.Error())
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		appError := errs.NewAppError("Error while calling PutItem: " + err.Error())
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}
	recNum := strconv.Itoa(response.RecommendNumber)
	logger.InfoLogger.Println("Successfully added '" + recNum + response.ExplainedWell + response.DiagnosisFeeling + "to table: " + tableName)

	return nil

}

func createTable(svc *dynamodb.DynamoDB) *errs.AppError {
	// Create table Movies
	tableName := "Responses"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PatientId"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Diagnosis"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("RecommendNumber"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("ExplainedWell"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("DiagnosisFeeling"),
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
	logger.InfoLogger.Println("Created the table", tableName)
	return nil
}

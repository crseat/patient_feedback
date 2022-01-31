package data

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	OpenDatabase()
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("non_existant_table"),
	})
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() != "ResourceNotFoundException" {
			t.Fatalf("Database service not initalized properly.: " + awsErr.Code())
		}
	}
}

func TestCreateTable(t *testing.T) {
	OpenDatabase()
	err := CreateTable("test_responses")
	if err != nil {
		t.Fatalf("Error while creating the table: " + err.Message)
	}
	table, describeErr := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("test_responses"),
	})
	if describeErr != nil {
		DeleteTable("test_responses")
		t.Fatalf("Error while getting the table information: " + describeErr.Error())
	}
	if *table.Table.TableStatus != "ACTIVE" {
		DeleteTable("test_responses")
		t.Fatalf("test_responses dynamoDB table not active:")
	}
	DeleteTable("test_responses")
}

func TestDeleteTable(t *testing.T) {
	OpenDatabase()
	err := CreateTable("test_responses")
	if err != nil {
		t.Fatalf("Error while creating the test_responses table: " + err.Message)
	}
	err = DeleteTable("test_responses")
	if err != nil {
		t.Fatalf("Error while deleting the test_responses table: " + err.Message)
	}
	_, descEerr := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("test_responses"),
	})
	if awsErr, ok := descEerr.(awserr.Error); ok {
		if awsErr.Code() != "ResourceNotFoundException" {
			t.Fatalf("Table was not deleted!")
		}
	}
}

func TestGetPatient(t *testing.T) {
	patient, err := GetPatient("./patient-feedback-raw-data-test.json")
	if err != nil {
		t.Fatalf("Error while getting patient info: " + err.Message)
	}
	if patient.Name != "Test" || patient.PatientId != "9999-9999-999-999" || patient.Diagnosis != "Spontaneous Combustion" || patient.Doctor != "Schmo" {
		t.Fatalf("Patient information is incorrect.")
	}
}

func TestSaveResponse(t *testing.T) {
	OpenDatabase()
	tableName := "test_responses"
	CreateTable(tableName)
	response := Response{
		Diagnosis:        "testD",
		PatientId:        "testPid",
		RecommendNumber:  "testRec",
		ExplainedWell:    "testExp",
		DiagnosisFeeling: "testFel",
	}
	err := SaveResponse(response, tableName)
	if err != nil {
		DeleteTable("test_responses")
		t.Fatalf("Error saving response " + err.Message)
	}

	result, itemErr := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Diagnosis": {
				S: aws.String("testD"),
			},
			"PatientId": {
				S: aws.String("testPid"),
			},
		},
	})

	if itemErr != nil {
		DeleteTable("test_responses")
		t.Fatalf("Error retrieving item from dynamoDB " + itemErr.Error())
	}
	if result.Item == nil {
		DeleteTable("test_responses")
		t.Fatalf("Test item not in table!")
	}
	check_response := Response{}
	marshErr := dynamodbattribute.UnmarshalMap(result.Item, &check_response)
	if err != nil {
		DeleteTable("test_responses")
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", marshErr))
	}
	if check_response != response {
		DeleteTable("test_responses")
		t.Fatalf("Test item does not match expected!")
	}
	DeleteTable("test_responses")
}

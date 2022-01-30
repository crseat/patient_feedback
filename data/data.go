package data

import (
	"database/sql"
	"github.com/crseat/patient_feedback/errs"
	"github.com/crseat/patient_feedback/logger"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDb() error {
	var err error

	db, err = sql.Open("sqlite3", "./patientResponseDb.db")
	if err != nil {
		return err
	}
	return db.Ping()
}

// CreateTable creates the responses table in the sqlLite db so we can record responses.
func CreateTable() {
	createTbl := `CREATE TABLE IF NOT EXISTS responses (
        "recommendNumber" INTEGER,
        "explainedWell" TEXT,
        "diagnosisFeeling" TEXT
      );`

	// Prepare our statement
	stmt, err := db.Prepare(createTbl)
	if err != nil {
		appError := errs.NewAppError("Error while preparing create statement for responses table")
		logger.ErrorLogger.Println(appError)
	}

	// Execute our statement
	_, err = stmt.Exec()
	if err != nil {
		appError := errs.NewAppError("Error while executing create statement for responses table")
		logger.ErrorLogger.Println(appError)
	}
	logger.InfoLogger.Println("Responses table created")
}

//Package logger defines logging infrastructure and log destination.
package logger

//source https://www.honeybadger.io/blog/golang-logging/

import (
	"io"
	"log"
	"os"
)

var (
	DebugLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// init opens a file for writing logs. Also defines a MultiWriter so we can write to stdout and the file.
func init() {
	// Normally I would output the logs to the cloud, but for this exercise I'll create a file.
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// We want to output to stdout and the file at the same time so we'll use a MultiWriter
	mw := io.MultiWriter(os.Stdout, file)

	DebugLogger = log.New(mw, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

/* Example usage:
	InfoLogger.Println("Starting the application...")
	InfoLogger.Println("Something noteworthy happened")
	WarningLogger.Println("There is something you should know about")
	ErrorLogger.Println("Something went wrong")
}
*/

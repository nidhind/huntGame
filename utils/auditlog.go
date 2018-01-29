// auditlog
package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var audit *log.Logger

func InitAuditLog() {
	log.Println("Initializing audit logger...")

	logBasePath := os.Getenv("AUDIT_LOG_BASE_PATH")
	fileName := os.Getenv("AUDIT_LOG_FIlE_NAME")

	if logBasePath == "" {
		// set location of log file
		logBasePath = filepath.Join("/", "tmp")
		log.Println("WARNING: Audit log base path not set, using /tmp as base path")
	}

	if fileName == "" {
		// set name of log file
		fileName = "audit-answer.log"
	}
	logPath := filepath.Join(logBasePath, fileName)
	// If the file doesn't exist, create it, or append to the file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Cannot initialize audit logger...", err)
	}
	audit = log.New(file, "", 0)
	audit.Println("Audit Logger initialized at ", time.Now().UTC())
	log.Println("Audit Logger initialized...")
	log.Println("Audit Log File: ", logPath)
}

func AuditAnswer(u, a string, l int, t time.Time) {
	audit.Printf("user:= %s level:= %d answer:= %s time:= %s", u, l, a, t.UTC())
}

func GetGinLogFilePath() string {

	logBasePath := os.Getenv("GIN_LOG_BASE_PATH")
	fileName := os.Getenv("GIN_LOG_FIlE_NAME")

	if logBasePath == "" {
		// set location of log file
		logBasePath = filepath.Join("/", "tmp")
		log.Println("WARNING: Audit log base path not set, using /tmp as base path")
	}

	if fileName == "" {
		// set name of log file
		fileName = "gin.log"
	}

	logPath := filepath.Join(logBasePath, fileName)
	return logPath
}

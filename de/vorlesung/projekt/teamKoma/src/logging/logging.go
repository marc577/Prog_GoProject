package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace       *log.Logger
	Info        *log.Logger
	Warning     *log.Logger
	Error       *log.Logger
	traceFile   *os.File
	infoFile    *os.File
	warningFile *os.File
	errorFile   *os.File
)

func initQueues(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) bool {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return true
}

func LogInit(logLoc string) bool {
	initQueues(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	createDirIfNotExist(logLoc)

	traceFile, traceErr := createLogfileIfNotExist(logLoc, "trace")
	if traceErr != nil {
		Error.Fatal("Could not create Logfile", traceErr)
		return false
	}
	traceWriter := io.MultiWriter(ioutil.Discard, traceFile)
	Trace.SetOutput(traceWriter)

	infoFile, infoErr := createLogfileIfNotExist(logLoc, "info")
	if infoErr != nil {
		Error.Fatal("Could not create Logfile", infoErr)
		return false
	}
	infoWriter := io.MultiWriter(os.Stdout, infoFile)
	Info.SetOutput(infoWriter)

	warningFile, warnErr := createLogfileIfNotExist(logLoc, "warning")
	if warnErr != nil {
		Error.Fatal("Could not create Logfile", warnErr)
		return false
	}
	warningWriter := io.MultiWriter(os.Stdout, warningFile)
	Warning.SetOutput(warningWriter)

	errorFile, err := createLogfileIfNotExist(logLoc, "error")
	if err != nil {
		Error.Fatal("Could not create Logfile", err)
		return false
	}
	errorWriter := io.MultiWriter(os.Stderr, errorFile)
	Error.SetOutput(errorWriter)

	Trace.Println("NEW TRACE LOG")
	Info.Println("NEW INFO LOG")
	Warning.Println("NEW WARNINGS LOG")
	Error.Println("NEW ERROR LOG")
	return true
}

func createDirIfNotExist(dir string) (success bool) {
	success = false
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			success = false
			Error.Panic("Could not create LogFolder: ", err)
			return success
		}
		success = true
	}

	return success

}

func createLogfileIfNotExist(dir string, file string) (*os.File, error) {

	logFile, fileErr := os.OpenFile(strings.Join([]string{dir, "/", file, ".log"}, ""),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		Error.Fatal("Could not create LogFile:", fileErr)
	}

	return logFile, fileErr
}

func closeLogFile(file *os.File) (success bool) {
	success = true
	err := file.Close()
	if err != nil {
		success = false
		Error.Println("Could not close file:", err)
		return success
	}
	return success
}

func ShutdownLogging() bool {
	if traceFile != nil && infoFile != nil && warningFile != nil && errorFile != nil {
		Trace.Println("Gracefully shutting down TraceLog")
		state := closeLogFile(traceFile)
		Info.Println("Gracefully shutting down InfoLog")
		state = closeLogFile(infoFile)
		Warning.Println("Gracefully shutting down WarningLog")
		state = closeLogFile(warningFile)
		Error.Println("Gracefully shutting down ErrorLog")
		state = closeLogFile(errorFile)
		return state
	}
	return true

}

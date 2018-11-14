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

func initQueues(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInit(logLoc string) {
	initQueues(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	createDirIfNotExist(logLoc)

	traceFile = createLogfileIfNotExist("trace")
	Trace.SetOutput(traceFile)
	infoFile = createLogfileIfNotExist("info")
	Info.SetOutput(infoFile)
	warningFile = createLogfileIfNotExist("warning")
	Warning.SetOutput(warningFile)
	errorFile = createLogfileIfNotExist("error")
	Error.SetOutput(errorFile)

	Trace.Println("NEW TRACE LOG")
	Info.Println("NEW INFO LOG")
	Warning.Println("NEW WARNINGS LOG")
	Error.Println("NEW ERROR LOG")
}

func createDirIfNotExist(dir string) (success bool) {
	success = true
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			success = false
			Error.Panic("Could not create LogFolder: ", err)
			return success
		}
	}
	return success
}

func createLogfileIfNotExist(file string) (logFile *os.File) {

	logFile, fileErr := os.OpenFile(strings.Join([]string{"../../log/", file, ".log"}, ""),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		Error.Fatal("Could not create LogFile:", fileErr)
	}

	return logFile
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

func ShutdownLogging() {
	Trace.Println("Gracefully shutting down TraceLog")
	closeLogFile(traceFile)
	Info.Println("Gracefully shutting down InfoLog")
	closeLogFile(infoFile)
	Warning.Println("Gracefully shutting down WarningLog")
	closeLogFile(warningFile)
	Error.Println("Gracefully shutting down ErrorLog")
	closeLogFile(errorFile)
}

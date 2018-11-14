package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
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

func LogInit() {
	initQueues(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	createDirIfNotExist("../../log")

	Trace.SetOutput(createLogfileIfNotExist("trace"))
	infoFile := createLogfileIfNotExist("info")
	Info.SetOutput(infoFile)
	Warning.SetOutput(createLogfileIfNotExist("warning"))
	Error.SetOutput(createLogfileIfNotExist("error"))

	Trace.Println("NEW TRACE LOG")
	Info.Println("NEW INFO LOG")
	Warning.Println("NEW WARNINGS LOG")
	Error.Println("NEW ERROR LOG")
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func createLogfileIfNotExist(file string) (logFile *os.File) {

	logFile, fileErr := os.OpenFile(joinStr("../../log/", file, ".log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		Error.Println(fileErr)
	}
	//defer logFile.Close()

	return logFile
}

func joinStr(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

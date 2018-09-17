package golog

import (
	"fmt"
	"io"
	"log"
	"path"
	"runtime"
	"strconv"
	"time"
)

func init() {

}

//Golog represents a golog instance with all the necessary options.
//LogLeve: the level of the log [0=off 1=Error, 2=Error&Info, 3=all]. Defaults to 3
//infoPrefix: the prefix used for info logs, defaults to INFO
//debugPrefix: the prefix used for debug logs, defaults to DEBUG
//errorPrefix: the prefix used for error logs, defaults to ERROR
//showTimestmp: show timestamp in the logs; defults to true
//showPrefix: show prefix in the logs; defaults to true
//showCallerInfo: show info about the caller in the logs; defaults to true
//out: output destination for the logs; defaults to stdout
//gologer: the logger instance
type Golog struct {
	LogLevel       int
	InfoPrefix     string
	DebugPrefix    string
	ErrorPrefix    string
	ShowTimestamp  bool
	ShowPrefix     bool
	ShowCallerInfo bool
	Out            io.Writer
	Gologger       *log.Logger
	InfoLogger     *log.Logger
	DebugLogger    *log.Logger
	ErrorLogger    *log.Logger
}

//New initializes a new Golog instance
func New(output io.Writer) *Golog {

	//Create a default logger having as output destination the writer passed in the constructor.
	//All loggers will use this writer unless it is excplicitly overwrited by set*Output() functions
	defaultLogger := log.New(output, "", 0)
	infoLogger := log.New(output, "", 0)
	debugLogger := log.New(output, "", 0)
	errorLogger := log.New(output, "", 0)

	return &Golog{
		InfoPrefix:     "INFO",
		DebugPrefix:    "DEBUG",
		ErrorPrefix:    "ERROR",
		ShowTimestamp:  true,
		ShowPrefix:     true,
		ShowCallerInfo: true,
		Out:            output,
		Gologger:       defaultLogger,
		InfoLogger:     infoLogger,
		ErrorLogger:    errorLogger,
		DebugLogger:    debugLogger,
		LogLevel:       3,
	}
}

//getCallerInfo returns the info about the function calling golog
func getCallerInfo(skip int) string {
	var callerInfo string
	//var callingFuncName string

	_, fullFilePath, lineNumber, ok := runtime.Caller(skip)

	if ok {
		//callingFuncName = runtime.FuncForPC(pc).Name()

		// Split the path and use only the last 2 elements (package and file name)
		dirPath, fileName := path.Split(fullFilePath)
		var moduleName string
		if dirPath != "" {
			dirPath = dirPath[:len(dirPath)-1]
			_, moduleName = path.Split(dirPath)
		}
		callerInfo = moduleName + "/" + fileName + " " + strconv.Itoa(lineNumber)
	}

	return callerInfo
}

//Println simply calls fmt.Println
func (g *Golog) Println(v ...interface{}) {

	g.Gologger.Println(v...)
}

func (g *Golog) buildPrefix(prefixType string) string {
	//init prefix values
	prefix := ""
	timestamp := ""
	callerInfo := ""

	if g.ShowPrefix {
		if prefixType == "info" {
			prefix = g.InfoPrefix
		}
		if prefixType == "debug" {
			prefix = g.DebugPrefix
		}
		if prefixType == "error" {
			prefix = g.ErrorPrefix
		}

	}
	if g.ShowTimestamp {
		timestamp = time.Now().Format("02-01-2006 15:04:05")
	}

	if g.ShowCallerInfo {
		callerInfo = getCallerInfo(3)
	}

	//build prefix
	prefix = fmt.Sprintf("%s %s %s => ", "[ "+timestamp+" ]", prefix, callerInfo)

	return prefix
}

//Info writes info messages to the established output
func (g *Golog) Info(format string, v ...interface{}) {

	//do not print Info logs  if log level is not 2 or 3
	if g.LogLevel != 2 && g.LogLevel != 3 {
		return
	}

	//build prefix
	prefix := g.buildPrefix("info")

	g.InfoLogger.SetPrefix(prefix)

	g.InfoLogger.Printf(format, v...)
}

//Error writes error messages to the established output
func (g *Golog) Error(format string, v ...interface{}) {

	//do not print errors only if log level = 0
	if g.LogLevel == 0 {
		return
	}

	//build prefix
	prefix := g.buildPrefix("error")

	g.ErrorLogger.SetPrefix(prefix)

	g.ErrorLogger.Printf(format, v...)
}

//Debug writes debug messages to the established output
func (g *Golog) Debug(format string, v ...interface{}) {

	//do not print Debug logs if level != 3
	if g.LogLevel != 3 {
		return
	}

	//build prefix
	prefix := g.buildPrefix("debug")

	g.DebugLogger.SetPrefix(prefix)

	g.DebugLogger.Printf(format, v...)
}

//SetErrorPrefix updates the prefix used for error logs
func (g *Golog) SetErrorPrefix(prefix string) {
	g.ErrorPrefix = prefix
}

//SetErrorOutput updates the  destination output for error logs
func (g *Golog) SetErrorOutput(out io.Writer) {
	g.ErrorLogger.SetOutput(out)
}

//SetInfoPrefix updates the prefix used for info logs
func (g *Golog) SetInfoPrefix(prefix string) {
	g.InfoPrefix = prefix
}

//SetInfoOutput updates the  destination output for info logs
func (g *Golog) SetInfoOutput(out io.Writer) {
	g.InfoLogger.SetOutput(out)
}

//SetOutput updates the  destination output for all logs
func (g *Golog) SetOutput(out io.Writer) {
	g.Gologger.SetOutput(out)
}

//SetDebugPrefix updates the prefix used for debug logs
func (g *Golog) SetDebugPrefix(prefix string) {
	g.DebugPrefix = prefix
}

//SetDebugOutput updates the  destination output for debug logs
func (g *Golog) SetDebugOutput(out io.Writer) {
	g.DebugLogger.SetOutput(out)
}

func (g *Golog) SetLogLevel(level int) {
	g.LogLevel = level
}

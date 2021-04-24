package golog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/paperclicks/golog/transporter"

	"github.com/paperclicks/golog/model"
)

var INFO, DEBUG, ERROR, EMERGENCY, ALERT, CRITICAL, NOTICE, WARNING int

func init() {
	EMERGENCY = 0
	ALERT = 1
	CRITICAL = 2
	ERROR = 3
	WARNING = 4
	NOTICE = 5
	INFO = 6
	DEBUG = 7
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

type CallerInfo struct {
	File   string
	Line   int
	Module string
}

func (c CallerInfo) String() string {

	return c.Module + "/" + c.File + " " + strconv.Itoa(c.Line)

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
		LogLevel:       DEBUG,
	}
}

//getCallerInfo returns the info about the function calling golog
func getCallerInfo(skip int) CallerInfo {
	var callerInfo CallerInfo
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
		callerInfo.Module = moduleName
		callerInfo.File = fileName
		callerInfo.Line = lineNumber

	}

	return callerInfo
}

//Println simply calls fmt.Println
func (g *Golog) Println(v ...interface{}) {

	g.Gologger.Println(v...)
}

func (g *Golog) buildPrefix(level int) string {
	//init prefix values
	prefix := ""
	timestamp := ""
	callerInfo := ""

	if g.ShowPrefix {

		switch level {
		case DEBUG:
			prefix = "DEBUG"
		case INFO:
			prefix = "INFO"
		case NOTICE:
			prefix = "NOTICE"
		case WARNING:
			prefix = "WARNING"
		case ERROR:
			prefix = "ERROR"
		case CRITICAL:
			prefix = "CRITICAL"
		case ALERT:
			prefix = "ALERT"
		case EMERGENCY:
			prefix = "EMERGENCY"
		}

	}
	if g.ShowTimestamp {
		timestamp = time.Now().Format("02-01-2006 15:04:05")
	}

	if g.ShowCallerInfo {
		callerInfo = getCallerInfo(3).String()
	}

	//build prefix
	prefix = fmt.Sprintf("%s %s %s => ", "[ "+timestamp+" ]", prefix, callerInfo)

	return prefix
}


func (g *Golog) Log(message interface{}, level int) {

	outType := reflect.ValueOf(g.Out).Type()
	messageType := reflect.ValueOf(message).Type()
	var m string

	//if the level of the current log is higher than the LogLevel, then do not print this log
	if level > g.LogLevel {
		return
	}

	//determine the type of message
	switch messageType {
	case reflect.TypeOf("s"):
		m = fmt.Sprintf("%s", message)

	case reflect.TypeOf(model.Greylog{}):

		ci := getCallerInfo(3)
		gl := reflect.ValueOf(message).Interface().(model.Greylog)
		gl.Level = level
		gl.CustomFields["file"] = ci.File
		gl.CustomFields["module"] = ci.Module
		gl.CustomFields["line"] = ci.Line
		m = gl.String()

	}

	//check the output type, and set eventual prefix
	switch outType {
	case reflect.TypeOf(os.Stdout):
		prefix := g.buildPrefix(level)

		g.Gologger.SetPrefix(prefix)

	case reflect.TypeOf(&transporter.FileTransporter{}):
		prefix := g.buildPrefix(level)

		g.Gologger.SetPrefix(prefix)
	case reflect.TypeOf(&transporter.AMQPTransporter{}):
		prefix := g.buildPrefix(level)
		g.Gologger.SetPrefix(prefix)

		g.Out.Write([]byte(m))
	}

	g.Gologger.Println(m)

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

package gologv2

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

	"github.com/paperclicks/golog/model"
	"github.com/paperclicks/golog/transporter"
)

type LogLevel int

type channelMessage struct {
	Message    string
	Prefix     string
	Level      LogLevel
	OutputType reflect.Type
}

var INFO, DEBUG, ERROR, EMERGENCY, ALERT, CRITICAL, NOTICE, WARNING LogLevel

var logger *log.Logger

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
	LogLevel       LogLevel
	InfoPrefix     string
	DebugPrefix    string
	ErrorPrefix    string
	ShowTimestamp  bool
	ShowPrefix     bool
	ShowCallerInfo bool
	Out            io.Writer
	Logger         *log.Logger
	Chann          chan channelMessage
	// Error          chan channelMessage
	// Info           chan channelMessage
	// Debug          chan channelMessage
}

type CallerInfo struct {
	File   string
	Line   int
	Module string
}

var infoChan, debugChan, errorChan, ch chan channelMessage
var stopChan chan struct{}
var golog *Golog

func (c CallerInfo) String() string {

	return c.Module + "/" + c.File + " " + strconv.Itoa(c.Line)

}

//New initializes a new Golog instance
func New(output io.Writer) *Golog {

	//Create a default logger having as output destination the writer passed in the constructor.
	//All loggers will use this writer unless it is excplicitly overwrited by set*Output() functions
	logger = log.New(output, "", 0)

	stopChan = make(chan struct{})

	ch = make(chan channelMessage, 10)

	// infoChan = make(chan channelMessage, 10)
	// debugChan = make(chan channelMessage, 10)
	// errorChan = make(chan channelMessage, 10)

	golog = &Golog{
		ShowTimestamp:  true,
		ShowPrefix:     true,
		ShowCallerInfo: true,
		Out:            output,
		Logger:         logger,
		LogLevel:       DEBUG,
		Chann:          ch,
		// Info:           infoChan,
		// Debug:          debugChan,
		// Error:          errorChan,
	}

	go golog.consume(stopChan)

	return golog
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

func (g *Golog) buildPrefix(level LogLevel) string {
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

func (g *Golog) Log(message interface{}, level LogLevel) {

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
		gl.Level = int(level)
		gl.CustomFields["_file"] = ci.File
		gl.CustomFields["_module"] = ci.Module
		gl.CustomFields["_line"] = ci.Line

		m = gl.String()
	}

	prefix := g.buildPrefix(level)

	//send message to chann
	chm := channelMessage{Message: m, Level: level, OutputType: outType, Prefix: prefix}
	g.Chann <- chm

}

func (g *Golog) consume(stopChan chan struct{}) {

	for {

		select {
		case m := <-g.Chann:
			outType := reflect.ValueOf(m.OutputType).Type()

			//determine the output destination for the log
			switch outType {
			case reflect.TypeOf(os.Stdout):

				g.Logger.SetPrefix(m.Prefix)

				g.Logger.Println(m.Message)

			case reflect.TypeOf(transporter.AMQPTransporter{}):

				g.Logger.SetPrefix("")

				g.Logger.Println(m.Message)
			}
		case <-stopChan:
			return
		}

	}

}

//SetErrorPrefix updates the prefix used for error logs
func (g *Golog) SetErrorPrefix(prefix string) {
	g.ErrorPrefix = prefix
}

//SetInfoPrefix updates the prefix used for info logs
func (g *Golog) SetInfoPrefix(prefix string) {
	g.InfoPrefix = prefix
}

//SetDebugPrefix updates the prefix used for debug logs
func (g *Golog) SetDebugPrefix(prefix string) {
	g.DebugPrefix = prefix
}

func (g *Golog) SetLogLevel(level LogLevel) {
	g.LogLevel = level
}

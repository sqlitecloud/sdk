//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud Web Server
//     ///             ///  ///         Version     : 0.2.0
//     //             ///   ///  ///    Date        : 2022/02/15
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description :
//   ////                ///  ///
//     ////     //////////   ///
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
)

type LogLevel int32

const (
	LogLevelPanic LogLevel = iota
	LogLevelFatal
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

const (
	calldepth = 2
)

// DefaultLogger is a default implementation of the Logger interface.
type Logger struct {
	*log.Logger
	LogLevel
}

var (
	StdErrLogger  = &Logger{Logger: log.New(os.Stderr, "", log.LstdFlags), LogLevel: LogLevelDebug}
	DiscardLogger = &Logger{Logger: log.New(ioutil.Discard, "", 0), LogLevel: LogLevelFatal}
)

func (l *Logger) Debug(v ...interface{}) {
	if l.LogLevel >= LogLevelDebug {
		l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.LogLevel >= LogLevelDebug {
		l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.LogLevel >= LogLevelInfo {
		l.Output(calldepth, header("INFO", fmt.Sprint(v...)))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.LogLevel >= LogLevelInfo {
		l.Output(calldepth, header("INFO", fmt.Sprintf(format, v...)))
	}
}

func (l *Logger) Warning(v ...interface{}) {
	if l.LogLevel >= LogLevelWarn {
		l.Output(calldepth, header("WARN", fmt.Sprint(v...)))
	}
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.LogLevel >= LogLevelWarn {
		l.Output(calldepth, header("WARN", fmt.Sprintf(format, v...)))
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.LogLevel >= LogLevelError {
		l.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.LogLevel >= LogLevelError {
		l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s %s: %s", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), lvl, msg)
}

// logger

func init() {
	initializeSQLiteWeb()
  }
  
func initLogger() {
	if ( len(SQLiteWeb.logFile)>0 ) {
		f, err := os.OpenFile(SQLiteWeb.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			SQLiteWeb.Logger.Fatal(err)
		}

		llevel := LogLevelDebug
		if ( len(SQLiteWeb.logLevel)>0 ) {
			switch strings.ToUpper(SQLiteWeb.logLevel) {
			case "PANIC": llevel = LogLevelPanic
			case "FATAL": llevel = LogLevelFatal
			case "ERROR": llevel = LogLevelError
			case "WARNING": llevel = LogLevelWarn
			case "INFO": llevel = LogLevelInfo
			case "DEBUG": llevel = LogLevelDebug
			}
		}

		SQLiteWeb.Logger = &Logger{ Logger: log.New(f, "", log.LstdFlags), LogLevel: llevel }
	}

	if ( len(SQLiteWeb.clfLogFile)>0 ) {
		f, err := os.OpenFile(SQLiteWeb.clfLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			SQLiteWeb.Logger.Fatal(err)
		}
		SQLiteWeb.CLFWriter = f
	}
}

// Web Server CommonLogFormat

func CommonLogFormatMiddleware(handler http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(SQLiteWeb.CLFWriter, handler)
}
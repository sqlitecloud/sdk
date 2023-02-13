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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
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
	// return fmt.Sprintf("%s %s: %s", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), lvl, msg)
	return fmt.Sprintf("%s: %s", lvl, msg)
}

// logger

func init() {
	initializeSQLiteWeb()
}

func initLogger() {
	if len(SQLiteWeb.logFile) > 0 {
		f, err := os.OpenFile(SQLiteWeb.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			SQLiteWeb.Logger.Fatal(err)
		}

		llevel := LogLevelDebug
		if len(SQLiteWeb.logLevel) > 0 {
			switch strings.ToUpper(SQLiteWeb.logLevel) {
			case "PANIC":
				llevel = LogLevelPanic
			case "FATAL":
				llevel = LogLevelFatal
			case "ERROR":
				llevel = LogLevelError
			case "WARNING":
				llevel = LogLevelWarn
			case "INFO":
				llevel = LogLevelInfo
			case "DEBUG":
				llevel = LogLevelDebug
			}
		}

		SQLiteWeb.Logger = &Logger{Logger: log.New(f, "", log.LstdFlags), LogLevel: llevel}
	}

	if len(SQLiteWeb.clfLogFile) > 0 {
		f, err := os.OpenFile(SQLiteWeb.clfLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			SQLiteWeb.Logger.Fatal(err)
		}
		SQLiteWeb.CLFWriter = f
	}
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := l.w.(http.Hijacker).Hijack()
	if err == nil && l.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		l.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}

func makeLogger(w http.ResponseWriter) (*responseLogger, http.ResponseWriter) {
	logger := &responseLogger{w: w, status: http.StatusOK}
	return logger, httpsnoop.Wrap(w, httpsnoop.Hooks{
		Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return logger.Write
		},
		WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return logger.WriteHeader
		},
	})
}

func LogReqWithDurationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger, w := makeLogger(w)

		next.ServeHTTP(w, r)

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		user := "-"
		userid, _ := SQLiteWeb.Auth.GetUserID(SQLiteWeb.Auth.getTokenFromAuthorization, r)
		if userid != -1 {
			user = strconv.Itoa(int(userid))
		}

		t := time.Since(startTime)
		SQLiteWeb.Logger.Infof("%s %s - \"%s %s\" %d %fms", host, user, r.Method, r.RequestURI, logger.status, float32(t.Nanoseconds())/1e6) // t.String()
	})
}

// Web Server CommonLogFormat

func CommonLogFormatMiddleware(handler http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(SQLiteWeb.CLFWriter, handler)
}

package logger

import (
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"strconv"
	"bytes"
	"time"
)

// Error logging level
const (
	LogLevelError   = iota
	LogLevelDefault
	LogLevelFromService
	LogLevelToService
)

type Logger struct {
	module string
}

// Log ... Append custom log identifier: [E!], [ ]
func (logger *Logger) Log(level int, s string) {
	t := time.Now()
	timeString := t.Format("2006/01/02 15:04:05.000")
	if level == LogLevelError {
		fmt.Println(timeString + " [E] " + logger.module + ": "  + s)
	} else if level == LogLevelDefault {
		fmt.Println(timeString + " [I] " + logger.module + ": "  + s)
	} else if level == LogLevelFromService {
		fmt.Println(timeString + " [<-] " + logger.module + ": "  + s)
	} else if level == LogLevelToService {
		fmt.Println(timeString + " [->] " + logger.module + ": "  + s)
	}
}

func New(module string) *Logger {
	return &Logger{module: module}
}

func ChildLogger(logger *Logger, module string) *Logger {
	return &Logger{module: logger.module + "->" + module}
}

// Logf ... Custom logging with format. Append custom log identifier: [E!], [ ]
func (logger *Logger) Logf(level int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	logger.Log(level, message)
}

// logResponse ... log response
func (logger *Logger) logResponse(res *http.Response, logBody bool, from bool) {
	headersString := ""
	for k, v := range res.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}

	message := fmt.Sprintf("Response for %s [%s]: Status: %s. Headers: %s.", res.Request.Method, res.Request.URL.String(), strconv.Itoa(res.StatusCode), headersString)
	if logBody && res.Body != nil {
		save, reader, err := readersFromReader(res.Body)
		if err != nil {
			logger.Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			logger.Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
		res.Body.Close()
		res.Body = ioutil.NopCloser(save)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if from {
			logger.Log(LogLevelFromService, message)
		} else {
			logger.Log(LogLevelToService, message)
		}
	} else {
		logger.Log(LogLevelError, message)
	}
}

func (logger *Logger) LogResponseToService(res *http.Response, logBody bool) {
	logger.logResponse(res, logBody, false)
}

func (logger *Logger) LogResponseFromService(res *http.Response, logBody bool) {
	logger.logResponse(res, logBody, true)
}

// logRequest ... log request
func (logger *Logger) logRequest(req *http.Request, logBody bool, from bool) {
	headersString := ""
	for k, v := range req.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}
	message := fmt.Sprintf("Request started: %s [%s]: Headers: %s.", req.Method, req.URL.String(), headersString)
	if logBody && req.Body != nil {
		save, reader, err := readersFromReader(req.Body)
		if err != nil {
			logger.Logf(LogLevelError, "Cant parse request %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			logger.Logf(LogLevelError, "Cant parse requesr %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
		req.Body.Close()
		req.Body = ioutil.NopCloser(save)
	}
	if from {
		logger.Log(LogLevelFromService, message)
	} else {
		logger.Log(LogLevelToService, message)
	}
}

func (logger *Logger) LogRequestToService(req *http.Request, logBody bool) {
	logger.logRequest(req, logBody, false)
}

func (logger *Logger) LogRequestFromService(req *http.Request, logBody bool) {
	logger.logRequest(req, logBody, true)
}


func readersFromReader(reader io.ReadCloser) (io.Reader, io.Reader, error) {
	b, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewBuffer(b), bytes.NewBuffer(b), nil
}

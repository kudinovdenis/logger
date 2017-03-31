package logger

import (
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"strconv"
	"bytes"
)

// Error logging level
const (
	LogLevelError   = iota
	LogLevelDefault
	LogLevelFromService
	LogLevelToService
)

// Log ... Append custom log identifier: [E!], [ ]
func Log(level int, s string) {
	if level == LogLevelError {
		fmt.Println("[E!] " + s)
	} else if level == LogLevelDefault {
		fmt.Println("[ ] " + s)
	} else if level == LogLevelFromService {
		fmt.Println("[<-] " + s)
	} else if level == LogLevelToService {
		fmt.Println("[->] " + s)
	}
}

// Logf ... Custom logging with format. Append custom log identifier: [E!], [ ]
func Logf(level int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Log(level, message)
}

// logResponse ... log response
func logResponse(res *http.Response, logBody bool, from bool) {
	headersString := ""
	for k, v := range res.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}

	message := fmt.Sprintf("Response for %s [%s]: Status: %s. Headers: %s.", res.Request.Method, res.Request.URL.String(), strconv.Itoa(res.StatusCode), headersString)
	if logBody && res.Body != nil {
		save, reader, err := readersFromReader(res.Body)
		if err != nil {
			Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
		res.Body.Close()
		res.Body = ioutil.NopCloser(save)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if from {
			Log(LogLevelFromService, message)
		} else {
			Log(LogLevelToService, message)
		}
	} else {
		Log(LogLevelError, message)
	}
}

func LogResponseToService(res *http.Response, logBody bool) {
	logResponse(res, logBody, false)
}

func LogResponseFromService(res *http.Response, logBody bool) {
	logResponse(res, logBody, true)
}

// logRequest ... log request
func logRequest(req *http.Request, logBody bool, from bool) {
	headersString := ""
	for k, v := range req.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}
	message := fmt.Sprintf("Request started: %s [%s]: Headers: %s.", req.Method, req.URL.String(), headersString)
	if logBody && req.Body != nil {
		save, reader, err := readersFromReader(req.Body)
		if err != nil {
			Logf(LogLevelError, "Cant parse request %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			Logf(LogLevelError, "Cant parse requesr %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
		req.Body.Close()
		req.Body = ioutil.NopCloser(save)
	}
	if from {
		Log(LogLevelFromService, message)
	} else {
		Log(LogLevelToService, message)
	}
}

func LogRequestToService(req *http.Request, logBody bool) {
	logRequest(req, logBody, false)
}

func LogRequestFromService(req *http.Request, logBody bool) {
	logRequest(req, logBody, true)
}


func readersFromReader(reader io.ReadCloser) (io.Reader, io.Reader, error) {
	b, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewBuffer(b), bytes.NewBuffer(b), nil
}

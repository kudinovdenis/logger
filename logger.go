package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"bytes"
)

// Error logging level
const (
	LogLevelError   = iota
	LogLevelDefault = iota
)

// Log ... Append custom log identifier: [E!], [ ]
func Log(level int, s string) {
	if level == LogLevelError {
		fmt.Println("[E!] " + s)
	} else if level == LogLevelDefault {
		fmt.Println("[ ] " + s)
	}
}

// Logf ... Custom logging with format. Append custom log identifier: [E!], [ ]
func Logf(level int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Log(level, message)
}

// LogResponse ... log response
func LogResponse(res *http.Response, logBody bool) {
	bodyString := ""
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	// Restore the io.ReadCloser to its original state
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	bodyString += string(bodyBytes)

	headersString := ""
	for k, v := range res.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}

	var message string
	if logBody && len(bodyString) > 0 {
		message = fmt.Sprintf("Response for [%s]: Status: %s. Headers: %s. Body: %s", res.Request.URL.String(), strconv.Itoa(res.StatusCode), headersString, bodyString)
	} else {
		message = fmt.Sprintf("Response for [%s]: Status: %s. Headers: %s.", res.Request.URL.String(), strconv.Itoa(res.StatusCode), headersString)
	}
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		Log(LogLevelDefault, message)
	} else {
		Log(LogLevelError, message)
	}
}

// LogRequest ... log request
func LogRequest(req *http.Request, logBody bool) {
	headersString := ""
	for k, v := range req.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}
	message := fmt.Sprintf("Request started: %s [%s]: Headers: %s.", req.Method, req.URL.String(), headersString)
	if logBody && req.Body != nil {
		reader, err := req.GetBody()
		if err != nil {
			Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
	}
	Log(LogLevelDefault, message)
}

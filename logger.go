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
	headersString := ""
	for k, v := range res.Header {
		headersString = headersString + fmt.Sprintf("[key:%s value:%s] ", k, v)
	}

	var message string
	save, reader, err := drainBody(res.Body)
	bodyBytes, _ := ioutil.ReadAll(reader)
	if logBody && len(bodyBytes) > 0 {
		if err != nil {
			Logf(LogLevelError, "Cant parse response %s.", err.Error())
		}
		bodyString := ""
		res.Body = save
		bodyString += string(bodyBytes)
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
		save, reader, err := drainBody(req.Body)
		if err != nil {
			Logf(LogLevelError, "Cant parse request %s.", err.Error())
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			Logf(LogLevelError, "Cant parse requesr %s.", err.Error())
		}
		bodyString := string(body[:])
		message += fmt.Sprintf(" Body: %s", bodyString)
		req.Body = save
	}
	Log(LogLevelDefault, message)
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

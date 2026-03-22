package logger

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetLogEntry(r *http.Request, status int, body []byte) string {
	var requestBody map[string]string

	err := json.Unmarshal(body, &requestBody)

	logEntry := APILogEntry{
		Method:      r.Method,
		URL:         r.URL.String(),
		IP:          r.RemoteAddr,
		Headers:     make(map[string]string),
		RequestBody: requestBody,
		Status:      status,
	}
	for key, values := range r.Header {
		logEntry.Headers[key] = values[0]
	}

	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(jsonBytes)
}

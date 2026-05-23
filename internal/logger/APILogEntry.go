package logger

type APILogEntry struct {
	Status      int
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	IP          string            `json:"ip"`
	Headers     map[string]string `json:"headers"`
	RequestBody map[string]string `json:"request_body"`
}

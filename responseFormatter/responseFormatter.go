package responseformatter

var FormatResponse = func(statusCode int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"data":    data,
	}
}

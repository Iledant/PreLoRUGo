package actions

// jsonError is used to embed JSON response for an error message
type jsonError struct {
	Error string `json:"error"`
}

// jsonMessage is used to embed JSON response for a message
type jsonMessage struct {
	Message string `json:"Message"`
}

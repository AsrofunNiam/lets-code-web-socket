package domain

// Message
type Message struct {
	Username  string      `json:"username"`
	Content   interface{} `json:"content"`
	Timestamp string      `json:"timestamp"`

	// Add other fields as needed
	Group bool
}

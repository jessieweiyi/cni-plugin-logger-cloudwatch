package logger

// Publisher interface for all log publishers
type Publisher interface {
	Publish([]byte) error
}

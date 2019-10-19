package logger

type Publisher interface {
	Publish([]byte) error
}

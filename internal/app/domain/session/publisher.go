package session

// Publisher Representation of a publisher interface
type Publisher interface {
	PublishResult(Result) error
}

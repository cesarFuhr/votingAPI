package vote

// DocValidator validates document numbers
type DocValidator interface {
	ValidateDocument(string) (bool, error)
}

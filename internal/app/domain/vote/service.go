package vote

// Service describes the agenda service interface
type Service interface {
	CreateVote(string, string, string) (Vote, error)
}

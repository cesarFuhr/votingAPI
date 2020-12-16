package agenda

// Service describes the agenda service interface
type Service interface {
	CreateAgenda(string) (Agenda, error)
	FindAgenda(string) (Agenda, error)
}

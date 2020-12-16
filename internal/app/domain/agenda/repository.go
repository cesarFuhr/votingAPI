package agenda

// Repository Persistency interface to serve the Agenda service
type Repository interface {
	FindAgenda(string) (Agenda, error)
	InsertAgenda(Agenda) error
}

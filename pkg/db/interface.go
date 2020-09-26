package db

// Interface defines the contractual methods for a db in the application
type Interface interface {
	Init()
	Connect() *Queries
}

package services

// Service contains details about provided Service.
// Execute runs the related action in the main thread.
// Next runs the next node(s)
type Service interface {
	Execute(attrs ...interface{})
	Next(attrs ...interface{})
}

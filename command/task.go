package command

// Single task interface.
type Task interface {
	// GetName gets a task name.
	GetName() string
	// Run runs a task.
	Run() error
}

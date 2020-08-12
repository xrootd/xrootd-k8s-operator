package reconciler

// StatusReconciler is called after reconciliation to update the status of CR instance
type StatusReconciler interface {
	UpdateStatus(instance resource) error
}

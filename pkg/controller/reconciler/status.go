package reconciler

type StatusReconciler interface {
	UpdateStatus(instance resource) error
}

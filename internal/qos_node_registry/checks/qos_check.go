package checks

type QosJob interface {
	PerformJob()
	ShouldRun() bool
}

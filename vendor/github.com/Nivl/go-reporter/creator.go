package reporter

//go:generate mockgen -destination implementations/mockreporter/creator.go -package mockreporter go de Creator

// Creator describes a struct used to create reporters
type Creator interface {
	New() (Reporter, error)
}

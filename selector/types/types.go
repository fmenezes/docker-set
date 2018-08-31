package types

type Driver interface {
	List() ([]EnvironmentEntry, error)
}

type EnvironmentEntry struct {
	Name   *string
	Active bool
	Source *string
	State  *string
	URL    *string
}

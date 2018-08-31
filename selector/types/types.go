package types

type Driver interface {
	Name() string
	Add(NewEnvironmentEntry) error
	List() ([]EnvironmentEntry, error)
}

type NewEnvironmentEntry struct {
	Name     string
	Driver   string
	Location string
}

type EnvironmentEntry struct {
	Name     string
	Driver   string
	Location string
	Active   bool
	State    *string
	URL      *string
}

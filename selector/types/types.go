package types

type Driver interface {
	Name() string
	Add(NewEnvironmentEntry) error
	Remove(NewEnvironmentEntry) error
	List() ([]EnvironmentEntry, error)
	Env(EnvironmentEntry) (map[string]*string, error)
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
	State    *string
}

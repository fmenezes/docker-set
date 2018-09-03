package common

type Driver interface {
	Name() string
	Add(EnvironmentEntry) error
	Remove(EnvironmentEntry) error
	List() ([]EnvironmentEntryWithState, error)
	Env(EnvironmentEntryWithState) (map[string]*string, error)
}

type EnvironmentEntry struct {
	Name     string
	Driver   string
	Location *string
}

type EnvironmentEntryWithState struct {
	EnvironmentEntry
	State *string
}

type Storage interface {
	Add(EnvironmentEntry) error
	Remove(EnvironmentEntry) error
	Load() ([]EnvironmentEntry, error)
	Save([]EnvironmentEntry) error
}

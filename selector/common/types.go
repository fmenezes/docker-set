// common package contains structs used in selector package
package common

// Holds all methods to manage environments
type Selector interface {
	// Adds driver to selector
	RegisterDriver(driver Driver)
	// Retrieves environment variables for given name,
	Env(entry string) (map[string]*string, error)
	// Retrieves a list of all environments,
	List() <-chan EnvironmentEntryWithState
	// Retrieves name of selected environment
	Selected() *string
	// Appends new environment into store
	Add(entry EnvironmentEntry) error
	// Removes from store the entry corresponding given name
	Remove(entry string) error
	// Starts the environment
	Start(entry string) error
	// Stops the environment
	Stop(entry string) error
}

// Interface defining a driver for docker environment
type Driver interface {
	// Returns the name of the driver
	Name() string
	// Starts the environment
	Start(EnvironmentEntry) error
	// Stops the environment
	Stop(EnvironmentEntry) error
	// Stores the envionment passed
	Add(EnvironmentEntry) error
	// Removes the envionment passed from the storage
	Remove(EnvironmentEntry) error
	// Lists all environments from this driver
	List() <-chan EnvironmentEntryWithState
	// Returns environment variables that need to be set
	Env(EnvironmentEntryWithState) (map[string]*string, error)
	// Checks if the driver is supported
	IsSupported() bool
}

// Struct defining an environment
type EnvironmentEntry struct {
	// Name of this environment
	Name string
	// Driver used in this environment
	Driver string
	// Location on where this environment is stored
	Location *string
}

// Extension of EnvironmentEntry struct to include state
type EnvironmentEntryWithState struct {
	// EnvironmentEntry struct
	EnvironmentEntry
	// State of this environment (e.g. running, stopped, etc...)
	State *string
}

// Interface defining how envionment entrys will be saved
type Storage interface {
	// Appends data into dataset
	Append(EnvironmentEntry) error
	// Removes data from dataset
	Remove(EnvironmentEntry) error
	// Retrieves data from dataset
	Load() ([]EnvironmentEntry, error)
	// Stores data in dataset
	Save([]EnvironmentEntry) error
}

// selector package contains logic to manage environments for docker
package selector

import "fmt"
import "sync"
import "os"

import "github.com/fmenezes/docker-set/selector/common"

// Holds all methods to manage environments
type Selector struct {
	drivers []common.Driver
}

// Returns a new instance of the selector,
// can fail on when trying to create storage.FileStorage
func NewSelector() *Selector {
	return &Selector{
		drivers: make([]common.Driver, 0),
	}
}

// Adds driver to selector
func (s *Selector) RegisterDriver(driver common.Driver) {
	s.drivers = append(s.drivers, driver)
}

func (s Selector) selectDriver(driver string) (common.Driver, error) {
	var selectedDriver *common.Driver = nil
	for _, item := range s.drivers {
		if item.Name() == driver {
			selectedDriver = &item
			break
		}
	}

	if selectedDriver == nil {
		return nil, fmt.Errorf("Driver '%s' is not supported", driver)
	}

	return *selectedDriver, nil
}

// Appends new environment into store,
// can fail when issuing driver add or adding the same name twice
func (s Selector) Add(entry common.EnvironmentEntry) error {
	selectedDriver, err := s.selectDriver(entry.Driver)
	if err != nil {
		return err
	}

	exists, err := s.existsEntry(entry.Name)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("'%s' already exists", entry.Name)
	}

	return selectedDriver.Add(entry)
}

// Starts the environment,
// can fail when issuing underlying driver command
func (s Selector) Start(entry string) error {
	found, err := s.findEntry(entry)
	if err != nil {
		return err
	}
	selectedDriver, err := s.selectDriver(found.Driver)
	if err != nil {
		return err
	}
	return selectedDriver.Start(found.EnvironmentEntry)
}

// Stops the environment,
// can fail when issuing underlying driver command
func (s Selector) Stop(entry string) error {
	found, err := s.findEntry(entry)
	if err != nil {
		return err
	}
	selectedDriver, err := s.selectDriver(found.Driver)
	if err != nil {
		return err
	}
	return selectedDriver.Stop(found.EnvironmentEntry)
}

func (s Selector) existsEntry(name string) (bool, error) {
	for item := range s.List() {
		if item.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (s Selector) findEntry(name string) (*common.EnvironmentEntryWithState, error) {
	for item := range s.List() {
		if item.Name == name {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("'%s' not found", name)
}

// Retrieves name of selected environment
func (s Selector) Selected() *string {
	var result *string = nil
	val, ok := os.LookupEnv("DOCKER_SET_MACHINE")
	if ok {
		result = &val
	}
	return result
}

// Retrieves environment variables for given name,
// can fail when issuing driver env
func (s Selector) Env(entry string) (map[string]*string, error) {
	found, err := s.findEntry(entry)
	if err != nil {
		return nil, err
	}
	selectedDriver, err := s.selectDriver(found.Driver)
	if err != nil {
		return nil, err
	}
	return selectedDriver.Env(*found)
}

// Removes from store the entry corresponding given name
// can fail when issuing driver remove
func (s Selector) Remove(entry string) error {
	found, err := s.findEntry(entry)
	if err != nil {
		return err
	}
	selectedDriver, err := s.selectDriver(found.Driver)
	if err != nil {
		return err
	}
	return selectedDriver.Remove(found.EnvironmentEntry)
}

func mergeChans(in ...<-chan common.EnvironmentEntryWithState) <-chan common.EnvironmentEntryWithState {
	out := make(chan common.EnvironmentEntryWithState)
	var wg sync.WaitGroup
	wg.Add(len(in))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, c := range in {
		go func(input <-chan common.EnvironmentEntryWithState) {
			for entry := range input {
				out <- entry
			}
			wg.Done()
		}(c)
	}

	return out
}

// Retrieves a list of all environments,
func (s Selector) List() <-chan common.EnvironmentEntryWithState {
	driverListChannels := make([]<-chan common.EnvironmentEntryWithState, 0)
	for _, driver := range s.drivers {
		driverListChannels = append(driverListChannels, driver.List())
	}

	return mergeChans(driverListChannels...)
}

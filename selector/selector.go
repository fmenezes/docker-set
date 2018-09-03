// selector package contains logic to manage environments for docker
package selector

import "fmt"
import "os"

import "github.com/fmenezes/docker-set/selector/common"
import "github.com/fmenezes/docker-set/selector/drivers/docker_for_mac"
import "github.com/fmenezes/docker-set/selector/drivers/docker_machine"
import "github.com/fmenezes/docker-set/selector/drivers/vagrant"
import "github.com/fmenezes/docker-set/selector/storage"

// Holds all methods to manage environments
type Selector struct {
	drivers []common.Driver
}

// Returns a new instance of the selector
// can fail on when trying to create storage.FileStorage
func NewSelector() (*Selector, error) {
	selector := Selector{
		drivers: make([]common.Driver, 0),
	}

	selector.drivers = append(selector.drivers, docker_for_mac.NewDriver())
	selector.drivers = append(selector.drivers, docker_machine.NewDriver())

	store, err := storage.NewFileStorage()
	if err != nil {
		return nil, err
	}
	selector.drivers = append(selector.drivers, vagrant.NewDriver(*store))

	return &selector, nil
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
		return nil, fmt.Errorf("Driver %s not supported", driver)
	}

	return *selectedDriver, nil
}

// Appends new environment into store
// can fail when issuing driver add
func (s Selector) Add(entry common.EnvironmentEntry) error {
	selectedDriver, err := s.selectDriver(entry.Driver)
	if err != nil {
		return err
	}

	return selectedDriver.Add(entry)
}

func (s Selector) findEntry(name string) (*common.EnvironmentEntryWithState, error) {
	list, err := s.List()
	if err != nil {
		return nil, err
	}
	for _, item := range list {
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

// Retrieves environment variables for given name
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

// Retrieves a list of all environments
// can fail when issuing driver list
func (s Selector) List() ([]common.EnvironmentEntryWithState, error) {
	list := make([]common.EnvironmentEntryWithState, 0)

	for _, driver := range s.drivers {
		entryList, err := driver.List()
		if err != nil {
			return nil, err
		}
		if entryList != nil {
			list = append(list, entryList...)
		}
	}

	return list, nil
}

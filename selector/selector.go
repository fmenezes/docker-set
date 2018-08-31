package selector

import "fmt"
import "os"

import "github.com/fmenezes/docker-set/selector/types"
import "github.com/fmenezes/docker-set/selector/drivers"

var driverList = make([]types.Driver, 0)

func init() {
	driverList = append(driverList, drivers.NewDockerForMacDriver())
	driverList = append(driverList, drivers.NewDockerMachineDriver())
	driverList = append(driverList, drivers.NewVagrantDriver())
}

func selectDriver(driver string) (types.Driver, error) {
	var selectedDriver *types.Driver = nil
	for _, item := range driverList {
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

func Add(entry types.NewEnvironmentEntry) error {
	selectedDriver, err := selectDriver(entry.Driver)
	if err != nil {
		return err
	}

	return selectedDriver.Add(entry)
}

func findEntry(name string) (types.EnvironmentEntry, error) {
	list, err := List()
	if err != nil {
		return types.EnvironmentEntry{}, err
	}
	for _, item := range list {
		if item.Name == name {
			return item, nil
		}
	}
	return types.EnvironmentEntry{}, fmt.Errorf("'%s' not found", name)
}

func Selected() *string {
	var result *string = nil
	val, ok := os.LookupEnv("DOCKER_SET_MACHINE")
	if ok {
		result = &val
	}
	return result
}

func Env(entry string) (map[string]*string, error) {
	found, err := findEntry(entry)
	if err != nil {
		return nil, err
	}
	selectedDriver, err := selectDriver(found.Driver)
	if err != nil {
		return nil, err
	}
	return selectedDriver.Env(found)
}

func Remove(entry string) error {
	found, err := findEntry(entry)
	if err != nil {
		return err
	}
	selectedDriver, err := selectDriver(found.Driver)
	if err != nil {
		return err
	}
	return selectedDriver.Remove(types.NewEnvironmentEntry{
		Name:     found.Name,
		Driver:   found.Driver,
		Location: found.Location,
	})
}

func List() ([]types.EnvironmentEntry, error) {
	list := make([]types.EnvironmentEntry, 0)

	for _, driver := range driverList {
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

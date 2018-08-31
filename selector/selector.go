package selector

import "fmt"

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

func Remove(entry string) error {
	list, err := List()
	if err != nil {
		return err
	}
	var found *types.NewEnvironmentEntry = nil
	for _, item := range list {
		if item.Name == entry {
			found = &types.NewEnvironmentEntry{
				Name:     item.Name,
				Driver:   item.Driver,
				Location: item.Location,
			}
			break
		}
	}
	if found == nil {
		return fmt.Errorf("'%s' not found", entry)
	}
	selectedDriver, err := selectDriver(found.Driver)
	if err != nil {
		return err
	}
	return selectedDriver.Remove(*found)
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

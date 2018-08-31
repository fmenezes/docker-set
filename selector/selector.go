package selector

import "github.com/fmenezes/docker-set/selector/types"
import "github.com/fmenezes/docker-set/selector/drivers"

var driverList = make([]types.Driver, 0)

func init() {
	driverList = append(driverList, drivers.NewDockerForMacDriver())
	driverList = append(driverList, drivers.NewDockerMachineDriver())
}

func List() ([]types.EnvironmentEntry, error) {
	list := make([]types.EnvironmentEntry, 0)

	for _, driver := range driverList {
		entryList, err := driver.List()
		if err != nil {
			return nil, err
		}
		list = append(list, entryList...)
	}

	return list, nil
}

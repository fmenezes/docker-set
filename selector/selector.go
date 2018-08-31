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

func Add(entry types.NewEnvironmentEntry) error {
  var selectedDriver *types.Driver = nil
  for _, driver := range driverList {
    if driver.Name() == *entry.Driver {
      selectedDriver = &driver
    }
  }
  if selectedDriver == nil {
    return fmt.Errorf("Driver %s not supported", *entry.Driver)
  }

  return (*selectedDriver).Add(entry)
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

package drivers

import "errors"
import "github.com/fmenezes/docker-set/selector/types"

const DOCKERFORMAC = "docker-for-mac"
const DOCKERFORMACURL = "unix:///var/run/docker.sock"

type DockerForMacDriver struct {
	name string
}

func (driver DockerForMacDriver) Name() string {
	return driver.name
}

func (driver DockerForMacDriver) Add(entry types.NewEnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerForMacDriver) List() ([]types.EnvironmentEntry, error) {
	list := make([]types.EnvironmentEntry, 0)
	list = append(list, driver.getDockerForMacEntry())
	return list, nil
}

func (driver DockerForMacDriver) getDockerForMacEntry() types.EnvironmentEntry {
	name := DOCKERFORMAC
	url := DOCKERFORMACURL

	return types.EnvironmentEntry{
		Name:   name,
		State:  nil,
		Driver: driver.name,
		Active: false,
		URL:    &url,
	}
}

func NewDockerForMacDriver() types.Driver {
	return DockerForMacDriver{
		name: DOCKERFORMAC,
	}
}

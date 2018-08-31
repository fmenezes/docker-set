package drivers

import "github.com/fmenezes/docker-set/selector/types"

const DOCKERFORMAC = "docker-for-mac"
const UNKNOWN = "Unknown"
const DOCKERFORMACURL = "unix:///var/run/docker.sock"

type DockerForMacDriver struct {
	name string
}

func (driver DockerForMacDriver) List() ([]types.EnvironmentEntry, error) {
	list := make([]types.EnvironmentEntry, 0)
	list = append(list, driver.getDockerForMacEntry())
	return list, nil
}

func (driver DockerForMacDriver) getDockerForMacEntry() types.EnvironmentEntry {
	name := DOCKERFORMAC
	state := UNKNOWN
	url := DOCKERFORMACURL

	return types.EnvironmentEntry{
		Name:   &name,
		State:  &state,
		Source: &driver.name,
		Active: false,
		URL:    &url,
	}
}

func NewDockerForMacDriver() types.Driver {
	return DockerForMacDriver{
		name: DOCKERFORMAC,
	}
}

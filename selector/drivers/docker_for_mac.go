package drivers

import "errors"
import "github.com/fmenezes/docker-set/selector/types"

const DOCKERFORMAC = "docker-for-mac"

type DockerForMacDriver struct {
	name string
}

func (driver DockerForMacDriver) Name() string {
	return driver.name
}

func (driver DockerForMacDriver) Env(entry types.EnvironmentEntry) (map[string]*string, error) {
	env := make(map[string]*string)

	env["DOCKER_SET_MACHINE"] = &entry.Name
	env["DOCKER_HOST"] = nil
	env["DOCKER_TLS_VERIFY"] = nil
	env["DOCKER_CERT_PATH"] = nil
	env["DOCKER_MACHINE_NAME"] = nil

	return env, nil
}

func (driver DockerForMacDriver) Remove(entry types.NewEnvironmentEntry) error {
	return errors.New("Not Supported")
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

	return types.EnvironmentEntry{
		Name:   name,
		State:  nil,
		Driver: driver.name,
	}
}

func NewDockerForMacDriver() types.Driver {
	return DockerForMacDriver{
		name: DOCKERFORMAC,
	}
}

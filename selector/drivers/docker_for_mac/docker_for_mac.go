package docker_for_mac

import "errors"
import "github.com/fmenezes/docker-set/selector/common"

const DOCKERFORMAC = "docker-for-mac"

type DockerForMacDriver struct {
	name string
}

func (driver DockerForMacDriver) Name() string {
	return driver.name
}

func (driver DockerForMacDriver) Env(entry common.EnvironmentEntryWithState) (map[string]*string, error) {
	env := make(map[string]*string)

	env["DOCKER_SET_MACHINE"] = &entry.Name
	env["DOCKER_HOST"] = nil
	env["DOCKER_TLS_VERIFY"] = nil
	env["DOCKER_CERT_PATH"] = nil
	env["DOCKER_MACHINE_NAME"] = nil

	return env, nil
}

func (driver DockerForMacDriver) Remove(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerForMacDriver) Add(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerForMacDriver) List() ([]common.EnvironmentEntryWithState, error) {
	list := make([]common.EnvironmentEntryWithState, 0)
	list = append(list, driver.getDockerForMacEntry())
	return list, nil
}

func (driver DockerForMacDriver) getDockerForMacEntry() common.EnvironmentEntryWithState {
	return common.EnvironmentEntryWithState{
		EnvironmentEntry: common.EnvironmentEntry{
			Name:   driver.name,
			Driver: driver.name,
		},
		State: nil,
	}
}

func NewDriver() common.Driver {
	return DockerForMacDriver{
		name: DOCKERFORMAC,
	}
}

// docker_for_mac package contains Docker for Mac driver implementation
package docker_for_mac

import "errors"
import "os"
import "github.com/fmenezes/docker-set/selector/common"

// Holds common.Driver interface implementation for Docker for Mac
type DockerForMacDriver struct {
	name string
}

// Returns "docker-for-mac" always
func (driver DockerForMacDriver) Name() string {
	return driver.name
}

// Checks if the driver is supported
func (driver DockerForMacDriver) IsSupported() bool {
	if _, err := os.Stat("/Applications/Docker.app"); os.IsNotExist(err) {
		return false
	}
	return true
}

// Returns environment variables for Docker for Mac
func (driver DockerForMacDriver) Env(entry common.EnvironmentEntryWithState) (map[string]*string, error) {
	env := make(map[string]*string)

	env["DOCKER_SET_MACHINE"] = &entry.Name
	env["DOCKER_HOST"] = nil
	env["DOCKER_TLS_VERIFY"] = nil
	env["DOCKER_CERT_PATH"] = nil
	env["DOCKER_MACHINE_NAME"] = nil

	return env, nil
}

// Not Supported
func (driver DockerForMacDriver) Remove(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

// Not Supported
func (driver DockerForMacDriver) Add(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

// Lists always one entry representing Docker for Mac
func (driver DockerForMacDriver) List() <-chan common.EnvironmentEntryWithState {
	list := make(chan common.EnvironmentEntryWithState)

	go func(c chan common.EnvironmentEntryWithState) {
		c <- driver.getDockerForMacEntry()
		close(c)
	}(list)

	return list
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

// Returns an instance of DockerForMacDriver struct
func NewDriver() common.Driver {
	return DockerForMacDriver{
		name: "docker-for-mac",
	}
}

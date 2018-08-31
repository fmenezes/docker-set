package drivers

import "os/exec"
import "encoding/csv"
import "bytes"
import "errors"
import "github.com/fmenezes/docker-set/selector/types"

const DOCKERMACHINE = "docker-machine"

type DockerMachineDriver struct {
	name string
}

func (driver DockerMachineDriver) Name() string {
	return driver.name
}

func (driver DockerMachineDriver) Env(entry types.EnvironmentEntry) (map[string]*string, error) {
	return nil, errors.New("Not Supported")
}

func (driver DockerMachineDriver) Add(entry types.NewEnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerMachineDriver) Remove(entry types.NewEnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerMachineDriver) List() ([]types.EnvironmentEntry, error) {
	cmd := exec.Command("docker-machine", "ls", "-f", "{{.Name}},{{.State}}")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bytes.NewReader(output))
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	machines := make([]types.EnvironmentEntry, 0)
	for _, record := range records {
		machines = append(machines, types.EnvironmentEntry{
			Name:   record[0],
			State:  &record[1],
			Driver: driver.name,
		})
	}

	return machines, nil
}

func NewDockerMachineDriver() types.Driver {
	return DockerMachineDriver{
		name: DOCKERMACHINE,
	}
}

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

func (driver DockerMachineDriver) Add(entry types.NewEnvironmentEntry) error {
	return errors.New("Not Supported")
}

func (driver DockerMachineDriver) List() ([]types.EnvironmentEntry, error) {
	cmd := exec.Command("docker-machine", "ls", "-f", "{{.Name}},{{.State}},{{.Active}},{{.URL}}")

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
			Active: record[2] == "*",
			URL:    &record[3],
		})
	}

	return machines, nil
}

func NewDockerMachineDriver() types.Driver {
	return DockerMachineDriver{
		name: DOCKERMACHINE,
	}
}

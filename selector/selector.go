package selector

import (
	"bytes"
	"encoding/csv"
	"os/exec"
)

type environmentEntry struct {
	Name   *string
	Active bool
	Source *string
	State  *string
	URL    *string
}

const DOCKERMACHINE = "docker-machine"
const DOCKERFORMAC = "docker-for-mac"
const UNKNOWN = "Unknown"
const DOCKERFORMACURL = "unix:///var/run/docker.sock"

func listDockerMachines() ([]environmentEntry, error) {
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

	machines := []environmentEntry{}
	source := DOCKERMACHINE
	for _, record := range records {
		machines = append(machines, environmentEntry{
			Name:   &record[0],
			State:  &record[1],
			Source: &source,
			Active: record[2] == "*",
			URL:    &record[3],
		})
	}

	return machines, nil
}

func getDockerForMacEntry() environmentEntry {
	name := DOCKERFORMAC
	state := UNKNOWN
	url := DOCKERFORMACURL

	return environmentEntry{
		Name:   &name,
		State:  &state,
		Source: &name,
		Active: false,
		URL:    &url,
	}
}

func List() ([]environmentEntry, error) {
	list, err := listDockerMachines()
	if err != nil {
		return nil, err
	}
	list = append(list, getDockerForMacEntry())
	return list, nil
}

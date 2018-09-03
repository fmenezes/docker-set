package vagrant

import "os"
import "strings"
import "fmt"
import "os/exec"
import "encoding/csv"
import "bytes"
import "errors"
import "path"

import "github.com/fmenezes/docker-set/selector/common"

const VAGRANT = "vagrant"

type VagrantDriver struct {
	name  string
	Store common.Storage
}

func (driver VagrantDriver) Name() string {
	return driver.name
}

func (driver VagrantDriver) Add(entry common.EnvironmentEntry) error {
	if entry.Location == nil {
		return fmt.Errorf("No file provided")
	}

	info, err := os.Stat(*entry.Location)
	if err != nil {
		return fmt.Errorf("Can not access %s", *entry.Location)
	}

	if info.IsDir() {
		return errors.New("Directories are not supported, pass the Vagrantfile's full path")
	}

	return driver.Store.Append(entry)
}

func (driver VagrantDriver) Remove(entry common.EnvironmentEntry) error {
	return driver.Store.Remove(entry)
}

func (driver VagrantDriver) Env(entry common.EnvironmentEntryWithState) (map[string]*string, error) {
	if *entry.State != "running" {
		return nil, fmt.Errorf("vm is not running")
	}

	if entry.Location == nil {
		return nil, fmt.Errorf("No file provided")
	}

	ip, err := getVagrantIp(*entry.Location)
	if err != nil {
		return nil, err
	}
	host := fmt.Sprintf("tcp://%s:2375", ip)

	env := make(map[string]*string)
	env["DOCKER_SET_MACHINE"] = &entry.Name
	env["DOCKER_HOST"] = &host
	env["DOCKER_TLS_VERIFY"] = nil
	env["DOCKER_CERT_PATH"] = nil
	env["DOCKER_MACHINE_NAME"] = nil

	return env, nil
}

func (driver VagrantDriver) List() ([]common.EnvironmentEntryWithState, error) {
	data, err := driver.Store.Load()
	if err != nil {
		return nil, err
	}

	list := make([]common.EnvironmentEntryWithState, 0)
	for _, item := range data {
		if item.Location == nil {
			return nil, fmt.Errorf("No file provided")
		}

		state, err := getVagrantState(*item.Location)
		if err != nil {
			return nil, err
		}

		list = append(list, common.EnvironmentEntryWithState{
			EnvironmentEntry: common.EnvironmentEntry{
				Name:     item.Name,
				Driver:   driver.name,
				Location: item.Location,
			},
			State: &state,
		})
	}

	return list, nil
}

func getVagrantIp(location string) (string, error) {
	cmd := exec.Command("vagrant", "ssh", "-c", "hostname -I | cut -d' ' -f2")
	cmd.Dir = path.Dir(location)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(string(output), "\n", "", 1), nil
}

func getVagrantState(location string) (string, error) {
	cmd := exec.Command("vagrant", "status", "--machine-readable")
	cmd.Dir = path.Dir(location)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	r := csv.NewReader(bytes.NewReader(output))
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	for _, record := range records {
		if record[2] == "state" {
			return record[3], nil
		}
	}

	return "none", nil
}

func NewDriver(store common.Storage) common.Driver {
	return &VagrantDriver{
		name:  VAGRANT,
		Store: store,
	}
}

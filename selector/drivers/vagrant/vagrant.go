// vagrant package contains Vagrant driver implementation
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

// Holds common.Driver interface implementation for Vagrant
type VagrantDriver struct {
	// Storage for Vagrant box list
	Store common.Storage

	name string
}

// Returns "vagrant" always
func (driver VagrantDriver) Name() string {
	return driver.name
}

// Stores a Vagrant box into the dataset
// can fail on any storage failure (e.g. disk failure)
func (driver VagrantDriver) Add(entry common.EnvironmentEntry) error {
	if entry.Location == nil || len(*entry.Location) == 0 {
		return fmt.Errorf("No location provided")
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

// Deletes a Vagrant box from the dataset
// can fail on any storage failure (e.g. disk failure)
func (driver VagrantDriver) Remove(entry common.EnvironmentEntry) error {
	return driver.Store.Remove(entry)
}

// Returns environment variables for Vagrant entry
// can fail if Vagrant box is not running, or while retrieving box's ip
func (driver VagrantDriver) Env(entry common.EnvironmentEntryWithState) (map[string]*string, error) {
	if *entry.State != "running" {
		return nil, fmt.Errorf("vm is not running")
	}

	if entry.Location == nil || len(*entry.Location) == 0 {
		return nil, fmt.Errorf("No location provided")
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

// Lists all Vagrant boxes you added
// can fail while retrieving Vagrant box status
func (driver VagrantDriver) List() ([]common.EnvironmentEntryWithState, error) {
	data, err := driver.Store.Load()
	if err != nil {
		return nil, err
	}

	list := make([]common.EnvironmentEntryWithState, 0)
	for _, item := range data {
		if item.Location == nil || len(*item.Location) == 0 {
			return nil, fmt.Errorf("No location provided")
		}

		state, err := getVagrantState(*item.Location)
		if err != nil {
			return nil, err
		}

		list = append(list, common.EnvironmentEntryWithState{
			EnvironmentEntry: item,
			State:            &state,
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

// Returns an instance of VagrantDriver struct
func NewDriver(store common.Storage) common.Driver {
	return &VagrantDriver{
		name:  "vagrant",
		Store: store,
	}
}

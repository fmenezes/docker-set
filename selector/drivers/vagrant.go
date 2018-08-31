package drivers

import "os"
import "fmt"
import "os/exec"
import "encoding/csv"
import "bytes"
import "errors"
import "path"

import "github.com/fmenezes/docker-set/selector/types"
import "github.com/fmenezes/docker-set/selector/storage"

const VAGRANT = "vagrant"

type VagrantDriver struct {
	name string
}

func (driver VagrantDriver) Name() string {
	return driver.name
}

func (driver VagrantDriver) Add(entry types.NewEnvironmentEntry) error {
	info, err := os.Stat(entry.Location)
	if err != nil {
		return fmt.Errorf("Can not access %s", entry.Location)
	}

	if info.IsDir() {
		return errors.New("Directories are not supported, pass the Vagrantfile's full path")
	}

	return storage.Save(storage.Entry{
		Name:     entry.Name,
		Driver:   entry.Driver,
		Location: entry.Location,
	})
}

func (driver VagrantDriver) List() ([]types.EnvironmentEntry, error) {
	data, err := storage.Load()
	if err != nil {
		return nil, err
	}

	list := make([]types.EnvironmentEntry, 0)
	for _, item := range data {
		state, err := getVagrantState(path.Dir(item.Location))
		if err != nil {
			return nil, err
		}

		list = append(list, types.EnvironmentEntry{
			Name:   item.Name,
			Active: false,
			Driver: driver.name,
			State:  &state,
			URL:    nil,
		})
	}

	return list, nil
}

func getVagrantState(dir string) (string, error) {
	cmd := exec.Command("vagrant", "status", "--machine-readable")
	cmd.Dir = dir

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

func NewVagrantDriver() types.Driver {
	return VagrantDriver{
		name: VAGRANT,
	}
}

// dockermachine package contains docker-machine driver implementation
package dockermachine

import "os/exec"
import "encoding/json"
import "encoding/csv"
import "fmt"
import "io"
import "strings"
import "errors"
import "github.com/fmenezes/docker-set/selector/common"

// Holds common.Driver interface implementation for docker-machine tool
type DockerMachineDriver struct {
	name string
}

type dockerMachineDetails struct {
	DockerTlsVerify bool    `json:"DOCKER_TLS_VERIFY"`
	DockerCertPath  *string `json:"DOCKER_CERT_PATH"`
}

type dockerMachineEnv struct {
	dockerMachineDetails
	DockerMachineName string `json:"DOCKER_MACHINE_NAME"`
	DockerHost        string `json:"DOCKER_HOST"`
}

// Returns "docker-machine" always
func (driver DockerMachineDriver) Name() string {
	return driver.name
}

// Checks if the driver is supported
func (driver DockerMachineDriver) IsSupported() bool {
	_, err := exec.LookPath("docker-machine")
	if err != nil {
		return false
	}
	return true
}

// Starts the docker-machine vm
func (driver DockerMachineDriver) Start(entry common.EnvironmentEntry) error {
	cmd := exec.Command("docker-machine", "start", entry.Name)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Stops the docker-machine vm
func (driver DockerMachineDriver) Stop(entry common.EnvironmentEntry) error {
	cmd := exec.Command("docker-machine", "stop", entry.Name)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getMachineDetails(machineName string) (dockerMachineDetails, error) {
	cmd := exec.Command("docker-machine", "inspect", machineName, "-f", "{\"DOCKER_TLS_VERIFY\":{{.HostOptions.EngineOptions.TlsVerify}},\"DOCKER_CERT_PATH\":\"{{.HostOptions.AuthOptions.StorePath}}\"}")
	output, err := cmd.Output()
	if err != nil {
		return dockerMachineDetails{}, err
	}

	var result dockerMachineDetails
	err = json.Unmarshal(output, &result)
	if err != nil {
		return dockerMachineDetails{}, err
	}

	return result, nil
}

func getMachineHost(machineName string) (string, error) {
	cmd := exec.Command("docker-machine", "ls", machineName, "-f", "{{.URL}}")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(string(output), "\n", "", -1), nil
}

func getMachineEnv(machineName string) (dockerMachineEnv, error) {
	details, err := getMachineDetails(machineName)
	if err != nil {
		return dockerMachineEnv{}, err
	}

	host, err := getMachineHost(machineName)
	if err != nil {
		return dockerMachineEnv{}, err
	}

	return dockerMachineEnv{
		dockerMachineDetails: dockerMachineDetails{
			DockerTlsVerify: details.DockerTlsVerify,
			DockerCertPath:  details.DockerCertPath,
		},
		DockerMachineName: machineName,
		DockerHost:        host,
	}, nil
}

// Returns environment variables for docker-machine entry, similar to 'docker-machine env name',
// can fail if docker-machine is not running
func (driver DockerMachineDriver) Env(entry common.EnvironmentEntryWithState) (map[string]*string, error) {
	if *entry.State != "Running" {
		return nil, fmt.Errorf("vm is not running")
	}

	vars, err := getMachineEnv(entry.Name)
	if err != nil {
		return nil, err
	}

	env := make(map[string]*string)
	env["DOCKER_SET_MACHINE"] = &entry.Name
	env["DOCKER_HOST"] = &vars.DockerHost
	tlsValue := "0"
	if vars.DockerTlsVerify {
		tlsValue = "1"
	}
	env["DOCKER_TLS_VERIFY"] = &tlsValue
	env["DOCKER_CERT_PATH"] = vars.DockerCertPath
	env["DOCKER_MACHINE_NAME"] = &vars.DockerMachineName

	return env, nil
}

// Not Supported
func (driver DockerMachineDriver) Add(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

// Not Supported
func (driver DockerMachineDriver) Remove(entry common.EnvironmentEntry) error {
	return errors.New("Not Supported")
}

// Lists docker-machine boxes, similar to 'docker-machine ls'
func (driver DockerMachineDriver) List() <-chan common.EnvironmentEntryWithState {
	list := make(chan common.EnvironmentEntryWithState)

	go func(c chan common.EnvironmentEntryWithState) {
		cmd := exec.Command("docker-machine", "ls", "-f", "{{.Name}},{{.State}}")

		defer close(c)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		r := csv.NewReader(stdout)
		r.FieldsPerRecord = -1
		for {
			record, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			c <- common.EnvironmentEntryWithState{
				EnvironmentEntry: common.EnvironmentEntry{
					Name:   record[0],
					Driver: driver.name,
				},
				State: &record[1],
			}
		}

		err = cmd.Wait()
		if err != nil {
			panic(err)
		}
	}(list)

	return list
}

// Returns an instance of DockerMachineDriver struct
func NewDriver() common.Driver {
	return DockerMachineDriver{
		name: "docker-machine",
	}
}

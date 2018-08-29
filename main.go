package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os/exec"
)

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

func getDockerMachineState(machine string) (string, error) {
	cmd := exec.Command("docker-machine", "ls", "-f", "{{.State}}", machine)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func main() {
	state, err := getVagrantState("/Users/fmenezes/Code/zendesk/docker-images")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("vagrant", state)

	state, err = getDockerMachineState("default")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("docker-machine", state)
}

// types of vms vagrant, docker-machine, docker-for-mac, remote
// add vagrant name, file path
// add remote name, $DOCKER_HOST
//
// ls (should say which one is activated in current terminal)
// -> docker-for-mac
// -> all docker machines
// -> all vagrant machines
// -> all remote machines
//
// rm name --only for vagrant or remote
//
// env name
// export DOCKER_HOST=XXXX

// package main

// import "fmt"

// func main() {
//  fmt.Println("export DOCKER_SET_HOST='docker-for-mac'")
//  fmt.Println("unset DOCKER_TLS_VERIFY")
//  fmt.Println("unset DOCKER_HOST")
//  fmt.Println("unset DOCKER_CERT_PATH")
//  fmt.Println("unset DOCKER_MACHINE_NAME")
//  fmt.Println("# eval $(docker-set env docker-for-mac)")
// }

// # docker-for-mac
// unset DOCKER_TLS_VERIFY
// unset DOCKER_HOST
// unset DOCKER_CERT_PATH
// unset DOCKER_MACHINE_NAME

// # docker-machine
// docker-machine ls -f '"DOCKER_MACHINE_NAME"="{{.Name}}","DOCKER_HOST"="{{.URL}}","STATE":"{{.State}}"'
// docker-machine inspect $DOCKER_MACHINE_NAME -f '"DOCKER_TLS_VERIFY"={{.HostOptions.EngineOptions.TlsVerify}},"DOCKER_CERT_PATH"="{{.HostOptions.AuthOptions.StorePath}}"'

// # vagrant
// unset DOCKER_TLS_VERIFY
// unset DOCKER_CERT_PATH
// unset DOCKER_MACHINE_NAME

// export DOCKER_HOST=
// tcp://<<IP>>:2375
// cd VAGRANTFILEPATH
// vagrant ssh -c "hostname -I | cut -d' ' -f2" 2>/dev/null

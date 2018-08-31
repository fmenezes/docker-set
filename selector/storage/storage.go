package storage

import "encoding/json"
import "io/ioutil"
import "os"
import "path"
import "os/user"

type Entry struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Location string `json:"location"`
}

func getStorageFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".docker-set"), nil
}

func Save(entry Entry) error {
	storageFile, err := getStorageFile()
	if err != nil {
		return err
	}

	list, err := Load()
	if err != nil {
		return err
	}

	list = append(list, entry)
	json, err := marshal(list)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(storageFile, json, 0644)
}

func marshal(data []Entry) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func unmarshal(file []byte) ([]Entry, error) {
	var result []Entry
	err := json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Load() ([]Entry, error) {
	storageFile, err := getStorageFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		return make([]Entry, 0), nil
	}

	file, err := ioutil.ReadFile(storageFile)
	if err != nil {
		return nil, err
	}
	return unmarshal(file)
}

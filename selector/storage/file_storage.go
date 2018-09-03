package storage

import "encoding/json"
import "io/ioutil"
import "os"
import "errors"
import "path"
import "os/user"

import "github.com/fmenezes/docker-set/selector/common"

type FileStorage struct {
	file string
}

func NewFileStorage() (*FileStorage, error) {
	file, err := getStorageFile()

	if err != nil {
		return nil, err
	}

	store := FileStorage{
		file: file,
	}

	return &store, nil
}

func getStorageFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".docker-set"), nil
}

func (s FileStorage) Add(entry common.EnvironmentEntry) error {
	list, err := s.Load()
	if err != nil {
		return err
	}

	list = append(list, entry)

	return s.Save(list)
}

func (s FileStorage) Remove(entry common.EnvironmentEntry) error {
	list, err := s.Load()
	if err != nil {
		return err
	}

	found := false

	newList := make([]common.EnvironmentEntry, 0)
	for _, item := range list {
		if entry.Name == item.Name && entry.Driver == item.Driver && *entry.Location == *item.Location {
			found = true
		} else {
			newList = append(newList, item)
		}
	}

	if !found {
		return errors.New("Entry not found")
	}

	return s.Save(newList)
}

func (s FileStorage) Save(list []common.EnvironmentEntry) error {
	data, err := marshal(list)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.file, data, 0644)
}

func marshal(data []common.EnvironmentEntry) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func unmarshal(file []byte) ([]common.EnvironmentEntry, error) {
	var result []common.EnvironmentEntry
	err := json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s FileStorage) Load() ([]common.EnvironmentEntry, error) {
	if _, err := os.Stat(s.file); os.IsNotExist(err) {
		return make([]common.EnvironmentEntry, 0), nil
	}

	file, err := ioutil.ReadFile(s.file)
	if err != nil {
		return nil, err
	}
	return unmarshal(file)
}

package storage

import (
	"encoding/json"
	"fmt"
	"github.com/jasongerard/remoteit-cli/client"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
)

type StorageFile string

const LoginFile StorageFile = "login"
const DeviceCacheFile StorageFile = "devices"

var configDir string

func Initialize () error {
	dir, err := homedir.Dir()

	if err != nil {
		return err
	}

	configDir = fmt.Sprintf("%s/.remoteit", dir)

	err = os.Mkdir(configDir, 0700)

	if os.IsExist(err) {
		return nil
	}

	return err
}

func getPath(name StorageFile) string {
	return fmt.Sprintf("%s/%s", configDir, name)
}

func WriteFile(name StorageFile, b []byte) error {
	p := getPath(name)

	return ioutil.WriteFile(p, b, 0600)
}

func ReadFile(name StorageFile) ([]byte, error) {
	p := getPath(name)

	return ioutil.ReadFile(p)
}

func CacheExists() bool {
	if _, err := os.Stat(getPath(DeviceCacheFile)); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDevicesFromCache() ([]*client.DeviceEntry, error) {
	b, err := ReadFile(DeviceCacheFile)

	if err != nil {
		return nil, err
	}

	devices := new(client.ListDevicesResponse)
	err = json.Unmarshal(b, devices)

	if err != nil {
		return nil, err
	}

	return devices.Devices, nil
}

func GetToken() (string, error) {
	b, err := ReadFile(LoginFile)

	if err != nil {
		return "", err
	}

	lr := new(client.LoginResponse)

	err = json.Unmarshal(b, lr)

	if err != nil {
		return "", err
	}

	return lr.Token, nil
}


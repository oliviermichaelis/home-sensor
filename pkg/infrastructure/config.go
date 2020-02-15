package infrastructure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type config map[string]string

var configMap = make(config)

func (c *config) register(key string, fallback string) (string, error) {
	if len(key) <= 0 {
		return "", errors.New("config: to be registered key is empty")
	}

	returned := getEnv(key, fallback)
	if len(returned) <= 0 {
		message := fmt.Sprintf("config: registered key: %s not found and no fallback provided", key)
		return "", errors.New(message)
	}

	(*c)[key] = returned
	return returned, nil
}

func (c *config) get(key string) (string, error) {
	if len(key) <= 0 {
		return "", errors.New("config: key is empty")
	}

	retrieved, ok := (*c)[key]
	if !ok {
		return "", errors.New("config: configMap doesn't contain the key:" + key)
	}
	if len(retrieved) <= 0 {
		return "", errors.New(fmt.Sprintf("config: value for key: %s is empty", key))
	}

	return retrieved, nil
}

// Returns an environment variable if it exists, else it returns a fallback variable passed as parameter
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// Registers a config key to be retrieved from ENV
func RegisterConfig(key string, fallback string) (string, error) {
	return configMap.register(key, fallback)
}

func GetConfig(key string) (string, error) {
	return configMap.get(key)
}

// TODO deprecated, use ReadConfig
// Returns the username read from a 'username' file under the specified parameter
func ReadUsername(secretPath string) (string, error) {
	u, err := ioutil.ReadFile(secretPath + "/username")
	return string(u), err
}

// TODO deprecated, use ReadConfig
// Returns the password read from a 'password' file under the specified parameter
func ReadPassword(secretPath string) (string, error) {
	p, err := ioutil.ReadFile(secretPath + "/password")
	return string(p), err
}

func ReadSecret(path string) (string, error) {
	c, err := ioutil.ReadFile(path)
	return string(c), err
}

package config

import (
	"path/filepath"
	"os"
	"encoding/json"
)


const (
	configFileName = ".gatorconfig.json"
	path = "/home/tspedtsberg/"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fullpath := filepath.Join(path, configFileName)
	// home/path/fgatorconfig.json
	file, err := os.Open(fullpath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	cfg := Config{}
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(UserName string) error {
	cfg.CurrentUserName = UserName
	return write(cfg)
}


func write(cfg *Config) error {
	fullpath := filepath.Join(path, configFileName)
	// home/path/fgatorconfig.json
	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}


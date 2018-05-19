package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"strings"
)

func loadConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	var cfg Config
	file := path.Join(usr.HomeDir, ".po-domaci.json")
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Save to %s:\n", file)
		//fmt.Println(jsonizePrettified(cfg))
		return nil, fmt.Errorf("error opening/reading %s: %s", file, err.Error())
	}

	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling %s: %s", file, err.Error())
	}

	if strings.HasPrefix(cfg.DatabaseFilename, "~") {
		cfg.DatabaseFilename = path.Join(usr.HomeDir, cfg.DatabaseFilename[1:])
	}

	return &cfg, nil
}

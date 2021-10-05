package clusterconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ClusterCfg struct {
	Name  string
	Nodes []Node
}

type Node struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	ExecCmd string `json:"command"`   //TODO - change this to execcmd to avoid confusion with type
	Args    string `json:"arguments"` //TODO - Ugh, fix this too. This needs to be an array. How we gonna input that?
	RootDir string `json:"rootdirectory"`
}

func WriteConfig(cfg ClusterCfg) error {
	file, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs("internal/set-configs/" + cfg.Name + ".json")
	if err != nil {
		return err
	}

	var dest *os.File

	//if file doesn't exist, create it
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		os.MkdirAll("internal/set-configs/", 0700)

		dest, err = os.Create(absPath)
		if err != nil {
			return err
		}
		defer dest.Close()

	} else { //otherwise append. Will probably change to overwrite
		dest, err = os.Open(absPath)
		if err != nil {
			log.Printf("error opening file: %s", err.Error())
			return err
		}
		defer dest.Close()
	}

	err = ioutil.WriteFile(absPath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(name string) (*ClusterCfg, error) {
	absPath, err := filepath.Abs("internal/set-configs/" + name + ".json")
	if err != nil {
		return nil, err
	}

	//check that file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteVal, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg ClusterCfg
	err = json.Unmarshal(byteVal, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil

}

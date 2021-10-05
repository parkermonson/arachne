package clusterconfig

//TODO - build constraints and a seperate file for windows. Because windows.
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func CreateCluster() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Creating New Cluster")

	// nodeSettings := []string{"name", "type", "command", "arguments", "rootdirectory"}

	finishedNodes := make([]Node, 0)

	//get the cluster name
	name, err := getClusterName(reader)
	if err != nil {
		log.Panic("error getting cluster name")
	}

	for {
		node := Node{}

		err := defineNode(&node, reader)
		if err != nil {
			log.Panicf("error getting input: %s\n", err.Error())
		}

		finishedNodes = append(finishedNodes, node)

		next, err := createNextNode(reader)
		if next && err == nil {
			continue
		} else if err != nil {
			log.Panicf("error handling next input: %s", err.Error)
		}
		break
	}

	file, err := json.MarshalIndent(ClusterCfg{Name: name, Nodes: finishedNodes}, "", "    ")
	if err != nil {
		log.Panicf("error encoding json file: %s\n", err.Error())
	}

	// err = ioutil.WriteFile("../cluster-configs/config.json", file, 0644)
	// if err != nil {
	// 	log.Panicf("error writing json file: %s\n", err.Error())
	// }

	// log.Println("success write?")
	// absDir, err := filepath.Abs("/internal/set-configs/")
	absPath, err := filepath.Abs("internal/set-configs/" + name + ".json")

	var dest *os.File

	//if file doesn't exist, create it
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		os.MkdirAll("internal/set-configs/", 0700)

		dest, err = os.Create(absPath)
		if err != nil {
			log.Printf("error creating file: %s", err.Error())
			return
		}
		defer dest.Close()

	} else { //otherwise append. Will probably change to overwrite
		dest, err = os.Open(absPath)
		if err != nil {
			log.Printf("error opening file: %s", err.Error())
			return
		}
		defer dest.Close()
	}

	err = ioutil.WriteFile(absPath, file, 0644)
	if err != nil {
		log.Panicf("error writing json file: %s\n", err.Error())
	}
}

func getClusterName(reader *bufio.Reader) (string, error) {
	fmt.Print("Cluster Name: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.Replace(text, "\n", "", -1)
	return text, nil
}

func defineNode(node *Node, reader *bufio.Reader) error {
	name, err := handleInput("name", reader)
	if err != nil {
		return err
	}
	node.Name = name

	nodeType, err := handleInput("type", reader)
	if err != nil {
		return err
	}
	node.Type = nodeType

	execCmd, err := handleInput("command", reader)
	if err != nil {
		return err
	}
	node.ExecCmd = execCmd

	args, err := handleInput("arguments", reader)
	if err != nil {
		return err
	}
	node.Args = args

	root, err := handleInput("rootdirectory", reader)
	if err != nil {
		return err
	}
	node.RootDir = root

	return nil
}

//TODO - validate input
func handleInput(settingKey string, reader *bufio.Reader) (string, error) {
	displayMsg := displayMessage(settingKey)
	fmt.Print(displayMsg)

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) //needs to be "\r\n" for windows

	return text, nil
}

func displayMessage(key string) string {
	switch key {
	case "name":
		return "Node Name: "
	case "type":
		return "Node Type [command/service/daemon]: "
	case "command":
		return "Execution Command: "
	case "arguments":
		return "Execution Arguments [x, y, z]: "
	case "rootdirectory":
		return "Root Directory: "
	}
	return "->"
}

func createNextNode(reader *bufio.Reader) (bool, error) {
	for {
		fmt.Print("Create another Node? [y/N]")
		text, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		text = strings.Replace(text, "\n", "", -1)

		if text != "y" && text != "N" {
			fmt.Println("please input a 'y' for yes or 'N' for no (case matters).")
			continue
		}
		return text == "y", nil
	}
}

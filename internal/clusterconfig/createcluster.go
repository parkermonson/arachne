package clusterconfig

//TODO - build constraints and a seperate file for windows. Because windows.
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parkermonson/arachne/internal/helpers"
)

func CreateCluster() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Creating New Cluster")

	// nodeSettings := []string{"name", "type", "command", "arguments", "rootdirectory"}

	finishedNodes := make([]helpers.Node, 0)

	//get the cluster name
	name, err := getClusterName(reader)
	if err != nil {
		log.Panic("error getting cluster name")
	}

	for {
		node := helpers.Node{}

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

	helpers.WriteConfig(helpers.ClusterCfg{Name: name, Nodes: finishedNodes})
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

func defineNode(node *helpers.Node, reader *bufio.Reader) error {
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

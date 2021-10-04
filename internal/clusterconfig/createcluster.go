package clusterconfig

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func CreateCluster() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Creating New Cluster")

	nodeSettings := []string{"name", "type", "command", "arguments", "rootdirectory"}

	finishedNodes := make([]map[string]string, 0)

	//get the cluster name
	name, err := getClusterName(reader)
	if err != nil {
		log.Panic("error getting cluster name")
	}

	for {

		newNode := map[string]string{}

		for _, settingKey := range nodeSettings {
			err := handleInput(settingKey, newNode, reader)
			if err != nil {
				//do the bad
			}
		}

		ok, err := createNextNode(reader)
		if !ok || err != nil {
			break
		}

		finishedNodes = append(finishedNodes, newNode)

	}

	fmt.Println("\n\nFinal cluster setup:")
	fmt.Println("Name: " + name)
	fmt.Printf("%+v\n", finishedNodes)

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

//TODO - handle bad input, maybe?
func handleInput(settingKey string, node map[string]string, reader *bufio.Reader) error {
	displayMsg := displayMessage(settingKey)
	fmt.Print(displayMsg)

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) //needs to be "\r\n" for windows

	node[settingKey] = text
	return nil
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

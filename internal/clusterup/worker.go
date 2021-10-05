package clusterup

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/parkermonson/arachne/internal/helpers"
)

type WorkerNode struct {
	//metadata to run the command
	Name    string
	Type    string
	Command string
	Args    []string
	RootDir string

	//out channels
	MessageChan chan MessageData
	StoppedChan chan struct{}

	//in channel
	stopChan chan struct{}
}

//messages on outchannel
type MessageData struct {
	CommandName string
	RawText     string
	ErrMessage  string
}

func CreateWorker(meta helpers.Node, messages chan MessageData, stopped chan struct{}, stopWorker chan struct{}) *WorkerNode {
	//set a name
	name := meta.Name
	if name == "" {
		name = meta.ExecCmd
	}

	return &WorkerNode{
		Name:        name,
		Type:        meta.Type,
		Command:     meta.ExecCmd,
		Args:        []string{meta.Args},
		RootDir:     meta.RootDir,
		MessageChan: messages,
		StoppedChan: stopped,
		stopChan:    stopWorker,
	}
}

func (wn *WorkerNode) WorkerUp() {
	if wn.Type == "command" {
		wn.runCmd()
		return
	}
	go func() {
		for {
			select {
			case <-wn.stopChan:
				wn.StoppedChan <- struct{}{}
				fmt.Printf("worker [%s] closing\n", wn.Name)
				return
			}
		}
	}()

	go watchRoot(wn.RootDir, wn.stopChan)
	go wn.runService()

	fmt.Println("returning from workerUp")

}

//commands are fire and forget, wait for it to finish
func (wn *WorkerNode) runCmd() error {
	cmd := exec.Command(wn.Command, wn.Args...)

	err := cmd.Run()
	if err != nil {
		return err //maybe this should be sent out in the error channel
	}

	//signal that we are done and the next worker can be spun up
	wn.MessageChan <- MessageData{
		CommandName: wn.Name,
		RawText:     "eagle has flopped", //TODO - encode this
	}
	wn.StoppedChan <- struct{}{}

	return nil
}

//services are expected to stay up
func (wn *WorkerNode) runService() {
	serviceCmd := exec.Command(wn.Command, wn.Args...)

	// stderr, err := serviceCmd.StderrPipe()//TODO - Later
	// if err != nil {

	// }

	stdout, err := serviceCmd.StdoutPipe()
	if err != nil {
		fmt.Println("error getting stdout pipe")
		wn.StoppedChan <- struct{}{}
		return
	}

	rd := bufio.NewReader(stdout)

	if err := serviceCmd.Start(); err != nil {
		fmt.Println("error starting service\n")
		wn.StoppedChan <- struct{}{}
		return
	}

	wn.MessageChan <- MessageData{
		CommandName: wn.Name,
		RawText:     "eagle has flopped",
	}

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			//TODO - this is erroring on an empty output, fix that
			fmt.Println("error reading output from command")
			fmt.Println("closing service " + wn.Name)
			wn.StoppedChan <- struct{}{}
			return
		}

		wn.MessageChan <- MessageData{
			CommandName: wn.Name,
			RawText:     str,
		}
	}

}

package clusterup

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/parkermonson/arachne/internal/helpers"
)

var StandbyQueue []helpers.Node

//this will be useless
type WorkerPool struct {
	woutChan    chan MessageData
	restartChan chan struct{}
	cancelFunc  context.CancelFunc
}

//this runs the pool, and reruns it in the event of a restart cluster
func Orchestrate(name string) {

	var wg sync.WaitGroup

	go func() {
		for {
			wg.Add(1)
			runPool(wg, name)
		}
	}()

	//this works
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan //block until interrupted

	wg.Done()
}

//this runs the commands and closes them on update or cluster exit
func runPool(wg sync.WaitGroup, name string) {
	var err error

	//set up queue of workers
	cluster, err := helpers.ReadConfig(name) //TODO - ugh, naming
	if err != nil {
		//log the bad, do the bad, this is bad
	}

	StandbyQueue = cluster.Nodes

	//this is how we do graceful shutdown, waiting to close channels until all workers goroutines have stopped
	totalWorkers := len(StandbyQueue)
	stoppedWorkers := 0

	messages := make(chan MessageData, len(StandbyQueue))
	stopped := make(chan struct{}, len(StandbyQueue))
	stopWorkers := make(chan struct{}, len(StandbyQueue))

	go func() {
		defer close(stopWorkers)

		defer func() {
			wg.Done()
		}()

		for {
			select {
			case msg := <-messages:
				if msg.RawText == "eagle has flopped" {
					//do the next worker
					if len(StandbyQueue) > 0 {
						next := CreateWorker(StandbyQueue[0], messages, stopped, stopWorkers)
						StandbyQueue = StandbyQueue[1:]

						next.WorkerUp()
					} else {
						//error, something is wrong we have spun up too many
					}
				} else {
					fmt.Printf("%s: %s\n", msg.CommandName, msg.RawText)
				}
			case <-stopped:
				fmt.Println("stopped message received")
				stoppedWorkers++
				if stoppedWorkers == totalWorkers {
					fmt.Println("all workers have closed, ready for restart")
					//close everything down
					return
				}

			}
		}
	}()

	//run the first worker, wait for a response to run the rest
	fmt.Printf("starting worker: %s\n", StandbyQueue[0].Name)
	nworker := CreateWorker(StandbyQueue[0], messages, stopped, stopWorkers)
	StandbyQueue = StandbyQueue[1:]
	nworker.WorkerUp()

	wg.Wait()

	// return "restart"

}

// //This gets pulled from the json
// type WorkerMetadata struct {
// 	Name    string `json:""`
// 	Type    string
// 	Command string
// 	Args    []string
// 	RootDir string
// }

// func getClusterConfig(name string) ([]WorkerMetadata, error) {
// 	runNodes := make([]WorkerMetadata, 0)

// 	node1 := WorkerMetadata{
// 		Name:    "test1",
// 		Type:    "command",
// 		Command: "/Users/parker/go/src/personal/arachne-test-services/test1/test1",
// 		RootDir: "/Users/parker/go/src/personal/arachne-test-services/test1/",
// 	}

// 	node2 := WorkerMetadata{
// 		Name:    "test2",
// 		Type:    "service",
// 		Command: "/Users/parker/go/src/personal/arachne-test-services/test2/test2",
// 		RootDir: "/Users/parker/go/src/personal/arachne-test-services/test2/",
// 	}

// 	runNodes = append(runNodes, node1)

// 	runNodes = append(runNodes, node2)

// 	return runNodes, nil
// }

// //maybe in another file?

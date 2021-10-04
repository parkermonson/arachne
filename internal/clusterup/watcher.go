package clusterup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func watchRoot(path string, stopSignal chan<- struct{}) error {
	fmt.Printf("watching directory: [%s]\n", path)
	initSize, err := getDirSize(path)
	if err != nil {
		fmt.Println("error getting directory size, sending stop signal")
		stopSignal <- struct{}{}
		return nil
	}

	for {
		sizeCheck, err := getDirSize(path)
		if err != nil {
			return err
		}

		if initSize != sizeCheck {
			fmt.Println("a size has changed")
			stopSignal <- struct{}{}
			return nil
		}

		time.Sleep(1 * time.Second)
	}

}

func getDirSize(dirpath string) (int64, error) {
	var dirSize int64
	var outSideError error

	filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// log.Fa tal(err.Error())
			outSideError = err
			return nil
		}
		// fmt.Printf("\n\nfile info: %+v\n\n", info)
		dirSize += info.Size()
		return nil
	})

	return dirSize, outSideError
}

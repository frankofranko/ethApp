package main

import (
	"ethApp/utils"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Indexer spawn error. Please provide start block number and worker number")
	}

	startBlockNum, err := strconv.ParseUint(os.Args[1], 10, 64)
	worker_num, err2 := strconv.Atoi(os.Args[2])
	if err != nil || err2 != nil {
		fmt.Println("Indexer spawn error. Please provide start block number and worker number")
		fmt.Println(err)
	} else {
		utils.SpawnWorkers(startBlockNum, worker_num)
	}
}

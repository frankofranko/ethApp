package main

import (
	"ethApp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()
	router.GET("/blocks/:id", getBlockByNumberView)
	router.GET("/blocks", getLatestNBlocksView)
	router.GET("/transaction/:txHash", getTransactionByHashView)
	router.Run("localhost:8080")
}

func getLatestNBlocksView(c *gin.Context) {
	numOfResults, err := strconv.ParseUint(c.Query("limit"), 10, 64)
	if err != nil {
		fmt.Println("Failed parsing request")
	}
	latestBlockNum := utils.GetLatestBlockNumber()
	result := make([]utils.Block, numOfResults)
	for i := uint64(0); i < numOfResults; i++ {
		blockRaw := utils.GetBlockRawByNumber(latestBlockNum - i)
		result[i] = utils.BlockRawToBlock(blockRaw)
	}
	c.IndentedJSON(http.StatusOK, result)
}

func getBlockByNumberView(c *gin.Context) {
	blockNum, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	blockRaw := utils.GetBlockRawByNumber(blockNum)
	c.IndentedJSON(http.StatusOK, utils.BlockRawToBlockAndTransactionHashes(blockRaw))
}

func getTransactionByHashView(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, utils.GetTransactionByHash(c.Param("txHash")))
}

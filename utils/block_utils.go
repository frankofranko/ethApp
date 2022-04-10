package utils

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var client, _ = ethclient.Dial("https://data-seed-prebsc-2-s3.binance.org:8545/")

// Block
type Block struct {
	BlockNum   uint64 `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}

type blockAndTransactionHashes struct {
	Block
	TransactionHashes []string `json:"transactions"`
}

type transactionLog struct {
	Index uint   `json:"index"`
	Data  string `json:"data"`
}

type transaction struct {
	TransactionHash string           `json:"tx_hash"`
	From            string           `json:"from"`
	To              string           `json:"to"`
	Nonce           uint64           `json:"nonce"`
	Data            string           `json:"data"`
	Value           string           `json:"value"`
	Logs            []transactionLog `json:"logs""`
}

func GetLatestBlockNumber() uint64 {
	latestBlockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return latestBlockNum
}

func GetBlockRawByNumber(blockNumber uint64) *types.Block {
	blockRaw, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		log.Fatal(err)
	}
	return blockRaw
}

func BlockRawToBlock(blockRaw *types.Block) Block {
	return Block{blockRaw.Number().Uint64(), blockRaw.Hash().Hex(), blockRaw.Time(), blockRaw.ParentHash().Hex()}
}

func BlockRawToBlockAndTransactionHashes(blockRaw *types.Block) blockAndTransactionHashes {
	transactionHashes := make([]string, blockRaw.Transactions().Len())

	for i, tx := range blockRaw.Transactions() {
		transactionHashes[i] = tx.Hash().Hex()
	}
	return blockAndTransactionHashes{Block{blockRaw.Number().Uint64(), blockRaw.Hash().Hex(),
		blockRaw.Time(), blockRaw.ParentHash().Hex()}, transactionHashes}
}

func GetTransactionByHash(transactionHash string) transaction {
	txHash := common.HexToHash(transactionHash)
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		fmt.Println("Failed retrieving transaction by transaction hash.")
	}
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		fmt.Println("Failed retrieving transaction receipt by transaction hash.")
	}
	logs := make([]transactionLog, len(receipt.Logs))
	var from string
	for i, txnLog := range receipt.Logs {
		logs[i] = transactionLog{Data: hex.EncodeToString(txnLog.Data), Index: txnLog.Index}
		sender, _ := client.TransactionSender(context.Background(), tx, txnLog.BlockHash, txnLog.Index)
		from = sender.Hex()
	}
	return transaction{TransactionHash: tx.Hash().Hex(), From: from, To: tx.To().Hex(), Nonce: tx.Nonce(),
		Data: hex.EncodeToString(tx.Data()), Value: tx.Value().String(), Logs: logs}
}

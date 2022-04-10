package utils

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "eth_db"
)

func SpawnWorkers(startBlockNum uint64, worker_num int) {
	pool := workerPool{worker_num, make(chan func(db *sql.DB), 100)}
	pool.Run()
	fmt.Printf("Workers are running.\n")
	for true {
		latestBlockNum := GetLatestBlockNumber()
		for blockNum := startBlockNum; blockNum <= latestBlockNum; blockNum++ {
			var blockNumToProcess = blockNum
			pool.AddDBTask(func(db *sql.DB) {
				block := BlockRawToBlock(GetBlockRawByNumber(blockNumToProcess))
				insertOrUpdateBlockToDB(block, db)
			})
		}
		startBlockNum = latestBlockNum + 1
	}
}

func insertOrUpdateBlockToDB(block Block, db *sql.DB) uint64 {
	insertSqlStatement := `INSERT INTO block (block_num, block_hash, block_time, parent_hash) VALUES ($1, $2, $3, $4) RETURNING block_num`
	var blockNum uint64
	err := db.QueryRow(insertSqlStatement, block.BlockNum, block.BlockHash, block.BlockTime, block.ParentHash).Scan(&blockNum)
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code.Name() == "unique_violation" {
			updateSqlStatement := `UPDATE block SET block_hash = $1, block_time = $2, parent_hash = $3 WHERE block_num = $4 RETURNING block_num`
			err = db.QueryRow(updateSqlStatement, block.BlockHash, block.BlockTime, block.ParentHash, block.BlockNum).Scan(&blockNum)
			fmt.Printf("Updated a single record %v\n", blockNum)
		} else {
			fmt.Printf("Unable to execute the query. %v\n", err)
		}
	} else {
		fmt.Printf("Inserted a single record %v\n", blockNum)
	}
	return blockNum
}

func createConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

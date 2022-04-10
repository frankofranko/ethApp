package utils

import "database/sql"

// T is a type alias to accept any type.
type T = interface{}

// WorkerPool is a contract for Worker Pool implementation
type WorkerPool interface {
	Run()
	AddDBTask(dbTask func(*sql.DB))
}

type workerPool struct {
	maxWorker     int
	queuedDBTaskC chan func(*sql.DB)
}

func (wp *workerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		//fmt.Println("2----")
		db := createConnection()
		go func(workerID int) {
			//fmt.Println("3----")
			for dbTask := range wp.queuedDBTaskC {
				//fmt.Println("4----")
				dbTask(db)
			}
		}(i + 1)
	}
}

func (wp *workerPool) AddDBTask(dbTask func(*sql.DB)) {
	wp.queuedDBTaskC <- dbTask
}

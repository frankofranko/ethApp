package utils

import "database/sql"

type T = interface{}

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
		db := createConnection()
		go func(workerID int) {
			for dbTask := range wp.queuedDBTaskC {
				dbTask(db)
			}
		}(i + 1)
	}
}

func (wp *workerPool) AddDBTask(dbTask func(*sql.DB)) {
	wp.queuedDBTaskC <- dbTask
}

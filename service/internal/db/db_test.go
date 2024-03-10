package db

import "testing"

func TestWorkerDB(t *testing.T) {
	workerdb := NewWorkerDB()

	workerdb.RegisterWorker("test")
	return
}

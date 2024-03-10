package db

import (
	postgrest_go "github.com/nedpals/postgrest-go/pkg"
	"github.com/nedpals/supabase-go"
	"github.com/rs/zerolog/log"
)

type WorkerDB struct {
	db *postgrest_go.Client
}

type WorkerRegistry struct {
	ID      string `json:"id"`
	MacAddr string `json:"mac_addr"`
}

func NewWorkerDB() WorkerDB {

	client := supabase.CreateClient("https://uqavovpznzxtraxxxmnh.supabase.co", "")

	return WorkerDB{db: client.DB}
}

func (r WorkerDB) getWorkerRegistriesTable() *postgrest_go.RequestBuilder {
	return r.db.From("workers")
}

func (r WorkerDB) RegisterWorker(macAddr string) (string, error) {
	if macAddr == "" {
		return "", ErrInvalidMacAddr
	}

	table := r.getWorkerRegistriesTable()

	var worker WorkerRegistry
	err := table.Select("*").Eq("test", "test").Execute(&worker)

	if err != nil {
		panic(err)
		log.Print(worker)
		return "", err
	}
	return "", err
}

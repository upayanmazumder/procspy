package store

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

func NewStore(connStr string) *Store {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return &Store{DB: db}
}

func (s *Store) RegisterMachine(name, osName, version string) (int, error) {
	var id int
	err := s.DB.QueryRow(
		`INSERT INTO machines(machine_name, os, agent_version) VALUES($1,$2,$3) RETURNING id`,
		name, osName, version).Scan(&id)
	return id, err
}

func (s *Store) SaveMetric(machineID int, cpu, ram float64) error {
	_, err := s.DB.Exec(
		`INSERT INTO metrics(machine_id, cpu_percent, ram_percent, collected_at) VALUES($1,$2,$3,$4)`,
		machineID, cpu, ram, time.Now())
	return err
}

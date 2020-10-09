package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Module struct {
	DB *sqlx.DB
	TX *sqlx.Tx
}

func New() *Module {
	return &Module{}
}

func (m *Module) InitConnection() {
	var err error
	m.DB, err = sqlx.Connect("postgres", "user=postgres"+" password=xyz"+" host=127.0.0.1"+" port=5432"+" dbname=benchmark"+" sslmode=disable")
	if err != nil {
		log.Fatal("Failed to init db connection, error:" + err.Error())
	}
	err = m.DB.Ping()
	if err != nil {
		log.Println("Failed to ping db", err.Error())
	}
}

func (m *Module) StartTransaction() (*sqlx.Tx, error) {
	tx, err := m.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("Failed to start database transaction. Error: %s", err.Error())
	}

	return tx, nil
}

func (m *Module) FinishTransaction(tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit database transaction. Error: %s", err.Error())
	}

	return nil
}

var BenchmarkStmt BenchmarkPreparedStatement

type BenchmarkPreparedStatement struct {
	InsertIntoBenchmark *sqlx.Stmt
	UpdateBenchmark     *sqlx.Stmt
}

func (m *Module) Preparex(query string) *sqlx.Stmt {
	statement, err := m.DB.Preparex(query)
	if err != nil {
		log.Fatal("Failed to prepared query")
	}
	return statement
}

var (
	insertQuery = "insert into benchmark(name, address, status) values($1, $2, $3)"
	updateQuery = "update benchmark set name=$1, address=$2, status=$3 where id=$4"
)

func (m *Module) InitPreparedStatement() {
	BenchmarkStmt = BenchmarkPreparedStatement{
		InsertIntoBenchmark: m.Preparex(insertQuery),
		UpdateBenchmark:     m.Preparex(updateQuery),
	}
}

type BenchmarkData struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Address   string    `db:"address"`
	Status    string    `db:"status"`
	CreatedOn time.Time `db:"created_on"`
	UpdatedOn time.Time `db:"updated_on"`
}

func (d BenchmarkData) InsertWithPreparedStatement(tx *sqlx.Tx) error {
	_, err := BenchmarkStmt.InsertIntoBenchmark.Exec(d.Name, d.Address, d.Status)
	if err != nil {
		return err
	}

	return nil
}

func (d BenchmarkData) InsertWithoutPreparedStatement(tx *sqlx.Tx) error {
	query := tx.Rebind(insertQuery)
	_, err := tx.Exec(query, d.Name, d.Address, d.Status)
	if err != nil {
		return err
	}
	return nil
}

func (d BenchmarkData) UpdateWithPreparedStatement(tx *sqlx.Tx) error {
	_, err := BenchmarkStmt.UpdateBenchmark.Exec(d.Name, d.Address, d.Status, d.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d BenchmarkData) UpdateWithoutPreparedStatement(tx *sqlx.Tx) error {
	query := tx.Rebind(updateQuery)
	_, err := tx.Exec(query, d.Name, d.Address, d.Status, d.ID)
	if err != nil {
		return err
	}
	return nil
}

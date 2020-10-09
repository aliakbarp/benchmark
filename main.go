package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	database "github.com/aliakbarp/benchmark/db"
)

func main() {
	modeStr := flag.String("mode", "", "for benchmarking")
	stmtStr := flag.Bool("stmt", false, "is use prepared statement?")
	flag.Parse()
	if *modeStr == "" || (*modeStr != "update" && *modeStr != "insert") {
		return
	}

	mode := *modeStr
	stmt := *stmtStr

	m := database.New()
	m.InitConnection()
	m.InitPreparedStatement()

	data := database.BenchmarkData{
		Name:    "Wahono",
		Address: "Priok",
		Status:  "Graduated",
	}

	var err error
	var function func(tx *sqlx.Tx) error
	if mode == "insert" {
		if stmt {
			log.Println("Benchmarking insert query with prepared statement")
			function = data.InsertWithPreparedStatement
		} else {
			log.Println("Benchmarking insert query without prepared statement")
			function = data.InsertWithoutPreparedStatement
		}
	} else if mode == "update" {
		if stmt {
			log.Println("Benchmarking update query with prepared statement")
			function = data.UpdateWithPreparedStatement
		} else {
			log.Println("Benchmarking update query without prepared statement")
			function = data.UpdateWithoutPreparedStatement
		}
	}

	then := time.Now()
	for i := int64(0); i < 2000; i++ {
		if mode == "update" {
			data.ID = i + 1
		}
		tx, _ := m.StartTransaction()
		err = function(tx)
		if err != nil {
			log.Println("Failed to insert, err: " + err.Error())
			break
		}
		m.FinishTransaction(tx)
	}
	diff := time.Now().Sub(then)
	fmt.Printf("%+v\n", diff)
}

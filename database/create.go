package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func init() {
	fmt.Println("main init")
}

func main() {

	url := "postgres://yiwovvbc:hpXxlhfd23PHotpuDQkHjmelSHBE3FBZ@rain.db.elephantsql.com/yiwovvbc"
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect err", err)
	}
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT );
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("err", err)
	}

	log.Println("okay")
}

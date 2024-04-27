package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect err", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO users (name, age) values ($1, $2) RETURNING id","Kitti", 24)
	var id int 
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("err")
	}
	fmt.Println("insert todo success id : ", id)
}
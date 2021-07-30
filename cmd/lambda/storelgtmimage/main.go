package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id   int
	name string
}

func Handler() {
	host := os.Getenv("DB_READER_HOSTNAME")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := user + ":" + password + "@tcp(" + host + ")/" + dbName
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`select id, name from user`)

	var userResult []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.id, &user.name); err != nil {
			log.Fatal(err)
		}
		userResult = append(userResult, user)
	}

	for _, u := range userResult {
		fmt.Println(u)
	}
}

func main() {
	lambda.Start(Handler)
}

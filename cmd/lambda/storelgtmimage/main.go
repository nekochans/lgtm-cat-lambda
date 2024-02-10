package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/nekochans/lgtm-cat-lambda/db/sqlc"
)

var q *db.Queries

func init() {
	host := os.Getenv("DB_HOSTNAME")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")

	var err error
	dataSourceName := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?tls=true"
	m, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	q = db.New(m)
}

func extractFilenameWithoutExt(objectKey string) string {
	base := filepath.Base(objectKey)

	return base[:len(base)-len(filepath.Ext(base))]
}

func Handler(event events.S3Event) {

	for _, record := range event.Records {
		objectKey := record.S3.Object.Key
		path := filepath.Dir(objectKey)
		filenameWithoutExt := extractFilenameWithoutExt(objectKey)

		// create LgtmImages
		ctx := context.Background()

		param := db.CreateLgtmImagesParams{
			Path:     path,
			Filename: filenameWithoutExt,
		}

		_, err := q.CreateLgtmImages(ctx, param)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	lambda.Start(Handler)
}

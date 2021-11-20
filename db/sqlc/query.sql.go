// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const createLgtmImages = `-- name: CreateLgtmImages :execresult
INSERT INTO lgtm_images (filename, path)
VALUES (?, ?)
`

type CreateLgtmImagesParams struct {
	Filename string
	Path     string
}

func (q *Queries) CreateLgtmImages(ctx context.Context, arg CreateLgtmImagesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createLgtmImages, arg.Filename, arg.Path)
}

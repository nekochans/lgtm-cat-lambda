-- name: CreateLgtmImages :execresult
INSERT INTO lgtm_images (filename, path)
VALUES (?, ?);

package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var downloader *manager.Downloader
var uploader *manager.Uploader
var destinationBucketName string

func init() {
	region := os.Getenv("REGION")
	destinationBucketName = os.Getenv("DESTINATION_BUCKET_NAME")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
	downloader = manager.NewDownloader(s3Client)
}

func downloadFromS3(
	ctx context.Context,
	downloader *manager.Downloader,
	bucket string,
	key string,
) (f *os.File, err error) {
	tmpFile, _ := os.CreateTemp("/tmp", "tmp_img_")

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
			log.Fatalln(err)
		}
	}()

	_, err = downloader.Download(
		ctx,
		tmpFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)

	if err != nil {
		return nil, err
	}

	return tmpFile, err
}

func uploadToS3(
	ctx context.Context,
	uploader *manager.Uploader,
	imgBytesBuffer *bytes.Buffer,
	bucket string,
	key string,
) error {
	contentType := decideS3ContentType(key)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Body:        imgBytesBuffer,
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	}

	_, err := uploader.Upload(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func extractImageExtension(fileName string) string {
	// 許可されている画像拡張子
	allowedImageExtList := [...]string{".jpg", ".jpeg", ".png"}

	ext := filepath.Ext(fileName)

	for _, v := range allowedImageExtList {
		if ext == v {
			return v
		}
	}

	return ""
}

func decideS3ContentType(s3Key string) string {
	ext := extractImageExtension(s3Key)

	contentType := ""

	switch ext {
	case ".png":
		contentType = "image/png"
	default:
		contentType = "image/jpeg"
	}

	return contentType
}

func Handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		img, err := downloadFromS3(ctx, downloader, bucket, key)
		if err != nil {
			return err
		}

		imgBuffer, err := io.ReadAll(img)
		if err != nil {
			return err
		}

		uploadKey := strings.ReplaceAll(key, "tmp/", "")

		imgBytesBuffer := new(bytes.Buffer)
		imgBytesBuffer.Write(imgBuffer)
		err = uploadToS3(ctx, uploader, imgBytesBuffer, destinationBucketName, uploadKey)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

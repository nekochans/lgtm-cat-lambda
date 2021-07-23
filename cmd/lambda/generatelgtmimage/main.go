package main

import (
	"bytes"
	"context"
	"image"
	"image/color"
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
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

const (
	lgtmFonts = 60
	eouFonts  = 30
	lgtmText  = "LGTM"
	meowText  = "eow"
)

var (
	downloader            *manager.Downloader
	uploader              *manager.Uploader
	destinationBucketName string
)

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
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Body:        imgBytesBuffer,
		ContentType: aws.String("image/png"),
		Key:         aws.String(key),
	}

	_, err := uploader.Upload(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func genLgtmImage(file *os.File) (b []byte, err error) {

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rct := img.Bounds()
	// Width > Height
	if rct.Dx() > rct.Dy() {
		// Resize the cropped image to width = 400px preserving the aspect ratio.
		img = imaging.Resize(img, 400, 0, imaging.Lanczos)
	} else {
		// Resize the cropped image to Height = 400px preserving the aspect ratio.
		img = imaging.Resize(img, 0, 400, imaging.Lanczos)
	}

	resizedRct := img.Bounds()
	dc := gg.NewContext(resizedRct.Dx(), resizedRct.Dy())
	dc.DrawImage(img, 0, 0)

	fontPath := filepath.Join("fonts", "MPLUSRounded1c-Medium.ttf")
	dc.SetColor(color.White)

	// for 'LGTM'
	if err := dc.LoadFontFace(fontPath, lgtmFonts); err != nil {
		return nil, err
	}
	textWidth, textHeight := dc.MeasureString(lgtmText)

	// for 'eow'
	if err := dc.LoadFontFace(fontPath, eouFonts); err != nil {
		return nil, err
	}
	eowTextWidth, _ := dc.MeasureString(meowText)

	x := (float64(dc.Width()) - (textWidth + eowTextWidth)) / 2
	y := (float64(dc.Height()) / 2) + (textHeight / 2)

	// for 'LGTM'
	if err := dc.LoadFontFace(fontPath, lgtmFonts); err != nil {
		return nil, err
	}
	dc.DrawString(lgtmText, x, y)

	// for 'eow'
	if err := dc.LoadFontFace(fontPath, eouFonts); err != nil {
		return nil, err
	}
	dc.DrawString(meowText, x+textWidth, y)

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func Handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		img, err := downloadFromS3(ctx, downloader, bucket, key)
		if err != nil {
			return err
		}

		imgByte, err := genLgtmImage(img)
		if err != nil {
			return err
		}
		imgBytesBuffer := new(bytes.Buffer)
		imgBytesBuffer.Write(imgByte)

		fileNameWithoutExt := getFileNameWithoutExt(strings.ReplaceAll(key, "tmp/", ""))
		uploadKey := fileNameWithoutExt + ".png"

		err = uploadToS3(ctx, uploader, imgBytesBuffer, destinationBucketName, uploadKey)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

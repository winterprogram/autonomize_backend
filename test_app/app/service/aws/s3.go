package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"test/test_app/app/service/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func (a *AWS) S3presignedUrl(ctx context.Context, id string) (url string, err error) {
	log := logger.Logger(ctx)

	bucket := s3.New(a.con)
	req, _ := bucket.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(a.BucketName),
		Key:    aws.String("id"),
	})

	url, err = req.Presign(60 * 150)
	if err != nil {
		log.Panic("failed to presign GetObjectRequest for key %q: %v", err)
		return
	}

	return
}

func (a *AWS) UploadFile(ctx context.Context, fileName string, file multipart.File) (url string, err error) {

	bucket := s3.New(a.con)

	buffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(buffer, file); err != nil {
		return
	}

	// create a unique file name for the file
	tempFileName := "attachments/" + bson.NewObjectId().Hex() + filepath.Ext(fileName)

	key := fmt.Sprintf("%s-%s", uuid.New().String(), tempFileName)
	typ := http.DetectContentType(buffer.Bytes())
	req := &s3.PutObjectInput{
		Body:   bytes.NewReader(buffer.Bytes()),
		Bucket: aws.String(a.BucketName),
		//ContentLength: &size,
		ContentType: &typ,
		Key:         aws.String(key),
	}

	_, err = bucket.PutObject(req)
	if err != nil {
		return
	}

	url = os.Getenv("APP_S3_URL") + key
	return
}

func (a *AWS) DeleteFile(ctx context.Context, fileName string, bucket string) (err error) {
	log := logger.Logger(ctx)
	svc := s3.New(a.con)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(fileName)})

	if err != nil {
		log.Panic("Unable to delete %q to %q, %v", fileName, bucket, err)
		return
	}
	return
}

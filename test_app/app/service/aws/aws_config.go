package aws

import (
	"context"
	"mime/multipart"
	"os"
	"test/test_app/app/service/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type IAwsService interface {
	S3presignedUrl(ctx context.Context, id string) (url string, err error)
	UploadFile(ctx context.Context, fileName string, file multipart.File) (url string, err error)
	DeleteFile(ctx context.Context, fileName string, bucket string) (err error)
}

type AWS struct {
	con        *session.Session
	BucketName string
}

// Init AWS
func InitAwsStr(ctx context.Context) (IAwsService, error) {
	log := logger.Logger(ctx)
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		log.Panic("Error at connecting to aws", err)
	}
	log.Info("Successfully connected to AWS SDK..")
	return &AWS{
		con:        sess,
		BucketName: os.Getenv("S3_BUCKET_NAME"),
	}, nil
}

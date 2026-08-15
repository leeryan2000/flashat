package db

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/leeryan2000/flashat/config"
)

// NewS3PresignClient builds an S3 presign client using the default AWS
// credential chain (an EC2 instance role in production, so no access keys
// are stored in .env).
func NewS3PresignClient(cfg config.Configuration) (*s3.PresignClient, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.AWS_REGION))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return s3.NewPresignClient(client), nil
}

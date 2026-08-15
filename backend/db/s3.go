package db

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/leeryan2000/flashat/config"
)

// NewS3Clients builds an S3 client and a presign client using the default
// AWS credential chain (an EC2 instance role in production, so no access
// keys are stored in .env). The plain client is used for direct
// backend-initiated calls (e.g. deleting an avatar); the presign client is
// used to hand out short-lived upload URLs to the browser.
func NewS3Clients(cfg config.Configuration) (*s3.Client, *s3.PresignClient, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.AWS_REGION))
	if err != nil {
		return nil, nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return client, s3.NewPresignClient(client), nil
}

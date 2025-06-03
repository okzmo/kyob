package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func GetAWSConfig() aws.Config {
	region := os.Getenv("AWS_REGION")
	keyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentialsProvider(keyId, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		panic(err)
	}

	return cfg
}
